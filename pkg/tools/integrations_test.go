package tools

import (
	"context"
	"testing"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/types"
)

// MockIntegrationsAPI is a mock implementation for unit testing
// You can extend this with more sophisticated behavior as needed.
type MockIntegrationsAPI struct {
	ListFunc      func(ctx context.Context) ([]types.Integration, error)
	GetFunc       func(ctx context.Context, id string) (*types.Integration, error)
	CreateFunc    func(ctx context.Context, config map[string]interface{}) (*types.Integration, error)
	UpdateFunc    func(ctx context.Context, id string, config map[string]interface{}) (*types.Integration, error)
	DeleteFunc    func(ctx context.Context, id string) error
	EnableFunc    func(ctx context.Context, id string) error
	DisableFunc   func(ctx context.Context, id string) error
	ListTypesFunc func(ctx context.Context) ([]types.IntegrationType, error)
	GetTypeFunc   func(ctx context.Context, id string) (*types.IntegrationType, error)
}

func (m *MockIntegrationsAPI) List(ctx context.Context) ([]types.Integration, error) {
	return m.ListFunc(ctx)
}
func (m *MockIntegrationsAPI) Get(ctx context.Context, id string) (*types.Integration, error) {
	return m.GetFunc(ctx, id)
}
func (m *MockIntegrationsAPI) Create(ctx context.Context, config map[string]interface{}) (*types.Integration, error) {
	return m.CreateFunc(ctx, config)
}
func (m *MockIntegrationsAPI) Update(ctx context.Context, id string, config map[string]interface{}) (*types.Integration, error) {
	return m.UpdateFunc(ctx, id, config)
}
func (m *MockIntegrationsAPI) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}
func (m *MockIntegrationsAPI) Enable(ctx context.Context, id string) error {
	return m.EnableFunc(ctx, id)
}
func (m *MockIntegrationsAPI) Disable(ctx context.Context, id string) error {
	return m.DisableFunc(ctx, id)
}
func (m *MockIntegrationsAPI) ListTypes(ctx context.Context) ([]types.IntegrationType, error) {
	return m.ListTypesFunc(ctx)
}
func (m *MockIntegrationsAPI) GetType(ctx context.Context, id string) (*types.IntegrationType, error) {
	return m.GetTypeFunc(ctx, id)
}

func NewMockIntegrationsAPI() *MockIntegrationsAPI {
	return &MockIntegrationsAPI{
		ListFunc: func(ctx context.Context) ([]types.Integration, error) {
			return []types.Integration{}, nil
		},
		GetFunc: func(ctx context.Context, id string) (*types.Integration, error) {
			return nil, nil
		},
		CreateFunc: func(ctx context.Context, config map[string]interface{}) (*types.Integration, error) {
			return &types.Integration{ID: "mockid"}, nil
		},
		UpdateFunc: func(ctx context.Context, id string, config map[string]interface{}) (*types.Integration, error) {
			return &types.Integration{ID: id}, nil
		},
		DeleteFunc: func(ctx context.Context, id string) error {
			return nil
		},
		EnableFunc: func(ctx context.Context, id string) error {
			return nil
		},
		DisableFunc: func(ctx context.Context, id string) error {
			return nil
		},
		ListTypesFunc: func(ctx context.Context) ([]types.IntegrationType, error) {
			return []types.IntegrationType{}, nil
		},
		GetTypeFunc: func(ctx context.Context, id string) (*types.IntegrationType, error) {
			return nil, nil
		},
	}
}

func TestIntegrationsTool_List_Create_Get_Update_Delete(t *testing.T) {
	api := NewMockIntegrationsAPI()
	_, handler := NewIntegrationsTool(api)
	ctx := context.Background()

	// Create
	createConfig := map[string]interface{}{"name": "test-integration", "type": "api"}
	createReq := mcp.CallToolRequest{
		Params: mcp.ToolCallParams{
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
		Params: mcp.ToolCallParams{
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
	var firstID string
	firstID = "mockid"
	getReq := mcp.CallToolRequest{
		Params: mcp.ToolCallParams{
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
		Params: mcp.ToolCallParams{
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

	// Enable
	enableReq := mcp.CallToolRequest{
		Params: mcp.ToolCallParams{
			Arguments: map[string]interface{}{"action": "enable", "id": firstID},
		},
	}
	enableRes, err := handler(ctx, enableReq)
	if err != nil {
		t.Fatalf("Enable failed: %v", err)
	}
	if enableRes == nil || len(enableRes.Content) == 0 {
		t.Errorf("Expected non-empty result for enable")
	}

	// Disable
	disableReq := mcp.CallToolRequest{
		Params: mcp.ToolCallParams{
			Arguments: map[string]interface{}{"action": "disable", "id": firstID},
		},
	}
	disableRes, err := handler(ctx, disableReq)
	if err != nil {
		t.Fatalf("Disable failed: %v", err)
	}
	if disableRes == nil || len(disableRes.Content) == 0 {
		t.Errorf("Expected non-empty result for disable")
	}

	// Delete
	deleteReq := mcp.CallToolRequest{
		Params: mcp.ToolCallParams{
			Arguments: map[string]interface{}{"action": "delete", "id": firstID},
		},
	}
	deleteRes, err := handler(ctx, deleteReq)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if deleteRes == nil || len(deleteRes.Content) == 0 {
		t.Errorf("Expected non-empty result for delete")
	}

	// Get (not found)
	getReqNotFound := mcp.CallToolRequest{
		Params: mcp.ToolCallParams{
			Arguments: map[string]interface{}{"action": "get", "id": "non-existent-id"},
		},
	}
	_, err = handler(ctx, getReqNotFound)
	if err == nil {
		t.Errorf("Expected error for get non-existent id")
	}
}

func TestIntegrationsTool_ListTypes_GetType(t *testing.T) {
	_, handler := NewIntegrationsMcpTool()
	ctx := context.Background()

	// ListTypes
	listTypesReq := mcp.CallToolRequest{
		Params: mcp.ToolCallParams{
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
		Params: mcp.ToolCallParams{
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
		Params: mcp.ToolCallParams{
			Arguments: map[string]interface{}{"action": "getType", "id": "notfound"},
		},
	}
	_, err = handler(ctx, getTypeReqNotFound)
	if err == nil {
		t.Errorf("Expected error for getType non-existent id")
	}
}
