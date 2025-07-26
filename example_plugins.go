package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GitHubPlugin implements fetching issues from GitHub API
type GitHubPlugin struct {
	id          string
	pluginType  string
	name        string
	version     string
	description string
	author      string
	apiToken    string
	repository  string
	client      *http.Client
	lastData    []GitHubIssue
}

// GitHubIssue represents a GitHub issue
type GitHubIssue struct {
	ID      int    `json:"id"`
	Number  int    `json:"number"`
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
	State   string `json:"state"`
	User    struct {
		Login string `json:"login"`
	} `json:"user"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Labels    []struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"labels"`
}

// NewGitHubPlugin creates a new GitHub plugin
func NewGitHubPlugin(apiToken, repository string) *GitHubPlugin {
	return &GitHubPlugin{
		id:          "github-issues",
		pluginType:  "issues",
		name:        "GitHub Issues",
		version:     "1.0.0",
		description: "Fetches issues from GitHub repository",
		author:      "GoDay Team",
		apiToken:    apiToken,
		repository:  repository,
		client:      &http.Client{Timeout: 10 * time.Second},
		lastData:    []GitHubIssue{},
	}
}

// GetID returns the plugin ID
func (gp *GitHubPlugin) GetID() string {
	return gp.id
}

// GetType returns the plugin type
func (gp *GitHubPlugin) GetType() string {
	return gp.pluginType
}

// Initialize sets up the plugin with configuration
func (gp *GitHubPlugin) Initialize(config map[string]interface{}) error {
	if apiToken, ok := config["api_token"].(string); ok {
		gp.apiToken = apiToken
	}
	if repository, ok := config["repository"].(string); ok {
		gp.repository = repository
	}
	return nil
}

// Fetch retrieves GitHub issues
func (gp *GitHubPlugin) Fetch(ctx context.Context) (interface{}, error) {
	if gp.repository == "" {
		return gp.lastData, fmt.Errorf("repository not configured")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/issues?state=open&per_page=10", gp.repository)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return gp.lastData, err
	}

	// Add GitHub API token if available
	if gp.apiToken != "" {
		req.Header.Set("Authorization", "token "+gp.apiToken)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := gp.client.Do(req)
	if err != nil {
		return gp.lastData, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return gp.lastData, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return gp.lastData, err
	}

	var issues []GitHubIssue
	if err := json.Unmarshal(body, &issues); err != nil {
		return gp.lastData, err
	}

	gp.lastData = issues
	return issues, nil
}

// GetMetadata returns plugin metadata
func (gp *GitHubPlugin) GetMetadata() PluginMetadata {
	return PluginMetadata{
		Name:        gp.name,
		Version:     gp.version,
		Description: gp.description,
		Author:      gp.author,
		Type:        gp.pluginType,
		Config: map[string]string{
			"api_token":  gp.apiToken,
			"repository": gp.repository,
		},
	}
}

// Cleanup performs cleanup
func (gp *GitHubPlugin) Cleanup() error {
	return nil
}

// ConvertToWidgetItems converts GitHub issues to widget items
func (gp *GitHubPlugin) ConvertToWidgetItems(data interface{}) []WidgetItem {
	if issues, ok := data.([]GitHubIssue); ok {
		var items []WidgetItem
		for _, issue := range issues {
			// Create status indicator based on labels
			status := ""
			if len(issue.Labels) > 0 {
				switch issue.Labels[0].Name {
				case "bug":
					status = "🐛"
				case "enhancement":
					status = "✨"
				case "documentation":
					status = "📚"
				case "help wanted":
					status = "🙏"
				default:
					status = "📋"
				}
			}

			items = append(items, WidgetItem{
				Title:    fmt.Sprintf("#%d %s", issue.Number, issue.Title),
				Subtitle: fmt.Sprintf("by %s", issue.User.Login),
				Status:   status,
				URL:      issue.HTMLURL,
			})
		}
		return items
	}
	return []WidgetItem{}
}

// CalendarPlugin is an example of another widget type plugin
type CalendarPlugin struct {
	id          string
	pluginType  string
	name        string
	version     string
	description string
	author      string
	apiKey      string
	calendarID  string
	client      *http.Client
	lastData    []CalendarEvent
}

// CalendarEvent represents a calendar event
type CalendarEvent struct {
	ID          string    `json:"id"`
	Title       string    `json:"summary"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start"`
	EndTime     time.Time `json:"end"`
	Location    string    `json:"location"`
	URL         string    `json:"htmlLink"`
}

// NewCalendarPlugin creates a new calendar plugin
func NewCalendarPlugin(apiKey, calendarID string) *CalendarPlugin {
	return &CalendarPlugin{
		id:          "google-calendar",
		pluginType:  "calendar",
		name:        "Google Calendar",
		version:     "1.0.0",
		description: "Fetches events from Google Calendar",
		author:      "GoDay Team",
		apiKey:      apiKey,
		calendarID:  calendarID,
		client:      &http.Client{Timeout: 10 * time.Second},
		lastData:    []CalendarEvent{},
	}
}

// Implement Plugin interface for CalendarPlugin
func (cp *CalendarPlugin) GetID() string   { return cp.id }
func (cp *CalendarPlugin) GetType() string { return cp.pluginType }

func (cp *CalendarPlugin) Initialize(config map[string]interface{}) error {
	if apiKey, ok := config["api_key"].(string); ok {
		cp.apiKey = apiKey
	}
	if calendarID, ok := config["calendar_id"].(string); ok {
		cp.calendarID = calendarID
	}
	return nil
}

func (cp *CalendarPlugin) Fetch(ctx context.Context) (interface{}, error) {
	// This is a placeholder implementation
	// In a real implementation, you would call the Google Calendar API
	mockEvents := []CalendarEvent{
		{
			ID:        "1",
			Title:     "Daily Standup",
			StartTime: time.Now().Add(2 * time.Hour),
			EndTime:   time.Now().Add(2*time.Hour + 30*time.Minute),
			URL:       "https://calendar.google.com/event?eid=1",
		},
		{
			ID:        "2",
			Title:     "Code Review",
			StartTime: time.Now().Add(4 * time.Hour),
			EndTime:   time.Now().Add(5 * time.Hour),
			URL:       "https://calendar.google.com/event?eid=2",
		},
	}

	cp.lastData = mockEvents
	return mockEvents, nil
}

func (cp *CalendarPlugin) GetMetadata() PluginMetadata {
	return PluginMetadata{
		Name:        cp.name,
		Version:     cp.version,
		Description: cp.description,
		Author:      cp.author,
		Type:        cp.pluginType,
		Config: map[string]string{
			"api_key":     cp.apiKey,
			"calendar_id": cp.calendarID,
		},
	}
}

func (cp *CalendarPlugin) Cleanup() error {
	return nil
}

// ConvertToWidgetItems converts calendar events to widget items
func (cp *CalendarPlugin) ConvertToWidgetItems(data interface{}) []WidgetItem {
	if events, ok := data.([]CalendarEvent); ok {
		var items []WidgetItem
		for _, event := range events {
			timeStr := event.StartTime.Format("15:04")
			items = append(items, WidgetItem{
				Title:    event.Title,
				Subtitle: timeStr,
				Status:   "📅",
				URL:      event.URL,
			})
		}
		return items
	}
	return []WidgetItem{}
}
