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

// BiDirectionalTrafficData represents traffic information for both directions
type BiDirectionalTrafficData struct {
	OriginToDestination TrafficData `json:"origin_to_destination"`
	DestinationToOrigin TrafficData `json:"destination_to_origin"`
	OriginName          string      `json:"origin_name"`
	DestinationName     string      `json:"destination_name"`
	Status              string      `json:"status"`
}

// LocationConfig represents either an address string or lat/lng coordinates
type LocationConfig struct {
	Address   string  `yaml:"address,omitempty"`
	Latitude  float64 `yaml:"latitude,omitempty"`
	Longitude float64 `yaml:"longitude,omitempty"`
	Name      string  `yaml:"name,omitempty"` // Optional display name
}

// OSRMTrafficPlugin implements traffic routing using OpenStreetMap data via OSRM
type OSRMTrafficPlugin struct {
	id          string
	origin      LocationConfig
	destination LocationConfig
	isReversed  bool
	client      *http.Client
}

// NewOSRMTrafficPlugin creates a new OSRM traffic plugin (no API key required)
func NewOSRMTrafficPlugin() *OSRMTrafficPlugin {
	return &OSRMTrafficPlugin{
		id:     "osrm_traffic",
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetID returns the plugin ID
func (o *OSRMTrafficPlugin) GetID() string {
	return o.id
}

// GetType returns the plugin type
func (o *OSRMTrafficPlugin) GetType() string {
	return "traffic"
}

// Initialize sets up the plugin with configuration
func (o *OSRMTrafficPlugin) Initialize(config map[string]interface{}) error {
	// Parse origin configuration
	if err := o.parseLocationConfig("origin", config, &o.origin); err != nil {
		return err
	}

	// Parse destination configuration
	if err := o.parseLocationConfig("destination", config, &o.destination); err != nil {
		return err
	}

	o.isReversed = false
	return nil
}

// parseLocationConfig parses location configuration from config map
func (o *OSRMTrafficPlugin) parseLocationConfig(key string, config map[string]interface{}, location *LocationConfig) error {
	if locationData, ok := config[key]; ok {
		switch v := locationData.(type) {
		case string:
			// Simple string address
			location.Address = v
		case map[string]interface{}:
			// Complex configuration with lat/lng or address
			if address, hasAddress := v["address"].(string); hasAddress {
				location.Address = address
			}
			if lat, hasLat := v["latitude"].(float64); hasLat {
				location.Latitude = lat
			}
			if lng, hasLng := v["longitude"].(float64); hasLng {
				location.Longitude = lng
			}
			if name, hasName := v["name"].(string); hasName {
				location.Name = name
			}

			// Validate that we have either address or lat/lng
			hasCoords := location.Latitude != 0 && location.Longitude != 0
			hasAddress := location.Address != ""
			if !hasCoords && !hasAddress {
				return fmt.Errorf("%s must have either 'address' or 'latitude'+'longitude'", key)
			}
		default:
			return fmt.Errorf("invalid %s configuration: must be string or object", key)
		}
		return nil
	}
	return fmt.Errorf("missing %s in config", key)
}

// OSRM API response structures
type OSRMResponse struct {
	Code   string `json:"code"`
	Routes []struct {
		Duration float64 `json:"duration"` // in seconds
		Distance float64 `json:"distance"` // in meters
		Legs     []struct {
			Duration float64 `json:"duration"`
			Distance float64 `json:"distance"`
		} `json:"legs"`
	} `json:"routes"`
}

type NominatimResponse []struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// geocodeLocation converts address to coordinates using Nominatim (free)
func (o *OSRMTrafficPlugin) geocodeLocation(location string) (lat, lon string, err error) {
	// Try multiple search strategies for better results
	searchQueries := []string{
		location, // Original query
		location + ", Bengaluru, Karnataka, India",                          // Add location context
		strings.Replace(location, "Pvt Ltd", "", -1) + ", Bengaluru, India", // Remove company suffixes
	}

	for i, query := range searchQueries {
		lat, lon, err := o.tryGeocoding(query)
		if err == nil {
			return lat, lon, nil
		}

		// Log the attempt for debugging
		if i == 0 {
			fmt.Printf("Geocoding attempt %d failed for '%s': %v\n", i+1, query, err)
		}
	}

	return "", "", fmt.Errorf("location not found after trying multiple queries: %s", location)
}

// tryGeocoding performs a single geocoding attempt
func (o *OSRMTrafficPlugin) tryGeocoding(location string) (lat, lon string, err error) {
	// Use Nominatim for geocoding (free OpenStreetMap service)
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", location)
	params.Add("format", "json")
	params.Add("limit", "5")                     // Get more results for better accuracy
	params.Add("addressdetails", "1")            // Get detailed address info
	params.Add("countrycodes", "in")             // Restrict to India for better results
	params.Add("bounded", "1")                   // Prefer results within viewbox
	params.Add("viewbox", "77.3,13.2,77.9,12.7") // Bengaluru bounding box

	apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("error creating geocoding request: %w", err)
	}

	// Add user agent as required by Nominatim
	req.Header.Set("User-Agent", "GoDay-Dashboard/1.0 (Contact: developer@goday.com)")

	resp, err := o.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error making geocoding request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	var geocodeResp NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&geocodeResp); err != nil {
		return "", "", fmt.Errorf("error decoding geocoding response: %w", err)
	}

	if len(geocodeResp) == 0 {
		return "", "", fmt.Errorf("no results found for: %s", location)
	}

	// Use the first (most relevant) result
	return geocodeResp[0].Lat, geocodeResp[0].Lon, nil
}

// getLocationCoordinates gets lat/lng coordinates from LocationConfig
// If coordinates are provided, uses them directly. Otherwise geocodes the address.
func (o *OSRMTrafficPlugin) getLocationCoordinates(location LocationConfig) (lat, lon string, err error) {
	// If coordinates are provided, use them directly
	if location.Latitude != 0 && location.Longitude != 0 {
		return fmt.Sprintf("%.6f", location.Latitude), fmt.Sprintf("%.6f", location.Longitude), nil
	}

	// Otherwise, geocode the address
	if location.Address != "" {
		return o.geocodeLocation(location.Address)
	}

	return "", "", fmt.Errorf("location has neither coordinates nor address")
}

// getLocationDisplayName gets a display name for the location
// Uses the custom name if provided, otherwise extracts from address, otherwise uses coordinates
func (o *OSRMTrafficPlugin) getLocationDisplayName(location LocationConfig) string {
	// Use custom name if provided
	if location.Name != "" {
		return location.Name
	}

	// Extract from address if available
	if location.Address != "" {
		return o.getLocationShortName(location.Address)
	}

	// Fall back to coordinates
	if location.Latitude != 0 && location.Longitude != 0 {
		return fmt.Sprintf("%.4f,%.4f", location.Latitude, location.Longitude)
	}

	return "Unknown Location"
}

// Fetch retrieves traffic data from OSRM for both directions
func (o *OSRMTrafficPlugin) Fetch(ctx context.Context) (interface{}, error) {
	// Get coordinates for both locations
	originLat, originLon, err := o.getLocationCoordinates(o.origin)
	if err != nil {
		return nil, fmt.Errorf("failed to get origin coordinates: %w", err)
	}

	destLat, destLon, err := o.getLocationCoordinates(o.destination)
	if err != nil {
		return nil, fmt.Errorf("failed to get destination coordinates: %w", err)
	}

	// Get route from origin to destination
	originToDestRoute, err := o.getRoute(ctx, originLon, originLat, destLon, destLat)
	if err != nil {
		return nil, fmt.Errorf("failed to get origin->destination route: %w", err)
	}

	// Get route from destination to origin
	destToOriginRoute, err := o.getRoute(ctx, destLon, destLat, originLon, originLat)
	if err != nil {
		return nil, fmt.Errorf("failed to get destination->origin route: %w", err)
	}

	// Get readable location names
	originName := o.getLocationDisplayName(o.origin)
	destName := o.getLocationDisplayName(o.destination)

	// Create traffic data for both directions
	originToDestData := TrafficData{
		Origin:      originName,
		Destination: destName,
		Duration:    o.formatDuration(int(originToDestRoute.Routes[0].Duration)),
		DurationSec: int(originToDestRoute.Routes[0].Duration),
		Distance:    fmt.Sprintf("%.1f km", originToDestRoute.Routes[0].Distance/1000),
		Status:      "OK",
		IsReversed:  false,
	}

	destToOriginData := TrafficData{
		Origin:      destName,
		Destination: originName,
		Duration:    o.formatDuration(int(destToOriginRoute.Routes[0].Duration)),
		DurationSec: int(destToOriginRoute.Routes[0].Duration),
		Distance:    fmt.Sprintf("%.1f km", destToOriginRoute.Routes[0].Distance/1000),
		Status:      "OK",
		IsReversed:  true,
	}

	return &BiDirectionalTrafficData{
		OriginToDestination: originToDestData,
		DestinationToOrigin: destToOriginData,
		OriginName:          originName,
		DestinationName:     destName,
		Status:              "OK",
	}, nil
}

// getRoute makes a single OSRM API call for a specific route
func (o *OSRMTrafficPlugin) getRoute(ctx context.Context, fromLon, fromLat, toLon, toLat string) (*OSRMResponse, error) {
	baseURL := "https://router.project-osrm.org/route/v1/driving"
	coordinates := fmt.Sprintf("%s,%s;%s,%s", fromLon, fromLat, toLon, toLat)
	apiURL := fmt.Sprintf("%s/%s?overview=false&alternatives=false&steps=false", baseURL, coordinates)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating route request: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making route request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OSRM API returned status %d", resp.StatusCode)
	}

	var osrmResp OSRMResponse
	if err := json.NewDecoder(resp.Body).Decode(&osrmResp); err != nil {
		return nil, fmt.Errorf("error decoding route response: %w", err)
	}

	if osrmResp.Code != "Ok" {
		return nil, fmt.Errorf("OSRM error: %s", osrmResp.Code)
	}

	if len(osrmResp.Routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	return &osrmResp, nil
}

// formatDuration converts seconds to readable format
func (o *OSRMTrafficPlugin) formatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%d sec", seconds)
	} else if seconds < 3600 {
		minutes := seconds / 60
		return fmt.Sprintf("%d min", minutes)
	} else {
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		if minutes == 0 {
			return fmt.Sprintf("%d hr", hours)
		}
		return fmt.Sprintf("%d hr %d min", hours, minutes)
	}
}

// getLocationShortName extracts a readable short name from full address
func (o *OSRMTrafficPlugin) getLocationShortName(address string) string {
	// Extract meaningful name from full address
	parts := strings.Split(address, ",")
	if len(parts) > 0 {
		firstPart := strings.TrimSpace(parts[0])

		// If the first part looks like a building/complex name, use it
		if len(firstPart) > 0 && !strings.Contains(strings.ToLower(firstPart), "road") &&
			!strings.Contains(strings.ToLower(firstPart), "street") &&
			!strings.Contains(strings.ToLower(firstPart), "avenue") {
			return firstPart
		}

		// If first part is a road/street, try to use a landmark from the second part
		if len(parts) > 1 {
			secondPart := strings.TrimSpace(parts[1])
			if len(secondPart) > 0 && !strings.Contains(strings.ToLower(secondPart), "bengaluru") &&
				!strings.Contains(strings.ToLower(secondPart), "karnataka") {
				return secondPart
			}
		}

		return firstPart
	}
	return address
}

// GetMetadata returns plugin metadata
func (o *OSRMTrafficPlugin) GetMetadata() PluginMetadata {
	return PluginMetadata{
		Name:        "OSRM Traffic",
		Version:     "1.0.0",
		Description: "Provides routing information using OpenStreetMap data via OSRM (no API key required)",
		Author:      "GoDay",
		Type:        "traffic",
		Config: map[string]string{
			"origin":      "Starting location",
			"destination": "Destination location",
		},
	}
}

// Cleanup performs any necessary cleanup
func (o *OSRMTrafficPlugin) Cleanup() error {
	return nil
}

// ToggleDirection switches between origin->destination and destination->origin
func (o *OSRMTrafficPlugin) ToggleDirection() {
	o.isReversed = !o.isReversed
}

// IsReversed returns whether the direction is currently reversed
func (o *OSRMTrafficPlugin) IsReversed() bool {
	return o.isReversed
}
