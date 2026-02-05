package auth

import (
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

		// TODO: Load secret from env
		userID, err := ValidateJWT(tokenString, os.Getenv("JWT_SECRET"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := h.DB.GetUserByID(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler(w, r, user)
	}
}
