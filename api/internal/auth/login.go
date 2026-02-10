package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	user, err := h.DB.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	match := CheckPassword(params.Password, user.PasswordHash)
	if !match {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "my-secret-key"
	}

	token, err := MakeJWT(user.ID, secret)
	if err != nil {
		log.Printf("Error generating token: %s", err)
		w.WriteHeader(500)
		return
	}

	type UserResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	type loginResponse struct {
		Token string       `json:"token"`
		User  UserResponse `json:"user"`
	}

	// Send back Token + User Info
	response := loginResponse{
		Token: token,
		User: UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}
