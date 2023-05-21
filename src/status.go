package main

import (
	"log"
	"os/exec"
	"strings"
	"net/http"
	"os"
	"fmt"
)

const titleOpen = "<title>"
const titleClose = "</title>"
const titleString = titleOpen + "%s" + titleClose

// statusHandler
func statusHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("sh", "-c", "echo q | htop | aha --black --line-fix")
	cmd.Env = append(os.Environ(), "TERM=xterm")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%v", err)
	}
	html := strings.Replace(string(out),
		fmt.Sprintf(titleString, "stdin"),
		fmt.Sprintf(titleString, "Status"),
		1,
	)
	_, err = w.Write([]byte(html))
	if err != nil {
		log.Printf("%v", err)
	}
	w.Header().Set("Content-Type", "text/html")
}