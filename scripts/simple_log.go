package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Create output directory
	logDir := "output/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		os.Exit(1)
	}

	// Create a log file
	logFile := filepath.Join(logDir, "simple.log")
	f, err := os.Create(logFile)
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Write to the log file
	_, err = f.WriteString("This is a test log entry\n")
	if err != nil {
		fmt.Printf("Failed to write to log file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Log file created at: %s\n", logFile)
}
