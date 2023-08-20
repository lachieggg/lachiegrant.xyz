package main

import (
	"fmt"
	"goservice/pkg/file"
	"goservice/pkg/xml"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const nfHtmlPath = "../public/neofetch.html"

// getStatusContent
func getStatusContent() (string, error) {
	htopHtml := ExecuteCmd("htop")

	nfHtml, err := file.ReadFile(nfHtmlPath)

	if err != nil {
		return "", err
	}

	// Merge
	merged, err := xml.MergeHTMLContents(htopHtml, nfHtml)
	return merged, err
}

// statusHandler
func statusHandler(w http.ResponseWriter, r *http.Request) {
	mergedOutput, err := getStatusContent()

	_, err = w.Write([]byte(mergedOutput))
	if err != nil {
		log.Printf("%v", err)
	}
	w.Header().Set("Content-Type", "text/html")
}

// ExecuteCmd
func ExecuteCmd(cmd string) string {
	cmdString := fmt.Sprintf("echo q | %s | aha --black --line-fix", cmd)
	htop := exec.Command("sh", "-c", cmdString)
	htop.Env = append(os.Environ(), "TERM=xterm")
	output, err := htop.CombinedOutput()
	if err != nil {
		log.Printf("%v", err)
	}

	content := string(output)
	return content
}
