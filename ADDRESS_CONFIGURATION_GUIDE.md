# Traffic Widget - Address Configuration Guide

## Supported Address Types

The traffic widget can handle various types of addresses, from general areas to specific street addresses. Here are examples of what works well:

## ✅ **Recommended Address Formats**

### 1. **Latitude/Longitude Coordinates** 🆕
```yaml
# Direct coordinates (most precise)
traffic:
  origin: 
    latitude: 12.9716
    longitude: 77.5946
    name: "UB City Mall"  # Optional display name
  destination:
    latitude: 12.9698
    longitude: 77.7500
    name: "Whitefield"

# Mix coordinates and address
traffic:
  origin: "Manyata Tech Park, Thanisandra Main Road, Bengaluru, Karnataka 560045, India"
  destination:
    latitude: 12.9698
    longitude: 77.7500
    name: "Custom Location"
```

### 2. **Specific Buildings/Complexes**
```yaml
# Tech Parks and Office Complexes
origin: "Manyata Tech Park, Thanisandra Main Road, Bengaluru, Karnataka 560045, India"
destination: "Brigade Metropolis, Whitefield Main Road, Bengaluru, Karnataka 560066, India"

# Residential Complexes
origin: "Prestige Shantiniketan, Whitefield, Bengaluru, Karnataka 560066, India"
destination: "Sobha City, Thanisandra, Bengaluru, Karnataka 560077, India"

# Shopping Malls
origin: "Phoenix MarketCity, Whitefield Road, Bengaluru, Karnataka 560048, India"
destination: "UB City Mall, Vittal Mallya Road, Bengaluru, Karnataka 560001, India"
```

### 3. **Street Addresses**
```yaml
# Complete street addresses
origin: "123 Whitefield Main Road, EPIP Zone, Whitefield, Bengaluru, Karnataka 560066"
destination: "45 Outer Ring Road, Marathahalli, Bengaluru, Karnataka 560037"

# With landmarks
origin: "Near Whitefield Railway Station, Whitefield, Bengaluru, Karnataka"
destination: "Opposite Brigade Gateway, Dr Rajkumar Road, Bengaluru, Karnataka"
```

### 4. **Major Landmarks**
```yaml
# Well-known places
origin: "Bangalore International Airport, Devanahalli, Bengaluru, Karnataka"
destination: "Cubbon Park, Bengaluru, Karnataka"

# Metro Stations
origin: "Whitefield Metro Station, Bengaluru, Karnataka"
destination: "MG Road Metro Station, Bengaluru, Karnataka"

# Universities
origin: "Indian Institute of Science, Bengaluru, Karnataka"
destination: "Christ University, Hosur Road, Bengaluru, Karnataka"
```

### 5. **Popular Bangalore Locations**
```yaml
# Business Districts
origin: "Electronic City Phase 1, Bengaluru, Karnataka"
destination: "Koramangala 5th Block, Bengaluru, Karnataka"

# IT Corridors
origin: "Outer Ring Road, Marathahalli, Bengaluru, Karnataka"
destination: "Sarjapur Road, Bengaluru, Karnataka"

# Residential Areas
origin: "Indiranagar, Bengaluru, Karnataka"
destination: "HSR Layout, Bengaluru, Karnataka"
```

## 🎯 **Best Practices**

### **Why Use Coordinates?** 🎯
- **Precision**: Exact location, no geocoding errors
- **Speed**: No geocoding delay, faster route calculation  
- **Reliability**: Always works, even for new/unlisted addresses
- **Custom Names**: Use any display name you want

### **When to Use Coordinates:**
- **New buildings** not yet in map databases
- **Precise locations** within large complexes
- **Custom points** like parking entrances
- **Performance** critical applications

### **Include These Elements:**
1. **Building/Area Name**: "Manyata Tech Park"
2. **Street/Road**: "Thanisandra Main Road"
3. **City**: "Bengaluru"
4. **State**: "Karnataka"
5. **Pincode** (optional but helpful): "560045"
6. **Country**: "India"

### **Address Format Template:**
```
[Building/Complex Name], [Street/Road Name], [Area], [City], [State] [Pincode], [Country]
```

## 📍 **Real Bangalore Examples**

### **Common Commute Routes:**

#### **Electronic City ↔ Whitefield**
```yaml
origin: "Electronic City Phase 1, Hosur Road, Bengaluru, Karnataka 560100, India"
destination: "ITPL, Whitefield, Bengaluru, Karnataka 560066, India"
```

#### **Koramangala ↔ Marathahalli**
```yaml
origin: "Koramangala 5th Block, Bengaluru, Karnataka 560095, India"
destination: "Marathahalli Bridge, Outer Ring Road, Bengaluru, Karnataka 560037, India"
```

#### **Indiranagar ↔ HSR Layout**
```yaml
origin: "Indiranagar 100 Feet Road, Bengaluru, Karnataka 560038, India"
destination: "HSR Layout Sector 1, Bengaluru, Karnataka 560102, India"
```

#### **Airport ↔ City Center**
```yaml
origin: "Kempegowda International Airport, Devanahalli, Bengaluru, Karnataka 560300, India"
destination: "MG Road, Bengaluru, Karnataka 560001, India"
```

## 🔍 **Testing Your Addresses**

Before using addresses in the config, test them:

1. **Search on OpenStreetMap**: Go to [openstreetmap.org](https://openstreetmap.org) and search
2. **Check Google Maps**: Verify the location exists and is correctly named
3. **Use Nominatim directly**: Test at [nominatim.openstreetmap.org](https://nominatim.openstreetmap.org)

## ⚠️ **Addresses to Avoid**

### **Too Vague:**
```yaml
# Don't use these - too general
origin: "Bengaluru"
destination: "Whitefield"
```

### **Non-existent Places:**
```yaml
# Don't use made-up or very new addresses
origin: "Random Building, Unknown Road, Bengaluru"
```

### **Incomplete Information:**
```yaml
# Missing important details
origin: "Tech Park"  # Which tech park?
destination: "Main Road"  # Which main road?
```

## 🚀 **Pro Tips**

### **For Better Accuracy:**
1. **Include pincode** when possible
2. **Use official building/complex names**
3. **Add major road names** (e.g., "Outer Ring Road", "Sarjapur Road")
4. **Include landmarks** when addresses are unclear

### **For Display Names:**
The widget automatically extracts short names for display:
```
"Manyata Tech Park, Thanisandra Main Road, Bengaluru..." 
→ Displays as: "Manyata Tech Park"
```

### **Common Bangalore Roads:**
- Outer Ring Road (ORR)
- Hosur Road
- Bannerghatta Road
- Whitefield Main Road
- Sarjapur Road
- Electronic City Elevated Expressway

## 🔧 **Configuration Example**

Here's a complete example for a typical Bangalore commute:

```yaml
widgets:
  traffic:
    ttl: 300s  # 5 minutes refresh
    origin: "Prestige Tech Park, Sarjapur-Marathahalli Road, Bengaluru, Karnataka 560103, India"
    destination: "RMZ Infinity, Old Madras Road, Bengaluru, Karnataka 560016, India"
```

### **Alternative: Using Coordinates**

For maximum precision, you can use latitude/longitude coordinates:

```yaml
widgets:
  traffic:
    ttl: 300s
    origin:
      latitude: 12.9352
      longitude: 77.6245
      name: "Home"  # Custom display name
    destination:
      latitude: 12.9698
      longitude: 77.7500
      name: "Office"
```

### **Mixed Configuration**

You can mix address strings and coordinates:

```yaml
widgets:
  traffic:
    ttl: 300s
    origin: "Bangalore International Airport, Devanahalli, Bengaluru, Karnataka"
    destination:
      latitude: 12.9716
      longitude: 77.5946
      name: "City Center"
```

## 📍 **Getting Coordinates**

To find latitude/longitude for any location:

1. **Google Maps**: Right-click → "What's here?" → Copy coordinates
2. **OpenStreetMap**: Click location → See coordinates in URL
3. **GPS Apps**: Use any GPS app to get precise coordinates
4. **Online Tools**: Use tools like [latlong.net](https://www.latlong.net/)

### **Common Bangalore Coordinates:**
```yaml
# Major landmarks with precise coordinates
bangalore_airport:
  latitude: 13.1986
  longitude: 77.7066

ub_city_mall:
  latitude: 12.9716
  longitude: 77.5946

electronic_city:
  latitude: 12.8456
  longitude: 77.6603

whitefield:
  latitude: 12.9698
  longitude: 77.7500

koramangala:
  latitude: 12.9279
  longitude: 77.6271
```

## 🔧 **Automatic Setup**

GoDay automatically handles configuration setup:

1. **First Run**: If `~/.goday/` doesn't exist, GoDay creates it
2. **Default Config**: If `config.yaml` doesn't exist, GoDay creates a default one
3. **User Feedback**: Shows where files are created
4. **Ready to Use**: Works immediately with default Bangalore locations

### **What GoDay Creates Automatically:**
```
~/.goday/                    # Config directory (created automatically)
├── config.yaml             # Default config (created automatically)
└── cache/                  # Cache directory (created as needed)
```

### **First Run Experience:**
```bash
$ ./goday
📁 Created config directory: /Users/yourname/.goday
📝 Created default config: /Users/yourname/.goday/config.yaml
💡 Edit the config file to customize your dashboard

[Dashboard starts with default settings]
```

The more specific your addresses, the more accurate your traffic estimates will be!

## 🚨 **Troubleshooting "Location Not Found" Errors**

If you're getting "location not found" errors, try these solutions:

### **1. Use Coordinates (Most Reliable)**
```yaml
traffic:
  origin:
    latitude: 12.9698
    longitude: 77.7500
    name: "Your Home Area"
  destination:
    latitude: 12.9279
    longitude: 77.6271
    name: "Your Office Area"
```

### **2. Simplify Address Names**
Instead of specific building names, use general area names:

❌ **Too Specific (may fail):**
```yaml
origin: "Disha Park West, Balagere Road, Panathur 560087"
destination: "Jfrog India Pvt Ltd, Green Glen Layout, Bellandur"
```

✅ **Better (more likely to work):**
```yaml
origin: "Panathur, Bengaluru, Karnataka, India"
destination: "Bellandur, Bengaluru, Karnataka, India"
```

### **3. Test Your Addresses First**
Before using addresses in config, test them:

```bash
# Test with OpenStreetMap Nominatim
curl -H "User-Agent: GoDay-Dashboard/1.0" \
  "https://nominatim.openstreetmap.org/search?q=Panathur,%20Bengaluru&format=json&limit=1&countrycodes=in"
```

### **4. Common Working Addresses**
These addresses typically work well:

```yaml
# Major Areas
origin: "Electronic City, Bengaluru, Karnataka, India"
destination: "Whitefield, Bengaluru, Karnataka, India"

# Metro Stations
origin: "MG Road Metro Station, Bengaluru, Karnataka, India"
destination: "Indiranagar Metro Station, Bengaluru, Karnataka, India"

# Well-known Places
origin: "Bangalore International Airport, Bengaluru, Karnataka, India"
destination: "Cubbon Park, Bengaluru, Karnataka, India"
```

### **5. Enable Debug Mode**
The traffic plugin will show you which geocoding attempts are failing:
```
Geocoding attempt 1 failed for 'Your Address': no results found
```

### **6. Quick Fix: Use Area Coordinates**
Here are coordinates for common Bangalore areas:

```yaml
# Copy these coordinates for your areas
panathur:
  latitude: 12.9698
  longitude: 77.7500

bellandur:
  latitude: 12.9279
  longitude: 77.6271

electronic_city:
  latitude: 12.8456
  longitude: 77.6603

whitefield:
  latitude: 12.9698
  longitude: 77.7500

koramangala:
  latitude: 12.9352
  longitude: 77.6245

indiranagar:
  latitude: 12.9784
  longitude: 77.6408

hsr_layout:
  latitude: 12.9082
  longitude: 77.6476

marathahalli:
  latitude: 12.9591
  longitude: 77.7074
```

# Traffic Widget - Address Configuration Guide

## Supported Address Types

The traffic widget can handle various types of addresses, from general areas to specific street addresses. Here are examples of what works well:

## ✅ **Recommended Address Formats**

### 1. **Latitude/Longitude Coordinates** 🆕
```yaml
# Direct coordinates (most precise)
traffic:
  origin: 
    latitude: 12.9716
    longitude: 77.5946
    name: "UB City Mall"  # Optional display name
  destination:
    latitude: 12.9698
    longitude: 77.7500
    name: "Whitefield"

# Mix coordinates and address
traffic:
  origin: "Manyata Tech Park, Thanisandra Main Road, Bengaluru, Karnataka 560045, India"
  destination:
    latitude: 12.9698
    longitude: 77.7500
    name: "Custom Location"
```

### 2. **Specific Buildings/Complexes**
```yaml
# Tech Parks and Office Complexes
origin: "Manyata Tech Park, Thanisandra Main Road, Bengaluru, Karnataka 560045, India"
destination: "Brigade Metropolis, Whitefield Main Road, Bengaluru, Karnataka 560066, India"

# Residential Complexes
origin: "Prestige Shantiniketan, Whitefield, Bengaluru, Karnataka 560066, India"
destination: "Sobha City, Thanisandra, Bengaluru, Karnataka 560077, India"

# Shopping Malls
origin: "Phoenix MarketCity, Whitefield Road, Bengaluru, Karnataka 560048, India"
destination: "UB City Mall, Vittal Mallya Road, Bengaluru, Karnataka 560001, India"
```

### 3. **Street Addresses**
```yaml
# Complete street addresses
origin: "123 Whitefield Main Road, EPIP Zone, Whitefield, Bengaluru, Karnataka 560066"
destination: "45 Outer Ring Road, Marathahalli, Bengaluru, Karnataka 560037"

# With landmarks
origin: "Near Whitefield Railway Station, Whitefield, Bengaluru, Karnataka"
destination: "Opposite Brigade Gateway, Dr Rajkumar Road, Bengaluru, Karnataka"
```

### 4. **Major Landmarks**
```yaml
# Well-known places
origin: "Bangalore International Airport, Devanahalli, Bengaluru, Karnataka"
destination: "Cubbon Park, Bengaluru, Karnataka"

# Metro Stations
origin: "Whitefield Metro Station, Bengaluru, Karnataka"
destination: "MG Road Metro Station, Bengaluru, Karnataka"

# Universities
origin: "Indian Institute of Science, Bengaluru, Karnataka"
destination: "Christ University, Hosur Road, Bengaluru, Karnataka"
```

### 5. **Popular Bangalore Locations**
```yaml
# Business Districts
origin: "Electronic City Phase 1, Bengaluru, Karnataka"
destination: "Koramangala 5th Block, Bengaluru, Karnataka"

# IT Corridors
origin: "Outer Ring Road, Marathahalli, Bengaluru, Karnataka"
destination: "Sarjapur Road, Bengaluru, Karnataka"

# Residential Areas
origin: "Indiranagar, Bengaluru, Karnataka"
destination: "HSR Layout, Bengaluru, Karnataka"
```

## 🎯 **Best Practices**

### **Why Use Coordinates?** 🎯
- **Precision**: Exact location, no geocoding errors
- **Speed**: No geocoding delay, faster route calculation  
- **Reliability**: Always works, even for new/unlisted addresses
- **Custom Names**: Use any display name you want

### **When to Use Coordinates:**
- **New buildings** not yet in map databases
- **Precise locations** within large complexes
- **Custom points** like parking entrances
- **Performance** critical applications

### **Include These Elements:**
1. **Building/Area Name**: "Manyata Tech Park"
2. **Street/Road**: "Thanisandra Main Road"
3. **City**: "Bengaluru"
4. **State**: "Karnataka"
5. **Pincode** (optional but helpful): "560045"
6. **Country**: "India"

### **Address Format Template:**
```
[Building/Complex Name], [Street/Road Name], [Area], [City], [State] [Pincode], [Country]
```

## 📍 **Real Bangalore Examples**

### **Common Commute Routes:**

#### **Electronic City ↔ Whitefield**
```yaml
origin: "Electronic City Phase 1, Hosur Road, Bengaluru, Karnataka 560100, India"
destination: "ITPL, Whitefield, Bengaluru, Karnataka 560066, India"
```

#### **Koramangala ↔ Marathahalli**
```yaml
origin: "Koramangala 5th Block, Bengaluru, Karnataka 560095, India"
destination: "Marathahalli Bridge, Outer Ring Road, Bengaluru, Karnataka 560037, India"
```

#### **Indiranagar ↔ HSR Layout**
```yaml
origin: "Indiranagar 100 Feet Road, Bengaluru, Karnataka 560038, India"
destination: "HSR Layout Sector 1, Bengaluru, Karnataka 560102, India"
```

#### **Airport ↔ City Center**
```yaml
origin: "Kempegowda International Airport, Devanahalli, Bengaluru, Karnataka 560300, India"
destination: "MG Road, Bengaluru, Karnataka 560001, India"
```

## 🔍 **Testing Your Addresses**

Before using addresses in the config, test them:

1. **Search on OpenStreetMap**: Go to [openstreetmap.org](https://openstreetmap.org) and search
2. **Check Google Maps**: Verify the location exists and is correctly named
3. **Use Nominatim directly**: Test at [nominatim.openstreetmap.org](https://nominatim.openstreetmap.org)

## ⚠️ **Addresses to Avoid**

### **Too Vague:**
```yaml
# Don't use these - too general
origin: "Bengaluru"
destination: "Whitefield"
```

### **Non-existent Places:**
```yaml
# Don't use made-up or very new addresses
origin: "Random Building, Unknown Road, Bengaluru"
```

### **Incomplete Information:**
```yaml
# Missing important details
origin: "Tech Park"  # Which tech park?
destination: "Main Road"  # Which main road?
```

## 🚀 **Pro Tips**

### **For Better Accuracy:**
1. **Include pincode** when possible
2. **Use official building/complex names**
3. **Add major road names** (e.g., "Outer Ring Road", "Sarjapur Road")
4. **Include landmarks** when addresses are unclear

### **For Display Names:**
The widget automatically extracts short names for display:
```
"Manyata Tech Park, Thanisandra Main Road, Bengaluru..." 
→ Displays as: "Manyata Tech Park"
```

### **Common Bangalore Roads:**
- Outer Ring Road (ORR)
- Hosur Road
- Bannerghatta Road
- Whitefield Main Road
- Sarjapur Road
- Electronic City Elevated Expressway

## 🔧 **Configuration Example**

Here's a complete example for a typical Bangalore commute:

```yaml
widgets:
  traffic:
    ttl: 300s  # 5 minutes refresh
    origin: "Prestige Tech Park, Sarjapur-Marathahalli Road, Bengaluru, Karnataka 560103, India"
    destination: "RMZ Infinity, Old Madras Road, Bengaluru, Karnataka 560016, India"
```

### **Alternative: Using Coordinates**

For maximum precision, you can use latitude/longitude coordinates:

```yaml
widgets:
  traffic:
    ttl: 300s
    origin:
      latitude: 12.9352
      longitude: 77.6245
      name: "Home"  # Custom display name
    destination:
      latitude: 12.9698
      longitude: 77.7500
      name: "Office"
```

### **Mixed Configuration**

You can mix address strings and coordinates:

```yaml
widgets:
  traffic:
    ttl: 300s
    origin: "Bangalore International Airport, Devanahalli, Bengaluru, Karnataka"
    destination:
      latitude: 12.9716
      longitude: 77.5946
      name: "City Center"
```

## 📍 **Getting Coordinates**

To find latitude/longitude for any location:

1. **Google Maps**: Right-click → "What's here?" → Copy coordinates
2. **OpenStreetMap**: Click location → See coordinates in URL
3. **GPS Apps**: Use any GPS app to get precise coordinates
4. **Online Tools**: Use tools like [latlong.net](https://www.latlong.net/)

### **Common Bangalore Coordinates:**
```yaml
# Major landmarks with precise coordinates
bangalore_airport:
  latitude: 13.1986
  longitude: 77.7066

ub_city_mall:
  latitude: 12.9716
  longitude: 77.5946

electronic_city:
  latitude: 12.8456
  longitude: 77.6603

whitefield:
  latitude: 12.9698
  longitude: 77.7500

koramangala:
  latitude: 12.9279
  longitude: 77.6271
```

## 🔧 **Automatic Setup**

GoDay automatically handles configuration setup:

1. **First Run**: If `~/.goday/` doesn't exist, GoDay creates it
2. **Default Config**: If `config.yaml` doesn't exist, GoDay creates a default one
3. **User Feedback**: Shows where files are created
4. **Ready to Use**: Works immediately with default Bangalore locations

### **What GoDay Creates Automatically:**
```
~/.goday/                    # Config directory (created automatically)
├── config.yaml             # Default config (created automatically)
└── cache/                  # Cache directory (created as needed)
```

### **First Run Experience:**
```bash
$ ./goday
📁 Created config directory: /Users/yourname/.goday
📝 Created default config: /Users/yourname/.goday/config.yaml
💡 Edit the config file to customize your dashboard

[Dashboard starts with default settings]
```

The more specific your addresses, the more accurate your traffic estimates will be!

## 🚨 **Troubleshooting "Location Not Found" Errors**

If you're getting "location not found" errors, try these solutions:

### **1. Use Coordinates (Most Reliable)**
```yaml
traffic:
  origin:
    latitude: 12.9698
    longitude: 77.7500
    name: "Your Home Area"
  destination:
    latitude: 12.9279
    longitude: 77.6271
    name: "Your Office Area"
```

### **2. Simplify Address Names**
Instead of specific building names, use general area names:

❌ **Too Specific (may fail):**
```yaml
origin: "Disha Park West, Balagere Road, Panathur 560087"
destination: "Jfrog India Pvt Ltd, Green Glen Layout, Bellandur"
```

✅ **Better (more likely to work):**
```yaml
origin: "Panathur, Bengaluru, Karnataka, India"
destination: "Bellandur, Bengaluru, Karnataka, India"
```

### **3. Test Your Addresses First**
Before using addresses in config, test them:

```bash
# Test with OpenStreetMap Nominatim
curl -H "User-Agent: GoDay-Dashboard/1.0" \
  "https://nominatim.openstreetmap.org/search?q=Panathur,%20Bengaluru&format=json&limit=1&countrycodes=in"
```

### **4. Common Working Addresses**
These addresses typically work well:

```yaml
# Major Areas
origin: "Electronic City, Bengaluru, Karnataka, India"
destination: "Whitefield, Bengaluru, Karnataka, India"

# Metro Stations
origin: "MG Road Metro Station, Bengaluru, Karnataka, India"
destination: "Indiranagar Metro Station, Bengaluru, Karnataka, India"

# Well-known Places
origin: "Bangalore International Airport, Bengaluru, Karnataka, India"
destination: "Cubbon Park, Bengaluru, Karnataka, India"
```

### **5. Enable Debug Mode**
The traffic plugin will show you which geocoding attempts are failing:
```
Geocoding attempt 1 failed for 'Your Address': no results found
```

### **6. Quick Fix: Use Area Coordinates**
Here are coordinates for common Bangalore areas:

```yaml
# Copy these coordinates for your areas
panathur:
  latitude: 12.9698
  longitude: 77.7500

bellandur:
  latitude: 12.9279
  longitude: 77.6271

electronic_city:
  latitude: 12.8456
  longitude: 77.6603

whitefield:
  latitude: 12.9698
  longitude: 77.7500

koramangala:
  latitude: 12.9352
  longitude: 77.6245

indiranagar:
  latitude: 12.9784
  longitude: 77.6408

hsr_layout:
  latitude: 12.9082
  longitude: 77.6476

marathahalli:
  latitude: 12.9591
  longitude: 77.7074
```

## 📍 **How to Get Precise Coordinates (Better than Geocoding APIs)**

You're right that geocoding APIs aren't always accurate. Here are the best ways to get exact coordinates:

### **Method 1: Google Maps (Most Accurate) 🎯**

1. **Open Google Maps** in your browser
2. **Right-click** on your exact location
3. **Click** the coordinates that appear
4. **Copy** the latitude, longitude values

**Example:**
```
Right-click on your office building
→ "12.9279, 77.6271"
→ Copy these numbers
```

### **Method 2: Google Maps Mobile App**
1. **Long press** on your location in Google Maps app
2. **Tap** on the coordinates at the bottom
3. **Share** or copy the coordinates

### **Method 3: What3Words (Super Precise)**
1. Go to [what3words.com](https://what3words.com)
2. Search for your location
3. Get the exact coordinates displayed
4. Very precise - accurate to 3m x 3m squares

### **Method 4: GPS from Your Phone**
Use any GPS app that shows coordinates:
- **iPhone**: Compass app shows coordinates
- **Android**: GPS apps like "GPS Coordinates"
- **Any navigation app** with coordinate display

### **Method 5: Plus Codes (Google's System)**
1. Search your location in Google Maps
2. Look for the **Plus Code** (e.g., "8Q6P+XM Bengaluru")
3. These convert to exact coordinates

### **Method 6: OpenStreetMap (Free but Less Accurate)**
1. Go to [openstreetmap.org](https://openstreetmap.org)
2. Search for your location
3. Right-click → "Show address"
4. Copy the coordinates from the URL

## 🎯 **Recommended: Use Google Maps Method**

**Step-by-step for your locations:**

1. **For your home/origin:**
   ```
   1. Open Google Maps
   2. Search "Disha Park West, Panathur, Bengaluru"
   3. Zoom in to the exact building
   4. Right-click on the building
   5. Copy coordinates (e.g., "12.9698, 77.7500")
   ```

2. **For your office/destination:**
   ```
   1. Search "JFrog India, Bellandur, Bengaluru"
   2. Right-click on the exact building
   3. Copy coordinates (e.g., "12.9279, 77.6271")
   ```

3. **Update your config:**
   ```yaml
   traffic:
     origin:
       latitude: 12.9698  # Your exact home coordinates
       longitude: 77.7500
       name: "Home"
     destination:
       latitude: 12.9279  # Your exact office coordinates
       longitude: 77.6271
       name: "Office"
   ```

## ⚡ **Why Coordinates Are Better Than Any API:**

- ✅ **100% Accurate** - No geocoding errors
- ✅ **Instant** - No API delays or rate limits
- ✅ **Free** - No API costs
- ✅ **Reliable** - Always works, no network dependency
- ✅ **Precise** - Down to the exact building entrance
- ✅ **Custom Names** - Display any name you want

## 🚀 **Quick Coordinate Lookup Tool**

I can also provide coordinates for common Bangalore locations if you tell me the specific areas/buildings you need!
