package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/internal/adapters"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/types"
)

// ResourcesTool provides methods for managing OpsRamp resources
type ResourcesTool struct {
	adapter *adapters.ResourcesAdapter
	logger  *common.CustomLogger
}

// NewResourcesMcpTool creates a new MCP tool for managing OpsRamp resources
func NewResourcesMcpTool() (mcp.Tool, server.ToolHandlerFunc) {
	// Get the OpsRamp client
	opsRampClient := client.GetOpsRampClient()

	// Create the adapter
	adapter := adapters.NewResourcesAdapter(opsRampClient)

	// Get the logger
	logger := common.GetLogger()

	// Create the tool
	resourcesTool := &ResourcesTool{
		adapter: adapter,
		logger:  logger,
	}

	// Define the tool schema
	tool := mcp.Tool{
		Name:        "resources",
		Description: "Manage OpsRamp resources (devices, groups, etc.)",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"description": "Action to perform on resources",
					"enum": []string{
						"list", "get", "getDetailed", "create", "update", "delete", "search",
						"bulkUpdate", "bulkDelete", "getResourceTypes", "changeState",
						"getMetrics", "getTags", "updateTags",
					},
				},
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Resource ID (for get, getDetailed, update, delete, changeState, getMetrics, getTags, updateTags)",
				},
				"config": map[string]interface{}{
					"type":        "object",
					"description": "Resource configuration (for create and update)",
				},
				"params": map[string]interface{}{
					"type":        "object",
					"description": "Search parameters (for search)",
				},
				"resourceIds": map[string]interface{}{
					"type":        "array",
					"description": "List of resource IDs (for bulkUpdate, bulkDelete)",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"updates": map[string]interface{}{
					"type":        "object",
					"description": "Updates to apply to resources (for bulkUpdate)",
				},
				"state": map[string]interface{}{
					"type":        "string",
					"description": "State to set for a resource (for changeState)",
				},
				"metricRequest": map[string]interface{}{
					"type":        "object",
					"description": "Metrics request parameters (for getMetrics)",
				},
				"tags": map[string]interface{}{
					"type":        "array",
					"description": "Tags to set for a resource (for updateTags)",
					"items": map[string]interface{}{
						"type": "object",
					},
				},
			},
			Required: []string{"action"},
		},
	}

	return tool, resourcesToolHandler(resourcesTool)
}

// resourcesToolHandler handles MCP tool requests for resources
func resourcesToolHandler(tool *ResourcesTool) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args := req.Params.Arguments
		action, _ := args["action"].(string)
		id, _ := args["id"].(string)
		config, _ := args["config"].(map[string]interface{})
		params, _ := args["params"].(map[string]interface{})
		resourceIdsRaw, _ := args["resourceIds"].([]interface{})
		updates, _ := args["updates"].(map[string]interface{})
		state, _ := args["state"].(string)
		metricRequest, _ := args["metricRequest"].(map[string]interface{})
		tagsRaw, _ := args["tags"].([]interface{})

		// Convert resourceIds from []interface{} to []string
		var resourceIds []string
		if resourceIdsRaw != nil {
			resourceIds = make([]string, len(resourceIdsRaw))
			for i, v := range resourceIdsRaw {
				if s, ok := v.(string); ok {
					resourceIds[i] = s
				}
			}
		}

		// Convert tags from []interface{} to []types.Tag
		var tags []types.Tag
		if tagsRaw != nil {
			tags = make([]types.Tag, 0, len(tagsRaw))
			for _, v := range tagsRaw {
				if tagMap, ok := v.(map[string]interface{}); ok {
					name, _ := tagMap["name"].(string)
					value, _ := tagMap["value"].(string)
					if name != "" {
						tags = append(tags, types.Tag{
							Name:  name,
							Value: value,
						})
					}
				}
			}
		}

		// Log the tool execution
		tool.logger.LogToolExecution("resources", action, args)

		var result interface{}
		var err error

		// Get the OpsRamp client
		opsRampClient := client.GetOpsRampClient()

		// Create the resources API
		resourcesAPI := NewOpsRampResourcesAPI(opsRampClient)

		// Execute the requested action
		switch action {
		case "list":
			tool.logger.Info("Executing List resources")
			// List is just a search with default parameters
			searchParams := types.ResourceSearchParams{
				PageSize: 100,
				PageNo:   1,
			}
			result, err = resourcesAPI.Search(ctx, searchParams)
		case "get":
			tool.logger.Info("Executing Get resource with ID: %s", id)
			result, err = resourcesAPI.Get(ctx, id)
		case "getDetailed":
			tool.logger.Info("Executing GetDetailed resource with ID: %s", id)
			result, err = resourcesAPI.GetDetailed(ctx, id)
		case "create":
			tool.logger.Info("Executing Create resource")
			// Convert config to ResourceCreateRequest
			var createRequest types.ResourceCreateRequest
			configJSON, _ := json.Marshal(config)
			if err := json.Unmarshal(configJSON, &createRequest); err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to parse create request: %v", err)}},
				}, nil
			}
			result, err = resourcesAPI.Create(ctx, createRequest)
		case "update":
			tool.logger.Info("Executing Update resource with ID: %s", id)
			// Convert config to ResourceUpdateRequest
			var updateRequest types.ResourceUpdateRequest
			configJSON, _ := json.Marshal(config)
			if err := json.Unmarshal(configJSON, &updateRequest); err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to parse update request: %v", err)}},
				}, nil
			}
			result, err = resourcesAPI.Update(ctx, id, updateRequest)
		case "delete":
			tool.logger.Info("Executing Delete resource with ID: %s", id)
			err = resourcesAPI.Delete(ctx, id)
		case "search":
			tool.logger.Info("Executing Search resources with parameters")
			// Convert params to ResourceSearchParams
			var searchParams types.ResourceSearchParams
			paramsJSON, _ := json.Marshal(params)
			if err := json.Unmarshal(paramsJSON, &searchParams); err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to parse search parameters: %v", err)}},
				}, nil
			}
			result, err = resourcesAPI.Search(ctx, searchParams)
		case "bulkUpdate":
			tool.logger.Info("Executing BulkUpdate resources")
			if len(resourceIds) == 0 {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "No resource IDs provided for bulk update"}},
				}, nil
			}
			if updates == nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "No updates provided for bulk update"}},
				}, nil
			}
			bulkUpdateRequest := types.ResourceBulkUpdateRequest{
				ResourceIDs: resourceIds,
				Updates:     updates,
			}
			err = resourcesAPI.BulkUpdate(ctx, bulkUpdateRequest)
		case "bulkDelete":
			tool.logger.Info("Executing BulkDelete resources")
			if len(resourceIds) == 0 {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "No resource IDs provided for bulk delete"}},
				}, nil
			}
			bulkDeleteRequest := types.ResourceBulkDeleteRequest{
				ResourceIDs: resourceIds,
			}
			err = resourcesAPI.BulkDelete(ctx, bulkDeleteRequest)
		case "getResourceTypes":
			tool.logger.Info("Executing GetResourceTypes")
			result, err = resourcesAPI.GetResourceTypes(ctx)
		case "changeState":
			tool.logger.Info("Executing ChangeState for resource %s to %s", id, state)
			if id == "" {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "No resource ID provided for change state"}},
				}, nil
			}
			if state == "" {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "No state provided for change state"}},
				}, nil
			}
			stateChangeRequest := types.ResourceStateChangeRequest{
				State: state,
			}
			err = resourcesAPI.ChangeState(ctx, id, stateChangeRequest)
		case "getMetrics":
			tool.logger.Info("Executing GetMetrics for resource %s", id)
			if id == "" {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "No resource ID provided for get metrics"}},
				}, nil
			}
			// Convert metricRequest to ResourceMetricsRequest
			var metricsRequest types.ResourceMetricsRequest
			metricRequestJSON, _ := json.Marshal(metricRequest)
			if err := json.Unmarshal(metricRequestJSON, &metricsRequest); err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: fmt.Sprintf("Failed to parse metrics request: %v", err)}},
				}, nil
			}
			result, err = resourcesAPI.GetMetrics(ctx, id, metricsRequest)
		case "getTags":
			tool.logger.Info("Executing GetTags for resource %s", id)
			if id == "" {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "No resource ID provided for get tags"}},
				}, nil
			}
			result, err = resourcesAPI.GetTags(ctx, id)
		case "updateTags":
			tool.logger.Info("Executing UpdateTags for resource %s", id)
			if id == "" {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "No resource ID provided for update tags"}},
				}, nil
			}
			if len(tags) == 0 {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{mcp.TextContent{Type: "text", Text: "No tags provided for update tags"}},
				}, nil
			}
			err = resourcesAPI.UpdateTags(ctx, id, tags)
		default:
			err = fmt.Errorf("unknown action: %s", action)
		}

		// Log the tool result
		tool.logger.LogToolResult("resources", action, result, err)

		// Handle errors
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
}
