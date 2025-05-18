package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/tools"
)

// TestIntegrationsRealAPI tests the integrations tool against the real OpsRamp API
// This test uses the credentials in config.yaml to make actual API calls
func TestIntegrationsRealAPI(t *testing.T) {
	// Load configuration from config.yaml
	// Use "../config.yaml" since we're running from the tests directory
	config, err := common.LoadConfig("../config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify that we have the required credentials
	if config.OpsRamp.TenantURL == "" || config.OpsRamp.AuthURL == "" ||
		config.OpsRamp.AuthKey == "" || config.OpsRamp.AuthSecret == "" {
		t.Skip("Skipping test: Missing OpsRamp credentials in config.yaml")
	}

	// Create the real OpsRamp client
	opsrampClient := client.NewOpsRampClient(config)

	// Create the real OpsRamp integrations API
	integrationsAPI := tools.NewOpsRampIntegrationsAPI(opsrampClient)

	// Create a custom integrations tool that uses the real API
	customIntTool := tools.NewIntegrationsTool(integrationsAPI)

	// Create a handler function for the tool
	intHandler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.IntegrationsToolHandler(ctx, req, customIntTool)
	}

	// Test cases - only testing list and get
	testCases := []struct {
		name   string
		action string
		id     string
	}{
		{
			name:   "List Integrations",
			action: "list",
		},
		// The Get test will be added dynamically after we get the list of integrations
	}

	// Run the List test first to get integration IDs
	t.Run(testCases[0].name, func(t *testing.T) {
		// Create request
		req := mcp.CallToolRequest{}
		req.Params.Name = "integrations"
		req.Params.Arguments = map[string]interface{}{
			"action": testCases[0].action,
		}

		// Set a timeout for the API call
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Call the handler
		result, err := intHandler(ctx, req)
		if err != nil {
			t.Fatalf("Error calling handler: %v", err)
		}

		// Check result
		if result.IsError {
			t.Fatalf("Handler returned error: %v", result.Content)
		}

		// Convert result to string for comparison
		resultText := ""
		for _, content := range result.Content {
			if textContent, ok := content.(mcp.TextContent); ok {
				resultText = textContent.Text
				break
			}
		}

		// Print the result
		fmt.Printf("List Integrations Result: %s\n", resultText)

		// If we got a result, try to extract an integration ID for the Get test
		if resultText != "" && resultText != "OK" {
			// Get the first integration from the list to use for the Get test
			integrations, err := integrationsAPI.List(ctx)
			if err != nil {
				t.Logf("Warning: Could not parse integrations list: %v", err)
			} else if len(integrations) > 0 {
				// Add a Get test case with the first integration ID
				integrationID := integrations[0].ID
				t.Logf("Found integration ID for Get test: %s", integrationID)

				// Run the Get test
				t.Run("Get Integration", func(t *testing.T) {
					// Create request
					req := mcp.CallToolRequest{}
					req.Params.Name = "integrations"
					req.Params.Arguments = map[string]interface{}{
						"action": "get",
						"id":     integrationID,
					}

					// Set a timeout for the API call
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()

					// Call the handler
					result, err := intHandler(ctx, req)
					if err != nil {
						t.Fatalf("Error calling handler: %v", err)
					}

					// Check result
					if result.IsError {
						t.Fatalf("Handler returned error: %v", result.Content)
					}

					// Convert result to string for comparison
					resultText := ""
					for _, content := range result.Content {
						if textContent, ok := content.(mcp.TextContent); ok {
							resultText = textContent.Text
							break
						}
					}

					// Print the result
					fmt.Printf("Get Integration Result: %s\n", resultText)

					// Verify we got a non-empty result
					if resultText == "" || resultText == "OK" {
						t.Errorf("Expected detailed integration data, got: %s", resultText)
					} else {
						t.Logf("Successfully retrieved integration details")
					}
				})
			} else {
				t.Logf("Warning: No integrations found to test Get operation")
			}
		}
	})
}
