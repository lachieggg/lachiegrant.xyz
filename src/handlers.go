package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// indexHandler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/app.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Title   string
		Header  string
		Content string
	}{
		Title:   "My Website",
		Header:  "Welcome to my website!",
		Content: "Thanks for visiting!",
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Home
func Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World!")
}
