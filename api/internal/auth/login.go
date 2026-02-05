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

	type response struct {
		Token string `json:"token"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response{Token: token})
}
