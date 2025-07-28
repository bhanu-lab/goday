#!/bin/bash

# GoDay Configuration Setup Script

echo "🚀 Setting up GoDay configuration..."

# Get the config directory
CONFIG_DIR="$HOME/.goday"
CONFIG_FILE="$CONFIG_DIR/config.yaml"

# Create config directory if it doesn't exist
if [ ! -d "$CONFIG_DIR" ]; then
    echo "📁 Creating config directory: $CONFIG_DIR"
    mkdir -p "$CONFIG_DIR"
fi

# Create default config if it doesn't exist
if [ ! -f "$CONFIG_FILE" ]; then
    echo "📝 Creating default config file: $CONFIG_FILE"
    cat > "$CONFIG_FILE" << 'EOF'
user:
  name: "Your Name"
  location: "Bengaluru,IN"

ui:
  layout: at_a_glance
  min_width: 100
  tile_height: 7

widgets:
  weather:
    ttl: 600s
    api_key: "YOUR_OWM_API_KEY"
  news:
    ttl: 600s
    tags: [golang, security, ai]
    provider: hn
  slack:
    ttl: 20s
  confluence:
    ttl: 300s
  jira:
    ttl: 45s
    log_work: true
  traffic:
    ttl: 300s
    # Option 1: Use addresses
    origin: "Electronic City Phase 1, Bengaluru, Karnataka, India"
    destination: "Whitefield, Bengaluru, Karnataka, India"
    
    # Option 2: Use coordinates (comment out above and uncomment below)
    # origin:
    #   latitude: 12.8456
    #   longitude: 77.6603
    #   name: "Electronic City"
    # destination:
    #   latitude: 12.9698
    #   longitude: 77.7500
    #   name: "Whitefield"
EOF
    echo "✅ Default config created!"
else
    echo "✅ Config file already exists: $CONFIG_FILE"
fi

echo ""
echo "📍 Your GoDay config is located at: $CONFIG_FILE"
echo "📝 Edit this file to customize your dashboard"
echo ""
echo "🔧 To edit your config:"
echo "   nano $CONFIG_FILE"
echo "   # or"
echo "   code $CONFIG_FILE"
echo ""
echo "🚀 Run GoDay with: ./goday"
