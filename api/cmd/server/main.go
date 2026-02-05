package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Tikkaaa3/t-learn/api/internal/auth"
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

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(http.ResponseWriter, *http.Request) {
		//
	})
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	// Protected Route
	mux.HandleFunc("GET /auth/me", authHandler.MiddlewareAuth(func(w http.ResponseWriter, r *http.Request, user database.User) {
		w.Write([]byte("Hello, " + user.Username))
	}))
	mux.HandleFunc("POST /auth/token", authHandler.MiddlewareAuth(authHandler.GenerateAPIKey))
	http.ListenAndServe(":8080", mux)
}
