package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/opsramp/or-mcp-v2/common"
)

func main() {
	// Create output directory
	logDir := "output/logs"
	logFile := "custom.log"
	
	// Initialize the logger
	customLogger, err := common.InitLogger(common.DEBUG, logDir, logFile)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer customLogger.Close()

	// Log test messages
	customLogger.Info("Starting logging test")
	customLogger.Debug("This is a debug message")
	customLogger.Info("This is an info message")
	customLogger.Warn("This is a warning message")
	customLogger.Error("This is an error message")

	// Print log file location
	logPath := filepath.Join(logDir, logFile)
	fmt.Printf("Log file created at: %s\n", logPath)
}
