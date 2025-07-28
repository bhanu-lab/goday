package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigAutoCreation(t *testing.T) {
	// Create a temporary directory to simulate user home
	tempHome := t.TempDir()

	// Set temporary home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Test GetConfigPath when directory doesn't exist
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	expectedPath := filepath.Join(tempHome, ".goday", "config.yaml")
	if configPath != expectedPath {
		t.Errorf("Expected config path %s, got %s", expectedPath, configPath)
	}

	// Test LoadConfigFromDefaultPath - should create directory and config
	cfg, err := LoadConfigFromDefaultPath()
	if err != nil {
		t.Fatalf("LoadConfigFromDefaultPath failed: %v", err)
	}

	// Verify config was loaded
	if cfg == nil {
		t.Fatal("Config is nil")
	}

	// Verify directory was created
	configDir := filepath.Join(tempHome, ".goday")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Config directory was not created: %s", configDir)
	}

	// Verify config file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created: %s", configPath)
	}

	// Verify config has expected defaults
	if cfg.User.Location != "Bengaluru,IN" {
		t.Errorf("Expected default location 'Bengaluru,IN', got '%s'", cfg.User.Location)
	}

	if cfg.Widgets.Traffic.TTL != "300s" {
		t.Errorf("Expected traffic TTL '300s', got '%s'", cfg.Widgets.Traffic.TTL)
	}
}

func TestConfigPathPriority(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// Change to temp directory
	os.Chdir(tempDir)

	// Create local config.yaml
	localConfig := "config.yaml"
	err := os.WriteFile(localConfig, []byte("test: local"), 0644)
	if err != nil {
		t.Fatalf("Failed to create local config: %v", err)
	}

	// Create temp home with ~/.goday/config.yaml
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	godayDir := filepath.Join(tempHome, ".goday")
	os.MkdirAll(godayDir, 0755)
	homeConfig := filepath.Join(godayDir, "config.yaml")
	err = os.WriteFile(homeConfig, []byte("test: home"), 0644)
	if err != nil {
		t.Fatalf("Failed to create home config: %v", err)
	}

	// GetConfigPath should prefer ~/.goday/config.yaml
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	if configPath != homeConfig {
		t.Errorf("Expected home config path %s, got %s", homeConfig, configPath)
	}

	// Remove home config, should fall back to local
	os.Remove(homeConfig)
	os.Remove(godayDir)

	configPath, err = GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	// Should still return home path (for creation), not local fallback in this case
	expectedPath := filepath.Join(tempHome, ".goday", "config.yaml")
	if configPath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, configPath)
	}
}
