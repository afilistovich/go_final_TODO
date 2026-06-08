package server

import (
	"log"
	"net/http"
	"time"

	"github.com/afilistovich/go_final_TODO/internal/api"
)

type Server struct {
	server *http.Server
	logger *log.Logger
}

func NewServer(port string, webDir string, logger *log.Logger) *Server {
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorLog:     logger,
	}

	return &Server{
		server: srv,
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Printf("Starting server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}
