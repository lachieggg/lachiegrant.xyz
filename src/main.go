package main

import (
	"goservice/pkg/loadenv"
	"net/http"
	"log"
	"os"
)

// main
func main() {
	// Configure the Go logger to write logs to stdout
	log.SetOutput(os.Stdout)
	log.Printf("Alive")

	// Define routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/code", githubHandler)
	http.HandleFunc("/status", statusHandler)
	// Blog
	http.HandleFunc("/blog", blogHandler)
	http.HandleFunc("/blog/", blogHandler) // Match /blog/{filename}

	// Load in key/values from environment file
	loadenv.LoadEnv(".env")

	http.ListenAndServe(":9000", nil)
}