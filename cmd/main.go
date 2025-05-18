package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	if err := os.MkdirAll(LogDir, 0755); err != nil {
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

	// Parse command line flags
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	// Load configuration
	config, err := common.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
		fmt.Printf("Using default configuration\n")

		// Get the logger
		logger := common.GetLogger()
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

	// Create OpsRamp client and store it globally
	opsRampClient := client.NewOpsRampClient(config)
	// Set the global client for use by tools
	client.SetGlobalClient(opsRampClient)

	// Create MCP server
	s := server.NewMCPServer("or-mcp-v2", "1.0.0")

	// Register the integrations tool as an MCP tool
	// Register all tools in alphabetical order
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

	// Register the resources tool
	resTool, resHandler := tools.NewResourcesMcpTool()
	s.AddTool(resTool, resHandler)

	// Start the server on stdio
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
