# Project HAL - Firmware

ESP8266 firmware for the HAL smart light switch, built with PlatformIO
and the Arduino framework.

## Hardware Requirements
- Wemos D1 Mini (ESP8266)
- HiLink HLK-PM01 (120V AC → 5V DC)
- Isolated relay module (active HIGH)
- Momentary SPDT center off switch (low voltage)
- RGB LED (common cathode)
- Panel mount momentary button (pairing)

## Wiring

### GPIO Pin Assignment
| Pin | GPIO | Role | Direction |
|-----|------|------|-----------|
| D3 | GPIO0  | Switch DOWN (always off) | Input |
| D1 | GPIO5 | Relay control | Output |
| D2 | GPIO4 | Switch UP (always on) | Input |
| D5 | GPIO14 | RGB Red | Output (PWM) |
| D6 | GPIO12 | RGB Green | Output (PWM) |
| D7 | GPIO13 | RGB Blue | Output (PWM) |
| D8 | GPIO15 | Pairing button | Input (pull down) |

### RGB LED
- R, G, B pins → respective GPIO via 220Ω resistor
- Common cathode → GND

### Relay Module
- IN → D1
- VCC → 5V
- GND → GND

### Switch
- Common → GND
- UP terminal → D2
- DOWN terminal → D3

### Pairing Button
- One terminal → D8
- Other terminal → 3.3V
- Pull down resistor required on D8

## Device Modes
| Mode | LED | Description |
|------|-----|-------------|
| Connected | Solid green | Connected, polling server |
| Connecting | Green blinking | Attempting WiFi connection |
| WiFi Offline | Red/blue blinking | WiFi connection lost |
| Server Unreachable | Red/green blinking | Server not responding |
| Pairing | Solid blue | AP mode, awaiting configuration (Phase 4) |

## Pairing
1. Hold pairing button for 3 seconds during normal operation
2. LED turns blue — device enters pairing mode
3. Connect to HAL access point via iOS app
4. Enter WiFi password when prompted
5. Device reboots into normal mode on success
6. LED returns to green when connected

## Development Setup
1. Install VS Code
2. Install PlatformIO extension
3. Open firmware/ directory in VS Code
4. PlatformIO will install dependencies automatically

## Building and Flashing
```bash
# Normal build and upload
pio run -t upload

# Debug build and upload
pio run -e d1_mini_debug -t upload

# Open serial monitor
pio device monitor
```

## Dependencies
None currently. WiFiManager (tzapu) will be added in Phase 4 for pairing mode.

## Memory Considerations
The ESP8266 has 80KB usable RAM. To conserve memory:
- Plain text query parameters used instead of JSON
- Arduino String class used sparingly
- If random reboots occur after extended runtime suspect memory
  pressure — see docs/architecture.md for upgrade path