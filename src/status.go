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

	html := string(htop) + string(nf)
	html = replacer(html)

	_, err := w.Write([]byte(html))
	if err != nil {
		log.Printf("%v", err)
	}
	w.Header().Set("Content-Type", "text/html")
}

// replacer
func replacer(input string) string {
	return strings.Replace(
		input,
		fmt.Sprintf(titleString, "stdin"),
		fmt.Sprintf(titleString, "Status"),
		1,
	)
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
