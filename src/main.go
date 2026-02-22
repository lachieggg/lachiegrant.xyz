package main

import (
	"fmt"
	"goservice/pkg/loadenv"
	"log"
	"net/http"
)

const (
	port    = 9000
	logPath = "log.out"
)

// main is the application entry point. It loads environment variables,
// initializes logging, registers HTTP routes, and starts the server.
func main() {
	// Load key/value pairs from .env file into the environment
	if err := loadenv.LoadEnv(".env"); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	initLogger(logPath)

	http.HandleFunc("/", middlewareFunc(indexHandler))
	http.HandleFunc("/code", middlewareFunc(githubHandler))
	http.HandleFunc("/blog", middlewareFunc(blogHandler))
	http.HandleFunc("/bookmarks", middlewareFunc(bookmarksHandler))
	http.HandleFunc("/resume", middlewareFunc(resumeHandler))
	http.HandleFunc("/blog/", middlewareFunc(blogHandler)) // Matches /blog/{post-name}

	logger.Printf("Server starting on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
