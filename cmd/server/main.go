package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/opsramp/or-mcp-v2/common"
	"github.com/opsramp/or-mcp-v2/pkg/handlers"
	"github.com/opsramp/or-mcp-v2/pkg/mcp"
	"github.com/opsramp/or-mcp-v2/pkg/tools"
)

const (
	// LogDir is the directory where logs will be stored
	LogDir = "output/logs"
	// LogFileName is the name of the log file
	LogFileName = "or-mcp.log"
	// DefaultPort is the default port to listen on
	DefaultPort = 8080
)

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port      int
	DebugMode bool
	Logger    *common.CustomLogger
	StartTime time.Time
}

// MCPServerComponents holds all MCP server components
type MCPServerComponents struct {
	MCPServer        *server.MCPServer
	SSEServer        *server.SSEServer
	InspectorHandler *mcp.InspectorHandler
	HTTPHandlers     *handlers.HTTPHandlers
	RegisteredTools  []string
}

func main() {
	// Initialize server configuration
	config, err := initializeServerConfig()
	if err != nil {
		fmt.Printf("Failed to initialize server config: %v\n", err)
		os.Exit(1)
	}
	defer config.Logger.Close()

	// Create MCP server components
	components, err := createMCPServerComponents(config)
	if err != nil {
		config.Logger.Fatal("Failed to create MCP server components: %v", err)
	}

	// Perform startup health check
	if err := performStartupHealthCheck(config.Logger); err != nil {
		config.Logger.Warn("Startup health check failed: %v", err)
		config.Logger.Info("Continuing with server startup despite health check failure")
	}

	// Start the HTTP server
	httpServer := createHTTPServer(config, components)
	startServer(config, httpServer)
}

// initializeServerConfig initializes the server configuration
func initializeServerConfig() (*ServerConfig, error) {
	startTime := time.Now()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(LogDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Initialize the logger
	logger, err := common.InitLogger(common.DEBUG, LogDir, LogFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Log server startup
	logger.Info("Starting HPE OpsRamp MCP server")
	logger.Info("Log file: %s", filepath.Join(LogDir, LogFileName))

	// Determine port from environment variable
	port := DefaultPort
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		if p, err := strconv.Atoi(portEnv); err == nil && p > 0 && p < 65536 {
			port = p
			logger.Info("Using port from environment: %d", port)
		} else {
			logger.Warn("Invalid PORT environment variable: %s, using default: %d", portEnv, port)
		}
	}

	// Check if debug mode is enabled
	debugMode := os.Getenv("DEBUG") == "true"
	if debugMode {
		logger.Info("*** DEBUG MODE ENABLED ***")
	}

	return &ServerConfig{
		Port:      port,
		DebugMode: debugMode,
		Logger:    logger,
		StartTime: startTime,
	}, nil
}

// createMCPServerComponents creates all MCP server components
func createMCPServerComponents(config *ServerConfig) (*MCPServerComponents, error) {
	// Create MCP server
	mcpServer := server.NewMCPServer("HPE OpsRamp MCP", "1.0.0")

	// Register tools
	registeredTools := make([]string, 0)

	// Register integrations tool
	integrationsTool, integrationsHandler := tools.NewIntegrationsMcpTool()
	mcpServer.AddTool(integrationsTool, integrationsHandler)
	registeredTools = append(registeredTools, integrationsTool.Name)
	config.Logger.Info("Registered tool: %s", integrationsTool.Name)

	// Register resources tool
	resourcesTool, resourcesHandler := tools.NewResourcesMcpTool()
	mcpServer.AddTool(resourcesTool, resourcesHandler)
	registeredTools = append(registeredTools, resourcesTool.Name)
	config.Logger.Info("Registered tool: %s", resourcesTool.Name)

	// Create SSE server with appropriate options for MCP
	sseOptions := []server.SSEOption{
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(30 * time.Second),
		server.WithMessageEndpoint("/mcp-message"),
		server.WithSSEEndpoint("/sse"),
		server.WithUseFullURLForMessageEndpoint(true),
		server.WithAppendQueryToMessageEndpoint(),
	}

	sseServer := server.NewSSEServer(mcpServer, sseOptions...)
	config.Logger.Debug("SSE server created with %d tools", len(registeredTools))
	config.Logger.Debug("Registered tools: %v", registeredTools)

	// Create MCP Inspector compatibility handler
	inspectorHandler := mcp.NewInspectorHandler(mcpServer, config.Logger)

	// Create HTTP handlers
	httpHandlers := handlers.NewHTTPHandlers(mcpServer, sseServer, config.Logger, config.StartTime, registeredTools)

	return &MCPServerComponents{
		MCPServer:        mcpServer,
		SSEServer:        sseServer,
		InspectorHandler: inspectorHandler,
		HTTPHandlers:     httpHandlers,
		RegisteredTools:  registeredTools,
	}, nil
}

// performStartupHealthCheck performs a real API call to verify connectivity
func performStartupHealthCheck(logger *common.CustomLogger) error {
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
	logger.Info("Startup health check passed: successfully listed %d integrations", len(integrations))
	return nil
}

// createHTTPServer creates and configures the HTTP server
func createHTTPServer(config *ServerConfig, components *MCPServerComponents) *http.Server {
	// Create HTTP mux to handle all endpoints
	mux := http.NewServeMux()

	// Register standard HTTP endpoints
	mux.HandleFunc("/health", components.HTTPHandlers.HealthHandler)
	mux.HandleFunc("/readiness", components.HTTPHandlers.ReadinessHandler)
	mux.HandleFunc("/debug", components.HTTPHandlers.DebugHandler)
	mux.HandleFunc("/mcp", components.HTTPHandlers.MCPHandler)

	// Register SSE endpoint (native MCP-Go implementation)
	mux.Handle("/sse", components.SSEServer)

	// Register native MCP-Go message endpoint (used by SSE server)
	mux.Handle("/mcp-message", components.SSEServer.MessageHandler())

	// Register MCP Inspector compatibility endpoint (for direct connections)
	mux.HandleFunc("/message", components.InspectorHandler.HandleMessage)

	config.Logger.Debug("HTTP routes configured")

	// Create HTTP server
	portString := fmt.Sprintf(":%d", config.Port)
	config.Logger.Info("Server listening on %s", portString)

	return &http.Server{
		Addr:    portString,
		Handler: mux,
		// Increase timeouts for long-running operations
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  240 * time.Second,
	}
}

// startServer starts the HTTP server and handles graceful shutdown
func startServer(config *ServerConfig, httpServer *http.Server) {
	// Start the server in a goroutine
	go func() {
		config.Logger.Info("Starting HTTP server on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			config.Logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	config.Logger.Info("Shutting down server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		config.Logger.Fatal("Server forced to shutdown: %v", err)
	}

	config.Logger.Info("Server exited gracefully")
}
