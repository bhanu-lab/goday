package main

import (
	"fmt"
	"testing"
)

func TestTrafficPlugin(t *testing.T) {
	plugin := NewGoogleMapsTrafficPlugin()

	// Test basic initialization
	config := map[string]interface{}{
		"api_key":     "test_key",
		"origin":      "Electronic City, Bengaluru",
		"destination": "Whitefield, Bengaluru",
	}

	err := plugin.Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	// Test metadata
	metadata := plugin.GetMetadata()
	if metadata.Name != "Google Maps Traffic" {
		t.Errorf("Expected plugin name 'Google Maps Traffic', got '%s'", metadata.Name)
	}

	// Test direction toggle
	if plugin.IsReversed() {
		t.Errorf("Expected initial direction to be false, got true")
	}

	plugin.ToggleDirection()
	if !plugin.IsReversed() {
		t.Errorf("Expected direction to be true after toggle, got false")
	}

	plugin.ToggleDirection()
	if plugin.IsReversed() {
		t.Errorf("Expected direction to be false after second toggle, got true")
	}

	fmt.Println("Traffic plugin tests passed!")
}

func TestUpdateTrafficWidget(t *testing.T) {
	wm := NewWidgetManager()
	wm.InitializeWidgets(nil)

	// Test with valid traffic data
	traffic := &TrafficData{
		Origin:      "Electronic City",
		Destination: "Whitefield",
		Duration:    "45 mins",
		DurationSec: 2700, // 45 minutes
		Distance:    "25.4 km",
		Status:      "OK",
		IsReversed:  false,
	}

	wm.UpdateTrafficWidget(traffic)

	widget := wm.Widgets["traffic"]
	if widget == nil {
		t.Fatal("Traffic widget not found")
	}

	if len(widget.Items) == 0 {
		t.Fatal("Traffic widget has no items")
	}

	item := widget.Items[0]
	expectedTitle := "Electronic City → Whitefield"
	if item.Title != expectedTitle {
		t.Errorf("Expected title '%s', got '%s'", expectedTitle, item.Title)
	}

	// Test with reversed direction
	traffic.IsReversed = true
	wm.UpdateTrafficWidget(traffic)
	item = widget.Items[0]
	expectedTitle = "Electronic City ← Whitefield"
	if item.Title != expectedTitle {
		t.Errorf("Expected reversed title '%s', got '%s'", expectedTitle, item.Title)
	}

	fmt.Println("Traffic widget update tests passed!")
}
