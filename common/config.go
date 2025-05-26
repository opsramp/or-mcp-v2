package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	OpsRamp OpsRampConfig `yaml:"opsramp"`
}

// OpsRampConfig holds the OpsRamp API configuration
type OpsRampConfig struct {
	TenantURL  string          `yaml:"tenant_url"`
	AuthURL    string          `yaml:"auth_url"`
	AuthKey    string          `yaml:"auth_key"`
	AuthSecret string          `yaml:"auth_secret"`
	TenantID   string          `yaml:"tenant_id"`
	Resources  ResourcesConfig `yaml:"resources"`
}

// ResourcesConfig holds resource management specific configuration
type ResourcesConfig struct {
	DefaultPageSize int  `yaml:"default_page_size"`
	MaxPageSize     int  `yaml:"max_page_size"`
	CacheTTL        int  `yaml:"cache_ttl"`
	EnableBulkOps   bool `yaml:"enable_bulk_operations"`
	MaxBulkSize     int  `yaml:"max_bulk_size"`
	RequestTimeout  int  `yaml:"request_timeout"`
	RetryAttempts   int  `yaml:"retry_attempts"`
	RetryDelay      int  `yaml:"retry_delay"`
	EnableMetrics   bool `yaml:"enable_metrics"`
	MetricsInterval int  `yaml:"metrics_interval"`
}

// LoadConfig loads configuration from environment or file
func LoadConfig(configPath string) (*Config, error) {
	// First try to load from file
	config, err := loadConfigFromFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from file: %w", err)
	}

	// Override with environment variables if they exist
	overrideConfigFromEnv(config)

	return config, nil
}

// loadConfigFromFile loads configuration from a YAML file
func loadConfigFromFile(configPath string) (*Config, error) {
	if configPath == "" {
		// Try default locations
		possiblePaths := []string{
			"config.yaml",
			"config.yml",
			filepath.Join("config", "config.yaml"),
			filepath.Join("config", "config.yml"),
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}

		if configPath == "" {
			return nil, fmt.Errorf("config file not found in default locations")
		}
	}

	// Sanitize config path for security
	cleanConfigPath := filepath.Clean(configPath)
	if strings.Contains(cleanConfigPath, "..") {
		return nil, fmt.Errorf("invalid config path: %s", configPath)
	}

	data, err := os.ReadFile(cleanConfigPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Apply defaults and validate
	applyResourceDefaults(&config.OpsRamp.Resources)
	if err := validateResourceConfig(&config.OpsRamp.Resources); err != nil {
		return nil, fmt.Errorf("resource configuration validation failed: %w", err)
	}

	return &config, nil
}

// overrideConfigFromEnv overrides configuration with environment variables
func overrideConfigFromEnv(config *Config) {
	// OpsRamp config
	if val := os.Getenv("OPSRAMP_TENANT_URL"); val != "" {
		config.OpsRamp.TenantURL = val
	}
	if val := os.Getenv("OPSRAMP_AUTH_URL"); val != "" {
		config.OpsRamp.AuthURL = val
	}
	if val := os.Getenv("OPSRAMP_AUTH_KEY"); val != "" {
		config.OpsRamp.AuthKey = val
	}
	if val := os.Getenv("OPSRAMP_AUTH_SECRET"); val != "" {
		config.OpsRamp.AuthSecret = val
	}
	if val := os.Getenv("OPSRAMP_TENANT_ID"); val != "" {
		config.OpsRamp.TenantID = val
	}
}

// GetEnvOrDefault gets an environment variable or returns a default value
func GetEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && strings.TrimSpace(value) != "" {
		return value
	}
	return defaultValue
}

// applyResourceDefaults applies default values to resource configuration
func applyResourceDefaults(config *ResourcesConfig) {
	if config.DefaultPageSize == 0 {
		config.DefaultPageSize = 50
	}
	if config.MaxPageSize == 0 {
		config.MaxPageSize = 1000
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 300 // 5 minutes
	}
	if config.MaxBulkSize == 0 {
		config.MaxBulkSize = 100
	}
	if config.RequestTimeout == 0 {
		config.RequestTimeout = 30
	}
	if config.RetryAttempts == 0 {
		config.RetryAttempts = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 1000 // 1 second
	}
	if config.MetricsInterval == 0 {
		config.MetricsInterval = 60
	}
}

// validateResourceConfig validates resource configuration values
func validateResourceConfig(config *ResourcesConfig) error {
	if config.DefaultPageSize < 1 || config.DefaultPageSize > config.MaxPageSize {
		return fmt.Errorf("default_page_size must be between 1 and %d", config.MaxPageSize)
	}

	if config.MaxPageSize < 1 || config.MaxPageSize > 10000 {
		return fmt.Errorf("max_page_size must be between 1 and 10000")
	}

	if config.CacheTTL < 0 || config.CacheTTL > 3600 {
		return fmt.Errorf("cache_ttl must be between 0 and 3600 seconds")
	}

	if config.MaxBulkSize < 1 || config.MaxBulkSize > 1000 {
		return fmt.Errorf("max_bulk_size must be between 1 and 1000")
	}

	if config.RequestTimeout < 1 || config.RequestTimeout > 300 {
		return fmt.Errorf("request_timeout must be between 1 and 300 seconds")
	}

	if config.RetryAttempts < 0 || config.RetryAttempts > 10 {
		return fmt.Errorf("retry_attempts must be between 0 and 10")
	}

	if config.RetryDelay < 100 || config.RetryDelay > 10000 {
		return fmt.Errorf("retry_delay must be between 100 and 10000 milliseconds")
	}

	if config.MetricsInterval < 10 || config.MetricsInterval > 3600 {
		return fmt.Errorf("metrics_interval must be between 10 and 3600 seconds")
	}

	return nil
}
