package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Tikkaaa3/t-learn/api/internal/database"
	"github.com/google/uuid"
)

type Handler struct {
	DB *database.Queries
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	hash, err := HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := h.DB.CreateUser(r.Context(), database.CreateUserParams{
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: hash,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Username  string    `json:"username"`
	}
	safeUser := response{
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Email,
		user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(safeUser)
}
