# Traffic Widget - Bidirectional Display

## New Feature: Both Directions Simultaneously

The traffic widget now shows both directions at once, giving you a complete view of your commute without needing to toggle.

## What You'll See

Instead of a single route with toggle functionality, you'll now see:

```
[Traffic]
Electronic City → Whitefield
42 min • 25.1 km • 🟡 Moderate

Whitefield → Electronic City  
38 min • 25.1 km • 🟢 Light
```

## Benefits

✅ **Complete Picture**: See both directions at once  
✅ **Compare Commutes**: Instantly compare morning vs evening commute times  
✅ **No Toggling**: No need to press keys to switch directions  
✅ **Better Planning**: Make informed decisions about when to leave  

## Traffic Indicators

- **🟢 Light**: < 30 minutes travel time
- **🟡 Moderate**: 30-60 minutes travel time  
- **🔴 Heavy**: > 60 minutes travel time

## Real-World Example for Bangalore

**Morning (9 AM)**:
```
Electronic City → Whitefield
52 min • 25.4 km • 🔴 Heavy

Whitefield → Electronic City
35 min • 25.4 km • 🟢 Light
```

**Evening (6 PM)**:
```
Electronic City → Whitefield  
38 min • 25.4 km • 🟡 Moderate

Whitefield → Electronic City
58 min • 25.4 km • 🔴 Heavy
```

This gives you immediate insight into traffic patterns - typically one direction is worse than the other depending on the time of day.

## Configuration

No changes needed to your existing config:

```yaml
widgets:
  traffic:
    ttl: 300s  # 5 minutes refresh
    origin: "Electronic City, Bengaluru, Karnataka, India"
    destination: "Whitefield, Bengaluru, Karnataka, India"
```

## Navigation

- **Tab/Shift+Tab**: Move focus to traffic widget
- **↑↓ or j/k**: Navigate between the two route items
- **r/R**: Refresh traffic data
- **Enter**: (Shows route details in bottom bar)

## What Changed

- **Removed**: Direction toggle (d/D key)
- **Added**: Bidirectional data fetching
- **Improved**: Better visual separation of routes
- **Enhanced**: Immediate comparison capabilities

## Performance

- **API Calls**: Now makes 2 API calls (one for each direction)
- **Refresh Rate**: Still 5 minutes by default
- **Free Service**: Still uses OSRM (no API key required)
- **Speed**: Minimal impact as calls are made concurrently

## Future Enhancements

Potential improvements:
- **Time-based highlighting**: Highlight the direction you're likely to take based on time
- **Historical comparison**: Show how current times compare to typical times
- **Smart notifications**: Alert when traffic is unusually heavy/light
- **Calendar integration**: Suggest departure times based on meetings
