package content

import (
	"encoding/json"
	"net/http"

	"github.com/Tikkaaa3/t-learn/api/internal/database"
	"github.com/google/uuid"
)

type Handler struct {
	DB *database.Queries
}

func (h *Handler) GetCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.DB.GetCourses(r.Context())
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func (h *Handler) GetLessons(w http.ResponseWriter, r *http.Request) {
	courseIDStr := r.PathValue("course_id")

	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid UUID format"))
		return
	}

	lessons, err := h.DB.GetLessonsByCourseID(r.Context(), courseID)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lessons)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	lessonIDStr := r.PathValue("lesson_id")

	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid UUID format"))
		return
	}

	task, err := h.DB.GetTaskByLessonID(r.Context(), lessonID)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
