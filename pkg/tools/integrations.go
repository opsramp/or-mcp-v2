package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/types"
)

// IntegrationsTool provides MCP methods for managing integrations

// IntegrationsAPI defines the contract for integration operations
// Implementations: real API client, mocks, etc.
type IntegrationsAPI interface {
	List(ctx context.Context) ([]types.Integration, error)
	Get(ctx context.Context, id string) (*types.Integration, error)
	GetDetailed(ctx context.Context, id string) (*types.DetailedIntegration, error)
	Create(ctx context.Context, config map[string]interface{}) (*types.Integration, error)
	Update(ctx context.Context, id string, config map[string]interface{}) (*types.Integration, error)
	Delete(ctx context.Context, id string) error
	Enable(ctx context.Context, id string) error
	Disable(ctx context.Context, id string) error
	ListTypes(ctx context.Context) ([]types.IntegrationType, error)
	GetType(ctx context.Context, id string) (*types.IntegrationType, error)
}

type IntegrationsTool struct {
	api    IntegrationsAPI
	logger *common.CustomLogger
}

// NewIntegrationsTool creates a new IntegrationsTool with the provided API implementation
func NewIntegrationsTool(api IntegrationsAPI) *IntegrationsTool {
	// Get the logger
	logger := common.GetLogger()

	return &IntegrationsTool{
		api:    api,
		logger: logger,
	}
}

// NewIntegrationsMcpTool returns the MCP tool definition and handler for integrations
func NewIntegrationsMcpTool() (mcp.Tool, server.ToolHandlerFunc) {
	// Get the logger
	logger := common.GetLogger()

	// Load configuration
	config, err := common.LoadConfig("")
	if err != nil {
		logger.Error("Failed to load config for OpsRamp Integrations API: %v", err)
		logger.Warn("Falling back to mock implementation")
		mockAPI := &MockIntegrationsAPI{}
		return createIntegrationsTool(mockAPI)
	}

	// Create and initialize the real API implementation
	api, err := NewOpsRampIntegrationsAPI(&config.OpsRamp)
	if err != nil {
		logger.Error("Failed to initialize OpsRamp Integrations API: %v", err)
		logger.Warn("Falling back to mock implementation")
		// Fall back to mock implementation if initialization fails
		mockAPI := &MockIntegrationsAPI{}
		return createIntegrationsTool(mockAPI)
	}

	logger.Info("Successfully initialized OpsRamp Integrations API")
	return createIntegrationsTool(api)
}

// createIntegrationsTool creates the MCP tool with the given API implementation
func createIntegrationsTool(api IntegrationsAPI) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.Tool{
			Name:        "integrations",
			Description: "Manage HPE OpsRamp integrations and their configurations.",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "Action to perform: list, get, getDetailed, create, update, delete, enable, disable, listTypes, getType",
					},
					"id": map[string]interface{}{
						"type":        "string",
						"description": "Integration ID (for get, update, delete, enable, disable, getType)",
					},
					"config": map[string]interface{}{
						"type":        "object",
						"description": "Integration configuration (for create and update)",
					},
				},
				Required: []string{"action"},
			},
		}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return IntegrationsToolHandler(ctx, req, api)
		}
}

// IntegrationsToolHandler routes requests to the correct method
// Exported for testing purposes
func IntegrationsToolHandler(ctx context.Context, req mcp.CallToolRequest, api IntegrationsAPI) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	action, _ := args["action"].(string)
	id, _ := args["id"].(string)
	config, _ := args["config"].(map[string]interface{})

	// Log the tool execution
	logger := common.GetLogger()
	logger.LogToolExecution("integrations", action, args)

	var err error
	var result interface{}

	switch action {
	case "list":
		logger.Info("Executing List integrations")
		result, err = api.List(ctx)
	case "get":
		logger.Info("Executing Get integration with ID: %s", id)
		result, err = api.Get(ctx, id)
	case "getDetailed":
		logger.Info("Executing GetDetailed integration with ID: %s", id)
		result, err = api.GetDetailed(ctx, id)
	case "create":
		logger.Info("Executing Create integration")
		result, err = api.Create(ctx, config)
	case "update":
		logger.Info("Executing Update integration with ID: %s", id)
		result, err = api.Update(ctx, id, config)
	case "delete":
		logger.Info("Executing Delete integration with ID: %s", id)
		err = api.Delete(ctx, id)
	case "enable":
		logger.Info("Executing Enable integration with ID: %s", id)
		err = api.Enable(ctx, id)
	case "disable":
		logger.Info("Executing Disable integration with ID: %s", id)
		err = api.Disable(ctx, id)
	case "listTypes":
		logger.Info("Executing List integration types")
		result, err = api.ListTypes(ctx)
	case "getType":
		logger.Info("Executing Get integration type with ID: %s", id)
		result, err = api.GetType(ctx, id)
	default:
		logger.Error("Unknown action: %s", action)
		err = server.ErrToolNotFound
	}

	// Log the result
	logger.LogToolResult("integrations", action, result, err)

	if err != nil {
		logger.Error("Error executing %s: %v", action, err)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}},
		}, nil
	}

	// Convert result to string if it exists
	resultText := "OK"
	if result != nil {
		// For large results, use JSON marshaling for better formatting
		jsonResult, jsonErr := json.Marshal(result)
		if jsonErr == nil {
			resultText = string(jsonResult)
			// Log the result (truncated if too large)
			if len(resultText) > 1000 {
				logger.Debug("Result (truncated): %s...", resultText[:1000])
			} else {
				logger.Debug("Result: %s", resultText)
			}
		} else {
			// Fallback to simple string conversion
			resultText = fmt.Sprintf("%v", result)
			logger.Debug("Result: %s", resultText)
		}
	}

	logger.Info("Successfully executed %s", action)
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: resultText}},
	}, nil
}

// MockIntegrationsAPI is a simple mock implementation of IntegrationsAPI
type MockIntegrationsAPI struct{}

func (m *MockIntegrationsAPI) List(ctx context.Context) ([]types.Integration, error) {
	// Return mock data
	return []types.Integration{
		{
			ID:     "int-001",
			Name:   "Mock Integration 1",
			Type:   "api",
			Status: "active",
		},
		{
			ID:     "int-002",
			Name:   "Mock Integration 2",
			Type:   "webhook",
			Status: "inactive",
		},
	}, nil
}

func (m *MockIntegrationsAPI) Get(ctx context.Context, id string) (*types.Integration, error) {
	return &types.Integration{
		ID:     id,
		Name:   "Mock Integration",
		Type:   "api",
		Status: "active",
	}, nil
}

func (m *MockIntegrationsAPI) Create(ctx context.Context, config map[string]interface{}) (*types.Integration, error) {
	name, _ := config["name"].(string)
	return &types.Integration{
		ID:     "new-int-001",
		Name:   name,
		Type:   "api",
		Status: "active",
	}, nil
}

func (m *MockIntegrationsAPI) Update(ctx context.Context, id string, config map[string]interface{}) (*types.Integration, error) {
	name, _ := config["name"].(string)
	return &types.Integration{
		ID:     id,
		Name:   name,
		Type:   "api",
		Status: "active",
	}, nil
}

func (m *MockIntegrationsAPI) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MockIntegrationsAPI) Enable(ctx context.Context, id string) error {
	return nil
}

func (m *MockIntegrationsAPI) Disable(ctx context.Context, id string) error {
	return nil
}

func (m *MockIntegrationsAPI) ListTypes(ctx context.Context) ([]types.IntegrationType, error) {
	return []types.IntegrationType{
		{
			ID:          "api",
			Name:        "API Integration",
			Description: "Integration with external APIs",
			Category:    "external",
		},
		{
			ID:          "webhook",
			Name:        "Webhook Integration",
			Description: "Integration with webhooks",
			Category:    "external",
		},
	}, nil
}

func (m *MockIntegrationsAPI) GetType(ctx context.Context, id string) (*types.IntegrationType, error) {
	switch id {
	case "api":
		return &types.IntegrationType{
			ID:          "api",
			Name:        "API Integration",
			Description: "Integration with external APIs",
			Category:    "external",
		}, nil
	case "webhook":
		return &types.IntegrationType{
			ID:          "webhook",
			Name:        "Webhook Integration",
			Description: "Integration with webhooks",
			Category:    "external",
		}, nil
	default:
		return nil, fmt.Errorf("integration type with ID %s not found", id)
	}
}

func (m *MockIntegrationsAPI) GetDetailed(ctx context.Context, id string) (*types.DetailedIntegration, error) {
	// Get the basic integration first
	integration, err := m.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create an extended integration
	extIntegration := types.ExtendedIntegration{
		Integration:     *integration,
		DisplayName:     "Mock Detailed Integration",
		App:             "mock-app",
		Version:         "1.0.0",
		Category:        "Mock Category",
		State:           "Active",
		InstalledBy:     "admin",
		InstalledTime:   time.Now().Format(time.RFC3339),
		ModifiedBy:      "admin",
		ModifiedTime:    time.Now().Format(time.RFC3339),
		UpdateAvailable: false,
	}

	// Create a detailed integration with mock data
	detailed := &types.DetailedIntegration{
		ExtendedIntegration: extIntegration,
		Resources: []types.IntegrationResource{
			{
				ID:           "res-001",
				Name:         "Mock Resource 1",
				Type:         "Server",
				Status:       "Up",
				DiscoveredAt: time.Now(),
			},
			{
				ID:           "res-002",
				Name:         "Mock Resource 2",
				Type:         "Database",
				Status:       "Down",
				DiscoveredAt: time.Now().Add(-24 * time.Hour),
			},
		},
		Metrics: []types.Metric{
			{
				Name:        "CPU Usage",
				Description: "CPU usage percentage",
				Unit:        "%",
				Type:        "Gauge",
			},
			{
				Name:        "Memory Usage",
				Description: "Memory usage percentage",
				Unit:        "%",
				Type:        "Gauge",
			},
		},
		Alerts: []types.Alert{
			{
				ID:          "alert-001",
				Name:        "High CPU Usage",
				Severity:    "Critical",
				Status:      "Active",
				CreatedTime: time.Now().Add(-1 * time.Hour),
			},
		},
		LastDiscoveryRun: &types.DiscoveryRunInfo{
			StartTime: time.Now().Add(-2 * time.Hour),
			EndTime:   time.Now().Add(-1 * time.Hour),
			Status:    "Completed",
			Resources: 2,
		},
	}

	return detailed, nil
}
