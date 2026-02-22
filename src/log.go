package main

import (
	"io"
	"log"
	"os"
)

const (
	logPrefix = ""
)

// logger is the application-wide logger instance used by middleware and handlers.
var logger *log.Logger

// initLogger configures the global logger to write to both os.Stdout and the specified file.
func initLogger(logFileName string) {
	f, err := os.OpenFile(
		logFileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Printf("Warning: Failed to open log file %s: %v", logFileName, err)
		logger = log.New(os.Stdout, logPrefix, log.LstdFlags)
		return
	}

	multi := io.MultiWriter(os.Stdout, f)
	logger = log.New(multi, logPrefix, log.LstdFlags)
}
