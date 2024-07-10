package main

import (
	"fmt"
	"goservice/pkg/loadenv"
	"log"
	"net/http"
)

const port = 9000

// main
func main() {
	initLogger("log.out")

	// Define routes
	http.HandleFunc("/", middlewareFunc(indexHandler))
	http.HandleFunc("/code", middlewareFunc(githubHandler))
	http.HandleFunc("/status", middlewareFunc(statusHandler))
	// Blog
	http.HandleFunc("/blog", middlewareFunc(blogHandler))
	http.HandleFunc("/blog/", middlewareFunc(blogHandler)) // Match /blog/{filename}
	// Bookmarks
	http.HandleFunc("/bookmarks", middlewareFunc(bookmarksHandler))

	// Load in key/values from environment file
	loadenv.LoadEnv(".env")

	log.Printf("Alive on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
