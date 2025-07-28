# GoDay Terminal Dashboard

A **single-view, zero-navigation** terminal dashboard that surfaces everything an engineer needs — tasks, code, builds, messages — and key personal context (user name, date, local weather) plus a bite-sized "Tech News" feed that can be filtered by tags.

## Features

- **Header Bar**: Shows user name, current date/time, and weather with live updates
- **Widget Grid**: Interactive 3x4 tile layout with all your essential tools
- **Tech News**: Real articles from Hacker News and Dev.to, filterable by tags
- **Plugin Architecture**: Extensible system for adding new data sources
- **Per-widget TTL**: Each widget refreshes at its own interval
- **Navigation**: Tab between widgets, arrow keys within widgets, Enter to open links
- **Live Data**: Real API integrations with fallback to cached data

## Widgets

- **JIRA**: Tasks with work-log shortcuts `[w]` (interactive)
- **PRs**: Pull requests with review status (interactive)
- **Builds**: CI/CD status with error indicators (interactive)
- **Commits**: Recent repository activity (interactive)
- **Calendar**: Gmail events with smart notifications and status indicators (Google Calendar API)
- **Slack**: Unread messages and channels (interactive)
- **Todos**: Personal task list (interactive)
- **Confluence**: Recent documentation updates (interactive)
- **PagerDuty**: On-call status (interactive)
- **Tech News**: Live articles from multiple sources with tag filtering (interactive)
- **Traffic**: Real-time commute information with direction toggle (interactive)

## Plugin Architecture

GoDay now uses a plugin-based architecture that makes it easy to add new data sources and widget types:

### Built-in Plugins

- **HackerNewsPlugin**: Fetches tech news from Hacker News
- **DevToPlugin**: Fetches articles from Dev.to
- **AggregateNewsPlugin**: Combines multiple news sources
- **WeatherPlugin**: Gets weather data from OpenWeatherMap
- **OSRMTrafficPlugin**: Real-time traffic data using OpenStreetMap (free)
- **GoogleCalendarPlugin**: Gmail account calendar events via Google Calendar API

### Adding New Plugins

The system supports easy extension with new plugins. See [PLUGIN_GUIDE.md](PLUGIN_GUIDE.md) for detailed instructions on:

- Creating news plugins (Reddit, Medium, etc.)
- Adding widget type plugins (GitHub, JIRA, Slack, etc.)
- Configuring and registering plugins
- Best practices and examples

## Installation

### Quick Setup

1. **Build the application:**
```bash
go build -o goday .
```

2. **Set up configuration:**
```bash
./setup-config.sh
```
This creates `~/.goday/config.yaml` with default settings.

3. **Run GoDay:**
```bash
./goday
```

### Manual Configuration

Your configuration file is located at `~/.goday/config.yaml`. If it doesn't exist, GoDay will create a default one on first run.

**Configuration Directory Structure:**
```
~/.goday/
├── config.yaml      # Main configuration
└── cache/           # Cached data (auto-created)
```

### Configuration File Location

GoDay looks for configuration in this order:
1. `~/.goday/config.yaml` (recommended)
2. `./config.yaml` (fallback for development)

The configuration file supports both address strings and precise coordinates for traffic widgets.

## Installation (Legacy)

If you prefer to manage the config file manually:
```bash
# Clone the repository
git clone <repository-url>
cd goday

# Install dependencies
go mod tidy

# Build the application
go build .

# Run the dashboard
./goday
```

## Configuration

Create a `config.yaml` file in the same directory:

```yaml
user:
  name: "Your Name"
  location: "City,Country"   # for weather API

ui:
  layout: at_a_glance
  min_width: 100
  tile_height: 7

widgets:
  weather:
    ttl: 600s
    api_key: "YOUR_OWM_API_KEY"  # Optional: OpenWeatherMap API key

  news:
    ttl: 600s
    tags: [golang, security, ai, javascript]  # News filter tags

  slack:
    ttl: 20s
  confluence:
    ttl: 300s
  jira:
    ttl: 45s
    log_work: true

# Plugin configurations (optional)
plugins:
  hackernews:
    tags: [golang, security, ai]
    current_tag: "all"
  devto:
    tags: [golang, javascript, webdev]
    current_tag: "all"
  openweathermap:
    api_key: "YOUR_OWM_API_KEY"
    city: "City,Country"
```

## Usage

### Keyboard Shortcuts

- `q` or `Ctrl+C`: Quit the application
- `Tab`/`Shift+Tab`: Navigate between widgets
- `↑↓` or `j/k`: Navigate within a widget
- `Enter`: Open selected item's URL in browser
- `t`: Cycle through news tags
- `T`: Reset news filter to "All"
- `r` or `R`: Refresh all widgets

### Navigation

The dashboard is fully interactive:

1. **Widget Focus**: Use Tab/Shift+Tab to move between widgets (focused widget has a blue border)
2. **Item Selection**: Use arrow keys or j/k to select items within a widget
3. **Open Links**: Press Enter to open the selected item's URL in your default browser
4. **News Filtering**: Press 't' to cycle through news tags, 'T' to reset to "All"

### Weather Setup

To get real weather data, sign up for a free API key at [OpenWeatherMap](https://openweathermap.org/api) and add it to your config:

```yaml
widgets:
  weather:
    api_key: "your_api_key_here"
```

Without an API key, the weather will show placeholder data.

## Architecture

The application follows a modern plugin-based architecture:

```
Plugin Manager ┌── WeatherPlugin ──────┐
              ├── HackerNewsPlugin ───┤
              ├── DevToPlugin ────────┤─► Widget Updates
              ├── AggregateNewsPlugin┤
              └── [Future Plugins]───┘

Plugin Registry ── Schedule Manager ── Widget Manager ── TUI Renderer
```

### Core Components

- **Plugin Manager**: Manages plugin lifecycle, registration, and execution
- **Plugin Registry**: Central registry for all plugins with type-based retrieval  
- **Plugin Scheduler**: Handles periodic plugin execution with configurable intervals
- **News Plugins**: Specialized plugins for news sources with tag filtering
- **Widget Manager**: Converts plugin data to UI-ready widget items

### Extensibility

The plugin system supports:
- **Multiple API sources**: Easily add new data sources
- **Custom widget types**: Create entirely new widget categories
- **Tag-based filtering**: Automatic filtering for categorized content
- **Graceful fallbacks**: Cached data when APIs are unavailable
- **Configurable refresh rates**: Per-plugin TTL settings

## Development

### Prerequisites

- Go 1.21 or later
- Terminal with support for Unicode characters

### Building

```bash
# Format code
go fmt .

# Build application
go build .

# Run tests (when implemented)
go test ./...
```

### Project Structure

```
goday/
├── main.go              # Application entry point and TUI logic
├── config_loader.go     # YAML configuration loading
├── providers.go         # Legacy providers (being phased out)
├── plugins.go           # Core plugin system interfaces
├── news_plugins.go      # News plugin implementations
├── weather_plugins.go   # Weather plugin implementation
├── example_plugins.go   # Example plugins for GitHub, Calendar, etc.
├── widgets.go           # Widget definitions and rendering
├── config.yaml          # User configuration
├── PLUGIN_GUIDE.md      # Detailed plugin development guide
├── go.mod               # Go module dependencies
└── README.md            # This file
```

## 📅 Google Calendar Integration

GoDay supports Gmail account integration via Google Calendar API:

### Quick Setup
1. **Enable Google Calendar API** in Google Cloud Console
2. **Create OAuth 2.0 credentials** (Desktop application)
3. **Download JSON credentials** to `~/.goday/google_calendar_credentials.json`
4. **Run GoDay** and complete OAuth flow
5. **View your events** in the Calendar widget!

### Features
- 🔴 **Live status indicators** (ongoing, starting soon, future)
- 📅 **Smart event display** (today's events show time, future events show date)
- 🔔 **Urgent event notifications** (widget title changes when events are imminent)
- ⚙️ **Configurable refresh** and event limits

### Documentation
- **[GOOGLE_CALENDAR_SETUP.md](GOOGLE_CALENDAR_SETUP.md)** - Complete setup guide
- **[ADDRESS_CONFIGURATION_GUIDE.md](ADDRESS_CONFIGURATION_GUIDE.md)** - Traffic widget configuration

### Configuration
```yaml
widgets:
  calendar:
    ttl: 300s        # Refresh every 5 minutes
    max_events: 10   # Maximum events to show
    days_ahead: 7    # Days ahead to fetch
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Ensure code is formatted with `go fmt`
5. Test the application compiles and runs
6. Submit a pull request

## License

[Add your license here]

## Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI
- Weather data from [OpenWeatherMap](https://openweathermap.org/)
- News data from [Hacker News](https://news.ycombinator.com/) # Git Integration Test
