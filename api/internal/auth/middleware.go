package auth

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/Tikkaaa3/t-learn/api/internal/database"
)

type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

func (h *Handler) MiddlewareAuth(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Try to look up User by API Key first (for CLI)
		user, err := h.DB.GetUserByAPIKey(r.Context(), sql.NullString{String: tokenString, Valid: true})
		if err == nil {
			handler(w, r, user)
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "my-secret-key"
		}

		// Validate JWT
		userID, err := ValidateJWT(tokenString, jwtSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err = h.DB.GetUserByID(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler(w, r, user)
	}
}

func (h *Handler) MiddlewareAdmin(handler AuthedHandler) http.HandlerFunc {
	return h.MiddlewareAuth(func(w http.ResponseWriter, r *http.Request, user database.User) {
		if user.Role != "admin" {
			w.WriteHeader(http.StatusForbidden) // 403 Forbidden
			w.Write([]byte("Access denied: Admins only"))
			return
		}
		handler(w, r, user)
	})
}
