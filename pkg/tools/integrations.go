package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/opsramp/or-mcp-v2/common"
	"github.com/opsramp/or-mcp-v2/pkg/types"
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
		integrationTypes, err := api.ListTypes(ctx)
		if err != nil {
			logger.Error("Error listing integration types: %v", err)
			return nil, err
		}

		// Log the results for debugging
		typesCount := len(integrationTypes)
		logger.Debug("Found %d integration types", typesCount)

		if typesCount > 0 {
			// Log a sample of the integration types
			sampleSize := min(3, typesCount)
			for i := 0; i < sampleSize; i++ {
				intType := integrationTypes[i]
				logger.Debug("Integration type %d: ID=%s, Name=%s, Category=%s", i, intType.ID, intType.Name, intType.Category)
			}
		} else {
			logger.Warn("No integration types found")
		}

		// Set the result
		result = integrationTypes
	case "getType":
		logger.Info("Executing Get integration type with ID: %s", id)
		result, err = api.GetType(ctx, id)
	default:
		logger.Error("Unknown action: %s", action)
		err = server.ErrToolNotFound
	}

	// Log the result
	logger.LogToolResult("integrations", action, result, err)

	// If there's an error, return it
	if err != nil {
		return nil, err
	}

	// Convert the result to a JSON string
	resultJSON, err := json.Marshal(result)
	if err != nil {
		logger.Error("Failed to marshal result to JSON: %v", err)
		return nil, err
	}

	// Log the JSON result
	logger.Debug("Result JSON: %s", string(resultJSON))

	// Create the MCP tool result
	toolResult := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}

	return toolResult, nil
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
