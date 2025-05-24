package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

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
)

func main() {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(LogDir, 0750); err != nil {
		log.Printf("Failed to create log directory: %v", err)
	}

	// Initialize the logger
	customLogger, err := common.InitLogger(common.DEBUG, LogDir, LogFileName)
	if err != nil {
		log.Printf("Failed to initialize logger: %v", err)
		log.Printf("Using default logger")
	} else {
		defer customLogger.Close()
		customLogger.Info("Starting OpsRamp MCP server")
		customLogger.Info("Log file: %s", filepath.Join(LogDir, LogFileName))
	}

	// Get the logger
	logger := common.GetLogger()

	// Parse command line flags
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	// Load configuration
	config, err := common.LoadConfig(*configPath)
	if err != nil {
		logger.Warn("Failed to load config: %v", err)
		logger.Info("Using default configuration")

		// Create a minimal default config
		config = &common.Config{
			OpsRamp: common.OpsRampConfig{
				TenantURL:  common.GetEnvOrDefault("OPSRAMP_TENANT_URL", "https://api.opsramp.com"),
				AuthURL:    common.GetEnvOrDefault("OPSRAMP_AUTH_URL", "https://api.opsramp.com/auth/oauth/token"),
				AuthKey:    common.GetEnvOrDefault("OPSRAMP_AUTH_KEY", ""),
				AuthSecret: common.GetEnvOrDefault("OPSRAMP_AUTH_SECRET", ""),
				TenantID:   common.GetEnvOrDefault("OPSRAMP_TENANT_ID", ""),
			},
		}
	}

	// Validate OpsRamp config
	logger.Info("Validating OpsRamp configuration...")
	if err := validateOpsRampConfig(&config.OpsRamp); err != nil {
		logger.Error("OpsRamp configuration validation failed: %v", err)
		logger.Warn("Some OpsRamp functionality may not work properly")
	} else {
		logger.Info("OpsRamp configuration is valid")
	}

	// Create OpsRamp client and store it globally
	opsRampClient := client.NewOpsRampClient(config)
	// Set the global client for use by tools
	client.SetGlobalClient(opsRampClient)

	// Test API connectivity
	logger.Info("Testing OpsRamp API connectivity...")
	if err := testApiConnectivity(opsRampClient); err != nil {
		logger.Error("OpsRamp API connectivity test failed: %v", err)
		logger.Warn("Some OpsRamp functionality may not work properly")
	} else {
		logger.Info("OpsRamp API connectivity test successful")
	}

	// Create MCP server
	s := server.NewMCPServer("or-mcp-v2", "1.0.0")

	// Register all tools in alphabetical order
	logger.Info("Registering MCP tools...")

	acctTool, acctHandler := tools.NewAccountsMcpTool()
	s.AddTool(acctTool, acctHandler)

	devTool, devHandler := tools.NewDevicesMcpTool()
	s.AddTool(devTool, devHandler)

	evtTool, evtHandler := tools.NewEventsMcpTool()
	s.AddTool(evtTool, evtHandler)

	intTool, intHandler := tools.NewIntegrationsMcpTool()
	s.AddTool(intTool, intHandler)

	jobsTool, jobsHandler := tools.NewJobsMcpTool()
	s.AddTool(jobsTool, jobsHandler)

	monTool, monHandler := tools.NewMonitoringMcpTool()
	s.AddTool(monTool, monHandler)

	polTool, polHandler := tools.NewPoliciesMcpTool()
	s.AddTool(polTool, polHandler)

	resTool, resHandler := tools.NewResourcesMcpTool()
	s.AddTool(resTool, resHandler)

	logger.Info("All tools registered successfully")

	// Start the server on stdio
	logger.Info("Starting MCP server on stdio...")
	if err := server.ServeStdio(s); err != nil {
		logger.Error("Failed to start server: %v", err)
		log.Fatalf("Failed to start server: %v", err)
	}
}

// validateOpsRampConfig validates the OpsRamp configuration
func validateOpsRampConfig(config *common.OpsRampConfig) error {
	// Check required fields
	if config.TenantURL == "" {
		return fmt.Errorf("tenant URL is missing")
	}
	if config.AuthURL == "" {
		return fmt.Errorf("auth URL is missing")
	}
	if config.AuthKey == "" {
		return fmt.Errorf("auth key is missing")
	}
	if config.AuthSecret == "" {
		return fmt.Errorf("auth secret is missing")
	}
	if config.TenantID == "" {
		return fmt.Errorf("tenant ID is missing")
	}

	// Check for placeholder values
	if strings.Contains(config.TenantURL, "your-tenant") {
		return fmt.Errorf("tenant URL contains placeholder value")
	}
	if strings.Contains(config.AuthKey, "YOUR_AUTH_KEY") || strings.Contains(config.AuthKey, "your-auth") {
		return fmt.Errorf("auth key contains placeholder value")
	}
	if strings.Contains(config.AuthSecret, "YOUR_AUTH_SECRET") || strings.Contains(config.AuthSecret, "your-secret") {
		return fmt.Errorf("auth secret contains placeholder value")
	}
	if strings.Contains(config.TenantID, "YOUR_TENANT_ID") || strings.Contains(config.TenantID, "your-tenant") {
		return fmt.Errorf("tenant ID contains placeholder value")
	}

	return nil
}

// testApiConnectivity tests connectivity to the OpsRamp API
func testApiConnectivity(client *client.OpsRampClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to make a simple API call to test connectivity
	// This will depend on what endpoints are available and don't require special permissions
	var response interface{}
	err := client.Get(ctx, "auth/ping", &response)
	if err != nil {
		// Try an alternative endpoint if the first one fails
		err = client.Get(ctx, "health", &response)
	}

	return err
}
