package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

// parseTemplates loads and parses all HTML templates from the templates directory.
// Returns nil and writes an error response if template parsing fails.
// NOTE: Templates are currently re-parsed on each request. For production,
// consider parsing once at startup and caching the result.
func parseTemplates(w http.ResponseWriter) *template.Template {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, fmt.Sprintf("Operating system error: %v", err), http.StatusInternalServerError)
		return nil
	}

	// Parse standard page templates
	standardFolder := filepath.Join(wd, "src", "templates", "*.html")
	t, err := template.ParseGlob(standardFolder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing templates: %v", err), http.StatusInternalServerError)
		return nil
	}

	// Parse blog post templates
	blogFolder := filepath.Join(wd, "src", "templates", "blog", "*.html")
	t, err = t.ParseGlob(blogFolder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing templates: %v", err), http.StatusInternalServerError)
		return nil
	}

	return t
}

// indexHandler serves the root path. It renders the home page for "/" and
// returns 404 for any other path (since "/" matches all unhandled routes).
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		logger.Printf("404: %s not found (matched index catch-all)", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	homeHandler(w, r)
}

// homeHandler renders the home page template.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	t := parseTemplates(w)
	if t == nil {
		return
	}

	if err := t.ExecuteTemplate(w, "home.html", nil); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

// githubHandler redirects users to the configured GitHub URL.
func githubHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, os.Getenv("GITHUB_URL"), http.StatusSeeOther)
}

// bookmarksHandler renders the bookmarks page template.
func bookmarksHandler(w http.ResponseWriter, r *http.Request) {
	t := parseTemplates(w)
	if t == nil {
		return
	}

	if err := t.ExecuteTemplate(w, "bookmarks.html", nil); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

// resumeHandler serves the resume PDF file.
func resumeHandler(w http.ResponseWriter, r *http.Request) {
	resumePath := os.Getenv("RESUME_PATH")
	if resumePath == "" {
		logger.Printf("Error: RESUME_PATH environment variable not set")
		http.Error(w, "Resume path not configured", http.StatusInternalServerError)
		return
	}

	absPath, err := filepath.Abs(resumePath)
	if err != nil {
		logger.Printf("Error resolving resume path: %v", err)
		http.Error(w, "Error resolving resume path", http.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		logger.Printf("404: Resume file not found at %s", absPath)
		http.NotFound(w, r)
		return
	}

	logger.Printf("Serving resume from: %s", absPath)
	http.ServeFile(w, r, absPath)
}
