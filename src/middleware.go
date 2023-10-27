package main

import (
	"fmt"
	"net/http"
)

// middlewareFunc is a general middleware function
// for all endpoints
func middlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		if logger == nil {
			fmt.Println("Logger is nil")
		} else if r == nil {
			logger.Println("Request is nil")
		} else if r.URL == nil {
			logger.Printf("URL is nil, remote address: %s\n", r.RemoteAddr)
		} else {
			logger.Printf("Request received for endpoint %s from %s\n",
				r.URL.Path,
				r.RemoteAddr,
			)
		}
		next(w, r)
	}
}
