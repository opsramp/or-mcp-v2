package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestResourcesTool_List(t *testing.T) {
	_, handler := NewResourcesMcpTool()
	req := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: map[string]interface{}{"action": "list"},
		},
	}
	res, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if res == nil || len(res.Content) == 0 {
		t.Errorf("Expected non-empty result for list")
	}
}

func TestResourcesTool_Get(t *testing.T) {
	_, handler := NewResourcesMcpTool()
	req := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: map[string]interface{}{"action": "get", "id": "test-id"},
		},
	}
	_, err := handler(context.Background(), req)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
}

func TestResourcesTool_Create(t *testing.T) {
	_, handler := NewResourcesMcpTool()
	config := map[string]interface{}{"name": "test-resource"}
	req := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: map[string]interface{}{"action": "create", "config": config},
		},
	}
	_, err := handler(context.Background(), req)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}
}

func TestResourcesTool_Update(t *testing.T) {
	_, handler := NewResourcesMcpTool()
	config := map[string]interface{}{"name": "updated-resource"}
	req := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: map[string]interface{}{"action": "update", "id": "test-id", "config": config},
		},
	}
	_, err := handler(context.Background(), req)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}
}

func TestResourcesTool_Delete(t *testing.T) {
	_, handler := NewResourcesMcpTool()
	req := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: map[string]interface{}{"action": "delete", "id": "test-id"},
		},
	}
	_, err := handler(context.Background(), req)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
}
