package api

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all API routes on the given router
func RegisterRoutes(r chi.Router) {

	// Public routes (no authentication required)
	r.Post("/api/signin", signInHandler)
	r.Get("/api/nextdate", nextDateHandler)

	// Protected routes (require JWT authentication)
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Get("/api/tasks", getTasksHandler)

		r.Get("/api/task", getTaskHandler)
		r.Post("/api/task", createTaskHandler)
		r.Put("/api/task", updateTaskHandler)
		r.Delete("/api/task", deleteTaskHandler)
		r.Post("/api/task/done", doneTaskHandler)
	})
}
