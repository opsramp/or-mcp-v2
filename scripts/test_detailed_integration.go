package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/tools"
)

const (
	// LogDir is the directory where logs will be stored
	LogDir = "output/logs"
	// LogFileName is the name of the log file
	LogFileName = "detailed-integration-test.log"
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
	customLogger.Info("Starting detailed integration test")

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

	// First, list integrations to get an ID
	customLogger.Info("Listing integrations to get an ID")
	listReq := mcp.CallToolRequest{}
	listReq.Params.Name = "integrations"
	listReq.Params.Arguments = map[string]interface{}{
		"action": "list",
	}

	// Call the handler to list integrations
	listResult, err := tools.IntegrationsToolHandler(context.Background(), listReq, integrationsTool)
	if err != nil {
		customLogger.Error("Error calling handler for list: %v", err)
		os.Exit(1)
	}

	// Check result
	if listResult.IsError {
		customLogger.Error("Handler returned error for list: %v", listResult.Content)
		os.Exit(1)
	}

	// Convert result to string for display
	listResultText := ""
	for _, content := range listResult.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			listResultText = textContent.Text
			break
		}
	}

	// Print the list result
	fmt.Printf("List Integrations Result:\n%s\n\n", listResultText)

	// Parse the list result to get the first integration ID
	var integrations []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(listResultText), &integrations); err != nil {
		customLogger.Error("Failed to parse integrations list: %v", err)
		os.Exit(1)
	}

	if len(integrations) == 0 {
		customLogger.Error("No integrations found")
		os.Exit(1)
	}

	// Get the first integration ID
	integrationID := integrations[0].ID
	customLogger.Info("Using integration ID: %s", integrationID)

	// Test getting detailed integration
	customLogger.Info("Testing GetDetailed Integration")
	req := mcp.CallToolRequest{}
	req.Params.Name = "integrations"
	req.Params.Arguments = map[string]interface{}{
		"action": "getDetailed",
		"id":     integrationID,
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
	fmt.Printf("GetDetailed Integration Result:\n%s\n", resultText)
	customLogger.Info("GetDetailed Integration completed successfully")

	// Print log file location
	logPath := filepath.Join(LogDir, LogFileName)
	fmt.Printf("\nLog file created at: %s\n", logPath)
}
