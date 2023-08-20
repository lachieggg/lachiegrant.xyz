package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

// parseTemplates
func parseTemplates(w http.ResponseWriter) *template.Template {
	wd, err := os.Getwd()

	// First folder
	standardFolder := filepath.Join(wd, "src", "templates", "*.html")
	t, err := template.ParseGlob(standardFolder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing templates: %v", err), http.StatusInternalServerError)
		return nil
	}

	// Second folder
	blogFolder := filepath.Join(wd, "src", "templates", "blog", "*.html")
	t, err = t.ParseGlob(blogFolder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing templates: %v", err), http.StatusInternalServerError)
		return nil
	}

	return t
}

// indexHandler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	} else {
		homeHandler(w, r)
	}
}

// homeHandler
func homeHandler(w http.ResponseWriter, r *http.Request) {
	t := parseTemplates(w)
	if t == nil {
		return
	}

	err := t.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

// githubHandler
func githubHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, os.Getenv("GITHUB_URL"), http.StatusSeeOther)
}
