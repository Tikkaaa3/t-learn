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
	AdminUser = "admin_seeder"
	AdminPass = "admin123"
)

func main() {
	// Setup Env & DB Connection
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env")
	}

	bootstrapAdmin()

	fmt.Println("Starting API-Based Seeder...")

	// Login to get the JWT/API Key
	token := loginAndGetToken()
	fmt.Println("Logged in as Admin")

	// Create Course
	courseID := createCourse(token)

	// Create Lesson
	lessonID := createLesson(token, courseID)

	// Create Task (With Steps)
	createTask(token, lessonID)

	fmt.Println("\nSeeding Complete via API! The backend is fully functional.")
}

// --- Bootstrap Helper (Raw SQL) ---

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
		// Create him if missing
		hashedPassword, _ := auth.HashPassword(AdminPass)
		_, err := queries.CreateUser(ctx, database.CreateUserParams{
			Username:     AdminUser,
			Email:        "admin@t-learn.com",
			PasswordHash: hashedPassword,
		})
		if err != nil {
			log.Fatal("Failed to create admin user:", err)
		}
		fmt.Println("Created 'admin_seeder' user.")
	}

	_, err = conn.Exec("UPDATE users SET role = 'admin' WHERE username = $1", AdminUser)
	if err != nil {
		log.Fatal("Failed to promote user to admin:", err)
	}
	fmt.Println("Promoted 'admin_seeder' to Admin Role.")
}

// --- API Client Helpers ---

func loginAndGetToken() string {
	payload := map[string]string{
		"username": AdminUser,
		"password": AdminPass,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(BaseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("Login failed:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Login returned non-200 status. Is the server running?")
	}

	var res struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Token
}

func createCourse(token string) string {
	payload := map[string]string{
		"title":       "Go HTTP Mastery",
		"description": "Learn how to build APIs with Go Standard Library.",
	}
	data, _ := json.Marshal(payload)

	var res struct {
		ID string `json:"id"`
	}
	doRequest("POST", "/admin/courses", token, data, &res)

	fmt.Printf("Created Course: %s\n", res.ID)
	return res.ID
}

func createLesson(token, courseID string) string {
	payload := map[string]interface{}{
		"title":    "The Handler Interface",
		"content":  "# Handlers\n\nEverything in Go is a handler...",
		"position": 1,
	}
	data, _ := json.Marshal(payload)

	var res struct {
		ID string `json:"id"`
	}
	path := fmt.Sprintf("/admin/courses/%s/lessons", courseID)
	doRequest("POST", path, token, data, &res)

	fmt.Printf("Created Lesson: %s\n", res.ID)
	return res.ID
}

func createTask(token, lessonID string) {
	type Step struct {
		Command        string `json:"command"`
		ExpectedOutput string `json:"expected_output"`
		Position       int    `json:"position"`
	}

	payload := map[string]interface{}{
		"description": "Create a hello.go file and run it.",
		"steps": []Step{
			{
				Position:       1,
				Command:        "echo 'package main; import \"fmt\"; func main() { fmt.Println(\"API Test\") }' > hello.go",
				ExpectedOutput: "",
			},
			{
				Position:       2,
				Command:        "go run hello.go",
				ExpectedOutput: "API Test",
			},
		},
	}
	data, _ := json.Marshal(payload)

	path := fmt.Sprintf("/admin/lessons/%s/task", lessonID)
	doRequest("POST", path, token, data, nil)

	fmt.Println("Created Multi-Step Task")
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
		log.Fatalf("API Error [%s]: %s %s", resp.Status, path, string(respBody))
	}

	if target != nil {
		json.NewDecoder(resp.Body).Decode(target)
	}
}
