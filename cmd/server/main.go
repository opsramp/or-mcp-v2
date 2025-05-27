package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/client"
	"github.com/vobbilis/codegen/or-mcp-v2/pkg/tools"
)

const (
	// LogDir is the directory where logs will be stored
	LogDir = "output/logs"
	// LogFileName is the name of the log file
	LogFileName = "or-mcp.log"
	// DefaultPort is the default port to listen on
	DefaultPort = 8080
)

var (
	// Global variables to track server state
	startTime       time.Time
	customLogger    *common.CustomLogger
	mcpServer       *server.MCPServer
	sseServer       *server.SSEServer
	registeredTools []string // Track registered tool names manually
)

// jsonRpcRequest represents a JSON-RPC 2.0 request
type jsonRpcRequest struct {
	JsonRpc string                 `json:"jsonrpc"`
	Id      string                 `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

// jsonRpcResponse represents a JSON-RPC 2.0 response
type jsonRpcResponse struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// jsonRpcError represents a JSON-RPC 2.0 error
type jsonRpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// healthHandler provides a simple health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime).String()

	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "hpe-opsramp-mcp",
		"uptime":    uptime,
		"tools":     registeredTools,
		"endpoints": map[string]string{
			"health":    "/health",
			"readiness": "/readiness",
			"sse":       "/sse",
			"message":   "/message",
			"debug":     "/debug",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response) // Ignore encoding errors
}

// readinessHandler provides a more detailed readiness check
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"ready":     true,
		"timestamp": time.Now().Format(time.RFC3339),
		"checks": map[string]interface{}{
			"server":   "ok",
			"sessions": "ok",
			"tools":    "ok",
		},
		"tools": registeredTools,
	}

	// Check if server is initialized
	if mcpServer == nil {
		response["ready"] = false
		response["checks"].(map[string]interface{})["server"] = "not initialized"
	}

	w.Header().Set("Content-Type", "application/json")
	if response["ready"].(bool) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	_ = json.NewEncoder(w).Encode(response) // Ignore encoding errors
}

// debugHandler provides detailed debug information about the server
func debugHandler(w http.ResponseWriter, r *http.Request) {
	debugInfo := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(startTime).String(),
		"tools":     registeredTools,
		"server": map[string]interface{}{
			"name":    "HPE OpsRamp MCP",
			"version": "1.0.0",
		},
	}

	// Include SSE info if session ID is provided
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID != "" {
		debugInfo["session"] = map[string]interface{}{
			"id":     sessionID,
			"exists": false, // Assume false until proven otherwise
		}

		// Simplify testing by accepting any session ID when in debug mode
		// Note: This should only be used for testing/debugging, not production
		debugMode := os.Getenv("DEBUG") == "true"
		if debugMode {
			// Print the debug information to the log
			customLogger.Info("Debug mode enabled, accepting any session ID: %s", sessionID)
			w.Header().Set("X-Accept-Any-Session", "true")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(debugInfo) // Ignore encoding errors
}

// messageHandler handles JSON-RPC requests for tools
func messageHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session ID from query parameter
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		jsonError(w, "Missing session ID", http.StatusBadRequest, "")
		return
	}

	// Debug mode: Accept any session ID for testing
	debugMode := os.Getenv("DEBUG") == "true"
	if !debugMode {
		// In production, verify session ID with the SSE server
		// For now, we'll assume it's valid since we're in development
		customLogger.Debug("Received message for session ID: %s", sessionID)
	} else {
		customLogger.Info("Debug mode: accepting any session ID: %s", sessionID)
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, "Failed to read request body", http.StatusBadRequest, "")
		return
	}
	defer r.Body.Close()

	// Parse JSON-RPC request
	var rpcRequest jsonRpcRequest
	if err := json.Unmarshal(body, &rpcRequest); err != nil {
		jsonError(w, "Invalid JSON-RPC request", http.StatusBadRequest, rpcRequest.Id)
		return
	}

	// Validate JSON-RPC version
	if rpcRequest.JsonRpc != "2.0" {
		jsonError(w, "Unsupported JSON-RPC version", http.StatusBadRequest, rpcRequest.Id)
		return
	}

	// Handle different methods
	var result interface{}
	var methodErr error

	customLogger.Debug("Received JSON-RPC request: method=%s, id=%s", rpcRequest.Method, rpcRequest.Id)

	// Handle different methods
	switch rpcRequest.Method {
	case "callTool":
		// Extract tool name and arguments
		toolName, ok := rpcRequest.Params["name"].(string)
		if !ok {
			jsonError(w, "Missing tool name", http.StatusBadRequest, rpcRequest.Id)
			return
		}

		// Extract arguments
		arguments, ok := rpcRequest.Params["arguments"].(map[string]interface{})
		if !ok {
			arguments = make(map[string]interface{})
		}

		customLogger.Debug("Calling tool: %s with arguments: %v", toolName, arguments)

		// Check if tool exists
		toolExists := false
		for _, registeredTool := range registeredTools {
			if registeredTool == toolName {
				toolExists = true
				break
			}
		}

		if !toolExists {
			jsonError(w, fmt.Sprintf("Tool not found: %s", toolName), http.StatusNotFound, rpcRequest.Id)
			return
		}

		// Get the action parameter
		action, ok := arguments["action"].(string)
		if !ok {
			jsonError(w, "Missing action parameter", http.StatusBadRequest, rpcRequest.Id)
			return
		}

		// Generate log entry for tool execution
		customLogger.Info("Tool Execution: %s, Action: %s, Args: %v", toolName, action, arguments)

		// Execute the tool based on the name
		if toolName == "integrations" {
			// Load configuration
			config, err := common.LoadConfig("")
			if err != nil {
				customLogger.Error("Failed to load config: %v", err)
				jsonError(w, "Failed to load configuration", http.StatusInternalServerError, rpcRequest.Id)
				return
			}

			// Create the integrations API
			integrationsAPI, err := tools.NewOpsRampIntegrationsAPI(&config.OpsRamp)
			if err != nil {
				customLogger.Error("Failed to create integrations API: %v", err)
				jsonError(w, "Failed to initialize integrations API", http.StatusInternalServerError, rpcRequest.Id)
				return
			}

			// Create request for the tool handler
			mcpRequest := mcp.CallToolRequest{
				Params: struct {
					Name      string    `json:"name"`
					Arguments any       `json:"arguments,omitempty"`
					Meta      *mcp.Meta `json:"_meta,omitempty"`
				}{
					Name:      "integrations",
					Arguments: arguments,
				},
			}

			// Log the specific action being called
			customLogger.Debug("Executing integrations action: %s", action)

			// Execute the integrations tool handler directly
			handlerResult, handlerErr := tools.IntegrationsToolHandler(r.Context(), mcpRequest, integrationsAPI)
			if handlerErr != nil {
				customLogger.Error("Error in integrations tool handler: %v", handlerErr)
				methodErr = handlerErr
			} else {
				customLogger.Debug("Integrations tool handler executed successfully")
				// Extract text content from the result
				if handlerResult != nil && len(handlerResult.Content) > 0 {
					customLogger.Debug("Handling result content with %d items", len(handlerResult.Content))
					for _, content := range handlerResult.Content {
						if textContent, ok := content.(mcp.TextContent); ok {
							customLogger.Debug("Processing text content: %s", textContent.Text)
							// Try to parse JSON result
							if err := json.Unmarshal([]byte(textContent.Text), &result); err != nil {
								// If not valid JSON, just use the text
								customLogger.Debug("Not valid JSON, using text directly: %s", err)
								result = textContent.Text
							}
							break
						}
					}
				} else {
					customLogger.Warn("Empty or nil result from integrations tool handler")
					// Return empty array for list operations to prevent null response
					if action == "list" || action == "listTypes" {
						customLogger.Info("Returning empty array for %s action", action)
						result = []interface{}{}
					}
				}
			}
		} else if toolName == "resources" {
			// Load configuration
			config, err := common.LoadConfig("")
			if err != nil {
				customLogger.Error("Failed to load config: %v", err)
				jsonError(w, "Failed to load configuration", http.StatusInternalServerError, rpcRequest.Id)
				return
			}

			// Create the resources API
			opsRampClient := client.NewOpsRampClient(config)
			resourcesAPI := tools.NewOpsRampResourcesAPI(opsRampClient)

			// Create request for the tool handler
			mcpRequest := mcp.CallToolRequest{
				Params: struct {
					Name      string    `json:"name"`
					Arguments any       `json:"arguments,omitempty"`
					Meta      *mcp.Meta `json:"_meta,omitempty"`
				}{
					Name:      "resources",
					Arguments: arguments,
				},
			}

			// Log the specific action being called
			customLogger.Debug("Executing resources action: %s", action)

			// Execute the resources tool handler directly
			handlerResult, handlerErr := tools.ResourcesToolHandler(r.Context(), mcpRequest, resourcesAPI)
			if handlerErr != nil {
				customLogger.Error("Error in resources tool handler: %v", handlerErr)
				methodErr = handlerErr
			} else {
				customLogger.Debug("Resources tool handler executed successfully")
				// Extract text content from the result
				if handlerResult != nil && len(handlerResult.Content) > 0 {
					customLogger.Debug("Handling result content with %d items", len(handlerResult.Content))
					for _, content := range handlerResult.Content {
						if textContent, ok := content.(mcp.TextContent); ok {
							customLogger.Debug("Processing text content: %s", textContent.Text)
							// Try to parse JSON result
							if err := json.Unmarshal([]byte(textContent.Text), &result); err != nil {
								// If not valid JSON, just use the text
								customLogger.Debug("Not valid JSON, using text directly: %s", err)
								result = textContent.Text
							}
							break
						}
					}
				} else {
					customLogger.Warn("Empty or nil result from resources tool handler")
					// Return empty array for list operations to prevent null response
					if action == "list" || action == "getMinimal" {
						customLogger.Info("Returning empty array for %s action", action)
						result = []interface{}{}
					}
				}
			}
		} else {
			jsonError(w, fmt.Sprintf("Tool not implemented: %s", toolName), http.StatusNotImplemented, rpcRequest.Id)
			return
		}
	case "tools/list":
		// Return list of available tools
		result = registeredTools
	default:
		jsonError(w, fmt.Sprintf("Method not supported: %s", rpcRequest.Method), http.StatusBadRequest, rpcRequest.Id)
		return
	}

	// Handle errors
	if methodErr != nil {
		customLogger.Error("Error executing method %s: %v", rpcRequest.Method, methodErr)
		jsonError(w, methodErr.Error(), http.StatusInternalServerError, rpcRequest.Id)
		return
	}

	// Create and send response
	response := jsonRpcResponse{
		JsonRpc: "2.0",
		Id:      rpcRequest.Id,
		Result:  result,
	}

	// Log the response payload for debugging
	responseJSON, _ := json.Marshal(response)
	customLogger.Debug("Sending response: %s", string(responseJSON))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response) // Ignore encoding errors
}

// jsonError sends a JSON-RPC error response
func jsonError(w http.ResponseWriter, message string, httpStatus int, id string) {
	// Log error
	customLogger.Error("JSON-RPC error: %s (HTTP %d)", message, httpStatus)

	// Create JSON-RPC error response
	response := jsonRpcResponse{
		JsonRpc: "2.0",
		Id:      id,
		Error: jsonRpcError{
			Code:    httpStatus,
			Message: message,
		},
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(response) // Ignore encoding errors
}

// startupHealthCheck performs a real API call to verify connectivity and config validity
func startupHealthCheck() error {
	// Load configuration
	config, err := common.LoadConfig("")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create the integrations API
	integrationsAPI, err := tools.NewOpsRampIntegrationsAPI(&config.OpsRamp)
	if err != nil {
		return fmt.Errorf("failed to create integrations API: %w", err)
	}

	// Make a real API call (e.g., list integrations)
	integrations, err := integrationsAPI.List(context.Background())
	if err != nil {
		return fmt.Errorf("startup health check failed: %w", err)
	}

	// Log success
	customLogger.Info("Startup health check passed: successfully listed %d integrations", len(integrations))
	return nil
}

// main is the entry point for the MCP server
func main() {
	// Record start time
	startTime = time.Now()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(LogDir, 0750); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize the logger
	var err error
	customLogger, err = common.InitLogger(common.DEBUG, LogDir, LogFileName)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer customLogger.Close()

	// Log server startup
	customLogger.Info("Starting HPE OpsRamp MCP server")
	customLogger.Info("Log file: %s", filepath.Join(LogDir, LogFileName))

	// Determine port from environment variable
	port := DefaultPort
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		if p, err := strconv.Atoi(portEnv); err == nil && p > 0 && p < 65536 {
			port = p
			customLogger.Info("Using port from environment: %d", port)
		} else {
			customLogger.Warn("Invalid PORT environment variable: %s, using default: %d", portEnv, port)
		}
	}

	// Check if debug mode is enabled
	debugMode := os.Getenv("DEBUG") == "true"
	if debugMode {
		customLogger.Info("*** DEBUG MODE ENABLED ***")
	}

	// Create MCP server
	mcpServer = server.NewMCPServer("HPE OpsRamp MCP", "1.0.0")

	// Initialize the registered tools slice
	registeredTools = make([]string, 0)

	// Register tools
	integrationsTool, integrationsHandler := tools.NewIntegrationsMcpTool()
	mcpServer.AddTool(integrationsTool, integrationsHandler)
	registeredTools = append(registeredTools, integrationsTool.Name)

	// Register resources tool
	resourcesTool, resourcesHandler := tools.NewResourcesMcpTool()
	mcpServer.AddTool(resourcesTool, resourcesHandler)
	registeredTools = append(registeredTools, resourcesTool.Name)

	// Log registered tools
	customLogger.Info("Registered tool: %s", integrationsTool.Name)
	customLogger.Info("Registered tool: %s", resourcesTool.Name)

	// Perform startup health check
	if err := startupHealthCheck(); err != nil {
		customLogger.Warn("Startup health check failed: %v", err)
		customLogger.Info("Continuing with server startup despite health check failure")
	}

	// Create SSE server with appropriate options for MCP
	var options []server.SSEOption
	options = append(options,
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30*time.Second),
		server.WithMessageEndpoint("/message"),
		server.WithSSEEndpoint("/sse"),
		server.WithUseFullURLForMessageEndpoint(true),
		server.WithAppendQueryToMessageEndpoint(),
	)

	// Create SSE server with tools
	sseServer = server.NewSSEServer(mcpServer, options...)
	customLogger.Debug("SSE server created with %d tools", len(registeredTools))
	customLogger.Debug("Registered tools: %v", registeredTools)

	// Create an HTTP mux to handle health, readiness, debug, and SSE endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/readiness", readinessHandler)
	mux.HandleFunc("/debug", debugHandler)
	mux.HandleFunc("/message", messageHandler)
	mux.Handle("/sse", sseServer) // SSE endpoint
	mux.Handle("/", sseServer)    // Fallback for all other paths
	customLogger.Debug("HTTP routes configured")

	// Create HTTP server with the mux
	portString := fmt.Sprintf(":%d", port)
	customLogger.Info("Server listening on %s", portString)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    portString,
		Handler: mux,
		// Increase timeouts
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  240 * time.Second,
	}

	// Start the server
	go func() {
		customLogger.Info("Starting HTTP server on %s", portString)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			customLogger.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	customLogger.Info("Shutting down server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		customLogger.Fatal("Server forced to shutdown: %v", err)
	}

	customLogger.Info("Server exited gracefully")
}
