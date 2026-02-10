package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Tikkaaa3/t-learn/api/internal/auth"
	"github.com/Tikkaaa3/t-learn/api/internal/database"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

const (
	BaseURL   = "http://localhost:8080"
	AdminUser = "admin"
	AdminPass = "admin"
)

// --- Helper Text to Append to Every Lesson ---
const cliHelper = `
---
### 🛠 CLI Setup Helper
**Install the t-cli if you don't have it yet:**

1. **Clone the repository:**
   ` + "```bash" + `
   git clone [https://github.com/Tikkaaa3/t-cli.git](https://github.com/Tikkaaa3/t-cli.git)
   cd t-cli
   ` + "```" + `

2. **Build and Install:**
   ` + "```bash" + `
   # Build the binary
   make build

   # Install it globally (requires sudo)
   sudo make install
   ` + "```" + `
`

// --- Seeding Data Structures ---

type TaskSeed struct {
	Description    string
	Command        string
	ExpectedOutput string
}

type LessonSeed struct {
	Title    string
	Content  string
	TaskSeed TaskSeed
}

type CourseSeed struct {
	Title       string
	Description string
	Lessons     []LessonSeed
}

func main() {
	// Setup Env & DB Connection
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env")
	}

	// Create Admin User directly in DB
	bootstrapAdmin()

	fmt.Println("Starting API-Based Seeder...")

	// Login to get the Token
	token := loginAndGetToken()
	fmt.Println("Logged in as Admin")

	// Define Curriculum
	curriculum := []CourseSeed{
		{
			Title:       "Python Basics",
			Description: "Start your journey with Python 3.",
			Lessons: []LessonSeed{
				{
					Title:   "Hello Python",
					Content: "# Hello World\nCreate a file named `main.py` that prints exactly `Hello Python`.",
					TaskSeed: TaskSeed{
						Description:    "Create main.py and print 'Hello Python'",
						Command:        "python3 main.py",
						ExpectedOutput: "Hello Python",
					},
				},
				{
					Title:   "Simple Math",
					Content: "# Variables\nCreate a file `math.py`. Print the result of `10 + 5`.",
					TaskSeed: TaskSeed{
						Description:    "Create math.py and print the integer 15",
						Command:        "python3 math.py",
						ExpectedOutput: "15",
					},
				},
				{
					Title:   "Loops",
					Content: "# Loops\nCreate `loop.py`. Write a loop that prints `Iteration` 3 times.",
					TaskSeed: TaskSeed{
						Description:    "Create loop.py that prints 'Iteration' on 3 separate lines",
						Command:        "python3 loop.py",
						ExpectedOutput: "Iteration\nIteration\nIteration",
					},
				},
			},
		},
		{
			Title:       "Go Basics",
			Description: "Master the fundamentals of Golang.",
			Lessons: []LessonSeed{
				{
					Title:   "Hello Go",
					Content: "# Go Setup\nCreate a `main.go` file with package main. Print `Hello Go`.",
					TaskSeed: TaskSeed{
						Description:    "Create main.go and print 'Hello Go'",
						Command:        "go run main.go",
						ExpectedOutput: "Hello Go",
					},
				},
				{
					Title:   "Integers",
					Content: "# Types\nCreate `nums.go`. Declare an integer `x = 42` and print it.",
					TaskSeed: TaskSeed{
						Description:    "Create nums.go and print 42",
						Command:        "go run nums.go",
						ExpectedOutput: "42",
					},
				},
				{
					Title:   "Functions",
					Content: "# Functions\nCreate `funcs.go`. Create a function `greet()` that prints `Greetings` and call it from main.",
					TaskSeed: TaskSeed{
						Description:    "Create funcs.go using a helper function",
						Command:        "go run funcs.go",
						ExpectedOutput: "Greetings",
					},
				},
			},
		},
		{
			Title:       "Rust Basics",
			Description: "Blazing fast memory safety.",
			Lessons: []LessonSeed{
				{
					Title:   "Hello Rust",
					Content: "# Rust Start\nCreate `main.rs`. Use `println!` to print `Hello Rust`.",
					TaskSeed: TaskSeed{
						Description:    "Create main.rs, compile it, and run it.",
						Command:        "rustc main.rs && ./main",
						ExpectedOutput: "Hello Rust",
					},
				},
				{
					Title:   "Variables",
					Content: "# Let\nCreate `vars.rs`. Define `let x = 100;` and print `{}`.",
					TaskSeed: TaskSeed{
						Description:    "Create vars.rs and print 100",
						Command:        "rustc vars.rs && ./vars",
						ExpectedOutput: "100",
					},
				},
				{
					Title:   "Mutability",
					Content: "# Mut\nCreate `mut.rs`. Define `let mut y = 1;`, change it to 2, and print it.",
					TaskSeed: TaskSeed{
						Description:    "Create mut.rs using mutable variables",
						Command:        "rustc mut.rs && ./mut",
						ExpectedOutput: "2",
					},
				},
			},
		},
	}

	// Loop and Seed
	for _, course := range curriculum {
		fmt.Printf("\n--- Seeding Course: %s ---\n", course.Title)
		courseID := createCourse(token, course.Title, course.Description)

		for i, lesson := range course.Lessons {
			// AUTOMATICALLY APPEND HELPER TEXT
			fullContent := lesson.Content + cliHelper

			// Position is i+1
			lessonID := createLesson(token, courseID, lesson.Title, fullContent, i+1)
			createTask(token, lessonID, lesson.TaskSeed)
		}
	}

	fmt.Println("\nSeeding Complete! Python, Go, and Rust courses are ready.")
}

// --- Helpers ---

func bootstrapAdmin() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	queries := database.New(conn)
	ctx := context.Background()

	_, err = queries.GetUserByUsername(ctx, AdminUser)
	if err != nil {
		hashedPassword, _ := auth.HashPassword(AdminPass)
		_, err := queries.CreateUser(ctx, database.CreateUserParams{
			Username:     AdminUser,
			Email:        "admin@t-learn.com",
			PasswordHash: hashedPassword,
		})
		if err != nil {
			log.Fatal("Failed to create admin user:", err)
		}
		fmt.Println("Created 'admin' user.")
	}

	_, err = conn.Exec("UPDATE users SET role = 'admin' WHERE username = $1", AdminUser)
	if err != nil {
		log.Fatal("Failed to promote user to admin:", err)
	}
}

func loginAndGetToken() string {
	payload := map[string]string{
		"username": AdminUser,
		"password": AdminPass,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(BaseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("Login failed (is server running?):", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("Login failed: %s", string(b))
	}

	var res struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Token
}

func createCourse(token, title, desc string) string {
	payload := map[string]string{
		"title":       title,
		"description": desc,
	}
	data, _ := json.Marshal(payload)

	var res struct {
		ID string `json:"id"`
	}
	doRequest("POST", "/admin/courses", token, data, &res)
	return res.ID
}

func createLesson(token, courseID, title, content string, position int) string {
	payload := map[string]interface{}{
		"title":    title,
		"content":  content,
		"position": position,
	}
	data, _ := json.Marshal(payload)

	var res struct {
		ID string `json:"id"`
	}
	path := fmt.Sprintf("/admin/courses/%s/lessons", courseID)
	doRequest("POST", path, token, data, &res)
	return res.ID
}

func createTask(token, lessonID string, task TaskSeed) {
	type Step struct {
		Command        string `json:"command"`
		ExpectedOutput string `json:"expected_output"`
		Position       int    `json:"position"`
	}

	payload := map[string]interface{}{
		"description": task.Description,
		"steps": []Step{
			{
				Position:       1,
				Command:        task.Command,
				ExpectedOutput: task.ExpectedOutput,
			},
		},
	}
	data, _ := json.Marshal(payload)

	path := fmt.Sprintf("/admin/lessons/%s/task", lessonID)
	doRequest("POST", path, token, data, nil)
}

func doRequest(method, path, token string, body []byte, target interface{}) {
	req, _ := http.NewRequest(method, BaseURL+path, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		respBody, _ := io.ReadAll(resp.Body)
		log.Fatalf("API Error [%s] %s: %s", method, path, string(respBody))
	}

	if target != nil {
		json.NewDecoder(resp.Body).Decode(target)
	}
}
