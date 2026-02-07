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
	// Setup Environment & Database
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

	// Create a Course
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

	// Create the Task (The Container)
	task, err := queries.CreateTask(ctx, database.CreateTaskParams{
		LessonID:    lesson.ID,
		Description: "Create a file named 'main.go' that prints 'Hello, World!' and run it.",
	})
	if err != nil {
		log.Fatal("Failed to create task:", err)
	}
	fmt.Printf("Created Task Container for Lesson: %s\n", lesson.Title)

	// Add Steps to the Task

	// Step 1: Create the file
	// We use a simple echo command to simulate the user writing code
	_, err = queries.CreateTaskStep(ctx, database.CreateTaskStepParams{
		TaskID:         task.ID,
		Position:       1,
		Command:        "echo 'package main; import \"fmt\"; func main() { fmt.Println(\"Hello, World!\") }' > main.go",
		ExpectedOutput: "", // Creating a file produces no output
	})
	if err != nil {
		log.Fatal("Failed to create step 1:", err)
	}
	fmt.Println("   🔹 Added Step 1: Create main.go")

	// Step 2: Run the file
	_, err = queries.CreateTaskStep(ctx, database.CreateTaskStepParams{
		TaskID:         task.ID,
		Position:       2,
		Command:        "go run main.go",
		ExpectedOutput: "Hello, World!",
	})
	if err != nil {
		log.Fatal("Failed to create step 2:", err)
	}
	fmt.Println("   🔹 Added Step 2: Run main.go")

	fmt.Println("Seeding complete!")
}
