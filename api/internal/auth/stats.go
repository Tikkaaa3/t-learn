package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Tikkaaa3/t-learn/api/internal/database"
)

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request, user database.User) {
	stats, err := h.DB.GetUserStats(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to get user stats", http.StatusInternalServerError)
		return
	}

	type statsResponse struct {
		Username       string `json:"username"`
		CompletedTasks int64  `json:"completed_tasks"`
	}

	response := statsResponse{
		Username:       stats.Username,
		CompletedTasks: stats.CompletedTasks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
