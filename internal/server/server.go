package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/afilistovich/go_final_TODO/internal/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	server *http.Server
}

func NewServer(port string, webDir string) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	api.RegisterRoutes(r)

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		server: srv,
	}
}

func (s *Server) Start() error {
	slog.Info("Server is listening", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}
