package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// OAuth2Config holds configuration for OAuth2.0
type OAuth2Config struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
	Scopes       []string
}

// TokenResponse represents the OAuth2.0 token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// AuthClient manages OAuth2.0 tokens
type AuthClient struct {
	Config      OAuth2Config
	token       string
	tokenExpiry time.Time
	mu          sync.Mutex
	httpClient  *http.Client
	logger      *CustomLogger
}

// NewAuthClient creates a new AuthClient
func NewAuthClient(config OAuth2Config) *AuthClient {
	// Get the logger
	logger := GetLogger()

	return &AuthClient{
		Config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
	}
}

// GetToken retrieves a valid OAuth2.0 token, refreshing if necessary
func (a *AuthClient) GetToken() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if we have a valid token
	if a.token != "" && time.Now().Before(a.tokenExpiry) {
		a.logger.Debug("Using cached token, valid until %s", a.tokenExpiry.Format(time.RFC3339))
		return a.token, nil
	}

	a.logger.Info("Token expired or not present, fetching new token")

	// We need to get a new token
	tokenResp, err := a.fetchNewToken()
	if err != nil {
		a.logger.Error("Failed to fetch token: %v", err)
		return "", fmt.Errorf("failed to fetch token: %w", err)
	}

	// Store the token and its expiry time
	a.token = tokenResp.AccessToken
	// Set expiry time with a small buffer to ensure we refresh before actual expiry
	a.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)

	a.logger.Info("Successfully obtained new token, valid until %s", a.tokenExpiry.Format(time.RFC3339))

	return a.token, nil
}

// fetchNewToken gets a new OAuth2.0 token from the authorization server
func (a *AuthClient) fetchNewToken() (*TokenResponse, error) {
	a.logger.Debug("Preparing token request to %s", a.Config.TokenURL)

	// Prepare the token request as form data (x-www-form-urlencoded)
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", a.Config.ClientID)
	formData.Set("client_secret", a.Config.ClientSecret)

	// Add scopes if specified
	if len(a.Config.Scopes) > 0 {
		scopes := strings.Join(a.Config.Scopes, " ")
		formData.Set("scope", scopes)
		a.logger.Debug("Added scopes: %s", scopes)
	}

	// Create the HTTP request with form data
	req, err := http.NewRequest("POST", a.Config.TokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		a.logger.Error("Failed to create token request: %v", err)
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	a.logger.Debug("Set Content-Type header to application/x-www-form-urlencoded")

	// Send the request
	a.logger.Info("Sending token request to %s", a.Config.TokenURL)
	startTime := time.Now()
	resp, err := a.httpClient.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		a.logger.Error("Token request failed: %v", err)
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	a.logger.Info("Token response received in %v with status code %d", duration, resp.StatusCode)

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		a.logger.Error("Token request returned non-OK status: %d", resp.StatusCode)
		return nil, fmt.Errorf("token request returned status %d", resp.StatusCode)
	}

	// Parse the response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		a.logger.Error("Failed to parse token response: %v", err)
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		a.logger.Error("Received empty access token")
		return nil, fmt.Errorf("received empty access token")
	}

	a.logger.Debug("Successfully parsed token response, token type: %s, expires in: %d seconds",
		tokenResp.TokenType, tokenResp.ExpiresIn)
	return &tokenResp, nil
}
