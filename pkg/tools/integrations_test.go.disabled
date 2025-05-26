package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestIntegrationsTool_List_Create_Get_Update_Delete(t *testing.T) {
	api := &MockIntegrationsAPI{}
	_, handler := createIntegrationsTool(api)
	ctx := context.Background()

	// Create
	createConfig := map[string]interface{}{"name": "test-integration", "type": "api"}
	createReq := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: map[string]interface{}{"action": "create", "config": createConfig},
		},
	}
	createRes, err := handler(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if createRes == nil || len(createRes.Content) == 0 {
		t.Fatalf("Expected non-empty result for create")
	}

	// List
	listReq := mcp.CallToolRequest{
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
	listRes, err := handler(ctx, listReq)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if listRes == nil || len(listRes.Content) == 0 {
		t.Errorf("Expected non-empty result for list")
	}

	// Get (using the first integration's ID from the mock)
	firstID := "mockid"
	getReq := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: map[string]interface{}{"action": "get", "id": firstID},
		},
	}
	getRes, err := handler(ctx, getReq)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if getRes == nil || len(getRes.Content) == 0 {
		t.Errorf("Expected non-empty result for get")
	}

	// Update
	updateConfig := map[string]interface{}{"name": "updated-integration"}
	updateReq := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: map[string]interface{}{"action": "update", "id": firstID, "config": updateConfig},
		},
	}
	updateRes, err := handler(ctx, updateReq)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updateRes == nil || len(updateRes.Content) == 0 {
		t.Errorf("Expected non-empty result for update")
	}
}

func TestIntegrationsTool_ListTypes_GetType(t *testing.T) {
	_, handler := NewIntegrationsMcpTool()
	ctx := context.Background()

	// ListTypes
	listTypesReq := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: map[string]interface{}{"action": "listTypes"},
		},
	}
	listTypesRes, err := handler(ctx, listTypesReq)
	if err != nil {
		t.Fatalf("ListTypes failed: %v", err)
	}
	if listTypesRes == nil || len(listTypesRes.Content) == 0 {
		t.Errorf("Expected non-empty result for listTypes")
	}

	// GetType (existing)
	getTypeReq := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: map[string]interface{}{"action": "getType", "id": "api"},
		},
	}
	getTypeRes, err := handler(ctx, getTypeReq)
	if err != nil {
		t.Fatalf("GetType failed: %v", err)
	}
	if getTypeRes == nil || len(getTypeRes.Content) == 0 {
		t.Errorf("Expected non-empty result for getType")
	}

	// GetType (not found)
	getTypeReqNotFound := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Arguments: map[string]interface{}{"action": "getType", "id": "notfound"},
		},
	}
	result, err := handler(ctx, getTypeReqNotFound)
	if err != nil {
		t.Fatalf("Handler returned Go error: %v", err)
	}
	if result == nil || !result.IsError {
		t.Errorf("Expected error result for getType non-existent id")
	} else {
		found := false
		for _, c := range result.Content {
			if text, ok := c.(mcp.TextContent); ok &&
				text.Text != "" &&
				(text.Text == "integration type with ID notfound not found" ||
					// allow substring match for flexibility
					contains(text.Text, "notfound") && contains(text.Text, "not found")) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected error message for getType non-existent id, got: %+v", result.Content)
		}
	}
}

// contains is a helper for substring matching
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (contains(s[1:], substr) || contains(s[:len(s)-1], substr))))
}
