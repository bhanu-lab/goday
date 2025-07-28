package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	User struct {
		Name     string `yaml:"name"`
		Location string `yaml:"location"`
	} `yaml:"user"`
	UI struct {
		Layout     string `yaml:"layout"`
		MinWidth   int    `yaml:"min_width"`
		TileHeight int    `yaml:"tile_height"`
	} `yaml:"ui"`
	Widgets struct {
		Weather struct {
			TTL    string `yaml:"ttl"`
			APIKey string `yaml:"api_key"`
		} `yaml:"weather"`
		News struct {
			TTL      string   `yaml:"ttl"`
			Tags     []string `yaml:"tags"`
			Provider string   `yaml:"provider"`
		} `yaml:"news"`
		Slack struct {
			TTL string `yaml:"ttl"`
		} `yaml:"slack"`
		Confluence struct {
			TTL string `yaml:"ttl"`
		} `yaml:"confluence"`
		Jira struct {
			TTL     string `yaml:"ttl"`
			LogWork bool   `yaml:"log_work"`
		} `yaml:"jira"`
		Traffic struct {
			TTL         string      `yaml:"ttl"`
			Origin      interface{} `yaml:"origin"`      // Can be string or LocationConfig
			Destination interface{} `yaml:"destination"` // Can be string or LocationConfig
		} `yaml:"traffic"`
		Calendar struct {
			TTL             string `yaml:"ttl"`
			CredentialsFile string `yaml:"credentials_file"`
			TokenFile       string `yaml:"token_file"`
			MaxEvents       int    `yaml:"max_events"`
			DaysAhead       int    `yaml:"days_ahead"`
		} `yaml:"calendar"`
	} `yaml:"widgets"`
}

// GetConfigPath returns the path to the config file, checking multiple locations
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get user home directory: %w", err)
	}

	// Preferred location: ~/.goday/config.yaml
	configPath := filepath.Join(homeDir, ".goday", "config.yaml")

	// Check if config exists at preferred location
	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	}

	// Fallback: check current directory (for development)
	localConfig := "config.yaml"
	if _, err := os.Stat(localConfig); err == nil {
		return localConfig, nil
	}

	// Return preferred path for creation (directory will be created as needed)
	return configPath, nil
}

// LoadConfig loads configuration from the specified path
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// LoadConfigFromDefaultPath loads config from the default location
func LoadConfigFromDefaultPath() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory %s: %w", configDir, err)
		}

		// Create default config file
		if err := CreateDefaultConfig(configPath); err != nil {
			return nil, fmt.Errorf("failed to create default config at %s: %w", configPath, err)
		}

		// Inform user that config was created
		fmt.Printf("📁 Created config directory: %s\n", configDir)
		fmt.Printf("📝 Created default config: %s\n", configPath)
		fmt.Printf("💡 Edit the config file to customize your dashboard\n\n")
	}

	return LoadConfig(configPath)
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig(path string) error {
	defaultConfig := `# GoDay Dashboard Configuration
# Config location: ~/.goday/config.yaml
# Edit this file to customize your dashboard

user:
  name: "Your Name"  # Change this to your name
  location: "Bengaluru,IN"  # Your location for weather

ui:
  layout: at_a_glance
  min_width: 100
  tile_height: 7

widgets:
  weather:
    ttl: 600s  # Refresh every 10 minutes
    api_key: "YOUR_OWM_API_KEY"  # Get from openweathermap.org
  news:
    ttl: 600s
    tags: [golang, security, ai]  # Filter tech news by these tags
    provider: hn  # hn (Hacker News) or devto (Dev.to)
  slack:
    ttl: 20s
  confluence:
    ttl: 300s
  jira:
    ttl: 45s
    log_work: true
  traffic:
    ttl: 300s  # Refresh every 5 minutes
    # Option 1: Use addresses (geocoded automatically)
    origin: "Electronic City Phase 1, Bengaluru, Karnataka, India"
    destination: "Whitefield, Bengaluru, Karnataka, India"
    
    # Option 2: Use precise coordinates (uncomment to use)
    # origin:
    #   latitude: 12.8456
    #   longitude: 77.6603
    #   name: "Electronic City"
    # destination:
    #   latitude: 12.9698
    #   longitude: 77.7500
    #   name: "Whitefield"
  calendar:
    ttl: 300s  # Refresh every 5 minutes
    max_events: 10  # Maximum events to show
    days_ahead: 7   # Days ahead to fetch events
    # credentials_file: ~/.goday/google_calendar_credentials.json  # Will be set automatically
    # token_file: ~/.goday/google_calendar_token.json             # Will be set automatically

# Calendar Setup:
# 1. Go to https://console.cloud.google.com/
# 2. Create/select a project and enable Google Calendar API
# 3. Create OAuth 2.0 credentials (Desktop application)
# 4. Download JSON and save as ~/.goday/google_calendar_credentials.json
# 5. Restart GoDay and follow OAuth flow

# For more configuration examples, see:
# - ADDRESS_CONFIGURATION_GUIDE.md (address formats)
# - COORDINATE_EXAMPLES.md (coordinate examples)
# - CONFIG_GUIDE.md (configuration guide)
`

	return os.WriteFile(path, []byte(defaultConfig), 0644)
}
