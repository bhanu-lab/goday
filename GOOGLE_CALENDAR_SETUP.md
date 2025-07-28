# Google Calendar Integration Guide

GoDay now supports Google Calendar integration to display your upcoming events directly in the dashboard! 📅

## 🚀 **Quick Setup**

### **Step 1: Enable Google Calendar API**

1. **Go to Google Cloud Console**
   - Visit [https://console.cloud.google.com/](https://console.cloud.google.com/)
   - Sign in with your Google account

2. **Create/Select a Project**
   - Create a new project or select an existing one
   - Name it something like "GoDay Dashboard"

3. **Enable Google Calendar API**
   - Go to "APIs & Services" > "Library"
   - Search for "Google Calendar API"
   - Click on it and press "Enable"

### **Step 2: Create OAuth Credentials**

1. **Go to Credentials**
   - Navigate to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth client ID"

2. **Configure OAuth Consent Screen** (if first time)
   - Choose "External" user type
   - Fill in required fields:
     - App name: "GoDay Dashboard"
     - User support email: your email
     - Developer email: your email
   - Save and continue through the steps

3. **Create OAuth Client ID**
   - Application type: "Desktop application"
   - Name: "GoDay Calendar"
   - Click "Create"

4. **Download Credentials**
   - Click the download button (JSON icon)
   - Save the file as `google_calendar_credentials.json`

### **Step 3: Install Credentials**

1. **Create GoDay directory** (if it doesn't exist):
   ```bash
   mkdir -p ~/.goday
   ```

2. **Copy credentials file**:
   ```bash
   cp ~/Downloads/google_calendar_credentials.json ~/.goday/
   ```

### **Step 4: Update Configuration**

Edit your `~/.goday/config.yaml` file to include calendar settings:

```yaml
widgets:
  calendar:
    ttl: 300s        # Refresh every 5 minutes
    max_events: 10   # Maximum events to show
    days_ahead: 7    # Days ahead to fetch events
    # credentials_file and token_file are auto-configured
```

### **Step 5: Run GoDay**

1. **Start GoDay**:
   ```bash
   ./goday
   ```

2. **First-time OAuth Flow**:
   - GoDay will display a URL in the console
   - Open the URL in your browser
   - Sign in to Google and authorize GoDay
   - Copy the authorization code and paste it in the terminal
   - GoDay will save the token automatically

3. **Enjoy Your Calendar** 📅
   - Events will appear in the Calendar widget
   - Shows today's events with times
   - Future events with dates
   - Status indicators for urgent events

## 📋 **Calendar Widget Features**

### **Event Display**
- **Today's Events**: Show time (e.g., "9:00-10:00")
- **Future Events**: Show date and time (e.g., "Jan 28 14:00")
- **All-day Events**: Show "All day"

### **Status Indicators**
- 🔴 **Currently happening** (event is ongoing)
- 🟡 **Starting soon** (within 30 minutes)
- 🟢 **Future event** (more than 30 minutes away)

### **Smart Notifications**
- Calendar widget title shows 🔔 when you have urgent events
- Shows up to 5 most relevant events
- Prioritizes today's and upcoming events

## ⚙️ **Configuration Options**

```yaml
widgets:
  calendar:
    ttl: 300s                    # How often to refresh (default: 5 minutes)
    max_events: 10               # Max events to fetch (default: 10)
    days_ahead: 7                # Days ahead to look (default: 7)
    credentials_file: "custom/path/credentials.json"  # Optional: custom path
    token_file: "custom/path/token.json"              # Optional: custom path
```

## 🔧 **Troubleshooting**

### **Problem: "Calendar Setup Required"**
**Solution**: 
1. Check that `~/.goday/google_calendar_credentials.json` exists
2. Make sure the file contains valid JSON from Google Cloud Console
3. Restart GoDay

### **Problem: "Calendar unavailable" or OAuth errors**
**Solution**:
1. Delete the token file: `rm ~/.goday/google_calendar_token.json`
2. Restart GoDay to redo the OAuth flow
3. Make sure you authorize the correct Google account

### **Problem: No events showing**
**Solution**:
1. Check that you have events in your Google Calendar
2. Verify the `days_ahead` setting (default: 7 days)
3. Check `max_events` setting (default: 10)

### **Problem: "Plugin not initialized"**
**Solution**:
1. Verify your credentials file is valid JSON
2. Check Google Cloud Console that Calendar API is enabled
3. Ensure OAuth consent screen is properly configured

## 🔒 **Privacy & Security**

- **OAuth 2.0**: Uses Google's secure OAuth 2.0 flow
- **Read-Only Access**: GoDay only reads your calendar, never modifies it
- **Local Storage**: All tokens stored locally in `~/.goday/`
- **No Data Sharing**: Your calendar data never leaves your machine

## 📁 **File Locations**

```
~/.goday/
├── config.yaml                          # Main configuration
├── google_calendar_credentials.json     # OAuth app credentials (from Google)
└── google_calendar_token.json          # OAuth tokens (auto-generated)
```

## 🎯 **Advanced Configuration**

### **Multiple Calendars**
Currently supports primary calendar only. Multiple calendar support coming soon!

### **Custom Time Ranges**
```yaml
widgets:
  calendar:
    days_ahead: 14    # Look 2 weeks ahead
    max_events: 20    # Show more events
```

### **Custom Refresh Rates**
```yaml
widgets:
  calendar:
    ttl: 60s     # Refresh every minute (for busy schedules)
    # or
    ttl: 600s    # Refresh every 10 minutes (for lighter use)
```

## 🚨 **Common Setup Issues**

### **"Project not found" or 403 errors**
- Make sure you're signed in to the correct Google account
- Verify the project has Google Calendar API enabled
- Check that OAuth consent screen is published (not in testing mode)

### **"Invalid client" errors**
- Re-download credentials from Google Cloud Console
- Make sure you selected "Desktop application" not "Web application"
- Verify the credentials file is valid JSON

### **Authorization code flow not working**
- Try using an incognito/private browser window
- Make sure to copy the ENTIRE authorization code
- Check that your browser can access the authorization URL

## 🎉 **You're All Set!**

Once configured, your Google Calendar events will automatically appear in GoDay's Calendar widget. The widget refreshes every 5 minutes by default and shows your most relevant upcoming events.

**Pro Tip**: Use the calendar widget alongside the traffic widget to plan your commute around your meetings! 🚗📅

Need help? Check the console output when starting GoDay for detailed error messages and setup instructions.
