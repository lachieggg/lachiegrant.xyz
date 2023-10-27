package main

import (
	"log"
	"os"
)

var logger *log.Logger

// initLogger
func initLogger(logFileName string) {
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Create the global logger
	logger = log.New(logFile, "", log.LstdFlags)
}
