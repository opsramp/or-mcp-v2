package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
)

const (
	// LogDir is the directory where logs will be stored
	LogDir = "output/logs"
	// LogFileName is the name of the log file
	LogFileName = "test-logging.log"
)

func main() {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(LogDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize the logger
	customLogger, err := common.InitLogger(common.DEBUG, LogDir, LogFileName)
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

	// Create a test config
	config := &common.Config{
		OpsRamp: common.OpsRampConfig{
			TenantURL:  "https://api.opsramp.com",
			AuthURL:    "https://api.opsramp.com/auth/oauth/token",
			AuthKey:    "test-key",
			AuthSecret: "test-secret",
			TenantID:   "test-tenant",
		},
	}

	// Create a client with logging
	opsrampClient := client.NewOpsRampClient(config)

	// Test client logging
	customLogger.Info("Testing client logging")
	
	// Print log file location
	logPath := filepath.Join(LogDir, LogFileName)
	customLogger.Info("Log file created at: %s", logPath)
	fmt.Printf("Log file created at: %s\n", logPath)
}
