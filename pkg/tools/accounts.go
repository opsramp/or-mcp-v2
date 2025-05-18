package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type AccountsTool struct{}

func NewAccountsMcpTool() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.Tool{
		Name:        "accounts",
		Description: "Manage OpsRamp accounts and their configurations.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"description": "Action to perform: list, get, create, update, delete, enable, disable, listTypes, getType",
				},
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Account ID (for get, update, delete, enable, disable, getType)",
				},
				"config": map[string]interface{}{
					"type":        "object",
					"description": "Account configuration (for create and update)",
				},
			},
			Required: []string{"action"},
		},
	}, accountsToolHandler
}

func accountsToolHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	action, _ := args["action"].(string)
	id, _ := args["id"].(string)
	config, _ := args["config"].(map[string]interface{})

	tool := &AccountsTool{}
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
func (at *AccountsTool) List(ctx context.Context) ([]interface{}, error) {
	// TODO: Implement list accounts
	return []interface{}{}, nil
}
func (at *AccountsTool) Get(ctx context.Context, id string) (interface{}, error) {
	// TODO: Implement get account
	return struct{}{}, nil
}
func (at *AccountsTool) Create(ctx context.Context, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement create account
	return struct{}{}, nil
}
func (at *AccountsTool) Update(ctx context.Context, id string, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement update account
	return struct{}{}, nil
}
func (at *AccountsTool) Delete(ctx context.Context, id string) error {
	// TODO: Implement delete account
	return nil
}
func (at *AccountsTool) Enable(ctx context.Context, id string) error {
	// TODO: Implement enable account
	return nil
}
func (at *AccountsTool) Disable(ctx context.Context, id string) error {
	// TODO: Implement disable account
	return nil
}
func (at *AccountsTool) ListTypes(ctx context.Context) ([]interface{}, error) {
	// TODO: Implement list account types
	return []interface{}{}, nil
}
func (at *AccountsTool) GetType(ctx context.Context, id string) (interface{}, error) {
	// TODO: Implement get account type
	return struct{}{}, nil
}
