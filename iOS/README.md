# Project HAL - iOS

Native iOS app for managing and controlling HAL smart light switches.

## Requirements
- Xcode 16 or later
- iOS 18 or later
- Swift 5.9 or later
- Active Apple Developer account (for device builds)

## Getting Started
1. Open ios/ in Xcode
2. Select your development team in project signing settings
3. Update Environment.swift with your local server IP for development
4. Build and run on device or simulator

## Features
- User account registration and login
- Device pairing via guided flow
- Real time device state display
- Toggle switch state remotely
- Device health and diagnostic history

## Architecture
MVVM pattern throughout:
- **Views** — declarative SwiftUI, no business logic
- **ViewModels** — ObservableObject, @Published properties
- **Models** — Codable structs matching server API
- **Network** — HALClient abstraction over URLSession
- **Storage** — Keychain for JWT, UserDefaults for preferences

## Project Structure
HAL/
├── App/
│   ├── HALApp.swift
│   └── AppState.swift
├── Features/
│   ├── Auth/
│   │   ├── LoginView.swift
│   │   ├── RegisterView.swift
│   │   └── AuthViewModel.swift
│   ├── Devices/
│   │   ├── DeviceListView.swift
│   │   ├── DeviceDetailView.swift
│   │   ├── DeviceViewModel.swift
│   │   └── DeviceRowView.swift
│   ├── Pairing/
│   │   ├── PairingView.swift
│   │   ├── PairingViewModel.swift
│   │   └── WiFiPasswordView.swift
│   └── Diagnostics/
│       ├── DiagnosticsView.swift
│       └── DiagnosticsViewModel.swift
├── Models/
│   ├── User.swift
│   ├── Device.swift
│   └── DiagnosticEvent.swift
├── Network/
│   ├── HALClient.swift
│   ├── Endpoints.swift
│   └── NetworkError.swift
├── Storage/
│   └── KeychainHelper.swift
└── Utilities/
├── Extensions/
└── Constants.swift

## Environment
```swift
enum Environment {
    case development
    case production

    var baseURL: String {
        switch self {
        case .development:
            return "http://192.168.1.105:8080"
        case .production:
            return "https://your-production-server.com"
        }
    }

    static var current: Environment {
        #if DEBUG
        return .development
        #else
        return .production
        #endif
    }
}
```

## Pairing Flow
1. Tap "Add Device" — server registers device and generates API key
2. Hold pairing button on HAL device for 3 seconds — LED turns blue
3. App prompts to join HAL access point — one tap, no need to leave app
4. Enter home WiFi password — SSID pre-filled from current network
5. App configures device and rejoins home network
6. Device appears active in device list

## Required Entitlements
- `com.apple.developer.networking.HotspotConfiguration`
- `com.apple.developer.networking.wifi-info`

## Storage
- **JWT** — Keychain
- **User preferences** — UserDefaults (non sensitive only)
- **Server URL** — hardcoded per environment, never stored
- Never store sensitive data in UserDefaults

## Security Notes
- JWT stored in Keychain, never UserDefaults
- All requests over HTTPS in production
- No sensitive data in UserDefaults

## Testing
- Unit tests for ViewModels and network layer
- UI tests for critical flows — login, pairing, device toggle
- Run tests with Cmd+U in Xcode

## HomeKit (Future)
HomeKit integration is planned as a future feature via the HAL server
acting as a HomeKit bridge using brutella/hap. No HomeKit code in the
iOS app directly — control will be through the native Home app.