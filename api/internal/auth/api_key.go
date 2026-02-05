package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/Tikkaaa3/t-learn/api/internal/database"
)

func generateRandomKey() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (h *Handler) GenerateAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	newKey, err := generateRandomKey()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, err = h.DB.UpdateAPIKey(r.Context(), database.UpdateAPIKeyParams{
		ID: user.ID,
		ApiKey: sql.NullString{
			String: newKey,
			Valid:  true,
		},
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"api_key": newKey})
}
