package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/tools"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/types"
)

// Test constants
const (
	// TestLogDir is the directory where logs will be stored
	TestLogDir = "output/logs"
	// TestLogFileName is the name of the log file
	TestLogFileName = "test-resources.log"
)

// main is the entry point for the test script
func main() {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(TestLogDir, 0750); err != nil {
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
	customLogger.Info("Starting resources API test script")
	customLogger.Info("Log file: %s", filepath.Join(TestLogDir, TestLogFileName))

	// Load configuration
	// Try to load from the project root directory
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
			customLogger.Warn("Failed to load config from %s: %v", configPath, err)
			fmt.Printf("Failed to load config from %s: %v\n", configPath, err)

			// Create a minimal default config
			customLogger.Info("Using default configuration")
			fmt.Println("Using default configuration")
			config = &common.Config{
				OpsRamp: common.OpsRampConfig{
					TenantURL:  common.GetEnvOrDefault("OPSRAMP_TENANT_URL", "https://api.opsramp.com"),
					AuthURL:    common.GetEnvOrDefault("OPSRAMP_AUTH_URL", "https://api.opsramp.com/auth/oauth/token"),
					AuthKey:    common.GetEnvOrDefault("OPSRAMP_AUTH_KEY", ""),
					AuthSecret: common.GetEnvOrDefault("OPSRAMP_AUTH_SECRET", ""),
					TenantID:   common.GetEnvOrDefault("OPSRAMP_TENANT_ID", ""),
				},
			}
		}
	}

	// Create OpsRamp client
	opsRampClient := client.NewOpsRampClient(config)
	client.SetGlobalClient(opsRampClient)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create the resources tool
	resourcesTool, _ := tools.NewResourcesMcpTool()
	customLogger.Info("Created resources tool: %s", resourcesTool.Name)

	// Test listing resources
	testListResources(ctx, customLogger)

	// If we have resources, test getting a specific resource
	testGetResource(ctx, customLogger)

	// Test searching for resources
	testSearchResources(ctx, customLogger)

	// Log script completion
	customLogger.Info("Resources API test script completed successfully")
	fmt.Println("Test completed successfully. Check the log file for details:", filepath.Join(TestLogDir, TestLogFileName))
}

func testListResources(ctx context.Context, logger *common.CustomLogger) {
	logger.Info("Testing List Resources operation")
	fmt.Println("Testing List Resources operation...")

	// Get the OpsRamp client
	opsRampClient := client.GetOpsRampClient()

	// Create the resources API
	resourcesAPI := tools.NewOpsRampResourcesAPI(opsRampClient)

	// List resources
	params := types.ResourceSearchParams{
		PageSize: 10, // Limit to 10 resources for the test
		PageNo:   1,
		// Add any additional parameters as needed
	}

	response, err := resourcesAPI.Search(ctx, params)
	if err != nil {
		logger.Error("Failed to list resources: %v", err)
		fmt.Printf("Failed to list resources: %v\n", err)
		return
	}

	logger.Info("Successfully listed %d resources", len(response.Results))
	fmt.Printf("Successfully listed %d resources\n", len(response.Results))

	// Print the first few resources
	for i, resource := range response.Results {
		if i >= 3 {
			break // Only show the first 3 resources
		}
		logger.Info("Resource %d: ID=%s, Name=%s, Type=%s", i+1, resource.ID, resource.Name, resource.Type)
		fmt.Printf("Resource %d: ID=%s, Name=%s, Type=%s\n", i+1, resource.ID, resource.Name, resource.Type)
	}
}

func testGetResource(ctx context.Context, logger *common.CustomLogger) {
	logger.Info("Testing Get Resource operation")
	fmt.Println("\nTesting Get Resource operation...")

	// Get the OpsRamp client
	opsRampClient := client.GetOpsRampClient()

	// Create the resources API
	resourcesAPI := tools.NewOpsRampResourcesAPI(opsRampClient)

	// First, list resources to get an ID
	params := types.ResourceSearchParams{
		PageSize: 1, // Just get one resource
		PageNo:   1,
	}

	response, err := resourcesAPI.Search(ctx, params)
	if err != nil {
		logger.Error("Failed to list resources: %v", err)
		fmt.Printf("Failed to list resources: %v\n", err)
		return
	}

	if len(response.Results) == 0 {
		logger.Warn("No resources found to test Get operation")
		fmt.Println("No resources found to test Get operation")
		return
	}

	// Get the first resource ID
	resourceID := response.Results[0].ID
	logger.Info("Testing Get with resource ID: %s", resourceID)
	fmt.Printf("Testing Get with resource ID: %s\n", resourceID)

	// Get the resource
	resource, err := resourcesAPI.Get(ctx, resourceID)
	if err != nil {
		logger.Error("Failed to get resource %s: %v", resourceID, err)
		fmt.Printf("Failed to get resource %s: %v\n", resourceID, err)
		return
	}

	logger.Info("Successfully retrieved resource: %s", resource.Name)
	fmt.Printf("Successfully retrieved resource: %s\n", resource.Name)
	fmt.Printf("Resource details: ID=%s, Name=%s, Type=%s, Status=%s\n",
		resource.ID, resource.Name, resource.Type, resource.Status)
}

func testSearchResources(ctx context.Context, logger *common.CustomLogger) {
	logger.Info("Testing Search Resources operation")
	fmt.Println("\nTesting Search Resources operation...")

	// Get the OpsRamp client
	opsRampClient := client.GetOpsRampClient()

	// Create the resources API
	resourcesAPI := tools.NewOpsRampResourcesAPI(opsRampClient)

	// Create search parameters
	params := types.ResourceSearchParams{
		PageSize:          5,
		PageNo:            1,
		IsDescendingOrder: true,
		SortName:          "name",
	}

	// Add a filter
	params.ResourceType = "Server" // Example filter
	logger.Info("Searching for resources of type: %s", params.ResourceType)
	fmt.Printf("Searching for resources of type: %s\n", params.ResourceType)

	response, err := resourcesAPI.Search(ctx, params)
	if err != nil {
		logger.Error("Failed to search resources: %v", err)
		fmt.Printf("Failed to search resources: %v\n", err)
		return
	}

	logger.Info("Successfully searched resources, found %d results", len(response.Results))
	fmt.Printf("Successfully searched resources, found %d results\n", len(response.Results))

	// Print the search results
	for i, resource := range response.Results {
		if i >= 3 {
			break // Only show the first 3 resources
		}
		logger.Info("Search Result %d: ID=%s, Name=%s, Type=%s", i+1, resource.ID, resource.Name, resource.Type)
		fmt.Printf("Search Result %d: ID=%s, Name=%s, Type=%s\n", i+1, resource.ID, resource.Name, resource.Type)
	}
}
