# Traffic Widget

The Traffic widget provides real-time commute information between two locations using Google Maps Distance Matrix API. This is particularly useful for Bangalore traffic where commute times can vary significantly.

## Features

- **Real-time traffic data**: Shows current travel time considering traffic conditions
- **Bidirectional routing**: Toggle between origin→destination and destination→origin
- **Traffic indicators**: Color-coded traffic levels (🟢 Light, 🟡 Moderate, 🔴 Heavy)
- **Distance and duration**: Shows both time and distance for the route
- **Auto-refresh**: Configurable refresh interval (default: 5 minutes)

## Configuration

Add the following to your `config.yaml`:

```yaml
widgets:
  traffic:
    ttl: 300s  # 5 minutes refresh interval
    api_key: "YOUR_GOOGLE_MAPS_API_KEY"
    origin: "Electronic City, Bengaluru, Karnataka, India"
    destination: "Whitefield, Bengaluru, Karnataka, India"
```

## Google Maps API Setup

1. **Get API Key**:
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing one
   - Enable "Distance Matrix API" and "Places API"
   - Create credentials → API Key
   - Restrict the key to Distance Matrix API for security

2. **API Pricing**:
   - Distance Matrix API: $5 per 1000 requests
   - With 5-minute refresh: ~288 requests/day = ~$0.14/day
   - Consider setting usage limits in Google Cloud Console

## Usage

### Navigation
- **Tab/Shift+Tab**: Move focus to Traffic widget
- **d/D**: Toggle traffic direction (origin ↔ destination)
- **r/R**: Refresh traffic data

### Widget Display
The traffic widget shows:
```
[Traffic]
Electronic City → Whitefield
45 mins • 25.4 km • 🟡 Moderate
```

When direction is toggled:
```
[Traffic]
Electronic City ← Whitefield
52 mins • 25.4 km • 🔴 Heavy
```

## Traffic Indicators

- **🟢 Light**: < 30 minutes travel time
- **🟡 Moderate**: 30-60 minutes travel time
- **🔴 Heavy**: > 60 minutes travel time

## Common Bangalore Routes

Popular route configurations for Bangalore:

```yaml
# Electronic City to Whitefield
origin: "Electronic City, Bengaluru, Karnataka, India"
destination: "Whitefield, Bengaluru, Karnataka, India"

# Koramangala to HSR Layout
origin: "Koramangala, Bengaluru, Karnataka, India"
destination: "HSR Layout, Bengaluru, Karnataka, India"

# Indiranagar to Brigade Road
origin: "Indiranagar, Bengaluru, Karnataka, India"
destination: "Brigade Road, Bengaluru, Karnataka, India"

# Marathahalli to Electronic City
origin: "Marathahalli, Bengaluru, Karnataka, India"
destination: "Electronic City, Bengaluru, Karnataka, India"
```

## Troubleshooting

### Common Issues

1. **"Traffic unavailable"**:
   - Check API key is valid and has Distance Matrix API enabled
   - Verify locations are formatted correctly
   - Check internet connectivity

2. **High API costs**:
   - Increase TTL value (e.g., 600s for 10-minute refresh)
   - Set usage limits in Google Cloud Console
   - Consider enabling billing alerts

3. **Invalid locations**:
   - Use full addresses with city and state
   - Test locations in Google Maps first
   - Use landmark names that Google recognizes

### Error Messages

- **"API key not configured"**: Add your Google Maps API key to config.yaml
- **"Plugin not found"**: Check that traffic plugin is properly registered
- **"No route data available"**: Verify origin and destination are valid locations
- **"API error: [status]"**: Check API key permissions and quotas

## Integration with Calendar

You can combine the traffic widget with calendar events for smart commute planning:

1. Check traffic before leaving for meetings
2. Use direction toggle to plan return journey
3. Factor in traffic when scheduling meetings

## Future Enhancements

Potential improvements for the traffic widget:

- [ ] Multiple route options
- [ ] Historical traffic patterns
- [ ] Integration with calendar for automatic route suggestions
- [ ] Toll road information
- [ ] Public transportation options
- [ ] Weather impact on traffic
- [ ] Customizable traffic threshold colors
