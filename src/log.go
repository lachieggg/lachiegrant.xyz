package main

import (
	"log"
	"os"
)

var logger *log.Logger

const (
	loggingKey = "LOGGING"
)

// initLogger
func initLogger(logFileName string) {
	_, logging := os.LookupEnv(loggingKey)
	if !logging {
		return
	}

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Create the global logger
	logger = log.New(logFile, "", log.LstdFlags)
}
