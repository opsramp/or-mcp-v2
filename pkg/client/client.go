package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
"strings"
	"net/url"
	"path"
	"time"

	"github.com/vobbilis/codegen/or-mcp-v2/common"
)

// OpsRampClient is the client for the OpsRamp API
type OpsRampClient struct {
	baseURL    string
	tenantID   string
	authClient *common.AuthClient
	httpClient *http.Client
	logger     *common.CustomLogger
}

// NewOpsRampClient creates a new OpsRamp API client
func NewOpsRampClient(config *common.Config) *OpsRampClient {
	// Create auth client
	authConfig := common.OAuth2Config{
		ClientID:     config.OpsRamp.AuthKey,
		ClientSecret: config.OpsRamp.AuthSecret,
		TokenURL:     config.OpsRamp.AuthURL,
	}
	authClient := common.NewAuthClient(authConfig)

	// Get the logger
	logger := common.GetLogger()

	return &OpsRampClient{
		baseURL:    config.OpsRamp.TenantURL,
		tenantID:   config.OpsRamp.TenantID,
		authClient: authClient,
		httpClient: &http.Client{Timeout: 60 * time.Second},
		logger:     logger,
	}
}

// Request makes an authenticated request to the OpsRamp API
func (c *OpsRampClient) Request(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	_, err := c.RequestWithStatusCode(ctx, method, endpoint, body, result)
	return err
}

// RequestWithStatusCode makes an authenticated request to the OpsRamp API and returns the status code
func (c *OpsRampClient) RequestWithStatusCode(ctx context.Context, method, endpoint string, body interface{}, result interface{}) (int, error) {
	// Log the request
	c.logger.Debug("API Request: %s %s", method, endpoint)

	// Build the full URL
// Build the full URL
u, err := url.Parse(c.baseURL)
if err != nil {
c.logger.Error("Invalid base URL: %v", err)
return 0, fmt.Errorf("invalid base URL: %w", err)
}

// Check if endpoint contains query parameters
endpointParts := strings.SplitN(endpoint, "?", 2)
u.Path = path.Join(u.Path, endpointParts[0])

// If there are query parameters, add them to the URL
if len(endpointParts) > 1 {
u.RawQuery = endpointParts[1]
}

// Log the full URL
c.logger.Debug("Full URL: %s", u.String())
	// Prepare request body if provided
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			c.logger.Error("Failed to marshal request body: %v", err)
			return 0, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
		c.logger.Debug("Request Body: %s", string(jsonBody))
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		c.logger.Error("Failed to create request: %v", err)
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get and set the auth token
	token, err := c.authClient.GetToken()
	if err != nil {
		c.logger.Error("Failed to get auth token: %v", err)
		return 0, fmt.Errorf("failed to get auth token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	c.logger.Debug("Auth token obtained and set")

	// Set tenant ID if provided
	if c.tenantID != "" {
		req.Header.Set("X-Tenant-ID", c.tenantID)
		c.logger.Debug("Tenant ID set: %s", c.tenantID)
	}

	// Log request details
	c.logger.Info("Sending %s request to %s", method, u.String())

	// Send the request
	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		c.logger.Error("Request failed: %v", err)
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Log response details
	c.logger.Info("Response received in %v with status code %d", duration, resp.StatusCode)

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to read error response
		errorBody, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, string(errorBody))
		c.logger.Error(errorMsg)
		return resp.StatusCode, fmt.Errorf(errorMsg)
	}

	// Parse the response if a result container was provided
	if result != nil {
		// Read the response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.logger.Error("Failed to read response body: %v", err)
			return resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
		}

		// Log the response body (truncated if too large)
		respBodyStr := string(respBody)
		if len(respBodyStr) > 1000 {
			c.logger.Debug("Response Body (truncated): %s...", respBodyStr[:1000])
		} else {
			c.logger.Debug("Response Body: %s", respBodyStr)
		}

		// Parse the response
		if err := json.Unmarshal(respBody, result); err != nil {
			c.logger.Error("Failed to parse response: %v", err)
			return resp.StatusCode, fmt.Errorf("failed to parse response: %w", err)
		}

		c.logger.Debug("Response successfully parsed")
	}

	c.logger.Info("Request completed successfully")
	return resp.StatusCode, nil
}

// Get makes a GET request to the OpsRamp API
func (c *OpsRampClient) Get(ctx context.Context, endpoint string, result interface{}) error {
	return c.Request(ctx, http.MethodGet, endpoint, nil, result)
}

// GetWithStatusCode makes a GET request to the OpsRamp API and returns the status code
func (c *OpsRampClient) GetWithStatusCode(ctx context.Context, endpoint string, result interface{}) (int, error) {
	return c.RequestWithStatusCode(ctx, http.MethodGet, endpoint, nil, result)
}

// Post makes a POST request to the OpsRamp API
func (c *OpsRampClient) Post(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, endpoint, body, result)
}

// PostWithStatusCode makes a POST request to the OpsRamp API and returns the status code
func (c *OpsRampClient) PostWithStatusCode(ctx context.Context, endpoint string, body interface{}, result interface{}) (int, error) {
	return c.RequestWithStatusCode(ctx, http.MethodPost, endpoint, body, result)
}

// Put makes a PUT request to the OpsRamp API
func (c *OpsRampClient) Put(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPut, endpoint, body, result)
}

// Delete makes a DELETE request to the OpsRamp API
func (c *OpsRampClient) Delete(ctx context.Context, endpoint string) error {
	return c.Request(ctx, http.MethodDelete, endpoint, nil, nil)
}

// Patch makes a PATCH request to the OpsRamp API
func (c *OpsRampClient) Patch(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPatch, endpoint, body, result)
}

// GetTenantID returns the tenant ID
func (c *OpsRampClient) GetTenantID() string {
	return c.tenantID
}

// Global client instance
var globalClient *OpsRampClient
var clientInitialized bool

// GetOpsRampClient returns the global OpsRampClient instance
func GetOpsRampClient() *OpsRampClient {
	if !clientInitialized {
		// Load configuration
		config, err := common.LoadConfig("")
		if err != nil {
			// Get the logger
			logger := common.GetLogger()
			logger.Error("Failed to load config for OpsRamp client: %v", err)

			// Create a minimal default config
			config = &common.Config{
				OpsRamp: common.OpsRampConfig{
					TenantURL:  common.GetEnvOrDefault("OPSRAMP_TENANT_URL", "https://api.opsramp.com"),
					AuthURL:    common.GetEnvOrDefault("OPSRAMP_AUTH_URL", "https://api.opsramp.com/auth/oauth/token"),
					AuthKey:    common.GetEnvOrDefault("OPSRAMP_AUTH_KEY", ""),
					AuthSecret: common.GetEnvOrDefault("OPSRAMP_AUTH_SECRET", ""),
					TenantID:   common.GetEnvOrDefault("OPSRAMP_TENANT_ID", ""),
				},
			}
		}

		// Create the client
		globalClient = NewOpsRampClient(config)
		clientInitialized = true
	}

	return globalClient
}

// SetGlobalClient sets the global OpsRampClient instance
func SetGlobalClient(client *OpsRampClient) {
	globalClient = client
	clientInitialized = true
}
