package main

import (
	"fmt"
	"goservice/pkg/loadenv"
	"log"
	"net/http"
)

const (
	port = 9000
)

// main is the application entry point. It loads environment variables,
// initializes logging, registers HTTP routes, and starts the server.
func main() {
	// Load key/value pairs from .env file into the environment
	if err := loadenv.LoadEnv(".env"); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	initLogger("log.out")

	// Register HTTP routes with logging middleware
	http.HandleFunc("/", middlewareFunc(indexHandler))
	http.HandleFunc("/code", middlewareFunc(githubHandler))
	http.HandleFunc("/status", middlewareFunc(statusHandler))
	http.HandleFunc("/blog", middlewareFunc(blogHandler))
	http.HandleFunc("/blog/", middlewareFunc(blogHandler)) // Matches /blog/{post-name}
	http.HandleFunc("/bookmarks", middlewareFunc(bookmarksHandler))
	http.HandleFunc("/resume", middlewareFunc(resumeHandler))

	log.Printf("Server starting on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
