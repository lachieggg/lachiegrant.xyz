package main

import (
	"log"
	"os"
)

var (
	logger  *log.Logger
	logFile *os.File
)

const (
	loggingKey = "LOGGING"
)

// initLogger
func initLogger(logFileName string) {
	_, logging := os.LookupEnv(loggingKey)
	if !logging {
		// No logging
		null, _ := os.Open("/dev/null")
		logger = log.New(null, "", log.LstdFlags)
		return
	}

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Create global logger
	logger = log.New(logFile, "", log.LstdFlags)
}
