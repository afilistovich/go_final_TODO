package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/afilistovich/go_final_TODO/internal/db"
)

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var id int64

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		// if JSON not valid or nil, return 400 Bad Request
		slog.Warn("Invalid JSON in request", "error", err)
		writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		slog.Warn("Task creation failed: empty title")
		writeError(w, "Title is required", http.StatusBadRequest)
		return
	}

	if err = normalizeTaskDate(&task); err != nil {
		slog.Warn("Invalid date in request", "error", err, "date", task.Date)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err = db.AddTask(&task)
	if err != nil {
		slog.Error("Database error while adding task", "error", err, "task_title", task.Title)
		writeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Info("Task created successfully", "task_id", id, "title", task.Title)

	err = writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
	if err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
	}
}

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")
	var parsedDate string
	if t, parseErr := time.Parse("02.01.2006", search); parseErr == nil {
		parsedDate = t.Format(DateLayout)
	}

	var tasks []*db.Task
	var err error

	switch {
	case search == "":
		tasks, err = db.Tasks(50)

	case parsedDate != "":
		tasks, err = db.TasksByDate(parsedDate, 50)

	default:
		tasks, err = db.TasksBySearch(search, 50)
	}
	if err != nil {
		slog.Error("Failed to get tasks", "error", err, "search", search)
		writeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Info("Tasks retrieved successfully", "count", len(tasks), "search", search)

	if err = writeJSON(w, http.StatusOK, TaskResp{Tasks: tasks}); err != nil {
		slog.Error("Failed to encode tasks response", "error", err)
	}
}
