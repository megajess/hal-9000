# Project HAL - Documentation

Technical documentation for Project HAL.

## Contents

### Architecture
[architecture.md](architecture.md) — Full system architecture, component
interactions, and data flow. Start here to understand how all the pieces
fit together.

### Wiring
[wiring_diagram.md](wiring_diagram.md) — Hardware wiring reference for
the HAL device including GPIO assignments, component connections, and
safety notes.

## Quick Reference

### System Overview
iOS App ──────────────────────────────────────┐
▼
Physical Switch → D1 Mini → polls → Go Server → Web UI
↓
Relay → Load (120V)

### Component Summary
| Component | Role |
|-----------|------|
| Wemos D1 Mini | MCU, WiFi, relay control |
| HiLink HLK-PM01 | 120V AC → 5V DC power supply |
| Isolated relay module | Switches 120V load |
| Momentary SPDT switch | Physical control input |
| RGB LED | Device mode indication |
| Panel mount button | Pairing mode trigger |

### Device Modes
| Mode | LED | Trigger |
|------|-----|---------|
| Normal | Green | Successful WiFi connection |
| Pairing | Blue | Button held 3 seconds |
| Debug | Red | Compiled with DEBUG flag |

### Communication Protocol
GET /poll?state=1&boot=42&reason=0&uptime=3600
X-API-Key: unique-per-device-key
Response: 0 (no change) or 1 (toggle)

### Security Model
- Users authenticate with JWT
- Devices authenticate with unique per device API keys
- Device has no concept of user accounts
- Server associates devices to users internally

## Related
- [Firmware README](../firmware/README.md)
- [Server README](../server/README.md)
- [iOS README](../ios/README.md)
- [Hardware README](../hardware/README.md)