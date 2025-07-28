package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// GoogleCalendarPlugin implements a plugin for Google Calendar integration
type GoogleCalendarPlugin struct {
	id          string
	pluginType  string
	name        string
	version     string
	description string
	author      string

	// Configuration
	credentialsFile string
	tokenFile       string
	maxEvents       int
	daysAhead       int

	// Internal state
	config      *oauth2.Config
	client      *http.Client
	service     *calendar.Service
	lastData    []GoogleCalendarEvent
	initialized bool
}

// GoogleCalendarEvent represents a calendar event from Google Calendar
type GoogleCalendarEvent struct {
	ID          string    `json:"id"`
	Title       string    `json:"summary"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start"`
	EndTime     time.Time `json:"end"`
	Location    string    `json:"location"`
	URL         string    `json:"htmlLink"`
	Status      string    `json:"status"`
	Attendees   []string  `json:"attendees"`
}

// NewGoogleCalendarPlugin creates a new Google Calendar plugin
func NewGoogleCalendarPlugin() *GoogleCalendarPlugin {
	return &GoogleCalendarPlugin{
		id:          "google-calendar",
		pluginType:  "calendar",
		name:        "Google Calendar",
		version:     "1.0.0",
		description: "Fetches events from Google Calendar using OAuth2",
		author:      "GoDay Team",
		maxEvents:   10,
		daysAhead:   7,
		lastData:    []GoogleCalendarEvent{},
	}
}

// Plugin interface implementation
func (gcp *GoogleCalendarPlugin) GetID() string   { return gcp.id }
func (gcp *GoogleCalendarPlugin) GetType() string { return gcp.pluginType }

func (gcp *GoogleCalendarPlugin) Initialize(config map[string]interface{}) error {
	// Set default file paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	godayDir := filepath.Join(homeDir, ".goday")
	gcp.credentialsFile = filepath.Join(godayDir, "google_calendar_credentials.json")
	gcp.tokenFile = filepath.Join(godayDir, "google_calendar_token.json")

	// Override with config values if provided
	if credFile, ok := config["credentials_file"].(string); ok {
		gcp.credentialsFile = credFile
	}
	if tokenFile, ok := config["token_file"].(string); ok {
		gcp.tokenFile = tokenFile
	}
	if maxEvents, ok := config["max_events"].(int); ok {
		gcp.maxEvents = maxEvents
	}
	if daysAhead, ok := config["days_ahead"].(int); ok {
		gcp.daysAhead = daysAhead
	}

	// Initialize OAuth2 configuration - don't fail if credentials are missing
	if err := gcp.initializeOAuth(); err != nil {
		// Don't fail initialization - just mark as needing setup
		gcp.initialized = false
		fmt.Printf("📅 Calendar setup needed: %v\n", err)
		return nil // Return success but mark as not initialized
	}

	// Get authenticated HTTP client - don't fail if OAuth flow is needed
	client, err := gcp.getClient()
	if err != nil {
		// Don't fail initialization - just mark as needing OAuth
		gcp.initialized = false
		fmt.Printf("📅 Calendar OAuth needed: %v\n", err)
		return nil // Return success but mark as not initialized
	}
	gcp.client = client

	// Initialize Calendar service
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		gcp.initialized = false
		fmt.Printf("📅 Calendar service error: %v\n", err)
		return nil // Return success but mark as not initialized
	}
	gcp.service = srv

	gcp.initialized = true
	fmt.Printf("📅 Calendar plugin initialized successfully\n")
	return nil
}

func (gcp *GoogleCalendarPlugin) initializeOAuth() error {
	// Read credentials file
	credBytes, err := ioutil.ReadFile(gcp.credentialsFile)
	if err != nil {
		return fmt.Errorf("unable to read client secret file %s: %w\n\n"+
			"To setup Google Calendar integration:\n"+
			"1. Go to https://console.cloud.google.com/\n"+
			"2. Create a new project or select existing one\n"+
			"3. Enable the Google Calendar API\n"+
			"4. Create credentials (OAuth 2.0 Client ID)\n"+
			"5. Download the JSON file\n"+
			"6. Save it as %s\n"+
			"7. Restart GoDay", gcp.credentialsFile, err, gcp.credentialsFile)
	}

	// Parse credentials
	config, err := google.ConfigFromJSON(credBytes, calendar.CalendarReadonlyScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	gcp.config = config
	return nil
}

// getClient retrieves a token, saves the token, then returns the generated client
func (gcp *GoogleCalendarPlugin) getClient() (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens.
	tok, err := gcp.tokenFromFile()
	if err != nil {
		// Don't automatically trigger OAuth flow - just return error
		return nil, fmt.Errorf("OAuth token not found. Run './setup-calendar.sh' to set up calendar integration")
	}
	return gcp.config.Client(context.Background(), tok), nil
}

// getTokenFromWeb requests a token from the web, then returns the retrieved token
func (gcp *GoogleCalendarPlugin) getTokenFromWeb() (*oauth2.Token, error) {
	authURL := gcp.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser and then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	fmt.Print("Enter authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %w", err)
	}

	tok, err := gcp.config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %w", err)
	}
	return tok, nil
}

// tokenFromFile retrieves a token from a local file
func (gcp *GoogleCalendarPlugin) tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(gcp.tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken saves a token to a file path
func (gcp *GoogleCalendarPlugin) saveToken(token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", gcp.tokenFile)
	f, err := os.OpenFile(gcp.tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (gcp *GoogleCalendarPlugin) Fetch(ctx context.Context) (interface{}, error) {
	if !gcp.initialized {
		// Return helpful setup information instead of failing
		return []GoogleCalendarEvent{
			{
				ID:        "setup",
				Title:     "Calendar Setup Required",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			},
		}, nil
	}

	// Calculate time range
	now := time.Now()
	timeMin := now.Format(time.RFC3339)
	timeMax := now.AddDate(0, 0, gcp.daysAhead).Format(time.RFC3339)

	// Fetch events from primary calendar
	events, err := gcp.service.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(timeMin).
		TimeMax(timeMax).
		MaxResults(int64(gcp.maxEvents)).
		OrderBy("startTime").
		Do()

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve user's events: %w", err)
	}

	// Convert to our GoogleCalendarEvent format
	var calendarEvents []GoogleCalendarEvent
	for _, item := range events.Items {
		event := GoogleCalendarEvent{
			ID:          item.Id,
			Title:       item.Summary,
			Description: item.Description,
			Location:    item.Location,
			URL:         item.HtmlLink,
			Status:      item.Status,
		}

		// Parse start time
		if item.Start.DateTime != "" {
			if startTime, err := time.Parse(time.RFC3339, item.Start.DateTime); err == nil {
				event.StartTime = startTime
			}
		} else if item.Start.Date != "" {
			if startTime, err := time.Parse("2006-01-02", item.Start.Date); err == nil {
				event.StartTime = startTime
			}
		}

		// Parse end time
		if item.End.DateTime != "" {
			if endTime, err := time.Parse(time.RFC3339, item.End.DateTime); err == nil {
				event.EndTime = endTime
			}
		} else if item.End.Date != "" {
			if endTime, err := time.Parse("2006-01-02", item.End.Date); err == nil {
				event.EndTime = endTime
			}
		}

		// Extract attendees
		for _, attendee := range item.Attendees {
			if attendee.Email != "" {
				event.Attendees = append(event.Attendees, attendee.Email)
			}
		}

		calendarEvents = append(calendarEvents, event)
	}

	gcp.lastData = calendarEvents
	return calendarEvents, nil
}

func (gcp *GoogleCalendarPlugin) GetMetadata() PluginMetadata {
	return PluginMetadata{
		Name:        gcp.name,
		Version:     gcp.version,
		Description: gcp.description,
		Author:      gcp.author,
		Type:        gcp.pluginType,
		Config: map[string]string{
			"credentials_file": "Path to Google OAuth2 credentials JSON file",
			"token_file":       "Path to store OAuth2 tokens",
			"max_events":       "Maximum number of events to fetch (default: 10)",
			"days_ahead":       "Number of days ahead to fetch events (default: 7)",
		},
	}
}

func (gcp *GoogleCalendarPlugin) Cleanup() error {
	// No cleanup needed for HTTP client
	return nil
}

// GetLastData returns the last fetched calendar events
func (gcp *GoogleCalendarPlugin) GetLastData() []GoogleCalendarEvent {
	return gcp.lastData
}

// FormatEventsForDisplay formats calendar events for display in the widget
func (gcp *GoogleCalendarPlugin) FormatEventsForDisplay() []WidgetItem {
	var items []WidgetItem

	// Handle setup case
	if !gcp.initialized && len(gcp.lastData) > 0 && gcp.lastData[0].ID == "setup" {
		return []WidgetItem{
			{
				Title:    "📅 Calendar Setup Required",
				Subtitle: "See GOOGLE_CALENDAR_SETUP.md",
				Status:   "🔧",
			},
			{
				Title:    "🌐 Enable Google Calendar API",
				Subtitle: "console.cloud.google.com",
				Status:   "1️⃣",
			},
			{
				Title:    "🔑 Download OAuth credentials",
				Subtitle: "Save as ~/.goday/credentials.json",
				Status:   "2️⃣",
			},
		}
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	for _, event := range gcp.lastData {
		// Skip past events (except for current ongoing events)
		if event.EndTime.Before(now) {
			continue
		}

		// Format time display
		var timeStr string
		eventDate := event.StartTime.Format("2006-01-02")

		if eventDate == today {
			// Today's events - show time only
			if event.StartTime.Format("15:04") == event.EndTime.Format("15:04") {
				// All-day event
				timeStr = "All day"
			} else {
				timeStr = event.StartTime.Format("15:04")
				if !event.EndTime.IsZero() {
					timeStr += "-" + event.EndTime.Format("15:04")
				}
			}
		} else {
			// Future events - show date and time
			timeStr = event.StartTime.Format("Jan 2")
			if event.StartTime.Format("15:04") != "00:00" {
				timeStr += " " + event.StartTime.Format("15:04")
			}
		}

		// Create status indicator
		var status string
		if event.StartTime.Before(now) && event.EndTime.After(now) {
			status = "🔴" // Currently happening
		} else if event.StartTime.Sub(now) < 30*time.Minute {
			status = "🟡" // Starting soon
		} else {
			status = "🟢" // Future event
		}

		items = append(items, WidgetItem{
			Title:    event.Title,
			Subtitle: timeStr,
			Status:   status,
			URL:      event.URL,
		})

		// Limit to reasonable number for display
		if len(items) >= 5 {
			break
		}
	}

	if len(items) == 0 {
		items = append(items, WidgetItem{
			Title:    "No upcoming events",
			Subtitle: "Your calendar is clear",
			Status:   "📅",
		})
	}

	return items
}

// SetupOAuth performs the OAuth flow for calendar setup
func (gcp *GoogleCalendarPlugin) SetupOAuth() error {
	if gcp.config == nil {
		return fmt.Errorf("OAuth config not initialized. Ensure credentials file exists")
	}

	tok, err := gcp.getTokenFromWeb()
	if err != nil {
		return fmt.Errorf("OAuth setup failed: %w", err)
	}

	gcp.saveToken(tok)

	// Initialize client and service after successful OAuth
	gcp.client = gcp.config.Client(context.Background(), tok)
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(gcp.client))
	if err != nil {
		return fmt.Errorf("failed to create calendar service: %w", err)
	}
	gcp.service = srv
	gcp.initialized = true

	fmt.Printf("✅ Calendar OAuth setup completed successfully!\n")
	return nil
}
