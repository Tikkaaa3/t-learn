package content

import (
	"encoding/json"
	"net/http"

	"github.com/Tikkaaa3/t-learn/api/internal/database"
	"github.com/google/uuid"
)

type CLIResponse struct {
	LessonID        string `json:"lesson_id"`
	LessonTitle     string `json:"lesson_title"`
	LessonContent   string `json:"lesson_content"`
	TaskID          string `json:"task_id"`
	TaskDescription string `json:"task_description"`
	Steps           []Step `json:"steps"`
}

type Step struct {
	Position       int32  `json:"position"`
	Command        string `json:"command"`
	ExpectedOutput string `json:"expected_output"`
}

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

func (h *Handler) GetLessons(w http.ResponseWriter, r *http.Request, user database.User) {
	courseIDStr := r.PathValue("course_id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	lessons, err := h.DB.GetLessonsWithStatus(r.Context(), database.GetLessonsWithStatusParams{
		CourseID: courseID,
		UserID:   user.ID, // <--- We pass the logged-in user's ID here
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Map response to JSON
	type LessonResponse struct {
		ID        uuid.UUID `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"` // <--- New JSON field
	}

	response := make([]LessonResponse, len(lessons))
	for i, l := range lessons {
		response[i] = LessonResponse{
			ID:        l.ID,
			Title:     l.Title,
			Completed: l.IsCompleted, // Map the boolean from SQL
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	lessonIDStr := r.PathValue("lesson_id")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	// Fetch the Lesson (Title & Content)
	lesson, err := h.DB.GetLesson(r.Context(), lessonID)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(`{"error": "Lesson not found"}`))
		return
	}

	// Fetch the Task
	task, err := h.DB.GetTaskByLessonID(r.Context(), lessonID)
	if err != nil {
		// It's okay if a lesson has no task (maybe it's just reading)
		// But for now, we return 404
		w.WriteHeader(404)
		w.Write([]byte(`{"error": "Task not found"}`))
		return
	}

	// Fetch the Steps
	steps, err := h.DB.GetStepsByTaskID(r.Context(), task.ID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Build the Response
	jsonSteps := []Step{}
	for _, s := range steps {
		jsonSteps = append(jsonSteps, Step{
			Position:       s.Position,
			Command:        s.Command,
			ExpectedOutput: s.ExpectedOutput,
		})
	}

	response := CLIResponse{
		LessonID:        lesson.ID.String(),
		LessonTitle:     lesson.Title,
		LessonContent:   lesson.Content, // The Markdown content
		TaskID:          task.ID.String(),
		TaskDescription: task.Description, // "Create a hello.go file..."
		Steps:           jsonSteps,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) CompleteTask(w http.ResponseWriter, r *http.Request, user database.User) {
	taskIDStr := r.PathValue("task_id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	err = h.DB.CompleteTask(r.Context(), database.CompleteTaskParams{
		UserID: user.ID,
		TaskID: taskID,
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(`{"status":"success"}`))
}

// Admin

func (h *Handler) CreateCourse(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(400)
		return
	}

	course, err := h.DB.CreateCourse(r.Context(), database.CreateCourseParams{
		Title:       params.Title,
		Description: params.Description,
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

func (h *Handler) CreateLesson(w http.ResponseWriter, r *http.Request, user database.User) {
	courseIDStr := r.PathValue("course_id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	type parameters struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Position int32  `json:"position"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(400)
		return
	}

	lesson, err := h.DB.CreateLesson(r.Context(), database.CreateLessonParams{
		CourseID: courseID,
		Title:    params.Title,
		Content:  params.Content,
		Position: params.Position,
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lesson)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request, user database.User) {
	lessonIDStr := r.PathValue("lesson_id")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	type StepRequest struct {
		Command        string `json:"command"`
		ExpectedOutput string `json:"expected_output"`
		Position       int32  `json:"position"`
	}
	type TaskRequest struct {
		Description string        `json:"description"`
		Steps       []StepRequest `json:"steps"`
	}

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(400)
		return
	}

	task, err := h.DB.CreateTask(r.Context(), database.CreateTaskParams{
		LessonID:    lessonID,
		Description: req.Description,
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}

	for _, step := range req.Steps {
		_, err := h.DB.CreateTaskStep(r.Context(), database.CreateTaskStepParams{
			TaskID:         task.ID,
			Position:       step.Position,
			Command:        step.Command,
			ExpectedOutput: step.ExpectedOutput,
		})
		if err != nil {
			w.WriteHeader(500)
			return
		}
	}

	w.WriteHeader(201)
	w.Write([]byte(`{"status":"created", "task_id":"` + task.ID.String() + `"}`))
}

func (h *Handler) DeleteCourse(w http.ResponseWriter, r *http.Request, user database.User) {
	id, err := uuid.Parse(r.PathValue("course_id"))
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if err := h.DB.DeleteCourse(r.Context(), id); err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204) // 204 No Content (Success)
}

func (h *Handler) DeleteLesson(w http.ResponseWriter, r *http.Request, user database.User) {
	id, err := uuid.Parse(r.PathValue("lesson_id"))
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if err := h.DB.DeleteLesson(r.Context(), id); err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request, user database.User) {
	id, err := uuid.Parse(r.PathValue("task_id"))
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if err := h.DB.DeleteTask(r.Context(), id); err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}
