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
		// ==========================================
		// PYTHON BASICS
		// ==========================================
		{
			Title:       "Python Basics",
			Description: "Start your journey with Python 3. Learn syntax, variables, and loops.",
			Lessons: []LessonSeed{
				{
					Title: "Hello Python",
					Content: `
### The Print Function
In Python, the ` + "`print()`" + ` function is used to send data to the standard output (your terminal). It is one of the most fundamental tools for debugging and user interaction.

#### Syntax
` + "```python" + `
print("Your message here")
` + "```" + `

Strings in Python can be enclosed in either single quotes ('...') or double quotes ("...").

#### Example
` + "```python" + `
print("Hello World")
print('Python is fun')
` + "```" + `
`,
					TaskSeed: TaskSeed{
						Description:    "Create a file named `main.py`. Inside it, write a script that outputs exactly `Hello Python`.",
						Command:        "python3 main.py",
						ExpectedOutput: "Hello Python",
					},
				},
				{
					Title: "Variables & Math",
					Content: `
### Variables
A variable is a container for storing data values. Unlike other languages, Python has no command for declaring a variable. A variable is created the moment you first assign a value to it.

#### Example
` + "```python" + `
x = 5
y = "John"
print(x)
` + "```" + `

### Arithmetic Operators
Python supports standard math operators:
* ` + "`+`" + ` Addition
* ` + "`-`" + ` Subtraction
* ` + "`*`" + ` Multiplication
* ` + "`/`" + ` Division

`,
					TaskSeed: TaskSeed{
						Description:    "Create a file named `math.py`. Calculate `10 + 5` and print the result (it should be 15).",
						Command:        "python3 math.py",
						ExpectedOutput: "15",
					},
				},
				{
					Title: "For Loops",
					Content: `
### The For Loop
A ` + "`for`" + ` loop is used for iterating over a sequence (that is either a list, a tuple, a dictionary, a set, or a string).

To loop a specific number of times, we can use the ` + "`range()`" + ` function.

#### Syntax
` + "```python" + `
for i in range(5):
    print(i)
` + "```" + `

*Note: ` + "`range(5)`" + ` generates numbers from 0 to 4 (5 is exclusive).*
`,
					TaskSeed: TaskSeed{
						Description:    "Create `loop.py`. Write a loop that prints the word `Iteration` exactly 3 times (on 3 separate lines).",
						Command:        "python3 loop.py",
						ExpectedOutput: "Iteration\nIteration\nIteration",
					},
				},
			},
		},

		// ==========================================
		// GO BASICS
		// ==========================================
		{
			Title:       "Go Basics",
			Description: "Master the fundamentals of Golang: Static typing, packages, and compilation.",
			Lessons: []LessonSeed{
				{
					Title: "Hello Go",
					Content: `
### Structure of a Go Program
Every runnable Go program starts with a package declaration.

1.  **Package Declaration:** ` + "`package main`" + ` tells the Go compiler that this file should compile as an executable program rather than a shared library.
2.  **Import:** ` + "`import \"fmt\"`" + ` brings in the formatting package (contains ` + "`Println`" + `).
3.  **Main Function:** ` + "`func main() { ... }`" + ` is the entry point where the program starts running.

#### Example
` + "```go" + `
package main
import "fmt"

func main() {
    fmt.Println("Hi there!")
}
` + "```" + `
`,
					TaskSeed: TaskSeed{
						Description:    "Create `main.go`. Write a valid Go program that prints `Hello Go` to the console.",
						Command:        "go run main.go",
						ExpectedOutput: "Hello Go",
					},
				},
				{
					Title: "Integers & Variables",
					Content: `
### Declaring Variables
Go is statically typed. You can declare variables in two ways:

**1. Long Syntax (var keyword)**
` + "```go" + `
var x int = 10
` + "```" + `

**2. Short Syntax (Type Inference)**
Inside functions, you can use ` + "`:=`" + ` to let Go guess the type.
` + "```go" + `
y := 20 // Go knows this is an int
` + "```" + `
`,
					TaskSeed: TaskSeed{
						Description:    "Create `nums.go`. Declare an integer variable equal to `42` and print it using `fmt.Println`.",
						Command:        "go run nums.go",
						ExpectedOutput: "42",
					},
				},
				{
					Title: "Functions",
					Content: `
### Functions
A function is a block of code which only runs when it is called. You can pass data, known as parameters, into a function.

#### Syntax
` + "```go" + `
func functionName() {
  // code to be executed
}
` + "```" + `

You must call the function from ` + "`main`" + ` for it to run.
`,
					TaskSeed: TaskSeed{
						Description:    "Create `funcs.go`. Define a function named `greet()` that prints `Greetings`. Call it from the main function.",
						Command:        "go run funcs.go",
						ExpectedOutput: "Greetings",
					},
				},
			},
		},

		// ==========================================
		// RUST BASICS
		// ==========================================
		{
			Title:       "Rust Basics",
			Description: "Learn memory safety and modern systems programming with Rust.",
			Lessons: []LessonSeed{
				{
					Title: "Hello Rust",
					Content: `
### Macros vs Functions
In Rust, ` + "`main`" + ` is the entry point.
You print to the screen using ` + "`println!`" + `.

Notice the **exclamation mark** (` + "`!`" + `)? This means it is a **Macro**, not a standard function. Macros write code for you at compile time.

#### Example
` + "```rust" + `
fn main() {
    println!("Hello World");
}
` + "```" + `
`,
					TaskSeed: TaskSeed{
						Description:    "Create `main.rs`. Write a program that prints `Hello Rust`.",
						Command:        "rustc main.rs && ./main",
						ExpectedOutput: "Hello Rust",
					},
				},
				{
					Title: "Variables (Let)",
					Content: `
### Variable Binding
In Rust, we use the ` + "`let`" + ` keyword to bind a value to a name.

` + "```rust" + `
fn main() {
    let x = 5;
    println!("{}", x); 
}
` + "```" + `

**Note:** The ` + "`{}`" + ` is a placeholder. Rust replaces it with the value of ` + "`x`" + `.
`,
					TaskSeed: TaskSeed{
						Description:    "Create `vars.rs`. Bind the value `100` to a variable named `x` and print it.",
						Command:        "rustc vars.rs && ./vars",
						ExpectedOutput: "100",
					},
				},
				{
					Title: "Mutability",
					Content: `
### Immutable by Default
In Rust, variables are **immutable** by default. Once a value is bound to a name, you can't change it.

To make a variable changeable, you must add ` + "`mut`" + `.

#### Example
` + "```rust" + `
let mut x = 5;
x = 6; // This is allowed because of 'mut'
` + "```" + `
`,
					TaskSeed: TaskSeed{
						Description:    "Create `mut.rs`. Define a mutable variable `y` starting at `1`. Change it to `2`, then print the final value.",
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
