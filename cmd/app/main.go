package main

import (
	"log/slog"
	"os"

	"github.com/afilistovich/go_final_TODO/internal/db"
	"github.com/afilistovich/go_final_TODO/internal/server"
)

const (
	envPort   = "TODO_PORT"
	webDir    = "./web"
	envDBFile = "TODO_DBFILE"
)

func main() {

	// Configure structured logger
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))

	// Get database path from environment or use default
	dbPath := os.Getenv(envDBFile)
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	// Initialize database connection
	err := db.Init(dbPath)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	// Close database on exit
	defer db.Close()

	// Get port from environment or use default
	port := os.Getenv(envPort)
	if port == "" {
		port = "7540"
	}

	// Create server with routes and middleware
	srv := server.NewServer(port, webDir)

	slog.Info("Starting server", "port", port)
	if err = srv.Start(); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
