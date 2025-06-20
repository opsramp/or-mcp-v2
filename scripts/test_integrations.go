package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/opsramp/or-mcp-v2/common"
	"github.com/opsramp/or-mcp-v2/pkg/client"
	"github.com/opsramp/or-mcp-v2/pkg/tools"
)

const (
	// LogDir is the directory where logs will be stored
	LogDir = "output/logs"
	// LogFileName is the name of the log file
	LogFileName = "integrations-test.log"
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

	// Log startup
	customLogger.Info("Starting integrations test")
	
	// Load configuration
	config, err := common.LoadConfig("config.yaml")
	if err != nil {
		customLogger.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	// Create OpsRamp client
	opsrampClient := client.NewOpsRampClient(config)
	customLogger.Info("Created OpsRamp client with tenant ID: %s", opsrampClient.GetTenantID())

	// Create the real OpsRamp integrations API
	integrationsAPI := tools.NewOpsRampIntegrationsAPI(opsrampClient)
	customLogger.Info("Created OpsRamp integrations API")

	// Create a custom integrations tool that uses the real API
	integrationsTool := tools.NewIntegrationsTool(integrationsAPI)
	customLogger.Info("Created integrations tool")

	// Test listing integrations
	customLogger.Info("Testing List Integrations")
	req := mcp.CallToolRequest{}
	req.Params.Name = "integrations"
	req.Params.Arguments = map[string]interface{}{
		"action": "list",
	}

	// Call the handler
	result, err := tools.IntegrationsToolHandler(context.Background(), req, integrationsTool)
	if err != nil {
		customLogger.Error("Error calling handler: %v", err)
		os.Exit(1)
	}

	// Check result
	if result.IsError {
		customLogger.Error("Handler returned error: %v", result.Content)
		os.Exit(1)
	}

	// Convert result to string for display
	resultText := ""
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			resultText = textContent.Text
			break
		}
	}

	// Print the result
	fmt.Printf("List Integrations Result:\n%s\n", resultText)
	customLogger.Info("List Integrations completed successfully")
	
	// Print log file location
	logPath := filepath.Join(LogDir, LogFileName)
	fmt.Printf("\nLog file created at: %s\n", logPath)
}
