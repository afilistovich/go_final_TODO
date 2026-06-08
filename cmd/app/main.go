package main

import (
	"log"
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

	dbPath := os.Getenv(envDBFile)
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	err := db.Init(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	port := os.Getenv(envPort)
	if port == "" {
		port = "7540"
	}

	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lshortfile)

	srv := server.NewServer(port, webDir, logger)

	if err = srv.Start(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
