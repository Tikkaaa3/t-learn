package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Tikkaaa3/t-learn/api/internal/auth"
	"github.com/Tikkaaa3/t-learn/api/internal/content"
	"github.com/Tikkaaa3/t-learn/api/internal/database"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found in .env")
	}

	dbConn, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	dbQueries := database.New(dbConn)

	authHandler := &auth.Handler{
		DB: dbQueries,
	}

	contentHandler := &content.Handler{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(http.ResponseWriter, *http.Request) {
		//
	})

	// Auth Routes
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("GET /auth/me", authHandler.MiddlewareAuth(func(w http.ResponseWriter, r *http.Request, user database.User) {
		w.Write([]byte("Hello, " + user.Username))
	}))
	mux.HandleFunc("POST /auth/token", authHandler.MiddlewareAuth(authHandler.GenerateAPIKey))

	// Content Routes
	mux.HandleFunc("GET /courses", contentHandler.GetCourses)
	mux.HandleFunc("GET /courses/{course_id}/lessons", authHandler.MiddlewareAuth(contentHandler.GetLessons))
	mux.HandleFunc("GET /lessons/{lesson_id}/task", contentHandler.GetTask)
	mux.HandleFunc("POST /tasks/{task_id}/complete", authHandler.MiddlewareAuth(contentHandler.CompleteTask))

	// Admin Routes
	mux.HandleFunc("POST /admin/courses", authHandler.MiddlewareAdmin(contentHandler.CreateCourse))
	mux.HandleFunc("POST /admin/courses/{course_id}/lessons", authHandler.MiddlewareAdmin(contentHandler.CreateLesson))
	mux.HandleFunc("POST /admin/lessons/{lesson_id}/task", authHandler.MiddlewareAdmin(contentHandler.CreateTask))

	mux.HandleFunc("DELETE /admin/courses/{course_id}", authHandler.MiddlewareAdmin(contentHandler.DeleteCourse))
	mux.HandleFunc("DELETE /admin/lessons/{lesson_id}", authHandler.MiddlewareAdmin(contentHandler.DeleteLesson))
	mux.HandleFunc("DELETE /admin/tasks/{task_id}", authHandler.MiddlewareAdmin(contentHandler.DeleteTask))

	log.Println("Server starting on :8080")
	err = http.ListenAndServe(":8080", enableCORS(mux))
	if err != nil {
		log.Fatal(err)
	}
}

// enableCORS adds headers to allow the React frontend to communicate with this server
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from specific frontend origin
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Vite default port

		// Allow specific methods (GET, POST, etc.)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow specific headers (Content-Type for JSON, Authorization for Tokens)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle "Preflight" requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass the request to the real handler
		next.ServeHTTP(w, r)
	})
}
