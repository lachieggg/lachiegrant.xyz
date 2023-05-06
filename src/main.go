package main

import (
	"net/http"
)

// main
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/files", filesHandler)
	http.HandleFunc("/github", githubHandler)
	

	http.ListenAndServe(":9000", nil)
}
