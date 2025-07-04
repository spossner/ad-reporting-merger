package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg == nil {
		t.Fatal("Config is nil")
	}

	groups := cfg.GetGroups()
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}

	// Check first group
	if groups[0].Prefix != "AdManager Reporting" {
		t.Errorf("Expected first group prefix to be 'AdManager Reporting', got '%s'", groups[0].Prefix)
	}

	if groups[0].Output != "raw.csv" {
		t.Errorf("Expected first group output to be 'raw.csv', got '%s'", groups[0].Output)
	}

	// Check second group
	if groups[1].Prefix != "Revenue per AdUnit" {
		t.Errorf("Expected second group prefix to be 'Revenue per AdUnit', got '%s'", groups[1].Prefix)
	}

	if groups[1].Output != "raw-revenue.csv" {
		t.Errorf("Expected second group output to be 'raw-revenue.csv', got '%s'", groups[1].Output)
	}
}

func TestGetWorkDir(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	workDir := cfg.GetWorkDir()
	if workDir != "~/Downloads" {
		t.Errorf("Expected work dir to be '~/Downloads', got '%s'", workDir)
	}
}

func TestGetWorkDirDefault(t *testing.T) {
	cfg := &Config{} // Empty config
	workDir := cfg.GetWorkDir()
	if workDir != "~/Downloads" {
		t.Errorf("Expected default work dir to be '~/Downloads', got '%s'", workDir)
	}
}