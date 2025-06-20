package tools

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/opsramp/or-mcp-v2/common"
	"github.com/opsramp/or-mcp-v2/pkg/client"
	"github.com/opsramp/or-mcp-v2/pkg/types"
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

	// GetMinimal retrieves minimal resource information for performance
	GetMinimal(ctx context.Context, id string) (*types.ResourceMinimal, error)
}

// OpsRampResourcesAPI implements the ResourcesAPI interface for OpsRamp
type OpsRampResourcesAPI struct {
	client *client.OpsRampClient
	logger *common.CustomLogger
	config *ResourcesAPIConfig
}

// ResourcesAPIConfig holds configuration for the Resources API client
type ResourcesAPIConfig struct {
	RetryAttempts  int           `json:"retry_attempts"`
	RetryDelay     time.Duration `json:"retry_delay"`
	RequestTimeout time.Duration `json:"request_timeout"`
	RateLimitDelay time.Duration `json:"rate_limit_delay"`
	CircuitBreaker bool          `json:"circuit_breaker"`
	MaxFailures    int           `json:"max_failures"`
	ResetTimeout   time.Duration `json:"reset_timeout"`
}

// NewOpsRampResourcesAPI creates a new OpsRamp resources API client
func NewOpsRampResourcesAPI(client *client.OpsRampClient) *OpsRampResourcesAPI {
	// Get the logger
	logger := common.GetLogger()

	// Default configuration
	config := &ResourcesAPIConfig{
		RetryAttempts:  3,
		RetryDelay:     1 * time.Second,
		RequestTimeout: 30 * time.Second,
		RateLimitDelay: 5 * time.Second,
		CircuitBreaker: true,
		MaxFailures:    5,
		ResetTimeout:   60 * time.Second,
	}

	return &OpsRampResourcesAPI{
		client: client,
		logger: logger,
		config: config,
	}
}

// NewOpsRampResourcesAPIWithConfig creates a new OpsRamp resources API client with custom configuration
func NewOpsRampResourcesAPIWithConfig(client *client.OpsRampClient, config *ResourcesAPIConfig) *OpsRampResourcesAPI {
	// Get the logger
	logger := common.GetLogger()

	return &OpsRampResourcesAPI{
		client: client,
		logger: logger,
		config: config,
	}
}

// Search searches for resources based on the provided parameters
func (api *OpsRampResourcesAPI) Search(ctx context.Context, params types.ResourceSearchParams) (*types.ResourceSearchResponse, error) {
	api.logger.Info("Searching for resources with parameters")

	// Build query parameters
	queryParams := url.Values{}

	// Add pagination parameters
	if params.PageNo > 0 {
		queryParams.Add("pageNo", strconv.Itoa(params.PageNo))
	}
	if params.PageSize > 0 {
		queryParams.Add("pageSize", strconv.Itoa(params.PageSize))
	}

	// Add sorting parameters
	if params.SortName != "" {
		queryParams.Add("sortName", params.SortName)
	}
	queryParams.Add("isDescendingOrder", strconv.FormatBool(params.IsDescendingOrder))

	// Add filter parameters
	if params.QueryString != "" {
		queryParams.Add("queryString", params.QueryString)
	}
	if params.HostName != "" {
		queryParams.Add("hostName", params.HostName)
	}
	if params.DNSName != "" {
		queryParams.Add("dnsName", params.DNSName)
	}
	if params.ResourceName != "" {
		queryParams.Add("resourceName", params.ResourceName)
	}
	if params.AliasName != "" {
		queryParams.Add("aliasName", params.AliasName)
	}
	if params.ID != "" {
		queryParams.Add("id", params.ID)
	}
	if params.SerialNumber != "" {
		queryParams.Add("serialNumber", params.SerialNumber)
	}
	if params.IPAddress != "" {
		queryParams.Add("ipAddress", params.IPAddress)
	}
	if params.SystemUID != "" {
		queryParams.Add("systemUID", params.SystemUID)
	}
	if params.State != "" {
		queryParams.Add("state", params.State)
	}
	if params.Type != "" {
		queryParams.Add("type", params.Type)
	}
	if params.DeviceType != "" {
		queryParams.Add("deviceType", params.DeviceType)
	}
	if params.ResourceType != "" {
		queryParams.Add("resourceType", params.ResourceType)
	}
	if params.StartCreationDate != "" {
		queryParams.Add("startCreationDate", params.StartCreationDate)
	}
	if params.EndCreationDate != "" {
		queryParams.Add("endCreationDate", params.EndCreationDate)
	}
	if params.StartUpdationDate != "" {
		queryParams.Add("startUpdationDate", params.StartUpdationDate)
	}
	if params.EndUpdationDate != "" {
		queryParams.Add("endUpdationDate", params.EndUpdationDate)
	}
	if params.Tags != "" {
		queryParams.Add("tags", params.Tags)
	}
	if params.Template != "" {
		queryParams.Add("template", params.Template)
	}
	if params.AgentProfile != "" {
		queryParams.Add("agentProfile", params.AgentProfile)
	}
	if params.GatewayProfile != "" {
		queryParams.Add("gatewayProfile", params.GatewayProfile)
	}
	if params.InstanceID != "" {
		queryParams.Add("instanceId", params.InstanceID)
	}
	if params.AccountNumber != "" {
		queryParams.Add("accountNumber", params.AccountNumber)
	}
	if params.AgentInstalled != nil {
		queryParams.Add("agentInstalled", strconv.FormatBool(*params.AgentInstalled))
	}
	if params.DeviceGroup != "" {
		queryParams.Add("deviceGroup", params.DeviceGroup)
	}
	if params.ServiceGroup != "" {
		queryParams.Add("serviceGroup", params.ServiceGroup)
	}
	if params.DeviceLocation != "" {
		queryParams.Add("deviceLocation", params.DeviceLocation)
	}
	if params.IsEquals != "" {
		queryParams.Add("isEquals", params.IsEquals)
	}
	// Add new filter parameters
	if params.AliasIP != "" {
		queryParams.Add("aliasIp", params.AliasIP)
	}
	if params.AppRoles != "" {
		queryParams.Add("appRoles", params.AppRoles)
	}
	if params.OSArchitecture != "" {
		queryParams.Add("osArchitecture", params.OSArchitecture)
	}
	if params.AssetManagedTime != "" {
		queryParams.Add("assetManagedTime", params.AssetManagedTime)
	}
	if params.FirstAssetManagedTime != "" {
		queryParams.Add("firstAssetManagedTime", params.FirstAssetManagedTime)
	}
	if params.Category != "" {
		queryParams.Add("category", params.Category)
	}
	if params.Make != "" {
		queryParams.Add("make", params.Make)
	}
	if params.Model != "" {
		queryParams.Add("model", params.Model)
	}
	if params.ProviderType != "" {
		queryParams.Add("providerType", params.ProviderType)
	}
	if params.ProviderUID != "" {
		queryParams.Add("providerUID", params.ProviderUID)
	}

	// Build the endpoint with query parameters
	// Build the endpoint without query parameters
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/search", api.client.GetTenantID())

	// Add query parameters separately to avoid URL encoding issues
	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())
	}

	api.logger.Debug("Using endpoint: %s", endpoint) // Make the request
	var response types.ResourceSearchResponse
	err := api.client.Get(ctx, endpoint, &response)
	if err != nil {
		api.logger.Error("Failed to search resources: %v", err)
		return nil, fmt.Errorf("failed to search resources: %w", err)
	}

	api.logger.Info("Successfully searched resources, found %d results", len(response.Results))
	return &response, nil
}

// Get retrieves a specific resource by ID
func (api *OpsRampResourcesAPI) Get(ctx context.Context, id string) (*types.Resource, error) {
	api.logger.Info("Getting resource with ID: %s", id)

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/%s", api.client.GetTenantID(), id)
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var resource types.Resource
	err := api.client.Get(ctx, endpoint, &resource)
	if err != nil {
		api.logger.Error("Failed to get resource %s: %v", id, err)
		return nil, fmt.Errorf("failed to get resource %s: %w", id, err)
	}

	api.logger.Info("Successfully retrieved resource: %s", resource.Name)
	return &resource, nil
}

// GetDetailed retrieves detailed information about a specific resource by ID
func (api *OpsRampResourcesAPI) GetDetailed(ctx context.Context, id string) (*types.DetailedResource, error) {
	api.logger.Info("Getting detailed resource with ID: %s", id)

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/%s", api.client.GetTenantID(), id)
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var detailedResource types.DetailedResource
	err := api.client.Get(ctx, endpoint, &detailedResource)
	if err != nil {
		api.logger.Error("Failed to get detailed resource %s: %v", id, err)
		return nil, fmt.Errorf("failed to get detailed resource %s: %w", id, err)
	}

	api.logger.Info("Successfully retrieved detailed resource: %s", detailedResource.Name)
	return &detailedResource, nil
}

// Create creates a new resource
func (api *OpsRampResourcesAPI) Create(ctx context.Context, resource types.ResourceCreateRequest) (*types.Resource, error) {
	api.logger.Info("Creating new resource of type: %s", resource.ResourceType)

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources", api.client.GetTenantID())
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var createdResource types.Resource
	err := api.client.Post(ctx, endpoint, resource, &createdResource)
	if err != nil {
		api.logger.Error("Failed to create resource: %v", err)
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	api.logger.Info("Successfully created resource with ID: %s", createdResource.ID)
	return &createdResource, nil
}

// Update updates an existing resource
func (api *OpsRampResourcesAPI) Update(ctx context.Context, id string, resource types.ResourceUpdateRequest) (*types.Resource, error) {
	api.logger.Info("Updating resource with ID: %s", id)

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/%s", api.client.GetTenantID(), id)
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var updatedResource types.Resource
	err := api.client.Post(ctx, endpoint, resource, &updatedResource)
	if err != nil {
		api.logger.Error("Failed to update resource %s: %v", id, err)
		return nil, fmt.Errorf("failed to update resource %s: %w", id, err)
	}

	api.logger.Info("Successfully updated resource: %s", updatedResource.Name)
	return &updatedResource, nil
}

// Delete deletes a resource by ID
func (api *OpsRampResourcesAPI) Delete(ctx context.Context, id string) error {
	api.logger.Info("Deleting resource with ID: %s", id)

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/%s", api.client.GetTenantID(), id)
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	err := api.client.Delete(ctx, endpoint)
	if err != nil {
		api.logger.Error("Failed to delete resource %s: %v", id, err)
		return fmt.Errorf("failed to delete resource %s: %w", id, err)
	}

	api.logger.Info("Successfully deleted resource with ID: %s", id)
	return nil
}

// BulkUpdate updates multiple resources at once
func (api *OpsRampResourcesAPI) BulkUpdate(ctx context.Context, request types.ResourceBulkUpdateRequest) error {
	api.logger.Info("Bulk updating %d resources", len(request.ResourceIDs))

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/bulk-update", api.client.GetTenantID())
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	err := api.client.Post(ctx, endpoint, request, nil)
	if err != nil {
		api.logger.Error("Failed to bulk update resources: %v", err)
		return fmt.Errorf("failed to bulk update resources: %w", err)
	}

	api.logger.Info("Successfully bulk updated %d resources", len(request.ResourceIDs))
	return nil
}

// BulkDelete deletes multiple resources at once
func (api *OpsRampResourcesAPI) BulkDelete(ctx context.Context, request types.ResourceBulkDeleteRequest) error {
	api.logger.Info("Bulk deleting %d resources", len(request.ResourceIDs))

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/bulk-delete", api.client.GetTenantID())
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	err := api.client.Post(ctx, endpoint, request, nil)
	if err != nil {
		api.logger.Error("Failed to bulk delete resources: %v", err)
		return fmt.Errorf("failed to bulk delete resources: %w", err)
	}

	api.logger.Info("Successfully bulk deleted %d resources", len(request.ResourceIDs))
	return nil
}

// GetResourceTypes retrieves all available resource types
func (api *OpsRampResourcesAPI) GetResourceTypes(ctx context.Context) ([]types.ResourceTypeInfo, error) {
	api.logger.Info("Getting resource types")

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/types", api.client.GetTenantID())
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var response struct {
		ResourceTypes []types.ResourceTypeInfo `json:"resourceTypes"`
	}
	err := api.client.Get(ctx, endpoint, &response)
	if err != nil {
		api.logger.Error("Failed to get resource types: %v", err)
		return nil, fmt.Errorf("failed to get resource types: %w", err)
	}

	api.logger.Info("Successfully retrieved %d resource types", len(response.ResourceTypes))
	return response.ResourceTypes, nil
}

// ChangeState changes the state of a resource
func (api *OpsRampResourcesAPI) ChangeState(ctx context.Context, id string, request types.ResourceStateChangeRequest) error {
	api.logger.Info("Changing state of resource %s to %s", id, request.State)

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/%s/state", api.client.GetTenantID(), id)
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	err := api.client.Post(ctx, endpoint, request, nil)
	if err != nil {
		api.logger.Error("Failed to change state of resource %s: %v", id, err)
		return fmt.Errorf("failed to change state of resource %s: %w", id, err)
	}

	api.logger.Info("Successfully changed state of resource %s to %s", id, request.State)
	return nil
}

// GetMetrics retrieves metrics for a resource
func (api *OpsRampResourcesAPI) GetMetrics(ctx context.Context, id string, request types.ResourceMetricsRequest) (*types.ResourceMetricsResponse, error) {
	api.logger.Info("Getting metrics for resource %s", id)

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/%s/metrics", api.client.GetTenantID(), id)
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var response types.ResourceMetricsResponse
	err := api.client.Post(ctx, endpoint, request, &response)
	if err != nil {
		api.logger.Error("Failed to get metrics for resource %s: %v", id, err)
		return nil, fmt.Errorf("failed to get metrics for resource %s: %w", id, err)
	}

	api.logger.Info("Successfully retrieved metrics for resource %s", id)
	return &response, nil
}

// GetTags retrieves all tags for a resource
func (api *OpsRampResourcesAPI) GetTags(ctx context.Context, id string) ([]types.Tag, error) {
	api.logger.Info("Getting tags for resource %s", id)

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/%s/tags", api.client.GetTenantID(), id)
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	var response struct {
		Tags []types.Tag `json:"tags"`
	}
	err := api.client.Get(ctx, endpoint, &response)
	if err != nil {
		api.logger.Error("Failed to get tags for resource %s: %v", id, err)
		return nil, fmt.Errorf("failed to get tags for resource %s: %w", id, err)
	}

	api.logger.Info("Successfully retrieved %d tags for resource %s", len(response.Tags), id)
	return response.Tags, nil
}

// UpdateTags updates the tags for a resource
func (api *OpsRampResourcesAPI) UpdateTags(ctx context.Context, id string, tags []types.Tag) error {
	api.logger.Info("Updating tags for resource %s", id)

	// Build the endpoint
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/resources/%s/tags", api.client.GetTenantID(), id)
	api.logger.Debug("Using endpoint: %s", endpoint)

	// Make the request
	request := struct {
		Tags []types.Tag `json:"tags"`
	}{
		Tags: tags,
	}
	err := api.client.Post(ctx, endpoint, request, nil)
	if err != nil {
		api.logger.Error("Failed to update tags for resource %s: %v", id, err)
		return fmt.Errorf("failed to update tags for resource %s: %w", id, err)
	}

	api.logger.Info("Successfully updated tags for resource %s", id)
	return nil
}

// ============================================================================
// RESILIENCE AND ERROR HANDLING METHODS (T3.3.1-T3.3.4)
// ============================================================================

// retryWithBackoff executes a function with retry logic and exponential backoff
func (api *OpsRampResourcesAPI) retryWithBackoff(ctx context.Context, operation string, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= api.config.RetryAttempts; attempt++ {
		if attempt > 0 {
			delay := time.Duration(attempt) * api.config.RetryDelay
			api.logger.Debug("Retrying %s after %v (attempt %d/%d)", operation, delay, attempt, api.config.RetryAttempts)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		lastErr = fn()
		if lastErr == nil {
			if attempt > 0 {
				api.logger.Info("Operation %s succeeded after %d retries", operation, attempt)
			}
			return nil
		}

		// Check if error is retryable
		if !isRetryableError(lastErr) {
			api.logger.Debug("Non-retryable error for %s: %v", operation, lastErr)
			break
		}

		// Check for rate limiting
		if isRateLimitError(lastErr) {
			api.logger.Warn("Rate limit hit for %s, waiting %v", operation, api.config.RateLimitDelay)
			select {
			case <-time.After(api.config.RateLimitDelay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return fmt.Errorf("operation %s failed after %d attempts: %w", operation, api.config.RetryAttempts, lastErr)
}

// isRetryableError determines if an error is retryable
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common retryable error patterns
	errStr := err.Error()
	retryablePatterns := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"temporary failure",
		"server unavailable",
		"502", "503", "504", // HTTP status codes
	}

	for _, pattern := range retryablePatterns {
		if containsSubstring(errStr, pattern) {
			return true
		}
	}

	return false
}

// isRateLimitError determines if an error is due to rate limiting
func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	rateLimitPatterns := []string{
		"rate limit",
		"too many requests",
		"429", // HTTP status code
	}

	for _, pattern := range rateLimitPatterns {
		if containsSubstring(errStr, pattern) {
			return true
		}
	}

	return false
}

// containsSubstring checks if a string contains a substring (case-insensitive)
func containsSubstring(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr ||
			str[:len(substr)] == substr ||
			str[len(str)-len(substr):] == substr ||
			hasSubstring(str, substr))
}

// hasSubstring checks if a string contains a substring
func hasSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// classifyError classifies errors into ResourceErrorType
func (api *OpsRampResourcesAPI) classifyError(err error) *types.ResourceError {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Check for specific error types
	if containsSubstring(errStr, "not found") || containsSubstring(errStr, "404") {
		return types.NewResourceError(types.ResourceErrorTypeNotFound, "RESOURCE_NOT_FOUND", err.Error())
	}

	if containsSubstring(errStr, "unauthorized") || containsSubstring(errStr, "401") {
		return types.NewResourceError(types.ResourceErrorTypePermission, "UNAUTHORIZED", err.Error())
	}

	if containsSubstring(errStr, "forbidden") || containsSubstring(errStr, "403") {
		return types.NewResourceError(types.ResourceErrorTypePermission, "FORBIDDEN", err.Error())
	}

	if isRateLimitError(err) {
		return types.NewResourceError(types.ResourceErrorTypeRateLimit, "RATE_LIMIT_EXCEEDED", err.Error())
	}

	if containsSubstring(errStr, "timeout") {
		return types.NewResourceError(types.ResourceErrorTypeTimeout, "REQUEST_TIMEOUT", err.Error())
	}

	if containsSubstring(errStr, "conflict") || containsSubstring(errStr, "409") {
		return types.NewResourceError(types.ResourceErrorTypeConflict, "RESOURCE_CONFLICT", err.Error())
	}

	if containsSubstring(errStr, "validation") || containsSubstring(errStr, "invalid") || containsSubstring(errStr, "400") {
		return types.NewResourceError(types.ResourceErrorTypeValidation, "VALIDATION_ERROR", err.Error())
	}

	// Default to server error
	return types.NewResourceError(types.ResourceErrorTypeServerError, "SERVER_ERROR", err.Error())
}

// GetMinimal retrieves minimal resource information for performance
func (api *OpsRampResourcesAPI) GetMinimal(ctx context.Context, id string) (*types.ResourceMinimal, error) {
	api.logger.Info("Getting minimal resource with ID: %s", id)

	var resource *types.Resource
	err := api.retryWithBackoff(ctx, "GetMinimal", func() error {
		var err error
		resource, err = api.Get(ctx, id)
		return err
	})

	if err != nil {
		return nil, api.classifyError(err)
	}

	// Convert to minimal resource
	minimal := &types.ResourceMinimal{
		ID:           resource.ID,
		HostName:     resource.HostName,
		IPAddress:    resource.IPAddress,
		Name:         resource.Name,
		ResourceName: resource.ResourceName,
		Type:         resource.Type,
		ResourceType: resource.ResourceType,
		State:        resource.State,
		Status:       resource.Status,
		Location:     resource.Location,
		Tags:         resource.Tags,
		UpdatedDate:  resource.UpdatedDate,
	}

	api.logger.Info("Successfully retrieved minimal resource: %s", minimal.Name)
	return minimal, nil
}
