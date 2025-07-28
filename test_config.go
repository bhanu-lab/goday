package main

import (
	"fmt"
	"log"
	"os"
)

func testConfig() {
	// Test the config path detection
	configPath, err := GetConfigPath()
	if err != nil {
		log.Fatalf("Error getting config path: %v", err)
	}

	fmt.Printf("Config path: %s\n", configPath)

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Config file does not exist - would be created on first run")
	} else {
		fmt.Println("Config file exists")
	}

	// Test loading config (this will create it if it doesn't exist)
	cfg, err := LoadConfigFromDefaultPath()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Printf("Config loaded successfully!\n")
	fmt.Printf("User name: %s\n", cfg.User.Name)
	fmt.Printf("User location: %s\n", cfg.User.Location)
	fmt.Printf("Traffic TTL: %s\n", cfg.Widgets.Traffic.TTL)
}
