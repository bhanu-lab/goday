package main

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file for testing
	configContent := `
user:
  name: "Test User"
  location: "Test City"

ui:
  layout: at_a_glance
  min_width: 100
  tile_height: 7

widgets:
  weather:
    ttl: 600s
    api_key: "test_key"
  news:
    ttl: 600s
    tags: [golang, test]
    provider: hn
  slack:
    ttl: 20s
  confluence:
    ttl: 300s
  jira:
    ttl: 45s
    log_work: true
`

	tmpFile, err := os.CreateTemp("", "test_config_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config content: %v", err)
	}
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.User.Name != "Test User" {
		t.Errorf("Expected user name 'Test User', got '%s'", cfg.User.Name)
	}

	if cfg.User.Location != "Test City" {
		t.Errorf("Expected location 'Test City', got '%s'", cfg.User.Location)
	}

	if len(cfg.Widgets.News.Tags) != 2 {
		t.Errorf("Expected 2 news tags, got %d", len(cfg.Widgets.News.Tags))
	}
}

func TestWeatherProvider(t *testing.T) {
	provider := NewWeatherProvider("test_key", "TestCity")

	// Test with invalid API key (should return fallback data)
	data, err := provider.Fetch()
	if err != nil {
		t.Fatalf("WeatherProvider.Fetch() failed: %v", err)
	}

	if data.Temperature != 30 {
		t.Errorf("Expected temperature 30, got %d", data.Temperature)
	}

	if data.Icon != "☁" {
		t.Errorf("Expected icon '☁', got '%s'", data.Icon)
	}
}

func TestNewsProvider(t *testing.T) {
	tags := []string{"golang", "test"}
	provider := NewNewsProvider(tags)

	// Test tag matching
	if !provider.matchesTags("golang programming language") {
		t.Error("Expected 'golang programming language' to match 'golang' tag")
	}

	if !provider.matchesTags("testing framework") {
		t.Error("Expected 'testing framework' to match 'test' tag")
	}

	if provider.matchesTags("Python programming") {
		t.Error("Expected 'Python programming' to not match any tags")
	}

	// Test with empty tags (should match everything)
	emptyProvider := NewNewsProvider([]string{})
	if !emptyProvider.matchesTags("Any title") {
		t.Error("Expected empty tags to match any title")
	}
}

func TestWidgetManager(t *testing.T) {
	wm := NewWidgetManager()

	// Test initialization
	if len(wm.Widgets) != 0 {
		t.Errorf("Expected empty widgets map, got %d widgets", len(wm.Widgets))
	}

	// Test tag cycling
	wm.NewsTags = []string{"golang", "security"}

	// Initial state should be "All"
	if wm.GetCurrentNewsTag() != "All" {
		t.Errorf("Expected initial tag 'All', got '%s'", wm.GetCurrentNewsTag())
	}

	// Cycle to first tag
	wm.CycleNewsTag()
	if wm.GetCurrentNewsTag() != "golang" {
		t.Errorf("Expected tag 'golang', got '%s'", wm.GetCurrentNewsTag())
	}

	// Cycle to second tag
	wm.CycleNewsTag()
	if wm.GetCurrentNewsTag() != "security" {
		t.Errorf("Expected tag 'security', got '%s'", wm.GetCurrentNewsTag())
	}

	// Cycle back to "All"
	wm.CycleNewsTag()
	if wm.GetCurrentNewsTag() != "All" {
		t.Errorf("Expected tag 'All', got '%s'", wm.GetCurrentNewsTag())
	}
}

func TestScheduler(t *testing.T) {
	scheduler := NewScheduler()

	// Test adding tasks
	scheduler.AddTask("test1", 10*time.Second, nil)
	scheduler.AddTask("test2", 5*time.Second, nil)

	if len(scheduler.GetTasks()) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(scheduler.GetTasks()))
	}

	// Test getting next task (should be the one with shorter interval)
	next := scheduler.GetNextTask()
	if next == nil {
		t.Fatal("Expected next task, got nil")
	}

	if next.ID != "test2" {
		t.Errorf("Expected next task ID 'test2', got '%s'", next.ID)
	}

	// Test updating task
	scheduler.UpdateTask("test1")

	// Test removing task
	scheduler.RemoveTask("test1")
	if len(scheduler.GetTasks()) != 1 {
		t.Errorf("Expected 1 task after removal, got %d", len(scheduler.GetTasks()))
	}
}

func TestParseTTL(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"600s", 600 * time.Second},
		{"20s", 20 * time.Second},
		{"", 600 * time.Second},        // Default
		{"invalid", 600 * time.Second}, // Default on error
	}

	for _, test := range tests {
		result := ParseTTL(test.input)
		if result != test.expected {
			t.Errorf("ParseTTL('%s') = %v, expected %v", test.input, result, test.expected)
		}
	}
}
