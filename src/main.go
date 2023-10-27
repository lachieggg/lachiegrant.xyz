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
	http.HandleFunc("/", middlewareFunc(indexHandler))
	http.HandleFunc("/code", middlewareFunc(githubHandler))
	http.HandleFunc("/status", middlewareFunc(statusHandler))
	// Blog
	http.HandleFunc("/blog", middlewareFunc(blogHandler))
	http.HandleFunc("/blog/", middlewareFunc(blogHandler)) // Match /blog/{filename}

	// Load in key/values from environment file
	loadenv.LoadEnv(".env")

	log.Printf("Alive on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
