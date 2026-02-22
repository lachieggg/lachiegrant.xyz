package main

import (
	"io"
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
// If LOGGING is set (to any value), logs are written to a file in addition to stdout.
const loggingKey = "LOGGING"

// initLogger configures the global logger. It always logs to os.Stdout.
// If the LOGGING environment variable is set, logs are also written to the specified file.
func initLogger(logFileName string) {
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	if _, logging := os.LookupEnv(loggingKey); logging {
		// Logging to file enabled - append to the specified file
		var err error
		logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Warning: Failed to open log file %s: %v", logFileName, err)
		} else {
			writers = append(writers, logFile)
		}
	}

	multi := io.MultiWriter(writers...)
	logger = log.New(multi, "", log.LstdFlags)
}
