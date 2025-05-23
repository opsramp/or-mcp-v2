package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/tools"
)

func TestResourcesTool(t *testing.T) {
	// Create a log directory in the project root output/logs
	logDir := filepath.Join("..", "output", "logs")
	logFile := "resources-test.log"

	// Clean up after the test
	defer os.Remove(filepath.Join(logDir, logFile))

	// Initialize the logger
	customLogger, err := common.InitLogger(common.DEBUG, logDir, logFile)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer customLogger.Close()

	// Load configuration
	config, err := common.LoadConfig("")
	if err != nil {
		customLogger.Warn("Failed to load config: %v", err)
		customLogger.Info("Using default configuration")
		// Create a minimal default config
		config = &common.Config{
			OpsRamp: common.OpsRampConfig{
				TenantURL:  "https://api.opsramp.com",
				AuthURL:    "https://api.opsramp.com/auth/oauth/token",
				AuthKey:    "test-key",
				AuthSecret: "test-secret",
				TenantID:   "test-tenant",
			},
		}
	}

	// Create OpsRamp client
	opsRampClient := client.NewOpsRampClient(config)
	client.SetGlobalClient(opsRampClient)

	// Create the resources tool
	tool, handler := tools.NewResourcesMcpTool()

	// Test that the tool was created successfully
	if tool.Name != "resources" {
		t.Errorf("Expected tool name to be 'resources', got '%s'", tool.Name)
	}

	// Test that the handler was created successfully
	if handler == nil {
		t.Errorf("Expected handler to be non-nil")
	}

	// Log test completion
	customLogger.Info("Resources tool test completed successfully")
}
