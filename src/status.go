package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const titleOpen = "<title>"
const titleClose = "</title>"
const titleString = titleOpen + "%s" + titleClose

// statusHandler
func statusHandler(w http.ResponseWriter, r *http.Request) {
	htop := ExecuteCmd("htop")
	nf := ExecuteCmd("neofetch")

	html := htop + nf

	_, err := w.Write([]byte(html))
	if err != nil {
		log.Printf("%v", err)
	}
	w.Header().Set("Content-Type", "text/html")
}

// replacer
func replacer(input []byte) string {
	return strings.Replace(
		string(input),
		fmt.Sprintf(titleString, "stdin"),
		fmt.Sprintf(titleString, "Status"),
		1,
	)
}

// ExecuteCmd executes a shell command and gets the formatted output
// for display
func ExecuteCmd(cmd string) string {
	cmdString := fmt.Sprintf("echo q | %s | aha --black --line-fix", cmd)
	htop := exec.Command("sh", "-c", cmdString)
	htop.Env = append(os.Environ(), "TERM=xterm")
	output, err := htop.CombinedOutput()
	if err != nil {
		log.Printf("%v", err)
	}

	return replacer(output)
}
