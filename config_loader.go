package main

import (
	"os"

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
	} `yaml:"widgets"`
}

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
