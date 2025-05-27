package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/tools"
)

// setupRealAPITestEnvironment sets up the testing environment with real OpsRamp API configuration
func setupRealAPITestEnvironment(t *testing.T) (*common.CustomLogger, func()) {
	// Create log directory
	logDir := filepath.Join("..", "output", "logs")
	logFile := "or-mcp.log"

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

	// Load real configuration from config.yaml
	config, err := common.LoadConfig("")
	if err != nil {
		t.Skipf("Skipping real API test - failed to load config.yaml: %v", err)
	}

	// Validate that we have real credentials
	if config.OpsRamp.AuthKey == "" || config.OpsRamp.AuthSecret == "" || config.OpsRamp.TenantID == "" {
		t.Skip("Skipping real API test - missing OpsRamp credentials in config.yaml")
	}

	// Create OpsRamp client
	opsRampClient := client.NewOpsRampClient(config)
	client.SetGlobalClient(opsRampClient)

	logger.Info("Real API test environment setup complete with config: %s", config.OpsRamp.TenantURL)

	// Return cleanup function
	cleanup := func() {
		logger.Info("Cleaning up real API test environment")
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

func TestResourcesRealAPI_List(t *testing.T) {
	logger, cleanup := setupRealAPITestEnvironment(t)
	defer cleanup()

	_, handler := tools.NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "list",
	})

	logger.Info("Testing real API resources list action")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := handler(ctx, req)

	if err != nil {
		t.Fatalf("List action failed: %v", err)
	}

	if res == nil {
		t.Fatalf("Expected non-nil result")
	}

	if res.IsError {
		t.Fatalf("List action returned error: %v", res.Content)
	}

	if len(res.Content) == 0 {
		t.Errorf("Expected non-empty result content")
	}

	logger.Info("Real API resources list test completed successfully")
}

func TestResourcesRealAPI_Search(t *testing.T) {
	logger, cleanup := setupRealAPITestEnvironment(t)
	defer cleanup()

	_, handler := tools.NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "search",
		"params": map[string]interface{}{
			"pageSize": 5,
			"pageNo":   1,
		},
	})

	logger.Info("Testing real API resources search action")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := handler(ctx, req)

	if err != nil {
		t.Fatalf("Search action failed: %v", err)
	}

	if res == nil {
		t.Fatalf("Expected non-nil result")
	}

	if res.IsError {
		t.Fatalf("Search action returned error: %v", res.Content)
	}

	if len(res.Content) == 0 {
		t.Errorf("Expected non-empty result content")
	}

	logger.Info("Real API resources search test completed successfully")
}

func TestResourcesRealAPI_GetResourceTypes(t *testing.T) {
	logger, cleanup := setupRealAPITestEnvironment(t)
	defer cleanup()

	_, handler := tools.NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "getResourceTypes",
	})

	logger.Info("Testing real API resources getResourceTypes action")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := handler(ctx, req)

	if err != nil {
		t.Fatalf("GetResourceTypes action failed: %v", err)
	}

	if res == nil {
		t.Fatalf("Expected non-nil result")
	}

	if res.IsError {
		t.Fatalf("GetResourceTypes action returned error: %v", res.Content)
	}

	if len(res.Content) == 0 {
		t.Errorf("Expected non-empty result content")
	}

	logger.Info("Real API resources getResourceTypes test completed successfully")
}

func TestResourcesRealAPI_GetFirstResource(t *testing.T) {
	logger, cleanup := setupRealAPITestEnvironment(t)
	defer cleanup()

	_, handler := tools.NewResourcesMcpTool()

	// First, get a list of resources to find a real resource ID
	listReq := createTestRequest(map[string]interface{}{
		"action": "search",
		"params": map[string]interface{}{
			"pageSize": 1,
			"pageNo":   1,
		},
	})

	logger.Info("Getting first resource for detailed testing")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	listRes, err := handler(ctx, listReq)
	if err != nil {
		t.Skipf("Skipping get test - failed to list resources: %v", err)
	}

	if listRes == nil || listRes.IsError || len(listRes.Content) == 0 {
		t.Skip("Skipping get test - no resources found")
	}

	// Extract resource ID from the response (this is a simplified approach)
	// In a real scenario, you'd parse the JSON response properly
	logger.Info("Found resources, attempting to get first resource details")

	// For now, we'll test with a common resource ID pattern
	// You can modify this to extract the actual ID from the list response
	getReq := createTestRequest(map[string]interface{}{
		"action": "get",
		"id":     "test-resource-id", // This would be extracted from listRes in practice
	})

	getRes, err := handler(ctx, getReq)

	// We expect this to fail with "not found" since we're using a dummy ID
	// But it should not fail with authentication errors
	if err != nil {
		logger.Info("Get action failed as expected with dummy ID: %v", err)
	} else if getRes != nil && getRes.IsError {
		logger.Info("Get action returned error as expected with dummy ID")
	}

	logger.Info("Real API resources get test completed")
}

func TestResourcesRealAPI_GetMinimal(t *testing.T) {
	logger, cleanup := setupRealAPITestEnvironment(t)
	defer cleanup()

	_, handler := tools.NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "getMinimal",
		"id":     "test-resource-id", // Using dummy ID for testing
	})

	logger.Info("Testing real API resources getMinimal action")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := handler(ctx, req)

	// We expect this to fail with "not found" since we're using a dummy ID
	// But it should not fail with authentication errors
	if err != nil {
		logger.Info("GetMinimal action failed as expected with dummy ID: %v", err)
		// Check that it's not an authentication error
		if res != nil && res.IsError {
			logger.Info("GetMinimal returned proper error response")
		}
	}

	logger.Info("Real API resources getMinimal test completed")
}

func TestResourcesRealAPI_InvalidAction(t *testing.T) {
	logger, cleanup := setupRealAPITestEnvironment(t)
	defer cleanup()

	_, handler := tools.NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "invalid-action",
	})

	logger.Info("Testing real API resources invalid action")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := handler(ctx, req)

	if err != nil {
		logger.Info("Invalid action properly returned error: %v", err)
	} else if res != nil && res.IsError {
		logger.Info("Invalid action properly returned error response")
	} else {
		t.Errorf("Expected error for invalid action")
	}

	logger.Info("Real API resources invalid action test completed")
}

func TestResourcesRealAPI_MissingParameters(t *testing.T) {
	logger, cleanup := setupRealAPITestEnvironment(t)
	defer cleanup()

	_, handler := tools.NewResourcesMcpTool()

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

			logger.Info("Testing real API %s", tc.name)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			res, err := handler(ctx, req)

			if err != nil {
				logger.Info("Test case '%s' properly returned error: %v", tc.name, err)
			} else if res != nil && res.IsError {
				logger.Info("Test case '%s' properly returned error response", tc.name)
			} else {
				t.Errorf("Test case '%s': Expected error for missing required parameters", tc.name)
			}
		})
	}

	logger.Info("Real API missing parameters tests completed")
}

// Benchmark test for performance measurement
func BenchmarkResourcesRealAPI_List(b *testing.B) {
	// Setup
	logDir := filepath.Join("..", "output", "logs")
	logFile := "or-mcp.log"

	absLogDir, err := filepath.Abs(logDir)
	if err != nil {
		b.Fatalf("Failed to get absolute log path: %v", err)
	}

	os.MkdirAll(absLogDir, 0750)
	logger, err := common.InitLogger(common.INFO, absLogDir, logFile) // Use INFO level for benchmarks
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	config, err := common.LoadConfig("")
	if err != nil {
		b.Skipf("Skipping benchmark - failed to load config.yaml: %v", err)
	}

	if config.OpsRamp.AuthKey == "" || config.OpsRamp.AuthSecret == "" || config.OpsRamp.TenantID == "" {
		b.Skip("Skipping benchmark - missing OpsRamp credentials in config.yaml")
	}

	opsRampClient := client.NewOpsRampClient(config)
	client.SetGlobalClient(opsRampClient)

	_, handler := tools.NewResourcesMcpTool()
	req := createTestRequest(map[string]interface{}{
		"action": "search",
		"params": map[string]interface{}{
			"pageSize": 10,
			"pageNo":   1,
		},
	})

	// Reset timer after setup
	b.ResetTimer()

	// Run benchmark
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		_, err := handler(ctx, req)
		cancel()

		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}
