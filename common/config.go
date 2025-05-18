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
	TenantURL  string `yaml:"tenant_url"`
	AuthURL    string `yaml:"auth_url"`
	AuthKey    string `yaml:"auth_key"`
	AuthSecret string `yaml:"auth_secret"`
	TenantID   string `yaml:"tenant_id"`
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

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
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
