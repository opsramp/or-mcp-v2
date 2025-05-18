package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/types"
)

// OpsRampIntegrationsAPI implements real API calls
type OpsRampIntegrationsAPI struct {
	client *client.OpsRampClient
	logger *common.CustomLogger
}

// NewOpsRampIntegrationsAPI creates a new OpsRamp integrations API client
func NewOpsRampIntegrationsAPI(client *client.OpsRampClient) *OpsRampIntegrationsAPI {
	// Get the logger
	logger := common.GetLogger()

	return &OpsRampIntegrationsAPI{
		client: client,
		logger: logger,
	}
}

// List returns all integrations
func (api *OpsRampIntegrationsAPI) List(ctx context.Context) ([]types.Integration, error) {
	api.logger.Info("Listing integrations")

	// The OpsRamp API returns a paginated response for integrations
	var response struct {
		Results []types.ExtendedIntegration `json:"results"`
	}

	// Use the correct OpsRamp API endpoint for listing installed integrations
	endpoint := "/api/v2/tenants/" + api.client.GetTenantID() + "/integrations/installed/search"
	api.logger.Debug("Using endpoint: %s", endpoint)

	err := api.client.Get(ctx, endpoint, &response)
	if err != nil {
		api.logger.Error("Failed to list integrations: %v", err)
		return nil, fmt.Errorf("failed to list integrations: %w", err)
	}

	// Convert ExtendedIntegration to Integration
	integrations := make([]types.Integration, len(response.Results))
	for i, extInt := range response.Results {
		integrations[i] = extInt.Integration
	}

	api.logger.Info("Successfully listed %d integrations", len(integrations))
	return integrations, nil
}

// Get returns a specific integration by ID
func (api *OpsRampIntegrationsAPI) Get(ctx context.Context, id string) (*types.Integration, error) {
	api.logger.Info("Getting integration with ID: %s", id)

	var extResult types.ExtendedIntegration

	// Use the correct OpsRamp API endpoint for getting an installed integration
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/integrations/installed/%s", api.client.GetTenantID(), id)
	api.logger.Debug("Using endpoint: %s", endpoint)

	err := api.client.Get(ctx, endpoint, &extResult)
	if err != nil {
		api.logger.Error("Failed to get integration %s: %v", id, err)
		return nil, fmt.Errorf("failed to get integration %s: %w", id, err)
	}

	api.logger.Info("Successfully retrieved integration: %s", extResult.DisplayName)
	return &extResult.Integration, nil
}

// GetDetailed returns detailed information about a specific integration by ID
func (api *OpsRampIntegrationsAPI) GetDetailed(ctx context.Context, id string) (*types.DetailedIntegration, error) {
	api.logger.Info("Getting detailed integration with ID: %s", id)

	// First get the extended integration info
	var extIntegration types.ExtendedIntegration

	// Use the correct OpsRamp API endpoint for getting an installed integration
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/integrations/installed/%s", api.client.GetTenantID(), id)
	api.logger.Debug("Using endpoint: %s", endpoint)

	err := api.client.Get(ctx, endpoint, &extIntegration)
	if err != nil {
		api.logger.Error("Failed to get integration %s: %v", id, err)
		return nil, fmt.Errorf("failed to get integration %s: %w", id, err)
	}

	// Create a detailed integration with the extended info
	detailed := &types.DetailedIntegration{
		ExtendedIntegration: extIntegration,
	}

	// Get resources for this integration
	api.logger.Debug("Getting resources for integration: %s", id)
	resourcesEndpoint := fmt.Sprintf("/api/v2/tenants/%s/integrations/installed/%s/resources",
		api.client.GetTenantID(), id)

	var resourcesResponse struct {
		Results []types.Resource `json:"results"`
	}

	err = api.client.Get(ctx, resourcesEndpoint, &resourcesResponse)
	if err != nil {
		api.logger.Warn("Failed to get resources for integration %s: %v", id, err)
		// Continue with other data even if resources fail
	} else {
		// Convert Resource to IntegrationResource
		integrationResources := make([]types.IntegrationResource, 0, len(resourcesResponse.Results))
		for _, res := range resourcesResponse.Results {
			integrationResources = append(integrationResources, types.IntegrationResource{
				ID:     res.ID,
				Name:   res.Name,
				Type:   res.Type,
				Status: res.Status,
				// Use current time as discovery time since it's not available in the API response
				DiscoveredAt: time.Now(),
			})
		}
		detailed.Resources = integrationResources
		api.logger.Info("Retrieved %d resources for integration %s", len(detailed.Resources), id)
	}

	// Get metrics for this integration
	api.logger.Debug("Getting metrics for integration: %s", id)
	metricsEndpoint := fmt.Sprintf("/api/v2/tenants/%s/integrations/installed/%s/metrics",
		api.client.GetTenantID(), id)

	var metricsResponse struct {
		Results []types.Metric `json:"results"`
	}

	err = api.client.Get(ctx, metricsEndpoint, &metricsResponse)
	if err != nil {
		api.logger.Warn("Failed to get metrics for integration %s: %v", id, err)
		// Continue with other data even if metrics fail
	} else {
		detailed.Metrics = metricsResponse.Results
		api.logger.Info("Retrieved %d metrics for integration %s", len(detailed.Metrics), id)
	}

	// Get alerts for this integration
	api.logger.Debug("Getting alerts for integration: %s", id)
	alertsEndpoint := fmt.Sprintf("/api/v2/tenants/%s/integrations/installed/%s/alerts",
		api.client.GetTenantID(), id)

	var alertsResponse struct {
		Results []types.Alert `json:"results"`
	}

	err = api.client.Get(ctx, alertsEndpoint, &alertsResponse)
	if err != nil {
		api.logger.Warn("Failed to get alerts for integration %s: %v", id, err)
		// Continue with other data even if alerts fail
	} else {
		detailed.Alerts = alertsResponse.Results
		api.logger.Info("Retrieved %d alerts for integration %s", len(detailed.Alerts), id)
	}

	api.logger.Info("Successfully retrieved detailed information for integration: %s", extIntegration.DisplayName)
	return detailed, nil
}

// Create creates a new integration
func (api *OpsRampIntegrationsAPI) Create(ctx context.Context, config map[string]interface{}) (*types.Integration, error) {
	var result types.Integration
	err := api.client.Post(ctx, "/api/integrations", config, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create integration: %w", err)
	}
	return &result, nil
}

// Update updates an existing integration
func (api *OpsRampIntegrationsAPI) Update(ctx context.Context, id string, config map[string]interface{}) (*types.Integration, error) {
	var result types.Integration
	err := api.client.Put(ctx, fmt.Sprintf("/api/integrations/%s", id), config, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to update integration %s: %w", id, err)
	}
	return &result, nil
}

// Delete deletes an integration
func (api *OpsRampIntegrationsAPI) Delete(ctx context.Context, id string) error {
	err := api.client.Delete(ctx, fmt.Sprintf("/api/integrations/%s", id))
	if err != nil {
		return fmt.Errorf("failed to delete integration %s: %w", id, err)
	}
	return nil
}

// Enable enables an integration
func (api *OpsRampIntegrationsAPI) Enable(ctx context.Context, id string) error {
	err := api.client.Post(ctx, fmt.Sprintf("/api/integrations/%s/enable", id), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to enable integration %s: %w", id, err)
	}
	return nil
}

// Disable disables an integration
func (api *OpsRampIntegrationsAPI) Disable(ctx context.Context, id string) error {
	err := api.client.Post(ctx, fmt.Sprintf("/api/integrations/%s/disable", id), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to disable integration %s: %w", id, err)
	}
	return nil
}

// ListTypes returns all integration types
func (api *OpsRampIntegrationsAPI) ListTypes(ctx context.Context) ([]types.IntegrationType, error) {
	var result []types.IntegrationType
	err := api.client.Get(ctx, "/api/integration-types", &result)
	if err != nil {
		return nil, fmt.Errorf("failed to list integration types: %w", err)
	}
	return result, nil
}

// GetType returns a specific integration type by ID
func (api *OpsRampIntegrationsAPI) GetType(ctx context.Context, id string) (*types.IntegrationType, error) {
	var result types.IntegrationType
	err := api.client.Get(ctx, fmt.Sprintf("/api/integration-types/%s", id), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration type %s: %w", id, err)
	}
	return &result, nil
}
