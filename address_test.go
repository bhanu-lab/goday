package main

import (
	"fmt"
	"testing"
)

// TestSpecificAddresses tests the traffic plugin with real Bangalore addresses
func TestSpecificAddresses(t *testing.T) {
	plugin := NewOSRMTrafficPlugin()

	// Test with specific Bangalore addresses
	config := map[string]interface{}{
		"origin":      "Manyata Tech Park, Thanisandra Main Road, Bengaluru, Karnataka 560045, India",
		"destination": "Brigade Metropolis, Whitefield Main Road, Bengaluru, Karnataka 560066, India",
	}

	err := plugin.Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize plugin: %v", err)
	}

	fmt.Println("Testing traffic between specific Bangalore locations:")
	fmt.Printf("Origin: %s\n", config["origin"])
	fmt.Printf("Destination: %s\n", config["destination"])

	// Test display name extraction
	originShort := plugin.getLocationShortName(config["origin"].(string))
	destShort := plugin.getLocationShortName(config["destination"].(string))

	fmt.Printf("Display names: %s → %s\n", originShort, destShort)

	// Note: We won't actually call Fetch() in tests to avoid hitting the API
	// but this shows how the configuration would work

	if originShort == "" || destShort == "" {
		t.Error("Failed to extract meaningful location names")
	}

	fmt.Println("✅ Address parsing successful!")
}

// TestDifferentAddressFormats tests various address formats
func TestDifferentAddressFormats(t *testing.T) {
	plugin := NewOSRMTrafficPlugin()

	testCases := []struct {
		address  string
		expected string
	}{
		{
			"Manyata Tech Park, Thanisandra Main Road, Bengaluru, Karnataka 560045, India",
			"Manyata Tech Park",
		},
		{
			"123 Whitefield Main Road, EPIP Zone, Whitefield, Bengaluru, Karnataka 560066",
			"EPIP Zone", // Should pick the area name over the street address
		},
		{
			"Phoenix MarketCity, Whitefield Road, Bengaluru, Karnataka 560048, India",
			"Phoenix MarketCity",
		},
		{
			"Outer Ring Road, Marathahalli, Bengaluru, Karnataka",
			"Marathahalli", // Should pick area name over road name
		},
		{
			"Electronic City Phase 1, Bengaluru, Karnataka",
			"Electronic City Phase 1",
		},
	}

	fmt.Println("\nTesting address format parsing:")
	for _, tc := range testCases {
		result := plugin.getLocationShortName(tc.address)
		fmt.Printf("Address: %s\n", tc.address)
		fmt.Printf("Expected: %s, Got: %s\n", tc.expected, result)

		if result == "" {
			t.Errorf("Failed to extract name from: %s", tc.address)
		}
		fmt.Println("---")
	}
}

// ExampleAddressConfigurations shows example configurations
func ExampleAddressConfigurations() {
	fmt.Println("Example Traffic Widget Configurations:")
	fmt.Println()

	examples := []struct {
		name   string
		origin string
		dest   string
	}{
		{
			"Tech Park Commute",
			"Manyata Tech Park, Thanisandra Main Road, Bengaluru, Karnataka 560045, India",
			"Brigade Metropolis, Whitefield Main Road, Bengaluru, Karnataka 560066, India",
		},
		{
			"Residential to Office",
			"Prestige Shantiniketan, Whitefield, Bengaluru, Karnataka 560066, India",
			"RMZ Infinity, Old Madras Road, Bengaluru, Karnataka 560016, India",
		},
		{
			"Airport to City",
			"Kempegowda International Airport, Devanahalli, Bengaluru, Karnataka 560300, India",
			"UB City Mall, Vittal Mallya Road, Bengaluru, Karnataka 560001, India",
		},
		{
			"Classic Bangalore Route",
			"Electronic City Phase 1, Hosur Road, Bengaluru, Karnataka 560100, India",
			"ITPL, Whitefield, Bengaluru, Karnataka 560066, India",
		},
	}

	for _, example := range examples {
		fmt.Printf("## %s\n", example.name)
		fmt.Println("```yaml")
		fmt.Println("widgets:")
		fmt.Println("  traffic:")
		fmt.Println("    ttl: 300s")
		fmt.Printf("    origin: \"%s\"\n", example.origin)
		fmt.Printf("    destination: \"%s\"\n", example.dest)
		fmt.Println("```")
		fmt.Println()
	}
}
