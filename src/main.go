package main

import (
	"fmt"
	"goservice/pkg/loadenv"
	"log"
	"net/http"
	"os"
)

const port = 9000

// main
func main() {
	// Configure the Go logger to write logs to stdout
	log.SetOutput(os.Stdout)

	// Define routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/code", githubHandler)
	http.HandleFunc("/status", statusHandler)
	// Blog
	http.HandleFunc("/blog", blogHandler)
	http.HandleFunc("/blog/", blogHandler) // Match /blog/{filename}

	// Load in key/values from environment file
	loadenv.LoadEnv(".env")

	log.Printf("Alive on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
