package api

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Get("/api/nextdate", nextDateHandler)

	r.Get("/api/tasks", getTasksHandler)

	r.Post("/api/task", createTaskHandler)

}
