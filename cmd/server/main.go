package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
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

	// Perform startup health check
	if err := startupHealthCheck(); err != nil {
		customLogger.Warn("Startup health check failed: %v", err)
		customLogger.Info("Continuing with server startup despite health check failure")
	}

	// Create SSE server with appropriate options for MCP
	var options []server.SSEOption

	// In debug mode, log more information
	if debugMode {
		customLogger.Info("Debug mode enabled, logging additional information")
	}

	// Create SSE server with tools
	sseServer = server.NewSSEServer(mcpServer, options...)
	customLogger.Debug("SSE server created with %d tools", len(registeredTools))
	customLogger.Debug("Registered tools: %v", registeredTools)

	// Create an HTTP mux to handle health, readiness, debug, and SSE endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/readiness", readinessHandler)
	mux.HandleFunc("/debug", debugHandler)
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
