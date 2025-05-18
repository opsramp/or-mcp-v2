package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vobbilis/codegen/or-mcp-v2/common"
)

func TestLogging(t *testing.T) {
	// Create a temporary log directory
	logDir := filepath.Join(os.TempDir(), "or-mcp-test-logs")
	logFile := "test.log"

	// Clean up after the test
	defer os.RemoveAll(logDir)

	// Initialize the logger
	customLogger, err := common.InitLogger(common.DEBUG, logDir, logFile)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer customLogger.Close()

	// Test logging at different levels
	customLogger.Debug("This is a debug message")
	customLogger.Info("This is an info message")
	customLogger.Warn("This is a warning message")
	customLogger.Error("This is an error message")

	// Verify that the log file exists
	logPath := filepath.Join(logDir, logFile)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("Log file was not created at %s", logPath)
	} else {
		t.Logf("Log file created successfully at %s", logPath)

		// Read the log file to verify content
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Errorf("Failed to read log file: %v", err)
		} else {
			t.Logf("Log file content: %s", string(content))
		}
	}
}
