package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	clockInterval   = 60 * time.Second
	weatherInterval = 600 * time.Second
	baseTileWidth   = 30
	baseTileHeight  = 8
)

type clockMsg string
type weatherMsg string
type newsMsg []NewsItem

// Commands that can access the model
type fetchWeatherCmd struct{}
type fetchNewsCmd struct{}
type fetchGitCommitsCmd struct{}
type fetchGitHubPRsCmd struct{}
type fetchTrafficCmd struct{}
type fetchCalendarCmd struct{}

func (fetchWeatherCmd) String() string    { return "fetch weather" }
func (fetchNewsCmd) String() string       { return "fetch news" }
func (fetchGitCommitsCmd) String() string { return "fetch git commits" }
func (fetchGitHubPRsCmd) String() string  { return "fetch github prs" }
func (fetchTrafficCmd) String() string    { return "fetch traffic" }
func (fetchCalendarCmd) String() string   { return "fetch calendar" }

// openURL opens a URL in the default browser
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// Widget item for list
type WidgetListItem struct {
	ItemTitle string
	Subtitle  string
	Status    string
	URL       string
}

func (i WidgetListItem) Title() string       { return i.ItemTitle }
func (i WidgetListItem) Description() string { return i.Subtitle }
func (i WidgetListItem) FilterValue() string { return i.ItemTitle }

// Widget tile model
type WidgetTile struct {
	title    string
	count    int
	hasError bool
	list     list.Model
	width    int
	height   int
}

func NewWidgetTile(title string, width, height int) WidgetTile {
	// Create list items for the widget
	items := []list.Item{
		WidgetListItem{ItemTitle: "Loading...", Subtitle: ""},
	}

	// Create list with proper sizing for content area
	l := list.New(items, list.NewDefaultDelegate(), width-6, height-4)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)

	return WidgetTile{
		title:  title,
		count:  0,
		width:  width,
		height: height,
		list:   l,
	}
}

func (wt *WidgetTile) UpdateItems(items []WidgetItem) {
	var listItems []list.Item
	if len(items) == 0 {
		listItems = []list.Item{
			WidgetListItem{ItemTitle: "No items available", Subtitle: ""},
		}
	} else {
		for _, item := range items {
			listItems = append(listItems, WidgetListItem{
				ItemTitle: item.Title,
				Subtitle:  item.Subtitle,
				Status:    item.Status,
				URL:       item.URL,
			})
		}
	}
	wt.list.SetItems(listItems)
	wt.count = len(items)
}

func (wt *WidgetTile) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229")).
		Align(lipgloss.Center).
		Width(wt.width - 2).
		Background(lipgloss.Color("235"))

	title := fmt.Sprintf("%s (%d)", wt.title, wt.count)
	if wt.hasError {
		title += " ❌"
	}

	// Get items directly from the list instead of using list.View()
	items := wt.list.Items()
	selectedIndex := wt.list.Index()
	var contentLines []string

	// Process each item to create readable content
	for i, item := range items {
		if widgetItem, ok := item.(WidgetListItem); ok {
			// Create a formatted line for each item
			line := widgetItem.ItemTitle
			if widgetItem.Subtitle != "" {
				line += " • " + widgetItem.Subtitle
			}
			if widgetItem.Status != "" {
				line += " " + widgetItem.Status
			}

			// Truncate if too long
			if len(line) > wt.width-4 {
				line = line[:wt.width-7] + "..."
			}

			// Highlight selected item
			if i == selectedIndex {
				selectedStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("0")).
					Background(lipgloss.Color("33")).
					Bold(true)
				line = selectedStyle.Render(line)
			}

			contentLines = append(contentLines, line)

			// Limit to prevent overflow
			if i >= wt.height-4 { // Leave space for title and borders
				remaining := len(items) - i - 1
				if remaining > 0 {
					contentLines = append(contentLines, fmt.Sprintf("+%d more…", remaining))
				}
				break
			}
		}
	}

	// Ensure we have some content
	if len(contentLines) == 0 {
		contentLines = []string{"No items"}
	}

	// Join content with proper spacing
	contentText := strings.Join(contentLines, "\n")

	// Create content area style
	contentStyle := lipgloss.NewStyle().
		Width(wt.width-2).
		Height(wt.height-2).
		Padding(0, 1).
		Align(lipgloss.Left)

	// Combine title and content
	fullContent := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(title),
		contentStyle.Render(contentText),
	)

	return fullContent
}

type Model struct {
	userName       string
	dateTime       string
	weather        string
	location       string
	config         *Config
	widgetManager  *WidgetManager
	pluginManager  *PluginManager
	scheduler      *Scheduler
	cancel         context.CancelFunc
	widgets        []WidgetTile
	focusedWidget  int
	terminalWidth  int
	terminalHeight int
}

func initialModel() Model {
	cfg, err := LoadConfigFromDefaultPath()
	userName := "Unknown User"
	location := "Bengaluru,IN"
	if err == nil && cfg != nil {
		userName = cfg.User.Name
		location = cfg.User.Location
	} else {
		// Log the error but continue with defaults
		fmt.Printf("Warning: Could not load config: %v\n", err)
	}

	widgetManager := NewWidgetManager()
	widgetManager.InitializeWidgets(cfg)
	// Create plugin manager
	pluginConfig := &PluginConfig{
		Plugins: make(map[string]map[string]interface{}),
	}

	if cfg != nil {
		// Configure weather plugin
		pluginConfig.Plugins["openweathermap"] = map[string]interface{}{
			"api_key": cfg.Widgets.Weather.APIKey,
			"city":    location,
		}

		// Configure news plugins
		pluginConfig.Plugins["hackernews"] = map[string]interface{}{
			"tags":        cfg.Widgets.News.Tags,
			"current_tag": "all",
		}
		pluginConfig.Plugins["devto"] = map[string]interface{}{
			"tags":        cfg.Widgets.News.Tags,
			"current_tag": "all",
		}
		pluginConfig.Plugins["aggregate-news"] = map[string]interface{}{
			"tags":        cfg.Widgets.News.Tags,
			"current_tag": "all",
		}

		// Configure traffic plugin (OSRM - no API key needed)
		pluginConfig.Plugins["osrm_traffic"] = map[string]interface{}{
			"origin":      cfg.Widgets.Traffic.Origin,
			"destination": cfg.Widgets.Traffic.Destination,
		}

		// Configure calendar plugin
		calendarConfig := map[string]interface{}{
			"max_events": cfg.Widgets.Calendar.MaxEvents,
			"days_ahead": cfg.Widgets.Calendar.DaysAhead,
		}
		// Add credentials_file and token_file if provided in config
		if cfg.Widgets.Calendar.CredentialsFile != "" {
			calendarConfig["credentials_file"] = cfg.Widgets.Calendar.CredentialsFile
		}
		if cfg.Widgets.Calendar.TokenFile != "" {
			calendarConfig["token_file"] = cfg.Widgets.Calendar.TokenFile
		}
		pluginConfig.Plugins["google-calendar"] = calendarConfig
	} else {
		// Default config when no config file is found
		defaultTags := []string{"golang", "security", "ai"}

		pluginConfig.Plugins["openweathermap"] = map[string]interface{}{
			"api_key": "YOUR_OWM_API_KEY",
			"city":    location,
		}

		pluginConfig.Plugins["hackernews"] = map[string]interface{}{
			"tags":        defaultTags,
			"current_tag": "all",
		}
		pluginConfig.Plugins["devto"] = map[string]interface{}{
			"tags":        defaultTags,
			"current_tag": "all",
		}
		pluginConfig.Plugins["aggregate-news"] = map[string]interface{}{
			"tags":        defaultTags,
			"current_tag": "all",
		}

		// Configure traffic plugin with defaults (OSRM - no API key needed)
		pluginConfig.Plugins["osrm_traffic"] = map[string]interface{}{
			"origin":      "Electronic City, Bengaluru, Karnataka, India",
			"destination": "Whitefield, Bengaluru, Karnataka, India",
		}

		// Configure calendar plugin with defaults
		pluginConfig.Plugins["google-calendar"] = map[string]interface{}{
			"max_events": 10,
			"days_ahead": 7,
		}
	}

	pluginManager := NewPluginManager(pluginConfig)

	// Register plugins - handle nil config gracefully
	var apiKey string
	if cfg != nil {
		apiKey = cfg.Widgets.Weather.APIKey
	}
	weatherPlugin := NewWeatherPlugin(apiKey, location)
	pluginManager.RegisterPlugin(weatherPlugin)

	// Create individual news plugins
	hackerNewsPlugin := NewHackerNewsPlugin()
	devToPlugin := NewDevToPlugin()
	hackernoonPlugin := NewHackernoonPlugin()
	pluginManager.RegisterPlugin(hackerNewsPlugin)
	pluginManager.RegisterPlugin(devToPlugin)
	pluginManager.RegisterPlugin(hackernoonPlugin)

	// Create aggregate news plugin with only tech-focused sources
	// Removed Hacker News as it includes general news articles
	aggregateNewsPlugin := NewAggregateNewsPlugin([]NewsPlugin{
		hackernoonPlugin,
		devToPlugin,
	})
	pluginManager.RegisterPlugin(aggregateNewsPlugin)

	// Create Git plugins
	gitCommitsPlugin := NewLocalGitCommitsPlugin()
	githubPRsPlugin := NewGitHubPRsPlugin()
	pluginManager.RegisterPlugin(gitCommitsPlugin)
	pluginManager.RegisterPlugin(githubPRsPlugin)

	// Create Traffic plugin (OSRM - no API key required)
	trafficPlugin := NewOSRMTrafficPlugin()
	pluginManager.RegisterPlugin(trafficPlugin)

	// Create Google Calendar plugin
	calendarPlugin := NewGoogleCalendarPlugin()
	pluginManager.RegisterPlugin(calendarPlugin)

	scheduler := NewScheduler()

	// Add scheduled tasks for each widget with their TTL
	if cfg != nil {
		scheduler.AddTask("weather", ParseTTL(cfg.Widgets.Weather.TTL), weatherPlugin)
		scheduler.AddTask("news", ParseTTL(cfg.Widgets.News.TTL), aggregateNewsPlugin)
		scheduler.AddTask("slack", ParseTTL(cfg.Widgets.Slack.TTL), nil)
		scheduler.AddTask("confluence", ParseTTL(cfg.Widgets.Confluence.TTL), nil)
		scheduler.AddTask("jira", ParseTTL(cfg.Widgets.Jira.TTL), nil)
		scheduler.AddTask("traffic", ParseTTL(cfg.Widgets.Traffic.TTL), trafficPlugin)
		scheduler.AddTask("calendar", ParseTTL(cfg.Widgets.Calendar.TTL), calendarPlugin)
	} else {
		// Default TTL values when no config
		scheduler.AddTask("weather", 600*time.Second, weatherPlugin)
		scheduler.AddTask("news", 600*time.Second, aggregateNewsPlugin)
		scheduler.AddTask("slack", 20*time.Second, nil)
		scheduler.AddTask("confluence", 300*time.Second, nil)
		scheduler.AddTask("jira", 45*time.Second, nil)
		scheduler.AddTask("traffic", 300*time.Second, trafficPlugin)
		scheduler.AddTask("calendar", 300*time.Second, calendarPlugin)
	}

	// Create widget tiles with fixed sizes
	widgets := []WidgetTile{
		NewWidgetTile("JIRA", baseTileWidth, baseTileHeight),
		NewWidgetTile("PRs", baseTileWidth, baseTileHeight),
		NewWidgetTile("Builds", baseTileWidth, baseTileHeight),
		NewWidgetTile("Commits", baseTileWidth, baseTileHeight),
		NewWidgetTile("Calendar", baseTileWidth, baseTileHeight),
		NewWidgetTile("Slack", baseTileWidth, baseTileHeight),
		NewWidgetTile("Todos", baseTileWidth, baseTileHeight),
		NewWidgetTile("Confluence", baseTileWidth, baseTileHeight),
		NewWidgetTile("PagerDuty", baseTileWidth, baseTileHeight),
		NewWidgetTile("Tech News", baseTileWidth, baseTileHeight),
		NewWidgetTile("Traffic", baseTileWidth, baseTileHeight),
	}

	// Populate widgets with data
	widgetNames := []string{"jira", "prs", "builds", "commits", "calendar", "slack", "todos", "confluence", "pagerduty", "news", "traffic"}
	for i, name := range widgetNames {
		if widget, exists := widgetManager.Widgets[name]; exists {
			var items []WidgetItem
			for _, item := range widget.Items {
				items = append(items, WidgetItem{
					Title:    item.Title,
					Subtitle: item.Subtitle,
					Status:   item.Status,
					URL:      item.URL,
				})
			}
			widgets[i].UpdateItems(items)
			widgets[i].hasError = widget.HasError
		}
	}

	return Model{
		userName:       userName,
		dateTime:       time.Now().Format("Mon 02 Jan 2006 15:04"),
		weather:        fmt.Sprintf("☁ N/A (%s)", location),
		location:       location,
		config:         cfg,
		widgetManager:  widgetManager,
		pluginManager:  pluginManager,
		scheduler:      scheduler,
		widgets:        widgets,
		focusedWidget:  0,
		terminalWidth:  100,
		terminalHeight: 24,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickClock(),
		tickWeather(),
		tickNews(),
		func() tea.Msg { return fetchNewsCmd{} }, // Immediate news fetch
		func() tea.Msg { return fetchWeatherCmd{} },    // Immediate weather fetch
		func() tea.Msg { return fetchGitCommitsCmd{} }, // Immediate git commits fetch
		func() tea.Msg { return fetchGitHubPRsCmd{} },  // Immediate GitHub PRs fetch
		func() tea.Msg { return fetchTrafficCmd{} },    // Immediate traffic fetch
		func() tea.Msg { return fetchCalendarCmd{} },   // Immediate calendar fetch
		tea.EnterAltScreen,
	)
}

func tickClock() tea.Cmd {
	return tea.Tick(clockInterval, func(t time.Time) tea.Msg {
		return clockMsg(t.Format("Mon 02 Jan 2006 15:04"))
	})
}

func tickWeather() tea.Cmd {
	return tea.Tick(weatherInterval, func(t time.Time) tea.Msg {
		return fetchWeatherCmd{}
	})
}

func tickNews() tea.Cmd {
	return tea.Tick(weatherInterval, func(t time.Time) tea.Msg {
		return fetchNewsCmd{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
		case "tab":
			m.focusedWidget = (m.focusedWidget + 1) % len(m.widgets)
			return m, nil
		case "shift+tab":
			m.focusedWidget = (m.focusedWidget - 1 + len(m.widgets)) % len(m.widgets)
			return m, nil
		case "up", "k":
			// Navigate up within the focused widget
			if m.focusedWidget < len(m.widgets) {
				var cmd tea.Cmd
				m.widgets[m.focusedWidget].list, cmd = m.widgets[m.focusedWidget].list.Update(msg)
				return m, cmd
			}
			return m, nil
		case "down", "j":
			// Navigate down within the focused widget
			if m.focusedWidget < len(m.widgets) {
				var cmd tea.Cmd
				m.widgets[m.focusedWidget].list, cmd = m.widgets[m.focusedWidget].list.Update(msg)
				return m, cmd
			}
			return m, nil
		case "t":
			m.widgetManager.CycleNewsTag()
			// Update the Tech News widget and refresh news
			m.updateNewsWidget()
			// Set the current tag in the news plugins
			currentTag := m.widgetManager.GetCurrentNewsTag()
			tagToSet := "all"
			if currentTag != "All" {
				tagToSet = strings.ToLower(currentTag)
			}

			// Update all news plugins
			newsPlugins := m.pluginManager.GetRegistry().GetAllNewsPlugins()
			for _, plugin := range newsPlugins {
				plugin.SetCurrentTag(tagToSet)
			}

			// Trigger immediate news refresh
			return m, func() tea.Msg { return fetchNewsCmd{} }
		case "T":
			m.widgetManager.NewsTagIndex = 0 // Reset to "All"
			// Update the Tech News widget and refresh news
			m.updateNewsWidget()

			// Set tag to "all" on all news plugins
			newsPlugins := m.pluginManager.GetRegistry().GetAllNewsPlugins()
			for _, plugin := range newsPlugins {
				plugin.SetCurrentTag("all")
			}

			// Trigger immediate news refresh
			return m, func() tea.Msg { return fetchNewsCmd{} }
		case "r", "R":
			// Refresh all widgets
			return m, tea.Batch(tickWeather(), tickNews())
		case "enter":
			// Open the selected item in the focused widget
			if m.focusedWidget < len(m.widgets) {
				selected := m.widgets[m.focusedWidget].list.SelectedItem()
				if item, ok := selected.(WidgetListItem); ok && item.URL != "" {
					// Open URL in browser
					go func() {
						if err := openURL(item.URL); err != nil {
							fmt.Printf("Error opening URL: %v\n", err)
						}
					}()
					// Show feedback message
					fmt.Printf("Opening: %s\n", item.URL)
				}
			}
			return m, nil
		}
	case clockMsg:
		m.dateTime = string(msg)
		return m, tickClock()
	case weatherMsg:
		m.weather = string(msg)
		return m, tickWeather()
	case newsMsg:
		// Update news widget with real data
		if len(msg) > 0 {
			var items []WidgetItem
			for _, news := range msg {
				// Format subtitle to include source
				subtitle := news.Author
				if news.Source == "hackernews" {
					subtitle = fmt.Sprintf("%s • HN", news.Author)
					if news.Points > 0 {
						subtitle = fmt.Sprintf("%s • %d pts", subtitle, news.Points)
					}
				} else if news.Source == "devto" {
					subtitle = fmt.Sprintf("%s • Dev.to", news.Author)
				}

				items = append(items, WidgetItem{
					Title:    news.Title,
					Subtitle: subtitle,
					URL:      news.URL,
				})
			}
			// Update the Tech News widget (index 9)
			if len(m.widgets) > 9 {
				m.widgets[9].UpdateItems(items)
			}
		}
		return m, tickNews()
	case fetchWeatherCmd:
		// Fetch real weather data using plugin
		weatherPlugin, exists := m.pluginManager.GetRegistry().GetPlugin("openweathermap")
		if !exists {
			return m, tea.Batch(
				tea.Tick(weatherInterval, func(t time.Time) tea.Msg { return fetchWeatherCmd{} }),
			)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		data, err := weatherPlugin.Fetch(ctx)
		if err != nil {
			return m, tea.Batch(
				tea.Tick(weatherInterval, func(t time.Time) tea.Msg { return fetchWeatherCmd{} }),
			)
		}

		if weatherData, ok := data.(*WeatherData); ok {
			return m, tea.Batch(
				tea.Tick(weatherInterval, func(t time.Time) tea.Msg { return fetchWeatherCmd{} }),
				func() tea.Msg {
					return weatherMsg(fmt.Sprintf("%s %d°C (%s)", weatherData.Icon, weatherData.Temperature, m.location))
				},
			)
		}

		return m, tea.Batch(
			tea.Tick(weatherInterval, func(t time.Time) tea.Msg { return fetchWeatherCmd{} }),
		)
	case fetchNewsCmd:
		// Fetch real news data using aggregate plugin
		newsPlugin, exists := m.pluginManager.GetRegistry().GetPlugin("aggregate-news")
		if !exists {
			// Update news widget to show error
			if len(m.widgets) > 9 {
				m.widgets[9].UpdateItems([]WidgetItem{
					{Title: "Plugin not found", Subtitle: "aggregate-news missing", Status: "❌"},
				})
			}
			return m, tea.Batch(
				tea.Tick(weatherInterval, func(t time.Time) tea.Msg { return fetchNewsCmd{} }),
			)
		}

		// Show fetching status
		if len(m.widgets) > 9 {
			m.widgets[9].UpdateItems([]WidgetItem{
				{Title: "Fetching news...", Subtitle: "Connecting to APIs", Status: "🔄"},
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		data, err := newsPlugin.Fetch(ctx)
		if err != nil {
			// Update news widget to show error
			if len(m.widgets) > 9 {
				m.widgets[9].UpdateItems([]WidgetItem{
					{Title: "Failed to fetch news", Subtitle: err.Error(), Status: "❌"},
				})
			}
			return m, tea.Batch(
				tea.Tick(weatherInterval, func(t time.Time) tea.Msg { return fetchNewsCmd{} }),
			)
		}

		if items, ok := data.([]NewsItem); ok {
			return m, tea.Batch(
				tea.Tick(weatherInterval, func(t time.Time) tea.Msg { return fetchNewsCmd{} }),
				func() tea.Msg { return newsMsg(items) },
			)
		} else {
			// Update news widget to show type error
			if len(m.widgets) > 9 {
				m.widgets[9].UpdateItems([]WidgetItem{
					{Title: "Data type error", Subtitle: fmt.Sprintf("Got %T", data), Status: "❌"},
				})
			}
		}

		return m, tea.Batch(
			tea.Tick(weatherInterval, func(t time.Time) tea.Msg { return fetchNewsCmd{} }),
		)
	case fetchGitCommitsCmd:
		// Fetch Git commits using local Git plugin
		gitPlugin, exists := m.pluginManager.GetRegistry().GetPlugin("local-git-commits")
		if exists {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			data, err := gitPlugin.Fetch(ctx)
			if err == nil {
				if commits, ok := data.([]GitCommit); ok {
					m.widgetManager.UpdateGitCommitsWidget(commits)
				}
			}
		}

		return m, tea.Batch(
			tea.Tick(5*time.Minute, func(t time.Time) tea.Msg { return fetchGitCommitsCmd{} }),
		)
	case fetchGitHubPRsCmd:
		// Fetch GitHub PRs using GitHub plugin
		githubPlugin, exists := m.pluginManager.GetRegistry().GetPlugin("github-prs")
		if exists {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			data, err := githubPlugin.Fetch(ctx)
			if err == nil {
				if prs, ok := data.([]GitPullRequest); ok {
					m.widgetManager.UpdateGitHubPRsWidget(prs)
				}
			}
		}

		return m, tea.Batch(
			tea.Tick(5*time.Minute, func(t time.Time) tea.Msg { return fetchGitHubPRsCmd{} }),
		)
	case fetchTrafficCmd:
		// Fetch traffic data using OSRM plugin
		trafficPlugin, exists := m.pluginManager.GetRegistry().GetPlugin("osrm_traffic")
		if exists {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			data, err := trafficPlugin.Fetch(ctx)
			if err == nil {
				if biTraffic, ok := data.(*BiDirectionalTrafficData); ok {
					m.widgetManager.UpdateBiDirectionalTrafficWidget(biTraffic)
					// Update the traffic widget (index 10)
					if len(m.widgets) > 10 {
						if widget, exists := m.widgetManager.Widgets["traffic"]; exists {
							var items []WidgetItem
							for _, item := range widget.Items {
								items = append(items, WidgetItem{
									Title:    item.Title,
									Subtitle: item.Subtitle,
									Status:   item.Status,
									URL:      item.URL,
								})
							}
							m.widgets[10].UpdateItems(items)
							m.widgets[10].hasError = widget.HasError
						}
					}
				} else if traffic, ok := data.(*TrafficData); ok {
					// Fallback for single direction traffic data
					m.widgetManager.UpdateTrafficWidget(traffic)
					// Update the traffic widget (index 10)
					if len(m.widgets) > 10 {
						if widget, exists := m.widgetManager.Widgets["traffic"]; exists {
							var items []WidgetItem
							for _, item := range widget.Items {
								items = append(items, WidgetItem{
									Title:    item.Title,
									Subtitle: item.Subtitle,
									Status:   item.Status,
									URL:      item.URL,
								})
							}
							m.widgets[10].UpdateItems(items)
							m.widgets[10].hasError = widget.HasError
						}
					}
				}
			} else {
				// Update traffic widget to show error
				if len(m.widgets) > 10 {
					m.widgets[10].UpdateItems([]WidgetItem{
						{Title: "Traffic unavailable", Subtitle: err.Error(), Status: "❌"},
					})
					m.widgets[10].hasError = true
				}
			}
		}

		return m, tea.Batch(
			tea.Tick(5*time.Minute, func(t time.Time) tea.Msg { return fetchTrafficCmd{} }),
		)
	case fetchCalendarCmd:
		// Fetch calendar data using Google Calendar plugin
		calendarPlugin, exists := m.pluginManager.GetRegistry().GetPlugin("google-calendar")
		if exists {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			data, err := calendarPlugin.Fetch(ctx)
			if err == nil {
				if events, ok := data.([]GoogleCalendarEvent); ok && len(events) > 0 {
					// Type assert to GoogleCalendarPlugin to access FormatEventsForDisplay
					if gcPlugin, ok := calendarPlugin.(*GoogleCalendarPlugin); ok {
						m.widgetManager.UpdateCalendarWidget(gcPlugin)
						// Update the calendar widget (index 4)
						if len(m.widgets) > 4 {
							if widget, exists := m.widgetManager.Widgets["calendar"]; exists {
								var items []WidgetItem
								for _, item := range widget.Items {
									items = append(items, WidgetItem{
										Title:    item.Title,
										Subtitle: item.Subtitle,
										Status:   item.Status,
										URL:      item.URL,
									})
								}
								m.widgets[4].UpdateItems(items)
								m.widgets[4].hasError = widget.HasError
							}
						}
					}
				}
			} else {
				// Update calendar widget to show error
				if len(m.widgets) > 4 {
					// Check if it's an OAuth error requiring setup
					errorMsg := err.Error()
					if strings.Contains(errorMsg, "credentials") || strings.Contains(errorMsg, "oauth") {
						m.widgets[4].UpdateItems([]WidgetItem{
							{Title: "Calendar Setup Required", Subtitle: "See ~/.goday/google_calendar_credentials.json", Status: "🔧"},
							{Title: "Setup Guide", Subtitle: "Check console.cloud.google.com", Status: "📋"},
						})
					} else {
						m.widgets[4].UpdateItems([]WidgetItem{
							{Title: "Calendar unavailable", Subtitle: errorMsg, Status: "❌"},
						})
					}
					m.widgets[4].hasError = true
				}
			}
		}

		return m, tea.Batch(
			tea.Tick(5*time.Minute, func(t time.Time) tea.Msg { return fetchCalendarCmd{} }),
		)
	}

	// Handle list updates for the focused widget
	if m.focusedWidget < len(m.widgets) {
		var cmd tea.Cmd
		m.widgets[m.focusedWidget].list, cmd = m.widgets[m.focusedWidget].list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	// Header styling with proper weather pill
	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("229")).
		Bold(true).
		Padding(0, 2).
		Width(m.terminalWidth - 4).
		Align(lipgloss.Left)

	weatherPill := lipgloss.NewStyle().
		Background(lipgloss.Color("24")).
		Foreground(lipgloss.Color("15")).
		Padding(0, 1).
		Bold(true)

	refreshPill := lipgloss.NewStyle().
		Background(lipgloss.Color("88")).
		Foreground(lipgloss.Color("15")).
		Padding(0, 1).
		Bold(true)

	headerContent := fmt.Sprintf("%s  •  %s  •  %s  •  %s",
		m.userName,
		m.dateTime,
		weatherPill.Render(m.weather),
		refreshPill.Render("R Refresh"),
	)

	header := headerStyle.Render(headerContent)

	grid := m.renderWidgetGrid()

	// Legend styling
	legendStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Italic(true).
		Padding(1, 2)

	legend := legendStyle.Render("Legend: [w] log work; Enter opens link; ↑↓/jk navigate items; Tab/Shift+Tab moves focus; t/T cycles news tags; r/R refresh")

	// Get selected item URL for display
	selectedURL := m.getSelectedItemURL()
	urlDisplay := ""
	if selectedURL != "" {
		urlStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Background(lipgloss.Color("236")).
			Padding(0, 2).
			Bold(true)
		urlDisplay = urlStyle.Render(m.formatURLDisplay(selectedURL))
	} else {
		// Show focused widget info even when no URL
		if m.focusedWidget < len(m.widgets) {
			title, subtitle, _ := m.getSelectedItemDetails()
			if title != "" {
				infoStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("245")).
					Background(lipgloss.Color("236")).
					Padding(0, 2).
					Italic(true)

				widgetName := m.widgets[m.focusedWidget].title
				info := fmt.Sprintf("[%s] %s", widgetName, title)
				if subtitle != "" {
					info += fmt.Sprintf(" • %s", subtitle)
				}
				urlDisplay = infoStyle.Render(info)
			}
		}
	}

	// Combine all parts without extra container
	var contentParts []string
	contentParts = append(contentParts, header, "", grid)

	if urlDisplay != "" {
		contentParts = append(contentParts, "", urlDisplay)
	}

	contentParts = append(contentParts, "", legend)

	content := lipgloss.JoinVertical(lipgloss.Left, contentParts...)

	return content
}

func (m Model) renderWidgetGrid() string {
	// Calculate tiles per row (3 for better readability)
	tilesPerRow := 3
	// Dynamic tile sizing based on terminal width
	tileWidth := baseTileWidth
	tileHeight := baseTileHeight

	// Make tiles much larger and use more screen space
	if m.terminalWidth > 120 {
		tileWidth = (m.terminalWidth - 10) / 3 // Use most of screen width
		tileHeight = baseTileHeight + 3
	} else if m.terminalWidth > 90 {
		tileWidth = baseTileWidth + 15
		tileHeight = baseTileHeight + 2
	}

	var rows []string

	for i := 0; i < len(m.widgets); i += tilesPerRow {
		var rowTiles []string
		for j := 0; j < tilesPerRow && i+j < len(m.widgets); j++ {
			tileIndex := i + j
			tile := m.widgets[tileIndex]

			// Update tile dimensions
			tile.width = tileWidth
			tile.height = tileHeight

			// Update the list dimensions to match new tile size
			tile.list.SetSize(tileWidth-6, tileHeight-4)

			// Create tile content
			tileContent := tile.View()

			// Apply border styling
			var borderStyle lipgloss.Style
			if tileIndex == m.focusedWidget {
				borderStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("33")).
					Width(tileWidth).
					Height(tileHeight).
					Bold(true).
					BorderStyle(lipgloss.DoubleBorder())
			} else {
				borderStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("240")).
					Width(tileWidth).
					Height(tileHeight)
			}

			styledTile := borderStyle.Render(tileContent)
			rowTiles = append(rowTiles, styledTile)

			// Update the original widget in the model
			m.widgets[tileIndex] = tile
		}

		// Join tiles horizontally with spacing
		row := lipgloss.JoinHorizontal(lipgloss.Top, rowTiles...)
		rows = append(rows, row)
	}

	// Join all rows vertically with spacing
	grid := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return grid
}

func (m *Model) updateNewsWidget() {
	currentTag := m.widgetManager.GetCurrentNewsTag()
	// Update the Tech News widget title to show current tag
	if len(m.widgets) > 9 {
		m.widgets[9].title = fmt.Sprintf("Tech News [%s]", currentTag)
	}
}

// getSelectedItemURL returns the URL of the currently selected item
func (m Model) getSelectedItemURL() string {
	if m.focusedWidget >= len(m.widgets) {
		return ""
	}

	selectedItem := m.widgets[m.focusedWidget].list.SelectedItem()
	if selectedItem == nil {
		return ""
	}

	if widgetItem, ok := selectedItem.(WidgetListItem); ok {
		return widgetItem.URL
	}

	return ""
}

// formatURLDisplay formats the URL for display at the bottom
func (m Model) formatURLDisplay(url string) string {
	if url == "" {
		return ""
	}

	// Get the focused widget name for context
	widgetName := ""
	if m.focusedWidget < len(m.widgets) {
		widgetName = m.widgets[m.focusedWidget].title
	}

	// Truncate URL if it's too long
	maxURLLength := m.terminalWidth - 30 // Leave space for prefix and widget name
	if len(url) > maxURLLength {
		url = url[:maxURLLength-3] + "..."
	}

	return fmt.Sprintf("[%s] → %s", widgetName, url)
}

// getSelectedItemDetails returns details of the currently selected item
func (m Model) getSelectedItemDetails() (title, subtitle, url string) {
	if m.focusedWidget >= len(m.widgets) {
		return "", "", ""
	}

	selectedItem := m.widgets[m.focusedWidget].list.SelectedItem()
	if selectedItem == nil {
		return "", "", ""
	}

	if widgetItem, ok := selectedItem.(WidgetListItem); ok {
		return widgetItem.ItemTitle, widgetItem.Subtitle, widgetItem.URL
	}

	return "", "", ""
}

func main() {
	// Check for command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config", "--config", "-c":
			configPath, err := GetConfigPath()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting config path: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Config file location: %s\n", configPath)

			// Check if config exists
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				fmt.Println("Config file does not exist. Run './setup-config.sh' to create it.")
			} else {
				fmt.Println("Config file exists and ready to use.")
			}
			return
		case "help", "--help", "-h":
			fmt.Println("GoDay Terminal Dashboard")
			fmt.Println("")
			fmt.Println("Usage:")
			fmt.Println("  goday              Start the dashboard")
			fmt.Println("  goday config       Show config file location")
			fmt.Println("  goday help         Show this help message")
			fmt.Println("")
			fmt.Println("Config file: ~/.goday/config.yaml")
			fmt.Println("Setup:       ./setup-config.sh")
			return
		}
	}

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
