package main

import (
	"fmt"
	"goservice/pkg/loadenv"
	"log"
	"net/http"
)

const (
	port    = 9000
	logPath = "log.out"
)

// main is the application entry point. It loads environment variables,
// initializes logging, registers HTTP routes, and starts the server.
func main() {
	// Load key/value pairs from .env file into the environment
	if err := loadenv.LoadEnv(".env"); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	initLogger(logPath)

	if v, err := publishVersionedAsset("public/css/app.css"); err != nil {
		logger.Printf("Warning: could not publish versioned CSS (falling back to /css/app.css): %v", err)
	} else {
		cssPath = v
	}
	if v, err := publishVersionedAsset("public/js/generated/home.js"); err != nil {
		logger.Printf("Warning: could not publish versioned home.js (falling back to /js/generated/home.js): %v", err)
	} else {
		homeJSPath = v
	}
	if v, err := publishVersionedAsset("public/images/rpi.png"); err != nil {
		logger.Printf("Warning: could not publish versioned logo (falling back to /images/rpi.png): %v", err)
	} else {
		logoPath = v
	}
	if v, err := publishVersionedAsset("public/images/debian.png"); err != nil {
		logger.Printf("Warning: could not publish versioned Debian logo (falling back to /images/debian.png): %v", err)
	} else {
		debianLogoPath = v
	}
	publishHomeImages()

	http.HandleFunc("/", middlewareFunc(indexHandler))
	http.HandleFunc("/code", middlewareFunc(codeHandler))
	http.HandleFunc("/blog", middlewareFunc(blogHandler))
	http.HandleFunc("/bookmarks", middlewareFunc(bookmarksHandler))
	http.HandleFunc("/resume", middlewareFunc(resumeHandler))
	http.HandleFunc("/blog/", middlewareFunc(blogHandler)) // Matches /blog/{post-name}

	logger.Printf("Server starting on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
