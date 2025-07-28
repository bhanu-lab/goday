# GoDay Configuration Guide

## Configuration File Location

GoDay uses a configuration file located in your home directory for better organization and user-specific settings.

### Default Location
```
~/.goday/config.yaml
```

### Setup Methods

#### Method 1: Automatic (Recommended) ✨
```bash
./goday  # Just run it - GoDay creates everything automatically!
```

GoDay automatically:
- Creates `~/.goday` directory if missing
- Creates default `config.yaml` if missing  
- Shows you where files are created
- Starts with sensible defaults

#### Method 2: Manual Setup Script
```bash
./setup-config.sh
```

#### Method 3: Manual Copy
```bash
mkdir -p ~/.goday
cp config.yaml ~/.goday/
```

## Configuration Lookup Order

GoDay searches for configuration in this order:

1. `~/.goday/config.yaml` (preferred user config)
2. `./config.yaml` (fallback for development)
3. **Auto-creates** default config at `~/.goday/config.yaml` if none found

## Automatic Directory Creation

GoDay handles setup automatically:

✅ **Creates directory**: `~/.goday` if it doesn't exist
✅ **Creates config**: Default `config.yaml` with examples
✅ **Creates cache**: `~/.goday/cache/` as needed  
✅ **User feedback**: Shows what was created
✅ **Works immediately**: No manual setup required

### **First Run Output:**
```
$ ./goday
📁 Created config directory: /Users/yourname/.goday
📝 Created default config: /Users/yourname/.goday/config.yaml
💡 Edit the config file to customize your dashboard

[Dashboard starts immediately]
```

## Command Line Options

### Show Config Location
```bash
./goday config
```

### Help
```bash
./goday help
```

## Configuration Structure

```yaml
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
  traffic:
    ttl: 300s
    # Address-based configuration
    origin: "Electronic City Phase 1, Bengaluru, Karnataka, India"
    destination: "Whitefield, Bengaluru, Karnataka, India"
    
    # OR coordinate-based configuration
    # origin:
    #   latitude: 12.8456
    #   longitude: 77.6603
    #   name: "Electronic City"
    # destination:
    #   latitude: 12.9698
    #   longitude: 77.7500
    #   name: "Whitefield"
```

## Benefits of ~/.goday Location

✅ **User-specific**: Each user has their own config
✅ **Persistent**: Config survives application updates
✅ **Standard**: Follows Unix/Linux conventions
✅ **Organized**: Keeps config separate from code
✅ **Backup-friendly**: Easy to include in dotfiles

## Migration from Local Config

If you have an existing `config.yaml` in your project directory:

```bash
# Move existing config to new location
mkdir -p ~/.goday
mv config.yaml ~/.goday/
```

## Troubleshooting

### Config Not Found
```bash
./goday config  # Shows expected location
./setup-config.sh  # Creates default config
```

### Permission Issues
```bash
chmod 755 ~/.goday
chmod 644 ~/.goday/config.yaml
```

### Reset to Defaults
```bash
rm ~/.goday/config.yaml
./setup-config.sh
```
