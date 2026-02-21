package main

import (
	"log"
	"os"
)

// logger is the application-wide logger instance used by middleware and handlers.
// logFile holds the file handle for cleanup (though currently not closed on shutdown).
var (
	logger  *log.Logger
	logFile *os.File
)

// loggingKey is the environment variable that enables file logging.
// If LOGGING is set (to any value), logs are written to a file.
// If unset, logs are discarded.
const loggingKey = "LOGGING"

// initLogger configures the global logger. If the LOGGING environment variable
// is set, logs are written to the specified file. Otherwise, logs are discarded.
func initLogger(logFileName string) {
	if _, logging := os.LookupEnv(loggingKey); !logging {
		// Logging disabled - send output to /dev/null
		null, _ := os.Open("/dev/null")
		logger = log.New(null, "", log.LstdFlags)
		return
	}

	// Logging enabled - append to the specified file
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger = log.New(logFile, "", log.LstdFlags)
}
