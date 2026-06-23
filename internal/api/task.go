package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/afilistovich/go_final_TODO/internal/db"
)

// createTaskHandler handles POST /api/task - creates new task
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

// getTaskHandler handles GET /api/task?id=X - returns task by ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		slog.Warn("Failed to parse id", "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, "Task not found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to get task from database", "error", err, "id", id)
		writeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err = writeJSON(w, http.StatusOK, task); err != nil {
		slog.Error("Failed to encode task response", "error", err)
	}
}

// updateTaskHandler handles PUT /api/task - updates existing task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		slog.Warn("Invalid JSON update request", "error", err)
		writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if task.ID <= 0 {
		writeError(w, "id must be a positive number", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "Title is required", http.StatusBadRequest)
		return
	}

	if err = normalizeTaskDate(&task); err != nil {
		slog.Warn("Invalid date in update request", "error", err, "date", task.Date)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = db.UpdateTask(&task); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, "Task not found", http.StatusNotFound)
			return
		}

		slog.Error("Database error while updating task", "error", err, "id", task.ID)
		writeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Info("Task updated successfully", "task_id", task.ID, "title", task.Title)

	if err = writeJSON(w, http.StatusOK, struct{}{}); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
	}
}

// doneTaskHandler handles POST /api/task/done?id=X - marks task as done
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		slog.Warn("Failed to parse id", "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = db.DoneTask(id); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, "Task not found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to mark task as done", "error", err, "id", id)
		writeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Info("Task marked as done", "task_id", id)
	writeJSON(w, http.StatusOK, struct{}{})
}

// deleteTaskHandler handles DELETE /api/task?id=X - deletes task
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		slog.Warn("Failed to parse id", "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = db.DeleteTask(id); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, "Task not found", http.StatusNotFound)
			return
		}
		slog.Error("Database error while deleting task", "error", err, "id", id)
		writeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Info("Task deleted successfully", "task_id", id)
	writeJSON(w, http.StatusOK, struct{}{})
}
