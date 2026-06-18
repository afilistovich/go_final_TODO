package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/afilistovich/go_final_TODO/internal/calc"
	"github.com/afilistovich/go_final_TODO/internal/db"
)

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")
	var parsedDate string
	if t, parseErr := time.Parse("02.01.2006", search); parseErr == nil {
		parsedDate = t.Format(calc.DateLayout)
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
