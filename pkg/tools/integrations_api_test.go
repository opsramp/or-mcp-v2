package tools

import (
	"context"
	"os"
	"testing"

	"github.com/vobbilis/codegen/or-mcp-v2/common"
	"gopkg.in/yaml.v2"
)

type opsrampConfig struct {
	Opsramp struct {
		TenantURL  string `yaml:"tenant_url"`
		AuthURL    string `yaml:"auth_url"`
		AuthKey    string `yaml:"auth_key"`
		AuthSecret string `yaml:"auth_secret"`
		TenantID   string `yaml:"tenant_id"`
	} `yaml:"opsramp"`
}

func loadOpsrampConfig() (*opsrampConfig, error) {
	f, err := os.Open("../../config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg opsrampConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func TestIntegrationsAPI_List_Create_Get(t *testing.T) {
	cfg, err := loadOpsrampConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Create OpsRamp config
	opsRampConfig := &common.OpsRampConfig{
		TenantURL:  cfg.Opsramp.TenantURL,
		AuthURL:    cfg.Opsramp.AuthURL,
		AuthKey:    cfg.Opsramp.AuthKey,
		AuthSecret: cfg.Opsramp.AuthSecret,
		TenantID:   cfg.Opsramp.TenantID,
	}

	// Create API client
	api, err := NewOpsRampIntegrationsAPI(opsRampConfig)
	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	ctx := context.Background()

	// List
	_, err = api.List(ctx)
	if err != nil {
		t.Errorf("List failed: %v", err)
	}

	// Create
	createCfg := map[string]interface{}{"name": "api-test-integration", "type": "api"}
	created, err := api.Create(ctx, createCfg)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("Expected created integration with ID")
	}

	// Get
	got, err := api.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Errorf("Get returned wrong integration: got %v, want %v", got.ID, created.ID)
	}
}
