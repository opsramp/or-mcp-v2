package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type DevicesTool struct{}

// NewDevicesMcpTool returns the MCP tool definition and handler for devices
func NewDevicesMcpTool() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.Tool{
		Name:        "devices",
		Description: "Manage HPE OpsRamp devices and their configurations.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"description": "Action to perform: list, get, create, update, delete, enable, disable, listTypes, getType",
				},
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Device ID (for get, update, delete, enable, disable, getType)",
				},
				"config": map[string]interface{}{
					"type":        "object",
					"description": "Device configuration (for create and update)",
				},
			},
			Required: []string{"action"},
		},
	}, devicesToolHandler
}

func devicesToolHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	action, _ := args["action"].(string)
	id, _ := args["id"].(string)
	config, _ := args["config"].(map[string]interface{})

	tool := &DevicesTool{}
	var err error
	var result interface{}

	switch action {
	case "list":
		result, err = tool.List(ctx)
	case "get":
		result, err = tool.Get(ctx, id)
	case "create":
		result, err = tool.Create(ctx, config)
	case "update":
		result, err = tool.Update(ctx, id, config)
	case "delete":
		err = tool.Delete(ctx, id)
	case "enable":
		err = tool.Enable(ctx, id)
	case "disable":
		err = tool.Disable(ctx, id)
	case "listTypes":
		result, err = tool.ListTypes(ctx)
	case "getType":
		result, err = tool.GetType(ctx, id)
	default:
		err = server.ErrToolNotFound
	}

	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}},
		}, nil
	}

	// Convert result to string if it exists
	resultText := "OK"
	if result != nil {
		resultText = fmt.Sprintf("%v", result)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}},
	}, nil
}

// Implementation stubs for actual OpsRamp logic
func (dt *DevicesTool) List(ctx context.Context) ([]interface{}, error) {
	// TODO: Implement list devices
	return []interface{}{}, nil
}
func (dt *DevicesTool) Get(ctx context.Context, id string) (interface{}, error) {
	// TODO: Implement get device
	return struct{}{}, nil
}
func (dt *DevicesTool) Create(ctx context.Context, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement create device
	return struct{}{}, nil
}
func (dt *DevicesTool) Update(ctx context.Context, id string, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement update device
	return struct{}{}, nil
}
func (dt *DevicesTool) Delete(ctx context.Context, id string) error {
	// TODO: Implement delete device
	return nil
}
func (dt *DevicesTool) Enable(ctx context.Context, id string) error {
	// TODO: Implement enable device
	return nil
}
func (dt *DevicesTool) Disable(ctx context.Context, id string) error {
	// TODO: Implement disable device
	return nil
}
func (dt *DevicesTool) ListTypes(ctx context.Context) ([]interface{}, error) {
	// TODO: Implement list device types
	return []interface{}{}, nil
}
func (dt *DevicesTool) GetType(ctx context.Context, id string) (interface{}, error) {
	// TODO: Implement get device type
	return struct{}{}, nil
}
