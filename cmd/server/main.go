package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/vobbilis/codegen/or-mcp-v2/common"
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

	json.NewEncoder(w).Encode(response)
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

	json.NewEncoder(w).Encode(response)
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
	json.NewEncoder(w).Encode(debugInfo)
}

func main() {
	// Record start time
	startTime = time.Now()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(LogDir, 0755); err != nil {
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

	// Log registered tools
	customLogger.Info("Registered tool: %s", integrationsTool.Name)

	// Create an HTTP server with the MCP server and a health check endpoint
	portString := fmt.Sprintf(":%d", port)
	customLogger.Info("Server listening on %s", portString)

	// Create SSE server with appropriate options for MCP
	var options []server.SSEOption

	// In debug mode, accept any session ID
	if debugMode {
		// Currently, we can't access the options API to modify session validation
		// This requires modifying the library itself
		customLogger.Info("Debug mode enabled, but note that session ID validation is controlled by the SSE server library")
	}

	sseServer = server.NewSSEServer(mcpServer, options...)

	// Set up HTTP mux for multiple endpoints
	mux := http.NewServeMux()

	// Add status endpoints
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/readiness", readinessHandler)
	mux.HandleFunc("/debug", debugHandler)

	// Add SSE server handler for all other paths
	mux.Handle("/", sseServer)

	// Start the HTTP server with our mux
	if err := http.ListenAndServe(portString, mux); err != nil {
		customLogger.Fatal("Failed to start server: %v", err)
	}
}
