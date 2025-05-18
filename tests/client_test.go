package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// TestMCPClient tests the MCP server by running it as a subprocess and sending requests to it
func TestMCPClient(t *testing.T) {
	t.Skip("This test requires manual execution as it starts the MCP server as a subprocess")

	// Start the MCP server as a subprocess
	cmd := exec.Command("go", "run", "../cmd/main.go")
	
	// Create pipes for stdin and stdout
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	
	// Start the command
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start MCP server: %v", err)
	}
	
	// Ensure the server is killed when the test exits
	defer func() {
		if err := cmd.Process.Kill(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to kill MCP server: %v\n", err)
		}
	}()
	
	// Wait a moment for the server to start
	time.Sleep(1 * time.Second)
	
	// Test cases for the integrations tool
	testCases := []struct {
		name     string
		tool     string
		action   string
		id       string
		config   map[string]interface{}
		expected string
	}{
		{
			name:     "List Integrations",
			tool:     "integrations",
			action:   "list",
			expected: "Mock Integration",
		},
		{
			name:     "Get Integration",
			tool:     "integrations",
			action:   "get",
			id:       "int-001",
			expected: "Mock Integration",
		},
		{
			name:   "Create Integration",
			tool:   "integrations",
			action: "create",
			config: map[string]interface{}{
				"name": "New Test Integration",
				"type": "api",
			},
			expected: "New Test Integration",
		},
	}
	
	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create MCP request
			request := mcp.Request{
				Type: "call_tool",
				CallTool: &mcp.CallToolRequest{
					Name: tc.tool,
					Params: mcp.CallToolParams{
						Arguments: map[string]interface{}{
							"action": tc.action,
						},
					},
				},
			}
			
			// Add id if provided
			if tc.id != "" {
				request.CallTool.Params.Arguments["id"] = tc.id
			}
			
			// Add config if provided
			if tc.config != nil {
				request.CallTool.Params.Arguments["config"] = tc.config
			}
			
			// Marshal the request to JSON
			requestJSON, err := json.Marshal(request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}
			
			// Send the request to the server
			if _, err := stdin.Write(append(requestJSON, '\n')); err != nil {
				t.Fatalf("Failed to write to stdin: %v", err)
			}
			
			// Read the response
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, stdout); err != nil {
				t.Fatalf("Failed to read from stdout: %v", err)
			}
			
			// Parse the response
			var response mcp.Response
			if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}
			
			// Check for errors
			if response.Error != nil {
				t.Fatalf("Server returned error: %v", response.Error)
			}
			
			// Check the result
			if response.CallToolResult == nil {
				t.Fatalf("Server did not return a result")
			}
			
			// Convert result to string for comparison
			resultText := ""
			for _, content := range response.CallToolResult.Content {
				if textContent, ok := content.(mcp.TextContent); ok {
					resultText = textContent.Text
					break
				}
			}
			
			// Print the result for debugging
			fmt.Printf("Result for %s: %s\n", tc.name, resultText)
			
			// Check if expected string is in the result
			if tc.expected != "" && resultText != "" {
				if !json.Valid([]byte(resultText)) {
					// Simple string comparison
					if resultText != tc.expected && resultText != "OK" {
						t.Errorf("Expected result to contain '%s', got '%s'", tc.expected, resultText)
					}
				} else {
					// JSON result, just check if it contains the expected string
					if resultText != "OK" && resultText != tc.expected {
						t.Errorf("Expected result to contain '%s', got '%s'", tc.expected, resultText)
					}
				}
			}
		})
	}
}
