package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Tikkaaa3/t-learn/api/internal/database"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer conn.Close()

	queries := database.New(conn)
	ctx := context.Background()

	fmt.Println("Seeding database...")

	//  Create a Course
	course, err := queries.CreateCourse(ctx, database.CreateCourseParams{
		Title:       "Go for Beginners",
		Description: "A gentle introduction to the Go programming language.",
	})
	if err != nil {
		log.Fatal("Failed to create course:", err)
	}
	fmt.Printf("Created Course: %s (%s)\n", course.Title, course.ID)

	// Create a Lesson
	lesson, err := queries.CreateLesson(ctx, database.CreateLessonParams{
		CourseID: course.ID,
		Title:    "Hello World",
		Content:  "# Hello World\n\nWelcome to your first Go program. We will print text to the console.",
		Position: 1,
	})
	if err != nil {
		log.Fatal("Failed to create lesson:", err)
	}
	fmt.Printf("Created Lesson: %s (%s)\n", lesson.Title, lesson.ID)

	// Create a Task for that Lesson
	task, err := queries.CreateTask(ctx, database.CreateTaskParams{
		LessonID:       lesson.ID,
		Description:    "Create a file named 'main.go' that prints 'Hello, World!'.",
		ExpectedOutput: "Hello, World!",
		Command:        "go run main.go",
	})
	if err != nil {
		log.Fatal("Failed to create task:", err)
	}
	fmt.Printf("Created Task for Lesson: %s\n", task.Description)

	fmt.Println("Seeding complete!")
}
