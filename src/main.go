package main

import (
	"fmt"
	"net/http"
)

// main
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			indexHandler(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Only GET requests are allowed")
		}
	})

	http.ListenAndServe(":9000", nil)
}
