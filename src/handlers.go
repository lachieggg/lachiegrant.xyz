package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	FeatureEnabled  = "enabled"
	FeatureDisabled = "disabled"

	EnvEnableBlog      = "ENABLE_BLOG"
	EnvEnableBookmarks = "ENABLE_BOOKMARKS"
	EnvGithubURL       = "GITHUB_URL"
	EnvResumePath      = "RESUME_PATH"

	TmplHome      = "home.html"
	TmplBookmarks = "bookmarks.html"
	TmplBlogIndex = "blog.html"
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

// PageData holds data passed to HTML templates.
type PageData struct {
	EnableBlog      bool
	EnableBookmarks bool
	Data            interface{}
}

// getPageData returns a PageData object with common feature flags populated.
func getPageData(data interface{}) PageData {
	blog := os.Getenv(EnvEnableBlog) == FeatureEnabled
	bookmarks := os.Getenv(EnvEnableBookmarks) == FeatureEnabled
	logger.Printf("Debug: PageData - Blog: %v, Bookmarks: %v", blog, bookmarks)
	return PageData{
		EnableBlog:      blog,
		EnableBookmarks: bookmarks,
		Data:            data,
	}
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

	if err := t.ExecuteTemplate(w, TmplHome, getPageData(nil)); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

// githubHandler redirects users to the configured GitHub URL.
func githubHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, os.Getenv(EnvGithubURL), http.StatusSeeOther)
}

// bookmarksHandler renders the bookmarks page template.
func bookmarksHandler(w http.ResponseWriter, r *http.Request) {
	if os.Getenv(EnvEnableBookmarks) != FeatureEnabled {
		logger.Printf("404: Bookmarks disabled via feature flag")
		http.NotFound(w, r)
		return
	}

	t := parseTemplates(w)
	if t == nil {
		return
	}

	if err := t.ExecuteTemplate(w, TmplBookmarks, getPageData(nil)); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

// resumeHandler serves the resume PDF file.
func resumeHandler(w http.ResponseWriter, r *http.Request) {
	resumePath := os.Getenv(EnvResumePath)
	if resumePath == "" {
		logger.Printf("Error: RESUME_PATH environment variable not set")
		http.Error(w, "Resume path not configured", http.StatusInternalServerError)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		logger.Printf("Error getting working directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Clean and resolve the absolute path relative to the 'public' directory
	publicDir := filepath.Join(wd, "public")
	absPath := filepath.Join(publicDir, filepath.Clean(resumePath))

	// Security Check: Ensure the resolved path is still within the 'public' directory
	if !strings.HasPrefix(absPath, publicDir) {
		logger.Printf("Security Warning: Unauthorized access attempt to %s", absPath)
		http.Error(w, "Forbidden", http.StatusForbidden)
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
