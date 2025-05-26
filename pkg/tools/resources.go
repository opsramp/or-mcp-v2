package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/types"
)

type ResourcesTool struct {
	api    ResourcesAPI
	logger *common.CustomLogger
}

// NewResourcesTool creates a new ResourcesTool with the provided API implementation
func NewResourcesTool(api ResourcesAPI) *ResourcesTool {
	// Get the logger
	logger := common.GetLogger()

	return &ResourcesTool{
		api:    api,
		logger: logger,
	}
}

// NewResourcesMcpTool returns the MCP tool definition and handler for resources
func NewResourcesMcpTool() (mcp.Tool, server.ToolHandlerFunc) {
	// Get the logger
	logger := common.GetLogger()

	// Load configuration
	config, err := common.LoadConfig("")
	if err != nil {
		logger.Error("Failed to load config for OpsRamp Resources API: %v", err)
		// Return error instead of falling back to mock
		return mcp.Tool{}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Configuration error: %v", err)}},
			}, nil
		}
	}

	// Create and initialize the real API implementation
	api := NewOpsRampResourcesAPI(client.NewOpsRampClient(config))

	logger.Info("Successfully initialized OpsRamp Resources API")
	return createResourcesTool(api)
}

// createResourcesTool creates the MCP tool with the given API implementation
func createResourcesTool(api ResourcesAPI) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.Tool{
			Name:        "resources",
			Description: "Manage HPE OpsRamp resources (devices, servers, network equipment, etc.)",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "Action to perform: list, get, getDetailed, getMinimal, create, update, delete, search, getResourceTypes",
					},
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Resource ID (for get, getDetailed, getMinimal, update, delete)",
					},
					"config": map[string]interface{}{
						"type":        "object",
						"description": "Resource configuration (for create and update)",
					},
					"params": map[string]interface{}{
						"type":        "object",
						"description": "Search parameters (for search)",
					},
				},
				Required: []string{"action"},
			},
		}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return ResourcesToolHandler(ctx, req, api)
		}
}

// ResourcesToolHandler routes requests to the correct method
// Exported for testing purposes
func ResourcesToolHandler(ctx context.Context, req mcp.CallToolRequest, api ResourcesAPI) (*mcp.CallToolResult, error) {
	// Extract arguments using the helper methods
	action := req.GetString("action", "")
	id := req.GetString("id", "")

	// Get arguments as a map
	args := req.GetArguments()

	// Extract config map if it exists
	var config map[string]interface{}
	if configArg, exists := args["config"]; exists && configArg != nil {
		if configMap, ok := configArg.(map[string]interface{}); ok {
			config = configMap
		}
	}

	// Extract params map if it exists
	var params map[string]interface{}
	if paramsArg, exists := args["params"]; exists && paramsArg != nil {
		if paramsMap, ok := paramsArg.(map[string]interface{}); ok {
			params = paramsMap
		}
	}

	// Log the tool execution
	logger := common.GetLogger()
	logger.LogToolExecution("resources", action, args)

	var err error
	var result interface{}

	switch action {
	case "list":
		logger.Info("Executing List resources")
		// List is just a search with default parameters
		searchParams := types.ResourceSearchParams{
			PageSize: 100,
			PageNo:   1,
		}
		result, err = api.Search(ctx, searchParams)
	case "get":
		logger.Info("Executing Get resource with ID: %s", id)
		if id == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "Resource ID is required for get action"}},
			}, nil
		}
		result, err = api.Get(ctx, id)
	case "getDetailed":
		logger.Info("Executing GetDetailed resource with ID: %s", id)
		if id == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "Resource ID is required for getDetailed action"}},
			}, nil
		}
		result, err = api.GetDetailed(ctx, id)
	case "getMinimal":
		logger.Info("Executing GetMinimal resource with ID: %s", id)
		if id == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "Resource ID is required for getMinimal action"}},
			}, nil
		}
		result, err = api.GetMinimal(ctx, id)
	case "create":
		logger.Info("Executing Create resource")
		if config == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "Configuration is required for create action"}},
			}, nil
		}
		// Convert config to ResourceCreateRequest
		var createRequest types.ResourceCreateRequest
		configJSON, _ := json.Marshal(config)
		if err := json.Unmarshal(configJSON, &createRequest); err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to parse create request: %v", err)}},
			}, nil
		}
		result, err = api.Create(ctx, createRequest)
	case "update":
		logger.Info("Executing Update resource with ID: %s", id)
		if id == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "Resource ID is required for update action"}},
			}, nil
		}
		if config == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "Configuration is required for update action"}},
			}, nil
		}
		// Convert config to ResourceUpdateRequest
		var updateRequest types.ResourceUpdateRequest
		configJSON, _ := json.Marshal(config)
		if err := json.Unmarshal(configJSON, &updateRequest); err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to parse update request: %v", err)}},
			}, nil
		}
		result, err = api.Update(ctx, id, updateRequest)
	case "delete":
		logger.Info("Executing Delete resource with ID: %s", id)
		if id == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "Resource ID is required for delete action"}},
			}, nil
		}
		err = api.Delete(ctx, id)
	case "search":
		logger.Info("Executing Search resources with parameters")
		// Convert params to ResourceSearchParams
		var searchParams types.ResourceSearchParams
		if params != nil {
			paramsJSON, _ := json.Marshal(params)
			if err := json.Unmarshal(paramsJSON, &searchParams); err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to parse search parameters: %v", err)}},
				}, nil
			}
		} else {
			// Default search parameters
			searchParams = types.ResourceSearchParams{
				PageSize: 100,
				PageNo:   1,
			}
		}
		result, err = api.Search(ctx, searchParams)
	case "getResourceTypes":
		logger.Info("Executing GetResourceTypes")
		result, err = api.GetResourceTypes(ctx)
	default:
		logger.Error("Unknown action: %s", action)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Unknown action: %s", action)}},
		}, nil
	}

	// Log the result
	logger.LogToolResult("resources", action, result, err)

	// If there's an error, return it
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}},
		}, nil
	}

	// Return the result
	if result != nil {
		// Convert the result to JSON
		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to marshal result: %v", err)}},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{mcp.TextContent{Type: "text", Text: string(resultJSON)}},
		}, nil
	}

	// Return a simple success message for actions that don't return a result
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "Operation completed successfully"}},
	}, nil
}
