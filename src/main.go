package main

import (
	"goservice/pkg/loadenv"
	"net/http"
)

// main
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/files", filesHandler)
	http.HandleFunc("/github", githubHandler)
	http.HandleFunc("/env", Env)

	loadenv.LoadEnv(".env")

	http.ListenAndServe(":9000", nil)
}
