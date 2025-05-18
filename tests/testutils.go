// Shared test utilities for integration and unit tests
package tests

import "github.com/vobbilis/codegen/or-mcp-v2/common"

// NewTestAuthClient creates a test AuthClient with the given credentials
func NewTestAuthClient(clientID, clientSecret, tokenURL string) *common.AuthClient {
	return &common.AuthClient{
		Config: common.OAuth2Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     tokenURL,
		},
	}
}
