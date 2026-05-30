package main

import (
	"net/http"
	"os"
)

const (
	envPort = "TODO_PORT"
	webDir  = "./web"
)

func main() {

	port := os.Getenv(envPort)
	if port == "" {
		port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.ListenAndServe(":"+port, nil)
}
