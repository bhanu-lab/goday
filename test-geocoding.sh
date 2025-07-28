#!/bin/bash

echo "🔍 Testing Address Geocoding..."

# Test addresses with different formats
addresses=(
    "Panathur, Bengaluru, Karnataka, India"
    "Bellandur, Bengaluru, Karnataka, India"
    "Electronic City, Bengaluru, Karnataka, India"
    "Whitefield, Bengaluru, Karnataka, India"
    "Koramangala, Bengaluru, Karnataka, India"
    "Indiranagar, Bengaluru, Karnataka, India"
    "HSR Layout, Bengaluru, Karnataka, India"
    "Marathahalli, Bengaluru, Karnataka, India"
)

echo ""
echo "📍 Testing common Bangalore locations..."

for address in "${addresses[@]}"; do
    echo ""
    echo "🔍 Testing: $address"
    
    # URL encode the address
    encoded_address=$(echo "$address" | sed 's/ /%20/g' | sed 's/,/%2C/g')
    
    # Test with Nominatim
    url="https://nominatim.openstreetmap.org/search?q=${encoded_address}&format=json&limit=1&countrycodes=in"
    
    # Use curl to test the geocoding
    result=$(curl -s -H "User-Agent: GoDay-Dashboard/1.0" "$url")
    
    if echo "$result" | grep -q '"lat"'; then
        lat=$(echo "$result" | grep -o '"lat":"[^"]*"' | cut -d'"' -f4)
        lon=$(echo "$result" | grep -o '"lon":"[^"]*"' | cut -d'"' -f4)
        echo "✅ Found: $lat, $lon"
    else
        echo "❌ Not found"
    fi
    
    # Rate limiting - be nice to Nominatim
    sleep 1
done

echo ""
echo "🎯 Recommendations:"
echo "1. Use coordinates for exact locations"
echo "2. Use general area names (Panathur, Bellandur) instead of specific building names"
echo "3. Always include 'Bengaluru, Karnataka, India' for better results"
