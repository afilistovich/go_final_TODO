package main

import (
	"log"
	"net/http"
	"os"

	"github.com/afilistovich/go_final_TODO/internal/db"
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

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
