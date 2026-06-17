package api

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Get("/api/nextdate", nextDateHandler)

	r.Get("/api/tasks", getTasksHandler)

	r.Get("/api/task", getTaskHandler)
	r.Post("/api/task", createTaskHandler)
	r.Put("/api/task", updateTaskHandler)

}
