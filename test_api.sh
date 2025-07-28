#!/bin/bash

# Test Google Maps API Key
# Replace YOUR_API_KEY with your actual API key

API_KEY="PASTE_YOUR_ACTUAL_GOOGLE_MAPS_API_KEY_HERE"
ORIGIN="Electronic City, Bengaluru, Karnataka, India"
DESTINATION="Whitefield, Bengaluru, Karnataka, India"

echo "Testing Google Maps Distance Matrix API..."
echo "Origin: $ORIGIN"
echo "Destination: $DESTINATION"
echo ""

URL="https://maps.googleapis.com/maps/api/distancematrix/json?origins=${ORIGIN}&destinations=${DESTINATION}&departure_time=now&traffic_model=best_guess&key=${API_KEY}"

echo "Making API request..."
curl -s "$URL" | jq '.' || curl -s "$URL"

echo ""
echo "If you see 'REQUEST_DENIED', check:"
echo "1. API key is valid"
echo "2. Distance Matrix API is enabled in Google Cloud Console"
echo "3. No IP/referrer restrictions on the API key"
echo "4. Billing is enabled for your Google Cloud project"
