package mcp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/server"
	"github.com/opsramp/or-mcp-v2/common"
)

// InspectorHandler handles MCP Inspector compatibility requirements
type InspectorHandler struct {
	mcpServer *server.MCPServer
	logger    *common.CustomLogger
}

// NewInspectorHandler creates a new MCP Inspector compatibility handler
func NewInspectorHandler(mcpServer *server.MCPServer, logger *common.CustomLogger) *InspectorHandler {
	return &InspectorHandler{
		mcpServer: mcpServer,
		logger:    logger,
	}
}

// jsonRpcRequest represents a JSON-RPC 2.0 request
type jsonRpcRequest struct {
	JsonRpc string                 `json:"jsonrpc"`
	Id      interface{}            `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// jsonRpcResponse represents a JSON-RPC 2.0 response
type jsonRpcResponse struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// jsonRpcError represents a JSON-RPC 2.0 error
type jsonRpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HandleMessage processes MCP Inspector messages with special compatibility handling
func (h *InspectorHandler) HandleMessage(w http.ResponseWriter, r *http.Request) {
	// Validate request method
	if r.Method != http.MethodPost {
		h.logger.Warn("Received non-POST request to /message endpoint: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate session and handle MCP Inspector compatibility
	if !h.validateSessionAndRoute(w, r) {
		return
	}

	// Read and validate request body
	body, rpcRequest, ok := h.readAndParseRequestBody(w, r)
	if !ok {
		return
	}

	// Handle response messages (acknowledgments)
	if h.handleResponseMessage(w, r, body) {
		return
	}

	// Validate JSON-RPC request
	if !h.validateJsonRpcRequest(w, rpcRequest) {
		return
	}

	// Handle MCP protocol methods
	if h.handleMCPProtocolMethods(w, r, rpcRequest) {
		return
	}

	// Handle tool calls by delegating to the MCP server
	if h.handleToolCalls(w, r, rpcRequest) {
		return
	}

	// If we get here, the method is not supported
	h.jsonError(w, "Method not found", http.StatusNotFound, rpcRequest.Id)
}

// validateSessionAndRoute validates session and routes MCP Inspector requests
func (h *InspectorHandler) validateSessionAndRoute(w http.ResponseWriter, r *http.Request) bool {
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		h.logger.Warn("Missing sessionId parameter in request")
		http.Error(w, "Missing sessionId parameter", http.StatusBadRequest)
		return false
	}

	// Debug mode accepts any session ID for MCP Inspector compatibility
	// Check environment variable for debug mode (like the original implementation)
	envDebugMode := os.Getenv("DEBUG") == "true"

	// Also check for MCP Inspector specific headers
	mcpInspectorMode := strings.Contains(r.Header.Get("User-Agent"), "MCP-Inspector") ||
		r.Header.Get("Accept") == "text/event-stream" ||
		strings.Contains(r.Header.Get("Accept"), "text/event-stream")

	if envDebugMode || mcpInspectorMode {
		h.logger.Info("Debug mode: accepting any session ID: %s", sessionID)
		return true
	}

	// In production mode, validate session exists
	// For now, we'll accept any session ID but log it
	h.logger.Info("Production mode: validating session ID: %s", sessionID)
	return true
}

// readAndParseRequestBody reads and parses the JSON-RPC request body
func (h *InspectorHandler) readAndParseRequestBody(w http.ResponseWriter, r *http.Request) ([]byte, *jsonRpcRequest, bool) {
	h.logger.Debug("About to read request body - Content-Length: %s", r.Header.Get("Content-Length"))

	// Read the request body properly
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var rawBody json.RawMessage
	if err := decoder.Decode(&rawBody); err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return nil, nil, false
	}

	body := []byte(rawBody)
	h.logger.Debug("Successfully read %d bytes from request body", len(body))
	h.logger.Debug("Raw request body: %s", string(body))
	h.logger.Debug("Request Content-Type: %s", r.Header.Get("Content-Type"))
	h.logger.Debug("Request method: %s, URL: %s", r.Method, r.URL.String())

	// Parse JSON-RPC request
	var rpcRequest jsonRpcRequest
	if err := json.Unmarshal(body, &rpcRequest); err != nil {
		h.logger.Error("Failed to parse JSON-RPC request: %v", err)
		http.Error(w, "Invalid JSON-RPC request", http.StatusBadRequest)
		return nil, nil, false
	}

	h.logger.Debug("Successfully parsed JSON-RPC request: method=%s, id=%v", rpcRequest.Method, rpcRequest.Id)
	return body, &rpcRequest, true
}

// handleResponseMessage handles JSON-RPC response messages (acknowledgments from MCP Inspector)
func (h *InspectorHandler) handleResponseMessage(w http.ResponseWriter, r *http.Request, body []byte) bool {
	var responseCheck map[string]interface{}
	if err := json.Unmarshal(body, &responseCheck); err != nil {
		return false
	}

	// Check if this is a response message (has "result" or "error" but no "method")
	if _, hasResult := responseCheck["result"]; hasResult {
		if _, hasMethod := responseCheck["method"]; !hasMethod {
			return h.handleResultResponse(w, r, responseCheck)
		}
	}

	if _, hasError := responseCheck["error"]; hasError {
		return h.handleErrorResponse(w, r)
	}
	return false
}

// handleResultResponse handles JSON-RPC result responses
func (h *InspectorHandler) handleResultResponse(w http.ResponseWriter, r *http.Request, responseCheck map[string]interface{}) bool {
	h.logger.Info("Received JSON-RPC response message (result) - this is likely an acknowledge from MCP Inspector")

	// Check if this is the acknowledgment after initialization
	if id, ok := responseCheck["id"]; ok {
		var idValue float64
		switch v := id.(type) {
		case float64:
			idValue = v
		case int:
			idValue = float64(v)
		case int64:
			idValue = float64(v)
		}

		// Accept acknowledgment with id=1 (typical pattern after initialize with id=0)
		if idValue == 1 {
			h.logger.Info("Received initialization acknowledgment (id=%v) - MCP Inspector is ready for operations", id)
			h.logger.Info("Handshake complete - sending 'initialized' notification")

			// Send the 'initialized' notification to complete the handshake
			initializedNotification := map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  "initialized",
			}

			initializedJSON, err := json.Marshal(initializedNotification)
			if err != nil {
				h.logger.Error("Error marshaling initialized notification: %v", err)
			} else {
				h.sendInitializedNotification(w, r, initializedJSON)
			}
			return true
		}
	}
	return false
}

// handleErrorResponse handles JSON-RPC error responses
func (h *InspectorHandler) handleErrorResponse(w http.ResponseWriter, r *http.Request) bool {
	h.logger.Info("Received JSON-RPC response message (error) - this is likely an acknowledge from MCP Inspector")

	// Check if client expects SSE format
	acceptHeader := r.Header.Get("Accept")
	if strings.Contains(acceptHeader, "text/event-stream") {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
	}
	w.WriteHeader(http.StatusOK)
	return true
}

// sendInitializedNotification sends the initialized notification to complete handshake
func (h *InspectorHandler) sendInitializedNotification(w http.ResponseWriter, r *http.Request, initializedJSON []byte) {
	// Check if client expects SSE format
	acceptHeader := r.Header.Get("Accept")
	if strings.Contains(acceptHeader, "text/event-stream") {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		sseResponse := fmt.Sprintf("event: message\ndata: %s\n\n", string(initializedJSON))
		h.logger.Info("Sending SSE initialized notification: %s", string(initializedJSON))
		w.Write([]byte(sseResponse))
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		h.logger.Info("Sending JSON initialized notification: %s", string(initializedJSON))
		w.Write(initializedJSON)
	}
}

// validateJsonRpcRequest validates the JSON-RPC request
func (h *InspectorHandler) validateJsonRpcRequest(w http.ResponseWriter, rpcRequest *jsonRpcRequest) bool {
	// Validate JSON-RPC version
	if rpcRequest.JsonRpc != "2.0" {
		h.jsonError(w, "Unsupported JSON-RPC version", http.StatusBadRequest, rpcRequest.Id)
		return false
	}

	// Handle empty method (this shouldn't happen with proper JSON-RPC, but let's be defensive)
	if rpcRequest.Method == "" {
		h.logger.Warn("Received JSON-RPC request with empty method - this might be a malformed response message")
		w.WriteHeader(http.StatusOK)
		return false
	}

	return true
}

// handleMCPProtocolMethods handles standard MCP protocol methods
func (h *InspectorHandler) handleMCPProtocolMethods(w http.ResponseWriter, r *http.Request, rpcRequest *jsonRpcRequest) bool {
	h.logger.Debug("Received JSON-RPC request: method=%s, id=%v", rpcRequest.Method, rpcRequest.Id)

	// Check if this is a standard MCP protocol method (for MCP Inspector compatibility)
	switch rpcRequest.Method {
	case "initialize":
		return h.handleInitializeMethod(w, r, rpcRequest)
	case "tools/list":
		return h.handleMCPToolsMethod(w, r, rpcRequest)
	}
	return false
}

// handleInitializeMethod handles the initialize method for MCP Inspector compatibility
func (h *InspectorHandler) handleInitializeMethod(w http.ResponseWriter, r *http.Request, rpcRequest *jsonRpcRequest) bool {
	h.logger.Info("Received initialize request - handling manually for protocol compatibility")

	// Extract the requested protocol version
	requestedVersion := "2024-11-05" // default
	if params, ok := rpcRequest.Params["protocolVersion"].(string); ok {
		requestedVersion = params
		h.logger.Info("MCP Inspector requested protocol version: %s", requestedVersion)
	}

	// Create a manual initialize response that matches MCP Inspector's expectations
	initResponse := jsonRpcResponse{
		JsonRpc: "2.0",
		Id:      rpcRequest.Id,
		Result: map[string]interface{}{
			"protocolVersion": requestedVersion,
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{
					"listChanged": true,
				},
				"logging": map[string]interface{}{},
				"resources": map[string]interface{}{
					"listChanged": true,
					"subscribe":   false,
				},
			},
			"serverInfo": map[string]interface{}{
				"name":    "HPE OpsRamp MCP",
				"version": "1.0.0",
			},
			"instructions": "HPE OpsRamp MCP Server providing access to OpsRamp integrations and resources. Use the 'integrations' tool to manage integrations and the 'resources' tool to access OpsRamp resources.",
		},
	}

	h.logger.Info("Sending manual initialize response with protocol version: %s", requestedVersion)
	h.sendMCPResponse(w, r, initResponse)
	return true
}

// handleMCPToolsMethod handles MCP tools/list method by delegating to the MCP server
func (h *InspectorHandler) handleMCPToolsMethod(w http.ResponseWriter, r *http.Request, rpcRequest *jsonRpcRequest) bool {
	h.logger.Info("Received standard MCP protocol method: %s - routing to MCP server", rpcRequest.Method)

	// Create a proper MCP protocol message and route it through the MCP server
	mcpMessage, err := json.Marshal(rpcRequest)
	if err != nil {
		h.logger.Error("Failed to marshal MCP message: %v", err)
		h.jsonError(w, "Failed to process MCP message", http.StatusInternalServerError, rpcRequest.Id)
		return true
	}

	h.logger.Debug("Sending to MCP server - method: %s, message: %s", rpcRequest.Method, string(mcpMessage))

	// Process through the MCP server directly
	mcpResponse := h.mcpServer.HandleMessage(r.Context(), json.RawMessage(mcpMessage))

	h.logger.Debug("MCP server response for method %s: %v", rpcRequest.Method, mcpResponse)

	if mcpResponse != nil {
		h.logger.Info("Sending MCP response for method %s", rpcRequest.Method)
		h.sendMCPResponse(w, r, mcpResponse)
	} else {
		h.logger.Warn("No response from MCP server for method %s", rpcRequest.Method)
		w.WriteHeader(http.StatusOK)
	}
	return true
}

// handleToolCalls handles tool/call requests by delegating to the MCP server
func (h *InspectorHandler) handleToolCalls(w http.ResponseWriter, r *http.Request, rpcRequest *jsonRpcRequest) bool {
	if rpcRequest.Method != "tools/call" {
		return false
	}

	h.logger.Info("Received tool call request - delegating to MCP server")

	// Create a proper MCP protocol message and route it through the MCP server
	mcpMessage, err := json.Marshal(rpcRequest)
	if err != nil {
		h.logger.Error("Failed to marshal tool call message: %v", err)
		h.jsonError(w, "Failed to process tool call", http.StatusInternalServerError, rpcRequest.Id)
		return true
	}

	h.logger.Debug("Sending tool call to MCP server: %s", string(mcpMessage))

	// Process through the MCP server directly
	mcpResponse := h.mcpServer.HandleMessage(r.Context(), json.RawMessage(mcpMessage))

	if mcpResponse != nil {
		h.logger.Info("Sending tool call response")
		h.sendMCPResponse(w, r, mcpResponse)
	} else {
		h.logger.Warn("No response from MCP server for tool call")
		w.WriteHeader(http.StatusOK)
	}
	return true
}

// sendMCPResponse sends MCP responses in the appropriate format (SSE or JSON)
func (h *InspectorHandler) sendMCPResponse(w http.ResponseWriter, r *http.Request, response interface{}) {
	// Check if client expects Server-Sent Events (like MCP Inspector)
	acceptHeader := r.Header.Get("Accept")
	if strings.Contains(acceptHeader, "text/event-stream") {
		h.logger.Info("Client expects SSE format - sending as event-stream")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		// Format response as SSE
		responseBytes, err := json.Marshal(response)
		if err != nil {
			h.logger.Error("Failed to marshal SSE response: %v", err)
			return
		}

		sseResponse := fmt.Sprintf("event: message\ndata: %s\n\n", string(responseBytes))
		h.logger.Debug("Sending SSE response: %s", string(responseBytes))
		w.Write([]byte(sseResponse))
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	} else {
		// Send as regular JSON
		h.logger.Info("Client expects JSON format - sending as application/json")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			h.logger.Error("Failed to encode JSON response: %v", err)
		}
	}
}

// jsonError sends a JSON-RPC error response
func (h *InspectorHandler) jsonError(w http.ResponseWriter, message string, httpStatus int, id interface{}) {
	h.logger.Error("JSON-RPC error: %s (HTTP %d)", message, httpStatus)

	response := jsonRpcResponse{
		JsonRpc: "2.0",
		Id:      id,
		Error: jsonRpcError{
			Code:    httpStatus,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(response)
}
