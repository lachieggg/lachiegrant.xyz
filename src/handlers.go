package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"io"
	"mime"
	"log"
	"os/exec"
)

// indexHandler
func indexHandler(w http.ResponseWriter, r *http.Request) {

	wd, err := os.Getwd()
	templateFolder := filepath.Join(wd, "src", "templates")
	templatePath := filepath.Join(templateFolder, "*.html")

	// Parse all templates in the specified directory
	t, err := template.ParseGlob(templatePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing templates: %v", err), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "app.html", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}

	// register the include function
	t.Funcs(template.FuncMap{
		"include": func(name string, data interface{}) (string, error) {
			// load the named template file
			tpl, err := template.ParseFiles(name)
			if err != nil {
				return "", err
			}

			// execute the named template with the provided data
			var buf bytes.Buffer
			if err := tpl.Execute(&buf, data); err != nil {
				return "", err
			}
			return buf.String(), nil
		},
	})
}

// filesHandler
func filesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		fmt.Fprintf(w, "%s\n", filepath.Join("/", file.Name()))
	}

	fmt.Fprint(w, "\n")
	files, err = os.ReadDir("/app")
	for _, file := range files {
		fmt.Fprintf(w, "%s\n", filepath.Join("/", file.Name()))
	}

	fmt.Fprint(w, "\n")
	wd, err := os.Getwd()
	fmt.Fprintf(w, "%s", wd)
}

// blogHandler handles the /blog route and specific files within /blog
func blogHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the requested file name from the URL path
	fileName := filepath.Base(r.URL.Path)
	log.Printf("Received request for file: %s", fileName)

	if fileName == "blog" {
		writeFile(w, "./public/blog/blog.html")
		return
	}

	// Construct the file path based on the requested file name
	filePath := filepath.Join("./public/blog/", fileName)

	log.Printf("Full file path: %s", filePath)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Printf("File not found: %s", filePath)
		http.NotFound(w, r)
		return
	} else if err != nil {
		log.Printf("Failed to check file existence: %v", err)
		http.Error(w, fmt.Sprintf("Failed to check file existence: %v", err), http.StatusInternalServerError)
		return
	}

	writeFile(w, filePath)
}

// topHandler
func topHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("sh", "-c", "echo q | htop | aha --black --line-fix")
	cmd.Env = append(os.Environ(), "TERM=xterm")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%v", err)
	}
	fmt.Fprintf(w, "%s\n", out)
}


// writeFile
func writeFile(w http.ResponseWriter, path string) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		http.Error(w, fmt.Sprintf("Failed to open file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the Content-Type header based on the file extension
	contentType := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", contentType)

	// Copy the file contents to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Failed to write file content to response: %v", err)
		http.Error(w, fmt.Sprintf("Failed to write file content to response: %v", err), http.StatusInternalServerError)
		return
	}
}

// githubHandler
func githubHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, os.Getenv("GITHUB_URL"), http.StatusSeeOther)
}