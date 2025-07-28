package main

import (
	"testing"
	"time"
)

func TestGoogleCalendarPlugin(t *testing.T) {
	// Test plugin creation
	plugin := NewGoogleCalendarPlugin()

	if plugin.GetID() != "google-calendar" {
		t.Errorf("Expected plugin ID 'google-calendar', got '%s'", plugin.GetID())
	}

	if plugin.GetType() != "calendar" {
		t.Errorf("Expected plugin type 'calendar', got '%s'", plugin.GetType())
	}

	// Test metadata
	metadata := plugin.GetMetadata()
	if metadata.Name != "Google Calendar" {
		t.Errorf("Expected plugin name 'Google Calendar', got '%s'", metadata.Name)
	}
}

func TestGoogleCalendarEventFormatting(t *testing.T) {
	plugin := NewGoogleCalendarPlugin()

	// Add some mock events
	now := time.Now()
	plugin.lastData = []GoogleCalendarEvent{
		{
			ID:        "1",
			Title:     "Test Meeting",
			StartTime: now.Add(30 * time.Minute),
			EndTime:   now.Add(90 * time.Minute),
			URL:       "https://calendar.google.com/event/1",
		},
		{
			ID:        "2",
			Title:     "Tomorrow Meeting",
			StartTime: now.Add(24 * time.Hour),
			EndTime:   now.Add(25 * time.Hour),
			URL:       "https://calendar.google.com/event/2",
		},
	}

	// Test formatting
	items := plugin.FormatEventsForDisplay()

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// First event should be today's event
	if items[0].Title != "Test Meeting" {
		t.Errorf("Expected first event title 'Test Meeting', got '%s'", items[0].Title)
	}

	// Should have a status indicator
	if items[0].Status == "" {
		t.Errorf("Expected first event to have status indicator")
	}
}

func TestCalendarWidgetUpdate(t *testing.T) {
	wm := NewWidgetManager()
	plugin := NewGoogleCalendarPlugin()

	// Add mock data
	now := time.Now()
	plugin.lastData = []GoogleCalendarEvent{
		{
			ID:        "1",
			Title:     "Test Event",
			StartTime: now.Add(time.Hour),
			EndTime:   now.Add(2 * time.Hour),
		},
	}

	// Update widget
	wm.UpdateCalendarWidget(plugin)

	// Check widget was created and updated
	if wm.Widgets["calendar"] == nil {
		t.Error("Calendar widget was not created")
	}

	if wm.Widgets["calendar"].Count != 1 {
		t.Errorf("Expected calendar widget count 1, got %d", wm.Widgets["calendar"].Count)
	}

	if len(wm.Widgets["calendar"].Items) != 1 {
		t.Errorf("Expected 1 calendar item, got %d", len(wm.Widgets["calendar"].Items))
	}

	if wm.Widgets["calendar"].Items[0].Title != "Test Event" {
		t.Errorf("Expected calendar item title 'Test Event', got '%s'", wm.Widgets["calendar"].Items[0].Title)
	}
}
