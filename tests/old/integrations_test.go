package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/opsramp/or-mcp-v2/common"
	"github.com/opsramp/or-mcp-v2/pkg/client"
	"github.com/opsramp/or-mcp-v2/pkg/tools"
)

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestIntegrationsTool(t *testing.T) {
	// Load configuration
	config, err := common.LoadConfig("")
	if err != nil {
		// Use default config for testing
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

	// Create OpsRamp client (not used in this test as we're using mocks)
	_ = client.NewOpsRampClient(config)

	// Create the integrations tool
	_, intHandler := tools.NewIntegrationsMcpTool()

	// Test cases
	testCases := []struct {
		name     string
		action   string
		id       string
		config   map[string]interface{}
		expected string
	}{
		{
			name:     "List Integrations",
			action:   "list",
			expected: "Mock Integration 1",
		},
		{
			name:     "Get Integration",
			action:   "get",
			id:       "int-001",
			expected: "Mock Integration",
		},
		{
			name:   "Create Integration",
			action: "create",
			config: map[string]interface{}{
				"name": "New Test Integration",
				"type": "api",
			},
			expected: "New Test Integration",
		},
		{
			name:     "List Integration Types",
			action:   "listTypes",
			expected: "API Integration",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req := mcp.CallToolRequest{}
			req.Params.Name = "integrations"
			req.Params.Arguments = map[string]interface{}{
				"action": tc.action,
			}

			// Add id if provided
			if tc.id != "" {
				req.Params.Arguments["id"] = tc.id
			}

			// Add config if provided
			if tc.config != nil {
				req.Params.Arguments["config"] = tc.config
			}

			// Call the handler
			result, err := intHandler(context.Background(), req)
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

			// Print the result for debugging
			fmt.Printf("Result for %s: %s\n", tc.name, resultText)

			// Check if expected string is in the result
			if tc.expected != "" && resultText != "" {
				// For this test, we just check if the expected string is contained in the result
				if !contains(resultText, tc.expected) {
					t.Errorf("Expected result to contain '%s', got '%s'", tc.expected, resultText)
				} else {
					t.Logf("Test passed: Result contains '%s'", tc.expected)
				}
			}
		})
	}
}
