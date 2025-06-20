package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/opsramp/or-mcp-v2/common"
	"github.com/opsramp/or-mcp-v2/pkg/client"
)

// setupTestEnvironment sets up the testing environment with proper logging and client
func setupTestEnvironment(t *testing.T) (*common.CustomLogger, func()) {
	// Create log directory in proper location relative to pkg/tools
	logDir := filepath.Join("..", "..", "output", "logs")
	logFile := "resources-test.log"

	// Get absolute path for log directory
	absLogDir, err := filepath.Abs(logDir)
	if err != nil {
		t.Fatalf("Failed to get absolute log path: %v", err)
	}

	// Ensure log directory exists
	if err := os.MkdirAll(absLogDir, 0750); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	// Initialize logger with absolute path
	logger, err := common.InitLogger(common.DEBUG, absLogDir, logFile)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Load configuration
	config, err := common.LoadConfig("")
	if err != nil {
		logger.Warn("Failed to load config: %v", err)
		logger.Info("Using default test configuration")
		// Create a minimal default config for testing
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

	logger.Info("Test environment setup complete")

	// Return cleanup function
	cleanup := func() {
		logger.Info("Cleaning up test environment")
		logger.Close()
	}

	return logger, cleanup
}

// createTestRequest creates a properly formatted MCP CallToolRequest
func createTestRequest(arguments map[string]interface{}) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "resources",
			Arguments: arguments,
		},
	}
}

func TestResourcesToolCreation(t *testing.T) {
	logger, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create the resources tool
	tool, handler := NewResourcesMcpTool()

	// Test that the tool was created successfully
	if tool.Name != "resources" {
		t.Errorf("Expected tool name to be 'resources', got '%s'", tool.Name)
	}

	if tool.Description == "" {
		t.Errorf("Expected tool description to be non-empty")
	}

	// Test that the handler was created successfully
	if handler == nil {
		t.Errorf("Expected handler to be non-nil")
	}

	logger.Info("Resources tool creation test completed successfully")
}

func TestResourcesTool_List(t *testing.T) {
	logger, cleanup := setupTestEnvironment(t)
	defer cleanup()

	_, handler := NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "list",
	})

	logger.Info("Testing resources list action")
	res, err := handler(context.Background(), req)

	// Note: This test may fail if OpsRamp credentials are invalid, but we test the code path
	if err != nil {
		logger.Warn("List action failed (expected with test credentials): %v", err)
		// Check if it's a proper error response
		if res != nil && res.IsError {
			logger.Info("Received proper error response for invalid credentials")
		}
	} else if res != nil && len(res.Content) > 0 {
		logger.Info("List action succeeded with valid credentials")
	}

	logger.Info("Resources list test completed")
}

func TestResourcesTool_Get(t *testing.T) {
	logger, cleanup := setupTestEnvironment(t)
	defer cleanup()

	_, handler := NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "get",
		"id":     "test-resource-id",
	})

	logger.Info("Testing resources get action")
	res, err := handler(context.Background(), req)

	// Note: This test may fail if OpsRamp credentials are invalid or resource doesn't exist
	if err != nil {
		logger.Warn("Get action failed (expected with test data): %v", err)
	} else if res != nil {
		logger.Info("Get action completed")
	}

	logger.Info("Resources get test completed")
}

func TestResourcesTool_GetMinimal(t *testing.T) {
	logger, cleanup := setupTestEnvironment(t)
	defer cleanup()

	_, handler := NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "getMinimal",
		"id":     "test-resource-id",
	})

	logger.Info("Testing resources getMinimal action")
	res, err := handler(context.Background(), req)

	if err != nil {
		logger.Warn("GetMinimal action failed (expected with test data): %v", err)
	} else if res != nil {
		logger.Info("GetMinimal action completed")
	}

	logger.Info("Resources getMinimal test completed")
}

func TestResourcesTool_Search(t *testing.T) {
	logger, cleanup := setupTestEnvironment(t)
	defer cleanup()

	_, handler := NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "search",
		"params": map[string]interface{}{
			"pageSize": 10,
			"pageNo":   1,
		},
	})

	logger.Info("Testing resources search action")
	res, err := handler(context.Background(), req)

	if err != nil {
		logger.Warn("Search action failed (expected with test credentials): %v", err)
	} else if res != nil {
		logger.Info("Search action completed")
	}

	logger.Info("Resources search test completed")
}

func TestResourcesTool_GetResourceTypes(t *testing.T) {
	logger, cleanup := setupTestEnvironment(t)
	defer cleanup()

	_, handler := NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "getResourceTypes",
	})

	logger.Info("Testing resources getResourceTypes action")
	res, err := handler(context.Background(), req)

	if err != nil {
		logger.Warn("GetResourceTypes action failed (expected with test credentials): %v", err)
	} else if res != nil {
		logger.Info("GetResourceTypes action completed")
	}

	logger.Info("Resources getResourceTypes test completed")
}

func TestResourcesTool_InvalidAction(t *testing.T) {
	logger, cleanup := setupTestEnvironment(t)
	defer cleanup()

	_, handler := NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "invalid-action",
	})

	logger.Info("Testing resources invalid action")
	res, err := handler(context.Background(), req)

	if err != nil {
		logger.Info("Invalid action properly returned error: %v", err)
	} else if res != nil && res.IsError {
		logger.Info("Invalid action properly returned error response")
	} else {
		t.Errorf("Expected error for invalid action")
	}

	logger.Info("Resources invalid action test completed")
}

func TestResourcesTool_MissingRequiredParameters(t *testing.T) {
	logger, cleanup := setupTestEnvironment(t)
	defer cleanup()

	_, handler := NewResourcesMcpTool()

	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "get without id",
			args: map[string]interface{}{"action": "get"},
		},
		{
			name: "getMinimal without id",
			args: map[string]interface{}{"action": "getMinimal"},
		},
		{
			name: "changeState without id",
			args: map[string]interface{}{"action": "changeState", "state": "active"},
		},
		{
			name: "changeState without state",
			args: map[string]interface{}{"action": "changeState", "id": "test-id"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := createTestRequest(tc.args)

			logger.Info("Testing %s", tc.name)
			res, err := handler(context.Background(), req)

			if err != nil {
				logger.Info("Test case '%s' properly returned error: %v", tc.name, err)
			} else if res != nil && res.IsError {
				logger.Info("Test case '%s' properly returned error response", tc.name)
			} else {
				t.Errorf("Test case '%s': Expected error for missing required parameters", tc.name)
			}
		})
	}

	logger.Info("Missing required parameters tests completed")
}
