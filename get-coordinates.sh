#!/bin/bash

echo "🗺️  GoDay Coordinate Helper"
echo "=========================="
echo ""

echo "📍 Common Bangalore Coordinates (copy these for quick setup):"
echo ""
echo "🏠 RESIDENTIAL AREAS:"
echo "   Panathur:        12.9698, 77.7500"
echo "   Bellandur:       12.9279, 77.6271"
echo "   HSR Layout:      12.9082, 77.6476"
echo "   Koramangala:     12.9352, 77.6245"
echo "   Indiranagar:     12.9784, 77.6408"
echo "   Jayanagar:       12.9165, 77.5833"
echo ""
echo "🏢 BUSINESS AREAS:"
echo "   Electronic City: 12.8456, 77.6603"
echo "   Whitefield:      12.9698, 77.7500"
echo "   Marathahalli:    12.9591, 77.7074"
echo "   Outer Ring Road: 12.9591, 77.7074"
echo "   Sarjapur Road:   12.9007, 77.6964"
echo ""
echo "🚇 LANDMARKS:"
echo "   MG Road:         12.9716, 77.5946"
echo "   Airport:         13.1986, 77.7066"
echo "   Cubbon Park:     12.9762, 77.5993"
echo "   UB City Mall:    12.9716, 77.5946"
echo ""

echo "🎯 TO GET EXACT COORDINATES FOR YOUR LOCATIONS:"
echo ""
echo "1. 📱 GOOGLE MAPS METHOD (Most Accurate):"
echo "   • Open https://maps.google.com"
echo "   • Search for your exact location"
echo "   • Right-click on the building/spot"
echo "   • Click the coordinates that appear"
echo "   • Copy the numbers (e.g., '12.9698, 77.7500')"
echo ""
echo "2. 📍 WHAT3WORDS METHOD (Super Precise):"
echo "   • Open https://what3words.com"
echo "   • Search for your location"
echo "   • Copy the coordinates shown"
echo ""
echo "3. 📱 MOBILE APPS:"
echo "   • Use any GPS app that shows coordinates"
echo "   • iPhone: Compass app"
echo "   • Android: GPS Coordinates app"
echo ""

echo "⚙️  EXAMPLE CONFIG:"
echo ""
echo "Once you have coordinates, add them to ~/.goday/config.yaml:"
echo ""
echo "widgets:"
echo "  traffic:"
echo "    ttl: 300s"
echo "    origin:"
echo "      latitude: 12.9698   # Your home coordinates"
echo "      longitude: 77.7500"
echo "      name: \"Home\""
echo "    destination:"
echo "      latitude: 12.9279   # Your office coordinates"
echo "      longitude: 77.6271"
echo "      name: \"Office\""
echo ""

echo "🚀 QUICK SETUP:"
echo ""
echo "Want to use approximate coordinates for your area? Enter the area name:"
read -p "Enter your home area (e.g., Panathur, Bellandur, etc.): " home_area
read -p "Enter your office area (e.g., Bellandur, Koramangala, etc.): " office_area

echo ""
echo "📋 Generated config based on your areas:"
echo ""

# Simple area to coordinate mapping
case ${home_area,,} in
    panathur) home_coords="latitude: 12.9698\n      longitude: 77.7500" ;;
    bellandur) home_coords="latitude: 12.9279\n      longitude: 77.6271" ;;
    koramangala) home_coords="latitude: 12.9352\n      longitude: 77.6245" ;;
    indiranagar) home_coords="latitude: 12.9784\n      longitude: 77.6408" ;;
    whitefield) home_coords="latitude: 12.9698\n      longitude: 77.7500" ;;
    "electronic city") home_coords="latitude: 12.8456\n      longitude: 77.6603" ;;
    marathahalli) home_coords="latitude: 12.9591\n      longitude: 77.7074" ;;
    "hsr layout") home_coords="latitude: 12.9082\n      longitude: 77.6476" ;;
    *) home_coords="latitude: 12.9698\n      longitude: 77.7500  # Update with exact coordinates" ;;
esac

case ${office_area,,} in
    panathur) office_coords="latitude: 12.9698\n      longitude: 77.7500" ;;
    bellandur) office_coords="latitude: 12.9279\n      longitude: 77.6271" ;;
    koramangala) office_coords="latitude: 12.9352\n      longitude: 77.6245" ;;
    indiranagar) office_coords="latitude: 12.9784\n      longitude: 77.6408" ;;
    whitefield) office_coords="latitude: 12.9698\n      longitude: 77.7500" ;;
    "electronic city") office_coords="latitude: 12.8456\n      longitude: 77.6603" ;;
    marathahalli) office_coords="latitude: 12.9591\n      longitude: 77.7074" ;;
    "hsr layout") office_coords="latitude: 12.9082\n      longitude: 77.6476" ;;
    *) office_coords="latitude: 12.9279\n      longitude: 77.6271  # Update with exact coordinates" ;;
esac

echo "widgets:"
echo "  traffic:"
echo "    ttl: 300s"
echo "    origin:"
echo -e "      $home_coords"
echo "      name: \"$home_area\""
echo "    destination:"
echo -e "      $office_coords"
echo "      name: \"$office_area\""
echo ""

echo "💡 To use this config:"
echo "1. Copy the above configuration"
echo "2. Edit ~/.goday/config.yaml"
echo "3. Replace the traffic section with the above"
echo "4. For exact coordinates, use Google Maps method described above"
echo ""
echo "🎯 This eliminates geocoding errors and gives you instant, accurate routes!"
