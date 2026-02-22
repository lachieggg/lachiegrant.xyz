package main

import (
	"fmt"
	"goservice/pkg/file"
	"goservice/pkg/xml"
	"net/http"
	"os"
	"os/exec"
)

// nfHtmlPath is the path to the pre-generated neofetch HTML output.
const nfHtmlPath = "./public/neofetch.html"

// getStatusContent generates status HTML by combining htop output with neofetch.
// It executes htop, converts it to HTML via aha, and merges it with the neofetch HTML.
func getStatusContent() (string, error) {
	htopHTML := ExecuteCmd("htop")

	nfHTML, err := file.ReadFile(nfHtmlPath)
	if err != nil {
		return "", err
	}

	// Merge htop and neofetch HTML into a single page
	merged, err := xml.MergeHTMLContents(htopHTML, nfHTML)
	if err != nil {
		return "", err
	}
	return merged, nil
}

// statusHandler serves the system status page showing htop and neofetch output.
func statusHandler(w http.ResponseWriter, r *http.Request) {
	mergedOutput, err := getStatusContent()
	if err != nil {
		logger.Printf("Failed to get status content: %v", err)
		http.Error(w, "Failed to generate status page", http.StatusInternalServerError)
		return
	}

	// Headers must be set before writing the response body
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(mergedOutput)); err != nil {
		logger.Printf("Failed to write response: %v", err)
	}
}

// ExecuteCmd runs a shell command and pipes it through aha to generate HTML output.
// It's primarily used to capture terminal output (like htop) as HTML.
func ExecuteCmd(cmd string) string {
	cmdString := fmt.Sprintf("echo q | %s | aha --black --line-fix", cmd)
	proc := exec.Command("sh", "-c", cmdString)
	proc.Env = append(os.Environ(), "TERM=xterm")

	output, err := proc.CombinedOutput()
	if err != nil {
		logger.Printf("Command execution failed: %v", err)
	}

	return string(output)
}
