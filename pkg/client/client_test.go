package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vobbilis/codegen/or-mcp-v2/common"
)

func TestOpsRampClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method and path
		if r.Method == http.MethodGet && r.URL.Path == "/api/test" {
			// Return a successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": true, "message": "Test successful"}`))
			return
		}

		// Return a 404 for any other request
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Create a test config
	config := &common.Config{
		OpsRamp: common.OpsRampConfig{
			TenantURL:  server.URL,
			AuthURL:    server.URL + "/auth/token",
			AuthKey:    "test-key",
			AuthSecret: "test-secret",
			TenantID:   "test-tenant",
		},
	}

	// Create the client
	client := NewOpsRampClient(config)

	// Test the Get method
	t.Run("Get", func(t *testing.T) {
		var result map[string]interface{}
		err := client.Get(context.Background(), "/api/test", &result)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Check the result
		if success, ok := result["success"].(bool); !ok || !success {
			t.Errorf("Expected success to be true, got %v", result["success"])
		}

		if message, ok := result["message"].(string); !ok || message != "Test successful" {
			t.Errorf("Expected message to be 'Test successful', got %v", result["message"])
		}
	})

	// Test with an invalid endpoint
	t.Run("Invalid Endpoint", func(t *testing.T) {
		var result map[string]interface{}
		err := client.Get(context.Background(), "/invalid", &result)
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})
}
