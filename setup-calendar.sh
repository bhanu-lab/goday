#!/bin/bash

# GoDay Calendar Setup Helper Script
echo "🗓️  GoDay Google Calendar Setup Helper"
echo "======================================"
echo ""

# Check if goday directory exists
GODAY_DIR="$HOME/.goday"
if [ ! -d "$GODAY_DIR" ]; then
    echo "📁 Creating GoDay config directory: $GODAY_DIR"
    mkdir -p "$GODAY_DIR"
else
    echo "✅ GoDay config directory exists: $GODAY_DIR"
fi

# Check for credentials file
CREDENTIALS_FILE="$GODAY_DIR/google_calendar_credentials.json"
if [ ! -f "$CREDENTIALS_FILE" ]; then
    echo ""
    echo "❌ Google Calendar credentials not found!"
    echo ""
    echo "📋 To set up Google Calendar integration:"
    echo ""
    echo "1. 🌐 Go to Google Cloud Console:"
    echo "   https://console.cloud.google.com/"
    echo ""
    echo "2. 🔧 Create/select a project and enable Google Calendar API"
    echo ""
    echo "3. 🔑 Create OAuth 2.0 credentials (Desktop application)"
    echo ""
    echo "4. 💾 Download the JSON file and save it as:"
    echo "   $CREDENTIALS_FILE"
    echo ""
    echo "5. 🚀 Run GoDay again to complete OAuth setup"
    echo ""
    echo "📖 For detailed instructions, see: GOOGLE_CALENDAR_SETUP.md"
    exit 1
else
    echo "✅ Google Calendar credentials found: $CREDENTIALS_FILE"
fi

# Check for token file
TOKEN_FILE="$GODAY_DIR/google_calendar_token.json"
if [ ! -f "$TOKEN_FILE" ]; then
    echo "⚠️  OAuth token not found (normal for first run)"
    echo "   GoDay will guide you through OAuth setup on first start"
else
    echo "✅ OAuth token exists: $TOKEN_FILE"
fi

# Check config file
CONFIG_FILE="$GODAY_DIR/config.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "⚠️  Config file not found"
    echo "   GoDay will create a default config on first run"
else
    echo "✅ Config file exists: $CONFIG_FILE"
    
    # Check if calendar is configured
    if grep -q "calendar:" "$CONFIG_FILE"; then
        echo "✅ Calendar widget is configured in config"
    else
        echo "⚠️  Calendar widget not configured yet"
        echo "   Add calendar section to config.yaml or restart GoDay to auto-add"
    fi
fi

echo ""
echo "🎉 Setup check complete!"
echo ""
echo "💡 Next steps:"
echo "   1. Make sure credentials are in place"
echo "   2. Run './goday' to start the dashboard"
echo "   3. Complete OAuth flow if needed"
echo ""
echo "🔧 Need help? Check GOOGLE_CALENDAR_SETUP.md for detailed instructions"
