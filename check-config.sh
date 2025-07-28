#!/bin/bash

echo "🔍 Checking GoDay config creation..."

# Check if ~/.goday exists
if [ -d "$HOME/.goday" ]; then
    echo "✅ ~/.goday directory exists"
    ls -la "$HOME/.goday"
else
    echo "❌ ~/.goday directory does not exist"
    echo "💡 This will be created when you first run GoDay"
fi

echo ""
echo "📍 Expected config location: $HOME/.goday/config.yaml"

# Check if config file exists
if [ -f "$HOME/.goday/config.yaml" ]; then
    echo "✅ Config file exists"
    echo "📄 First few lines of config:"
    head -10 "$HOME/.goday/config.yaml"
else
    echo "❌ Config file does not exist"
    echo "💡 Run './goday' to create it automatically"
fi

echo ""
echo "🚀 To create config automatically:"
echo "   ./goday"
echo ""
echo "🔧 To see config location:"
echo "   ./goday config"
