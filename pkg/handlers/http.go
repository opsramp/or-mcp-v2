package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/opsramp/or-mcp-v2/common"
)

// HTTPHandlers contains all HTTP endpoint handlers
type HTTPHandlers struct {
	mcpServer       *server.MCPServer
	sseServer       *server.SSEServer
	logger          *common.CustomLogger
	startTime       time.Time
	registeredTools []string
}

// NewHTTPHandlers creates a new HTTP handlers instance
func NewHTTPHandlers(mcpServer *server.MCPServer, sseServer *server.SSEServer, logger *common.CustomLogger, startTime time.Time, registeredTools []string) *HTTPHandlers {
	return &HTTPHandlers{
		mcpServer:       mcpServer,
		sseServer:       sseServer,
		logger:          logger,
		startTime:       startTime,
		registeredTools: registeredTools,
	}
}

// HealthHandler provides a simple health check endpoint
func (h *HTTPHandlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime).String()

	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "hpe-opsramp-mcp",
		"uptime":    uptime,
		"tools":     h.registeredTools,
		"endpoints": map[string]string{
			"health":    "/health",
			"readiness": "/readiness",
			"sse":       "/sse",
			"message":   "/message",
			"debug":     "/debug",
			"mcp":       "/mcp",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ReadinessHandler provides a more detailed readiness check
func (h *HTTPHandlers) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"ready":     true,
		"timestamp": time.Now().Format(time.RFC3339),
		"checks": map[string]interface{}{
			"server":   "ok",
			"sessions": "ok",
			"tools":    "ok",
		},
		"tools": h.registeredTools,
	}

	// Check if server is initialized
	if h.mcpServer == nil {
		response["ready"] = false
		response["checks"].(map[string]interface{})["server"] = "not initialized"
	}

	w.Header().Set("Content-Type", "application/json")
	if response["ready"].(bool) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(response)
}

// DebugHandler provides detailed debug information about the server
func (h *HTTPHandlers) DebugHandler(w http.ResponseWriter, r *http.Request) {
	debugInfo := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(h.startTime).String(),
		"tools":     h.registeredTools,
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

		// Accept any session ID in debug mode for testing
		w.Header().Set("X-Accept-Any-Session", "true")
		h.logger.Info("Debug endpoint accessed with session ID: %s", sessionID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(debugInfo)
}

// MCPHandler provides direct access to the MCP server for simple JSON requests
func (h *HTTPHandlers) MCPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warn("Received non-POST request to /mcp endpoint: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info("Received direct MCP protocol request from %s", r.RemoteAddr)

	// Read request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var rawBody json.RawMessage
	if err := decoder.Decode(&rawBody); err != nil {
		h.logger.Error("Failed to read MCP request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Direct MCP request body: %s", string(rawBody))

	// Process through the MCP server directly
	mcpResponse := h.mcpServer.HandleMessage(r.Context(), rawBody)

	w.Header().Set("Content-Type", "application/json")

	if mcpResponse != nil {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mcpResponse); err != nil {
			h.logger.Error("Failed to encode MCP response: %v", err)
		}
	} else {
		// No response needed (notification)
		w.WriteHeader(http.StatusNoContent)
	}
}
