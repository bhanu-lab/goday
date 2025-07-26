# GoDay Plugin System

The GoDay dashboard uses a plugin-based architecture that allows easy extension with new data sources and widget types. This document explains how to create and register new plugins.

## Plugin Architecture

### Core Interfaces

#### Plugin Interface
All plugins must implement the base `Plugin` interface:

```go
type Plugin interface {
    GetID() string                                    // Unique plugin identifier
    GetType() string                                  // Plugin type (e.g., "news", "weather", "calendar")
    Initialize(config map[string]interface{}) error  // Setup with configuration
    Fetch(ctx context.Context) (interface{}, error)  // Retrieve data from source
    GetMetadata() PluginMetadata                      // Return plugin metadata
    Cleanup() error                                   // Perform cleanup
}
```

#### NewsPlugin Interface
News plugins implement additional functionality for tag-based filtering:

```go
type NewsPlugin interface {
    Plugin
    SetTags(tags []string)           // Configure filter tags
    GetCurrentTag() string           // Get active tag
    SetCurrentTag(tag string)        // Set active tag
    GetSupportedTags() []string      // List all supported tags
}
```

## Creating New Plugins

### 1. News Plugins

To add a new news source (e.g., Reddit, Medium):

```go
// RedditPlugin fetches posts from Reddit
type RedditPlugin struct {
    *BaseNewsPlugin
    subreddit string
    apiKey    string
}

func NewRedditPlugin(subreddit string) *RedditPlugin {
    base := NewBaseNewsPlugin(
        "reddit",
        "Reddit",
        "1.0.0",
        "Fetches posts from Reddit subreddits",
        "Your Name",
    )
    base.supportedTags = []string{"all", "programming", "golang", "tech"}
    
    return &RedditPlugin{
        BaseNewsPlugin: base,
        subreddit:      subreddit,
    }
}

func (rp *RedditPlugin) Initialize(config map[string]interface{}) error {
    if subreddit, ok := config["subreddit"].(string); ok {
        rp.subreddit = subreddit
    }
    if apiKey, ok := config["api_key"].(string); ok {
        rp.apiKey = apiKey
    }
    return nil
}

func (rp *RedditPlugin) Fetch(ctx context.Context) (interface{}, error) {
    // Implement Reddit API calls here
    // Return []NewsItem
    return []NewsItem{}, nil
}
```

### 2. Other Widget Types

For non-news widgets (e.g., GitHub issues, JIRA tickets):

```go
type JIRAPlugin struct {
    id          string
    pluginType  string
    // ... other fields
}

func (jp *JIRAPlugin) GetID() string { return jp.id }
func (jp *JIRAPlugin) GetType() string { return jp.pluginType }

func (jp *JIRAPlugin) Fetch(ctx context.Context) (interface{}, error) {
    // Implement JIRA API calls
    // Return any data structure
    return []JIRAIssue{}, nil
}

// Convert plugin data to widget items for display
func (jp *JIRAPlugin) ConvertToWidgetItems(data interface{}) []WidgetItem {
    if issues, ok := data.([]JIRAIssue); ok {
        var items []WidgetItem
        for _, issue := range issues {
            items = append(items, WidgetItem{
                Title:    issue.Summary,
                Subtitle: issue.Status,
                URL:      issue.URL,
                Status:   "🎫",
            })
        }
        return items
    }
    return []WidgetItem{}
}
```

## Registering Plugins

### 1. In initialModel() function

Add your plugin registration in the `initialModel()` function in `main.go`:

```go
// Create and register your plugin
redditPlugin := NewRedditPlugin("programming")
pluginManager.RegisterPlugin(redditPlugin)

// Add to aggregate news plugin if it's a news source
newsPlugins := []NewsPlugin{
    hackerNewsPlugin,
    devToPlugin,
    redditPlugin,  // Add your news plugin here
}
aggregateNewsPlugin := NewAggregateNewsPlugin(newsPlugins)
```

### 2. Configure in config.yaml

Add plugin configuration to your `config.yaml`:

```yaml
plugins:
  reddit:
    subreddit: "programming"
    api_key: "your_reddit_api_key"
  jira:
    api_key: "your_jira_token"
    base_url: "https://yourcompany.atlassian.net"
    project: "PROJ"
```

### 3. Schedule Plugin Execution

Add your plugin to the scheduler for periodic updates:

```go
// Schedule the plugin to run every 5 minutes
pluginManager.SchedulePlugin("reddit", 5*time.Minute)
```

## Best Practices

### 1. Error Handling
- Always handle API errors gracefully
- Return cached data when API calls fail
- Use context timeouts for HTTP requests

### 2. Rate Limiting
- Respect API rate limits
- Implement exponential backoff for failed requests
- Cache data to reduce API calls

### 3. Configuration
- Use the config map for all plugin settings
- Provide sensible defaults
- Validate configuration in Initialize()

### 4. Testing
- Mock external API calls in tests
- Test error conditions
- Verify data transformation

## Example: Adding a Slack Plugin

Here's a complete example of adding a Slack plugin:

1. **Create the plugin** (`slack_plugin.go`):

```go
type SlackPlugin struct {
    id         string
    pluginType string
    token      string
    channels   []string
    client     *http.Client
    lastData   []SlackMessage
}

func NewSlackPlugin(token string, channels []string) *SlackPlugin {
    return &SlackPlugin{
        id:         "slack",
        pluginType: "messaging",
        token:      token,
        channels:   channels,
        client:     &http.Client{Timeout: 10 * time.Second},
    }
}

func (sp *SlackPlugin) Fetch(ctx context.Context) (interface{}, error) {
    // Implement Slack API calls using sp.token
    // Return []SlackMessage
    return sp.lastData, nil
}
```

2. **Register in main.go**:

```go
slackPlugin := NewSlackPlugin(cfg.Widgets.Slack.Token, cfg.Widgets.Slack.Channels)
pluginManager.RegisterPlugin(slackPlugin)
pluginManager.SchedulePlugin("slack", ParseTTL(cfg.Widgets.Slack.TTL))
```

3. **Update widget handling** to use the plugin data:

```go
case fetchSlackCmd:
    slackPlugin, exists := m.pluginManager.GetRegistry().GetPlugin("slack")
    if exists {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        data, err := slackPlugin.Fetch(ctx)
        if err == nil {
            // Update widget with data
            // Transform data to WidgetItems and update UI
        }
    }
```

## Available Plugin Examples

The codebase includes several example plugins:

- **HackerNewsPlugin**: Fetches tech news from Hacker News
- **DevToPlugin**: Fetches articles from Dev.to
- **WeatherPlugin**: Gets weather data from OpenWeatherMap
- **GitHubPlugin**: Fetches issues from GitHub repositories
- **CalendarPlugin**: Gets events from Google Calendar

These examples demonstrate different patterns for API integration, data transformation, and error handling.

## Future Extensibility

The plugin system is designed to support:

- **Dynamic plugin loading**: Load plugins from external files
- **Plugin marketplace**: Download and install community plugins
- **Real-time updates**: WebSocket and SSE support for live data
- **Custom widget types**: Create entirely new widget categories
- **Plugin dependencies**: Plugins that depend on other plugins
- **Plugin configuration UI**: Web interface for plugin management

## Troubleshooting

### Common Issues

1. **Plugin not found**: Ensure plugin ID matches registration
2. **API errors**: Check API keys and network connectivity
3. **Data format errors**: Verify data transformation logic
4. **Memory leaks**: Implement proper cleanup in Cleanup() method

### Debugging

Enable debug logging by setting environment variable:
```bash
export GODAY_DEBUG=true
```

This will show plugin execution times, error details, and API response information.
