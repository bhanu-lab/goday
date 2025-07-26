package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WeatherPlugin implements weather fetching from OpenWeatherMap
type WeatherPlugin struct {
	id          string
	pluginType  string
	name        string
	version     string
	description string
	author      string
	apiKey      string
	city        string
	client      *http.Client
	lastData    *WeatherData
}

// NewWeatherPlugin creates a new weather plugin
func NewWeatherPlugin(apiKey, city string) *WeatherPlugin {
	return &WeatherPlugin{
		id:          "openweathermap",
		pluginType:  "weather",
		name:        "OpenWeatherMap",
		version:     "1.0.0",
		description: "Fetches weather data from OpenWeatherMap API",
		author:      "GoDay Team",
		apiKey:      apiKey,
		city:        city,
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

// GetID returns the plugin ID
func (wp *WeatherPlugin) GetID() string {
	return wp.id
}

// GetType returns the plugin type
func (wp *WeatherPlugin) GetType() string {
	return wp.pluginType
}

// Initialize sets up the plugin with configuration
func (wp *WeatherPlugin) Initialize(config map[string]interface{}) error {
	if apiKey, ok := config["api_key"].(string); ok {
		wp.apiKey = apiKey
	}
	if city, ok := config["city"].(string); ok {
		wp.city = city
	}
	return nil
}

// Fetch retrieves weather data
func (wp *WeatherPlugin) Fetch(ctx context.Context) (interface{}, error) {
	if wp.apiKey == "" || wp.apiKey == "YOUR_OWM_API_KEY" {
		// Return mock data for demo
		return &WeatherData{
			Temperature: 30,
			Condition:   "Clouds",
			Icon:        "☁",
		}, nil
	}

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", wp.city, wp.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return wp.lastData, err
	}

	resp, err := wp.client.Do(req)
	if err != nil {
		return wp.lastData, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return wp.lastData, err
	}

	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return wp.lastData, err
	}

	// Return fallback data if the response is invalid
	if weatherResp.Main.Temp == 0 {
		return &WeatherData{
			Temperature: 30,
			Condition:   "Clouds",
			Icon:        "☁",
		}, nil
	}

	icon := "☁"
	condition := "Clouds"
	if len(weatherResp.Weather) > 0 {
		icon = getWeatherIcon(weatherResp.Weather[0].ID)
		condition = weatherResp.Weather[0].Main
	}

	data := &WeatherData{
		Temperature: int(weatherResp.Main.Temp),
		Condition:   condition,
		Icon:        icon,
	}
	wp.lastData = data
	return data, nil
}

// GetMetadata returns plugin metadata
func (wp *WeatherPlugin) GetMetadata() PluginMetadata {
	return PluginMetadata{
		Name:        wp.name,
		Version:     wp.version,
		Description: wp.description,
		Author:      wp.author,
		Type:        wp.pluginType,
		Config: map[string]string{
			"api_key": wp.apiKey,
			"city":    wp.city,
		},
	}
}

// Cleanup performs cleanup
func (wp *WeatherPlugin) Cleanup() error {
	return nil
}
