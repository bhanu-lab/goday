# Traffic Widget - Coordinate Configuration Examples

## Example 1: Both locations using coordinates
```yaml
widgets:
  traffic:
    ttl: 300s
    origin:
      latitude: 12.9716
      longitude: 77.5946
      name: "UB City Mall"
    destination:
      latitude: 12.9698
      longitude: 77.7500
      name: "Whitefield Tech Park"
```

## Example 2: Mixed address and coordinates
```yaml
widgets:
  traffic:
    ttl: 300s
    origin: "Bangalore International Airport, Devanahalli, Bengaluru, Karnataka"
    destination:
      latitude: 12.9279
      longitude: 77.6271
      name: "Koramangala"
```

## Example 3: Precise building locations
```yaml
widgets:
  traffic:
    ttl: 300s
    origin:
      latitude: 12.8456
      longitude: 77.6603
      name: "Electronic City Gate"
    destination:
      latitude: 13.0358
      longitude: 77.5970
      name: "Manyata Entrance"
```

## Testing Configuration

To test if your coordinates work:

1. Copy your coordinates
2. Paste them into Google Maps as: "12.9716, 77.5946"
3. Verify the location is correct
4. Use the configuration in your config.yaml

## Advantages:

✅ **Exact precision** - No geocoding ambiguity
✅ **Faster loading** - No geocoding delay  
✅ **Custom names** - Use any display name
✅ **Reliable** - Works even for new locations
✅ **Mixed mode** - Combine with address strings
