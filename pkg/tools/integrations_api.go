package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/types"
)

// OpsRampIntegrationsAPI implements the real OpsRamp API client
type OpsRampIntegrationsAPI struct {
	config     *common.OpsRampConfig
	httpClient *http.Client
	baseURL    string
	authURL    string
	authToken  string
	tokenExp   time.Time
	logger     *common.CustomLogger
}

// NewOpsRampIntegrationsAPI creates a new client for accessing the OpsRamp API
func NewOpsRampIntegrationsAPI(config *common.OpsRampConfig) (*OpsRampIntegrationsAPI, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Validate required config fields
	if config.TenantURL == "" || config.AuthURL == "" ||
		config.AuthKey == "" || config.AuthSecret == "" ||
		config.TenantID == "" {
		return nil, fmt.Errorf("invalid OpsRamp configuration: missing required fields")
	}

	// Check if using placeholder values
	if strings.Contains(config.TenantURL, "your-tenant-instance") ||
		strings.Contains(config.AuthKey, "YOUR_AUTH_KEY") ||
		strings.Contains(config.AuthSecret, "YOUR_AUTH_SECRET") ||
		strings.Contains(config.TenantID, "YOUR_TENANT_ID") {
		return nil, fmt.Errorf("invalid OpsRamp configuration: contains placeholder values")
	}

	api := &OpsRampIntegrationsAPI{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		config:  config,
		baseURL: config.TenantURL,
		logger:  common.GetLogger(),
	}

	// Authenticate to verify credentials immediately
	err := api.authenticate(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with OpsRamp API: %w", err)
	}

	return api, nil
}

// authenticate obtains a new OAuth token from OpsRamp
func (a *OpsRampIntegrationsAPI) authenticate(ctx context.Context) error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", a.config.AuthKey)
	data.Set("client_secret", a.config.AuthSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", a.config.AuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating auth request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error during auth request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading auth response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &authResp); err != nil {
		return fmt.Errorf("error unmarshaling auth response: %w", err)
	}

	a.authToken = authResp.AccessToken
	// Set expiry time with a small buffer to ensure we refresh before actual expiry
	a.tokenExp = time.Now().Add(time.Duration(authResp.ExpiresIn-60) * time.Second)

	return nil
}

// ensureAuth ensures a valid authentication token is available
func (a *OpsRampIntegrationsAPI) ensureAuth(ctx context.Context) error {
	if a.authToken == "" || time.Now().After(a.tokenExp) {
		return a.authenticate(ctx)
	}
	return nil
}

// makeRequest makes an authenticated request to the OpsRamp API
func (a *OpsRampIntegrationsAPI) makeRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	// Ensure we have a valid auth token
	if err := a.ensureAuth(ctx); err != nil {
		return nil, err
	}

	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
	}

	// Format according to OpsRamp API documentation
	// The URL should be in the format: {baseURL}/api/v2/tenants/{tenantId}/integrations/{path}
	fullURL := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/%s", a.baseURL, a.config.TenantID, path)
	a.logger.Debug("Making API request to URL: %s", fullURL)
	a.logger.Debug("Request method: %s, path: %s", method, path)
	a.logger.Debug("Base URL: %s, Tenant ID: %s", a.baseURL, a.config.TenantID)

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.authToken))

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		a.logger.Error("API request failed with status %d: %s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// List returns all integrations
func (a *OpsRampIntegrationsAPI) List(ctx context.Context) ([]types.Integration, error) {
	// Based on OpsRamp API docs: /api/v2/tenants/{tenantId}/integrations/installed/search
	respBody, err := a.makeRequest(ctx, "GET", "installed/search", nil)
	if err != nil {
		return nil, fmt.Errorf("error listing integrations: %w", err)
	}

	// Log the raw response for debugging
	a.logger.Debug("Raw response: %s", string(respBody))

	// Try to parse as a structured response first
	var structuredResp struct {
		Results []types.Integration `json:"results"`
	}

	if err := json.Unmarshal(respBody, &structuredResp); err != nil {
		// If that fails, try parsing as a direct array
		var directArray []types.Integration
		if err2 := json.Unmarshal(respBody, &directArray); err2 != nil {
			return nil, fmt.Errorf("error unmarshaling integration list: %w (array error: %w)", err, err2)
		}
		return directArray, nil
	}

	// If we got a structured response with results field
	if len(structuredResp.Results) > 0 {
		return structuredResp.Results, nil
	}

	// Final fallback: try to parse as a direct array
	var directArray []types.Integration
	if err := json.Unmarshal(respBody, &directArray); err != nil {
		return nil, fmt.Errorf("error unmarshaling integration list: %w", err)
	}

	return directArray, nil
}

// Get returns a specific integration by ID
func (a *OpsRampIntegrationsAPI) Get(ctx context.Context, id string) (*types.Integration, error) {
	// Based on OpsRamp API docs: /api/v2/tenants/{tenantId}/integrations/installed/{installedIntgId}
	respBody, err := a.makeRequest(ctx, "GET", fmt.Sprintf("installed/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error getting integration %s: %w", id, err)
	}

	// Log the raw response for debugging
	a.logger.Debug("Raw response: %s", string(respBody))

	var integration types.Integration
	if err := json.Unmarshal(respBody, &integration); err != nil {
		// Check if this might be a response with a nested field
		var wrappedResp map[string]types.Integration
		if err2 := json.Unmarshal(respBody, &wrappedResp); err2 != nil {
			return nil, fmt.Errorf("error unmarshaling integration: %w", err)
		}

		// Check common wrapper fields
		for _, field := range []string{"integration", "result", "data"} {
			if wrapped, ok := wrappedResp[field]; ok {
				return &wrapped, nil
			}
		}

		return nil, fmt.Errorf("error unmarshaling integration: %w", err)
	}

	// Ensure ID is set
	if integration.ID == "" {
		integration.ID = id
	}

	return &integration, nil
}

// GetDetailed returns detailed information about an integration
func (a *OpsRampIntegrationsAPI) GetDetailed(ctx context.Context, id string) (*types.DetailedIntegration, error) {
	// Using same endpoint as Get with additional processing if needed
	respBody, err := a.makeRequest(ctx, "GET", fmt.Sprintf("installed/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error getting detailed integration %s: %w", id, err)
	}

	// Log the raw response for debugging
	a.logger.Debug("Raw response: %s", string(respBody))

	var integration types.DetailedIntegration
	if err := json.Unmarshal(respBody, &integration); err != nil {
		// Check if this might be a response with a nested field
		var wrappedResp map[string]types.DetailedIntegration
		if err2 := json.Unmarshal(respBody, &wrappedResp); err2 != nil {
			return nil, fmt.Errorf("error unmarshaling detailed integration: %w", err)
		}

		// Check common wrapper fields
		for _, field := range []string{"integration", "result", "data"} {
			if wrapped, ok := wrappedResp[field]; ok {
				return &wrapped, nil
			}
		}

		return nil, fmt.Errorf("error unmarshaling detailed integration: %w", err)
	}

	// Ensure ID is set
	if integration.ID == "" {
		integration.ID = id
	}

	return &integration, nil
}

// Create creates a new integration
func (a *OpsRampIntegrationsAPI) Create(ctx context.Context, config map[string]interface{}) (*types.Integration, error) {
	// Get the integration name from the config
	intgName, ok := config["name"].(string)
	if !ok || intgName == "" {
		return nil, fmt.Errorf("integration name is required")
	}

	// Based on OpsRamp API docs: /api/v2/tenants/{tenantId}/integrations/install/{uniqueName}
	respBody, err := a.makeRequest(ctx, "POST", fmt.Sprintf("install/%s", intgName), config)
	if err != nil {
		return nil, fmt.Errorf("error creating integration: %w", err)
	}

	var integration types.Integration
	if err := json.Unmarshal(respBody, &integration); err != nil {
		return nil, fmt.Errorf("error unmarshaling created integration: %w", err)
	}

	return &integration, nil
}

// Update updates an existing integration
func (a *OpsRampIntegrationsAPI) Update(ctx context.Context, id string, config map[string]interface{}) (*types.Integration, error) {
	// Based on OpsRamp API docs: /api/v2/tenants/{tenantId}/integrations/installed/{installedIntgId}
	respBody, err := a.makeRequest(ctx, "POST", fmt.Sprintf("installed/%s", id), config)
	if err != nil {
		return nil, fmt.Errorf("error updating integration %s: %w", id, err)
	}

	var integration types.Integration
	if err := json.Unmarshal(respBody, &integration); err != nil {
		return nil, fmt.Errorf("error unmarshaling updated integration: %w", err)
	}

	return &integration, nil
}

// Delete removes an integration
func (a *OpsRampIntegrationsAPI) Delete(ctx context.Context, id string) error {
	// Based on OpsRamp API docs: /api/v2/tenants/{tenantId}/integrations/installed/{installedIntgId}
	_, err := a.makeRequest(ctx, "DELETE", fmt.Sprintf("installed/%s", id), nil)
	if err != nil {
		return fmt.Errorf("error deleting integration %s: %w", id, err)
	}

	return nil
}

// Enable enables an integration
func (a *OpsRampIntegrationsAPI) Enable(ctx context.Context, id string) error {
	// Based on OpsRamp API docs: /api/v2/tenants/{tenantId}/integrations/installed/{installedIntgId}/{actions}
	// where actions is 'enable'
	_, err := a.makeRequest(ctx, "POST", fmt.Sprintf("installed/%s/enable", id), nil)
	if err != nil {
		return fmt.Errorf("error enabling integration %s: %w", id, err)
	}

	return nil
}

// Disable disables an integration
func (a *OpsRampIntegrationsAPI) Disable(ctx context.Context, id string) error {
	// Based on OpsRamp API docs: /api/v2/tenants/{tenantId}/integrations/installed/{installedIntgId}/{actions}
	// where actions is 'disable'
	_, err := a.makeRequest(ctx, "POST", fmt.Sprintf("installed/%s/disable", id), nil)
	if err != nil {
		return fmt.Errorf("error disabling integration %s: %w", id, err)
	}

	return nil
}

// ListTypes returns all integration types
func (a *OpsRampIntegrationsAPI) ListTypes(ctx context.Context) ([]types.IntegrationType, error) {
	// Based on OpsRamp API docs: /api/v2/tenants/{tenantId}/integrations/available/search
	respBody, err := a.makeRequest(ctx, "GET", "available/search", nil)
	if err != nil {
		return nil, fmt.Errorf("error listing integration types: %w", err)
	}

	// Log the raw response for debugging
	a.logger.Debug("Raw response: %s", string(respBody))

	// Try to parse as a structured response first with results
	var structuredResp struct {
		Results []map[string]interface{} `json:"results"`
	}

	// Direct array of integration types (alternate response format)
	var directTypes []types.IntegrationType

	// Try the structured format first
	err1 := json.Unmarshal(respBody, &structuredResp)
	if err1 == nil && len(structuredResp.Results) > 0 {
		// Extract unique types from structured response
		typeMap := make(map[string]types.IntegrationType)
		for _, integration := range structuredResp.Results {
			typeInfo, ok := integration["type"].(map[string]interface{})
			if !ok {
				continue
			}

			id, ok := typeInfo["id"].(string)
			if !ok {
				continue
			}

			name, ok := typeInfo["name"].(string)
			if !ok {
				continue
			}

			if _, exists := typeMap[id]; !exists {
				typeMap[id] = types.IntegrationType{
					ID:   id,
					Name: name,
				}
			}
		}

		// Convert map to slice
		var types []types.IntegrationType
		for _, integrationType := range typeMap {
			types = append(types, integrationType)
		}

		return types, nil
	}

	// Try direct array format
	err2 := json.Unmarshal(respBody, &directTypes)
	if err2 == nil {
		return directTypes, nil
	}

	// Both parsing attempts failed
	return nil, fmt.Errorf("error unmarshaling integration types: structured error: %w, direct error: %v", err1, err2)
}

// GetType returns a specific integration type
func (a *OpsRampIntegrationsAPI) GetType(ctx context.Context, id string) (*types.IntegrationType, error) {
	// There's no specific endpoint for getting integration type by ID
	// So we'll get all types and filter by ID
	types, err := a.ListTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting integration types: %w", err)
	}

	for _, integrationType := range types {
		if integrationType.ID == id {
			return &integrationType, nil
		}
	}

	return nil, fmt.Errorf("integration type with ID %s not found", id)
}
