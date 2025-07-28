package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TrafficData represents traffic information between two locations
type TrafficData struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Duration    string `json:"duration"`
	DurationSec int    `json:"duration_seconds"`
	Distance    string `json:"distance"`
	Status      string `json:"status"`
	IsReversed  bool   `json:"is_reversed"`
}

// GoogleMapsTrafficPlugin implements the Plugin interface for Google Maps traffic data
type GoogleMapsTrafficPlugin struct {
	id          string
	apiKey      string
	origin      string
	destination string
	isReversed  bool
	client      *http.Client
}

// NewGoogleMapsTrafficPlugin creates a new Google Maps traffic plugin
func NewGoogleMapsTrafficPlugin() *GoogleMapsTrafficPlugin {
	return &GoogleMapsTrafficPlugin{
		id:     "googlemaps_traffic",
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetID returns the plugin ID
func (g *GoogleMapsTrafficPlugin) GetID() string {
	return g.id
}

// GetType returns the plugin type
func (g *GoogleMapsTrafficPlugin) GetType() string {
	return "traffic"
}

// Initialize sets up the plugin with configuration
func (g *GoogleMapsTrafficPlugin) Initialize(config map[string]interface{}) error {
	if apiKey, ok := config["api_key"].(string); ok {
		g.apiKey = apiKey
	} else {
		return fmt.Errorf("missing or invalid api_key in config")
	}

	if origin, ok := config["origin"].(string); ok {
		g.origin = origin
	} else {
		return fmt.Errorf("missing or invalid origin in config")
	}

	if destination, ok := config["destination"].(string); ok {
		g.destination = destination
	} else {
		return fmt.Errorf("missing or invalid destination in config")
	}

	g.isReversed = false
	return nil
}

// Google Maps Distance Matrix API response structure
type DistanceMatrixResponse struct {
	Status string `json:"status"`
	Rows   []struct {
		Elements []struct {
			Status   string `json:"status"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			DurationInTraffic struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration_in_traffic"`
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
		} `json:"elements"`
	} `json:"rows"`
	OriginAddresses      []string `json:"origin_addresses"`
	DestinationAddresses []string `json:"destination_addresses"`
}

// Fetch retrieves traffic data from Google Maps
func (g *GoogleMapsTrafficPlugin) Fetch(ctx context.Context) (interface{}, error) {
	if g.apiKey == "" || g.apiKey == "YOUR_GOOGLE_MAPS_API_KEY" {
		return nil, fmt.Errorf("Google Maps API key not configured. Please update config.yaml with a valid API key")
	}

	// Determine origin and destination based on current direction
	origin := g.origin
	destination := g.destination
	if g.isReversed {
		origin = g.destination
		destination = g.origin
	}

	// Build API URL
	baseURL := "https://maps.googleapis.com/maps/api/distancematrix/json"
	params := url.Values{}
	params.Add("origins", origin)
	params.Add("destinations", destination)
	params.Add("departure_time", "now")
	params.Add("traffic_model", "best_guess")
	params.Add("key", g.apiKey)

	apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Debug: log the request (without API key for security)
	fmt.Printf("Traffic API Request: %s → %s\n", origin, destination)

	// Make API request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse response
	var distanceResp DistanceMatrixResponse
	if err := json.NewDecoder(resp.Body).Decode(&distanceResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if distanceResp.Status != "OK" {
		switch distanceResp.Status {
		case "REQUEST_DENIED":
			return nil, fmt.Errorf("API key invalid or restrictions applied. Check: 1) API key is valid, 2) Distance Matrix API is enabled, 3) No IP/referrer restrictions")
		case "OVER_DAILY_LIMIT":
			return nil, fmt.Errorf("API quota exceeded for today")
		case "OVER_QUERY_LIMIT":
			return nil, fmt.Errorf("API rate limit exceeded")
		case "INVALID_REQUEST":
			return nil, fmt.Errorf("invalid request parameters")
		default:
			return nil, fmt.Errorf("API error: %s", distanceResp.Status)
		}
	}

	if len(distanceResp.Rows) == 0 || len(distanceResp.Rows[0].Elements) == 0 {
		return nil, fmt.Errorf("no route data available")
	}

	element := distanceResp.Rows[0].Elements[0]
	if element.Status != "OK" {
		return nil, fmt.Errorf("route error: %s", element.Status)
	}

	// Use traffic duration if available, otherwise use normal duration
	duration := element.Duration.Text
	durationSec := element.Duration.Value
	if element.DurationInTraffic.Text != "" {
		duration = element.DurationInTraffic.Text
		durationSec = element.DurationInTraffic.Value
	}

	// Get readable location names
	originName := g.getLocationShortName(origin)
	destName := g.getLocationShortName(destination)

	return &TrafficData{
		Origin:      originName,
		Destination: destName,
		Duration:    duration,
		DurationSec: durationSec,
		Distance:    element.Distance.Text,
		Status:      "OK",
		IsReversed:  g.isReversed,
	}, nil
}

// getLocationShortName extracts a readable short name from full address
func (g *GoogleMapsTrafficPlugin) getLocationShortName(address string) string {
	// Extract area name from full address (e.g., "Electronic City" from "Electronic City, Bengaluru, Karnataka, India")
	parts := strings.Split(address, ",")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return address
}

// GetMetadata returns plugin metadata
func (g *GoogleMapsTrafficPlugin) GetMetadata() PluginMetadata {
	return PluginMetadata{
		Name:        "Google Maps Traffic",
		Version:     "1.0.0",
		Description: "Provides real-time traffic information between two locations using Google Maps",
		Author:      "GoDay",
		Type:        "traffic",
		Config: map[string]string{
			"api_key":     "Google Maps API key",
			"origin":      "Starting location",
			"destination": "Destination location",
		},
	}
}

// Cleanup performs any necessary cleanup
func (g *GoogleMapsTrafficPlugin) Cleanup() error {
	return nil
}

// ToggleDirection switches between origin->destination and destination->origin
func (g *GoogleMapsTrafficPlugin) ToggleDirection() {
	g.isReversed = !g.isReversed
}

// IsReversed returns whether the direction is currently reversed
func (g *GoogleMapsTrafficPlugin) IsReversed() bool {
	return g.isReversed
}
