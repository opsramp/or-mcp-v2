package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/types"
)

// ResourcesAPI defines the contract for resource operations
type ResourcesAPI interface {
	// Search searches for resources based on the provided parameters
	Search(ctx context.Context, params types.ResourceSearchParams) (*types.ResourceSearchResponse, error)

	// Get retrieves a specific resource by ID
	Get(ctx context.Context, id string) (*types.Resource, error)

	// GetDetailed retrieves detailed information about a specific resource by ID
	GetDetailed(ctx context.Context, id string) (*types.DetailedResource, error)

	// Create creates a new resource
	Create(ctx context.Context, resource types.ResourceCreateRequest) (*types.Resource, error)

	// Update updates an existing resource
	Update(ctx context.Context, id string, resource types.ResourceUpdateRequest) (*types.Resource, error)

	// Delete deletes a resource by ID
	Delete(ctx context.Context, id string) error

	// BulkUpdate updates multiple resources at once
	BulkUpdate(ctx context.Context, request types.ResourceBulkUpdateRequest) error

	// BulkDelete deletes multiple resources at once
	BulkDelete(ctx context.Context, request types.ResourceBulkDeleteRequest) error

	// GetResourceTypes retrieves all available resource types
	GetResourceTypes(ctx context.Context) ([]types.ResourceTypeInfo, error)

	// ChangeState changes the state of a resource
	ChangeState(ctx context.Context, id string, request types.ResourceStateChangeRequest) error

	// GetMetrics retrieves metrics for a resource
	GetMetrics(ctx context.Context, id string, request types.ResourceMetricsRequest) (*types.ResourceMetricsResponse, error)

	// GetTags retrieves all tags for a resource
	GetTags(ctx context.Context, id string) ([]types.Tag, error)

	// UpdateTags updates the tags for a resource
	UpdateTags(ctx context.Context, id string, tags []types.Tag) error
}

// ResourcesAdapter handles mapping between MCP and OpsRamp API for resources
type ResourcesAdapter struct {
	api    ResourcesAPI
	logger *common.CustomLogger
	client *client.OpsRampClient
}

// NewResourcesAdapter creates a new ResourcesAdapter
func NewResourcesAdapter(client *client.OpsRampClient) *ResourcesAdapter {
	logger := common.GetLogger()

	// We'll create the API implementation in the resources.go file
	// For now, we'll use a placeholder that will be set by the caller
	return &ResourcesAdapter{
		logger: logger,
		client: client,
	}
}

// SetAPI sets the ResourcesAPI implementation
func (a *ResourcesAdapter) SetAPI(api ResourcesAPI) {
	a.api = api
}

// List returns a list of resources
func (a *ResourcesAdapter) List(ctx context.Context) ([]types.Resource, error) {
	a.logger.Info("ResourcesAdapter: Listing resources")

	// Build the endpoint for resources search
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/search?pageSize=100&pageNo=1", a.client.GetTenantID())
	a.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var response types.ResourceSearchResponse
	err := a.client.Get(ctx, endpoint, &response)
	if err != nil {
		a.logger.Error("ResourcesAdapter: Failed to list resources: %v", err)
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	a.logger.Info("ResourcesAdapter: Successfully listed %d resources", len(response.Results))
	return response.Results, nil
}

// Get returns a specific resource by ID
func (a *ResourcesAdapter) Get(ctx context.Context, id string) (*types.Resource, error) {
	a.logger.Info("ResourcesAdapter: Getting resource with ID: %s", id)

	// First, search for the resource by ID
	searchParams := types.ResourceSearchParams{
		PageSize: 1,
		PageNo:   1,
		// Add a filter for the resource ID if the API supports it
		// For now, we'll get all resources and filter client-side
	}

	// Build the endpoint
	queryParams := url.Values{}
	queryParams.Add("pageSize", fmt.Sprintf("%d", searchParams.PageSize))
	queryParams.Add("pageNo", fmt.Sprintf("%d", searchParams.PageNo))

	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/search?%s", a.client.GetTenantID(), queryParams.Encode())
	a.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var response types.ResourceSearchResponse
	err := a.client.Get(ctx, endpoint, &response)
	if err != nil {
		a.logger.Error("ResourcesAdapter: Failed to search for resource %s: %v", id, err)
		return nil, fmt.Errorf("failed to search for resource %s: %w", id, err)
	}

	// Find the resource with the matching ID
	for _, resource := range response.Results {
		if resource.ID == id {
			a.logger.Info("ResourcesAdapter: Successfully retrieved resource: %s", resource.Name)
			return &resource, nil
		}
	}

	// If we get here, the resource was not found
	a.logger.Error("ResourcesAdapter: Resource %s not found", id)
	return nil, fmt.Errorf("resource %s not found", id)
}

// GetDetailed returns detailed information about a specific resource by ID
func (a *ResourcesAdapter) GetDetailed(ctx context.Context, id string) (*types.DetailedResource, error) {
	a.logger.Info("ResourcesAdapter: Getting detailed resource with ID: %s", id)

	// First, get the basic resource
	resource, err := a.Get(ctx, id)
	if err != nil {
		a.logger.Error("ResourcesAdapter: Failed to get resource %s: %v", id, err)
		return nil, fmt.Errorf("failed to get resource %s: %w", id, err)
	}

	// Convert the basic resource to a detailed resource
	detailedResource := &types.DetailedResource{
		Resource: *resource,
		// Add any additional fields that would be in the detailed view
		// For now, we'll just return the basic resource information
	}

	a.logger.Info("ResourcesAdapter: Successfully retrieved detailed resource: %s", detailedResource.Name)
	return detailedResource, nil
}

// Create creates a new resource
func (a *ResourcesAdapter) Create(ctx context.Context, config map[string]interface{}) (*types.Resource, error) {
	a.logger.Info("ResourcesAdapter: Creating new resource")

	// Convert the config map to a ResourceCreateRequest
	var createRequest types.ResourceCreateRequest
	configJSON, err := json.Marshal(config)
	if err != nil {
		a.logger.Error("ResourcesAdapter: Failed to marshal resource config: %v", err)
		return nil, fmt.Errorf("failed to marshal resource config: %w", err)
	}

	if err := json.Unmarshal(configJSON, &createRequest); err != nil {
		a.logger.Error("ResourcesAdapter: Failed to unmarshal resource config: %v", err)
		return nil, fmt.Errorf("failed to unmarshal resource config: %w", err)
	}

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources", a.client.GetTenantID())
	a.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var resource types.Resource
	err = a.client.Post(ctx, endpoint, createRequest, &resource)
	if err != nil {
		a.logger.Error("ResourcesAdapter: Failed to create resource: %v", err)
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	a.logger.Info("ResourcesAdapter: Successfully created resource with ID: %s", resource.ID)
	return &resource, nil
}

// Update updates an existing resource
func (a *ResourcesAdapter) Update(ctx context.Context, id string, config map[string]interface{}) (*types.Resource, error) {
	a.logger.Info("ResourcesAdapter: Updating resource with ID: %s", id)

	// Convert the config map to a ResourceUpdateRequest
	var updateRequest types.ResourceUpdateRequest
	configJSON, err := json.Marshal(config)
	if err != nil {
		a.logger.Error("ResourcesAdapter: Failed to marshal resource config: %v", err)
		return nil, fmt.Errorf("failed to marshal resource config: %w", err)
	}

	if err := json.Unmarshal(configJSON, &updateRequest); err != nil {
		a.logger.Error("ResourcesAdapter: Failed to unmarshal resource config: %v", err)
		return nil, fmt.Errorf("failed to unmarshal resource config: %w", err)
	}

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/%s", a.client.GetTenantID(), id)
	a.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var resource types.Resource
	err = a.client.Put(ctx, endpoint, updateRequest, &resource)
	if err != nil {
		a.logger.Error("ResourcesAdapter: Failed to update resource %s: %v", id, err)
		return nil, fmt.Errorf("failed to update resource %s: %w", id, err)
	}

	a.logger.Info("ResourcesAdapter: Successfully updated resource: %s", resource.Name)
	return &resource, nil
}

// Delete deletes a resource by ID
func (a *ResourcesAdapter) Delete(ctx context.Context, id string) error {
	a.logger.Info("ResourcesAdapter: Deleting resource with ID: %s", id)

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/%s", a.client.GetTenantID(), id)
	a.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	err := a.client.Delete(ctx, endpoint)
	if err != nil {
		a.logger.Error("ResourcesAdapter: Failed to delete resource %s: %v", id, err)
		return fmt.Errorf("failed to delete resource %s: %w", id, err)
	}

	a.logger.Info("ResourcesAdapter: Successfully deleted resource with ID: %s", id)
	return nil
}

// Search searches for resources based on the provided parameters
func (a *ResourcesAdapter) Search(ctx context.Context, params map[string]interface{}) (*types.ResourceSearchResponse, error) {
	a.logger.Info("ResourcesAdapter: Searching for resources")

	// Convert the params map to a ResourceSearchParams
	var searchParams types.ResourceSearchParams
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		a.logger.Error("ResourcesAdapter: Failed to marshal search params: %v", err)
		return nil, fmt.Errorf("failed to marshal search params: %w", err)
	}

	if err := json.Unmarshal(paramsJSON, &searchParams); err != nil {
		a.logger.Error("ResourcesAdapter: Failed to unmarshal search params: %v", err)
		return nil, fmt.Errorf("failed to unmarshal search params: %w", err)
	}

	// Build the endpoint with query parameters
	queryParams := url.Values{}
	if searchParams.PageSize > 0 {
		queryParams.Add("pageSize", fmt.Sprintf("%d", searchParams.PageSize))
	}
	if searchParams.PageNo > 0 {
		queryParams.Add("pageNo", fmt.Sprintf("%d", searchParams.PageNo))
	}
	if searchParams.IsDescendingOrder {
		queryParams.Add("isDescendingOrder", "true")
	}
	if searchParams.SortName != "" {
		queryParams.Add("sortName", searchParams.SortName)
	}
	if searchParams.ResourceType != "" {
		queryParams.Add("resourceType", searchParams.ResourceType)
	}

	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/search?%s", a.client.GetTenantID(), queryParams.Encode())
	a.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var response types.ResourceSearchResponse
	err = a.client.Get(ctx, endpoint, &response)
	if err != nil {
		a.logger.Error("ResourcesAdapter: Failed to search resources: %v", err)
		return nil, fmt.Errorf("failed to search resources: %w", err)
	}

	a.logger.Info("ResourcesAdapter: Successfully searched resources, found %d results", len(response.Results))
	return &response, nil
}

// These methods are not applicable to resources but are kept for interface compatibility
func (a *ResourcesAdapter) Enable(ctx context.Context, id string) error {
	a.logger.Warn("ResourcesAdapter: Enable operation not supported for resources")
	return fmt.Errorf("enable operation not supported for resources")
}

func (a *ResourcesAdapter) Disable(ctx context.Context, id string) error {
	a.logger.Warn("ResourcesAdapter: Disable operation not supported for resources")
	return fmt.Errorf("disable operation not supported for resources")
}

func (a *ResourcesAdapter) ListTypes(ctx context.Context) ([]types.IntegrationType, error) {
	a.logger.Warn("ResourcesAdapter: ListTypes operation not implemented yet")
	return []types.IntegrationType{}, fmt.Errorf("list types operation not implemented yet")
}

func (a *ResourcesAdapter) GetType(ctx context.Context, id string) (*types.IntegrationType, error) {
	a.logger.Warn("ResourcesAdapter: GetType operation not implemented yet")
	return nil, fmt.Errorf("get type operation not implemented yet")
}
