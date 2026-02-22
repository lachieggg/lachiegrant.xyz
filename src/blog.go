package main

import (
	"fmt"
	"net/http"
	"strings"
)

// blogHandler serves blog content. For /blog it renders the blog index page.
// For /blog/{post-name} it renders the corresponding post template.
func blogHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(path, "/")

	// Extract the post name from the URL path (e.g., /blog/my-post -> "my-post")
	postName := segments[len(segments)-1]

	t := parseTemplates(w)
	if t == nil {
		return
	}

	var templateName string
	if postName == "" {
		templateName = "blog.html"
	} else {
		templateName = fmt.Sprintf("%s.html", postName)
	}

	err := t.ExecuteTemplate(w, templateName, nil)
	if err != nil {
		if strings.Contains(err.Error(), "incomplete or empty template") || strings.Contains(err.Error(), "is not defined") {
			logger.Printf("404: Blog template %s not found", templateName)
			http.NotFound(w, r)
			return
		}
		logger.Printf("Error executing blog template %s: %v", templateName, err)
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}
