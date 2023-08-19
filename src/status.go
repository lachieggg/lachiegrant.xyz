package main

import (
	"fmt"
	"goservice/pkg/xml"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// statusHandler
func statusHandler(w http.ResponseWriter, r *http.Request) {
	htop := ExecuteCmd("htop")
	nf := ExecuteCmd("neofetch")

	// Merge XML outputs
	mergedOutput := xml.MergeXML(htop, nf, "Status")

	_, err := w.Write([]byte(mergedOutput))
	if err != nil {
		log.Printf("%v", err)
	}
	w.Header().Set("Content-Type", "text/html")
}

// ExecuteCmd executes a shell command and gets the formatted output
// for display
func ExecuteCmd(cmd string) []byte {
	cmdString := fmt.Sprintf("echo q | %s | aha --black --line-fix", cmd)
	htop := exec.Command("sh", "-c", cmdString)
	htop.Env = append(os.Environ(), "TERM=xterm")
	output, err := htop.CombinedOutput()
	if err != nil {
		log.Printf("%v", err)
	}

	return output
}
