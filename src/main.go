package main

import (
	"goservice/pkg/loadenv"
	"net/http"
	"log"
	"os"
)

// main
func main() {
	log.Printf("Alive")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/files/", filesHandler)
	http.HandleFunc("/github", githubHandler)

	// Top
	http.HandleFunc("/top", topHandler)

	// Blog
	http.HandleFunc("/blog/", blogHandler) // Match /blog/{filename}
	http.HandleFunc("/blog", blogHandler)

	loadenv.LoadEnv(".env")

	http.ListenAndServe(":9000", nil)
}



func init() {
	// Configure the Go logger to write logs to stdout
	log.SetOutput(os.Stdout)
}
