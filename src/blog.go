package main

import (
	"fmt"
	"log"
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

	log.Printf("Blog request for post: %s", postName)

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
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}
