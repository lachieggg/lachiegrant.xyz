package main

import (
	"net/http"
	"log"
	"strings"
	"fmt"
)
// blogHandler handles the /blog route and specific files within /blog
func blogHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(path, "/")

	// The post name should be the last segment
	postName := segments[len(segments)-1]

	log.Printf("%s", postName)

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
