package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
)

// Test constants
const (
	// TestLogDir is the directory where logs will be stored
	TestLogDir = "output/logs"
	// TestLogFileName is the name of the log file
	TestLogFileName = "test-api-methods.log"
)

// Endpoint variants to test
var resourceEndpoints = []string{
	"/resources",
	"/resources/search",
	"/discovery/resources",
	"/discovery/resources/search",
	"/inventory/resources",
	"/inventory/resources/search",
	"/monitoring/resources",
	"/monitoring/resources/search",
}

// HTTP methods to test
var httpMethods = []string{
	"GET",
	"POST",
}

// Search parameters for POST requests
var searchParams = map[string]interface{}{
	"pageSize":          10,
	"pageNo":            1,
	"isDescendingOrder": false,
	"sortName":          "name",
}

func main() {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(TestLogDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize the logger
	customLogger, err := common.InitLogger(common.DEBUG, TestLogDir, TestLogFileName)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer customLogger.Close()

	// Log script start
	customLogger.Info("Starting API methods test script")
	customLogger.Info("Log file: %s", filepath.Join(TestLogDir, TestLogFileName))

	// Load configuration
	configPath := "../../config.yaml"
	customLogger.Info("Trying to load config from: %s", configPath)
	config, err := common.LoadConfig(configPath)
	if err != nil {
		customLogger.Warn("Failed to load config from %s: %v", configPath, err)
		fmt.Printf("Failed to load config from %s: %v\n", configPath, err)
		
		// Try alternate locations
		configPath = "../config.yaml"
		customLogger.Info("Trying to load config from: %s", configPath)
		config, err = common.LoadConfig(configPath)
		if err != nil {
			customLogger.Error("Failed to load config: %v", err)
			fmt.Printf("Failed to load config: %v\n", err)
			os.Exit(1)
		}
	}

	// Create OpsRamp client
	opsRampClient := client.NewOpsRampClient(config)
	client.SetGlobalClient(opsRampClient)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test each endpoint with each method
	for _, endpoint := range resourceEndpoints {
		for _, method := range httpMethods {
			testEndpoint(ctx, customLogger, opsRampClient, endpoint, method)
		}
	}

	// Log script completion
	customLogger.Info("API methods test script completed successfully")
	fmt.Println("Test completed successfully. Check the log file for details:", filepath.Join(TestLogDir, TestLogFileName))
}

func testEndpoint(ctx context.Context, logger *common.CustomLogger, client *client.OpsRampClient, endpoint string, method string) {
	// Build the full endpoint
	fullEndpoint := fmt.Sprintf("/api/v2/tenants/%s%s", client.GetTenantID(), endpoint)
	logger.Info("Testing %s %s", method, fullEndpoint)
	fmt.Printf("Testing %s %s\n", method, fullEndpoint)

	var err error
	var statusCode int
	var responseBody interface{}

	// Execute the request based on the method
	switch method {
	case "GET":
		statusCode, responseBody, err = executeGet(ctx, client, fullEndpoint)
	case "POST":
		statusCode, responseBody, err = executePost(ctx, client, fullEndpoint, searchParams)
	}

	// Log the result
	if err != nil {
		logger.Error("Request failed: %v", err)
		fmt.Printf("  Result: FAILED - %v\n", err)
	} else {
		if statusCode >= 200 && statusCode < 300 {
			logger.Info("Request succeeded with status code %d", statusCode)
			fmt.Printf("  Result: SUCCESS (%d)\n", statusCode)
			
			// Print response summary
			responseJSON, _ := json.MarshalIndent(responseBody, "", "  ")
			logger.Info("Response: %s", string(responseJSON))
			fmt.Printf("  Response: %s\n", summarizeResponse(responseBody))
		} else {
			logger.Warn("Request returned non-success status code %d", statusCode)
			fmt.Printf("  Result: FAILED (%d)\n", statusCode)
		}
	}
	fmt.Println()
}

func executeGet(ctx context.Context, client *client.OpsRampClient, endpoint string) (int, interface{}, error) {
	var response interface{}
	statusCode, err := client.GetWithStatusCode(ctx, endpoint, &response)
	return statusCode, response, err
}

func executePost(ctx context.Context, client *client.OpsRampClient, endpoint string, body interface{}) (int, interface{}, error) {
	var response interface{}
	statusCode, err := client.PostWithStatusCode(ctx, endpoint, body, &response)
	return statusCode, response, err
}

func summarizeResponse(response interface{}) string {
	if response == nil {
		return "empty response"
	}

	// Try to extract common fields from the response
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return "non-map response"
	}

	// Check for results array
	if results, ok := responseMap["results"].([]interface{}); ok {
		return fmt.Sprintf("%d results found", len(results))
	}

	// Check for common fields
	fields := []string{"id", "name", "type", "count", "total"}
	summary := make(map[string]interface{})
	for _, field := range fields {
		if value, ok := responseMap[field]; ok {
			summary[field] = value
		}
	}

	if len(summary) > 0 {
		summaryJSON, _ := json.Marshal(summary)
		return string(summaryJSON)
	}

	// If we can't extract meaningful summary, return the first few keys
	keys := make([]string, 0, len(responseMap))
	for k := range responseMap {
		keys = append(keys, k)
		if len(keys) >= 3 {
			break
		}
	}
	return fmt.Sprintf("keys: %v", keys)
}
