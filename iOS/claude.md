# Project HAL - iOS Context

## Overview
The HAL iOS app is a native Swift/SwiftUI application responsible for:
- User account management (register, login)
- Device pairing flow
- Viewing and controlling all registered devices
- Displaying device health and diagnostic history
- HomeKit integration via HAP (nice to have, future feature)

## Development Environment
- **Language:** Swift
- **UI Framework:** SwiftUI
- **Minimum iOS target:** iOS 18+
- **Xcode** for building and signing
- **Package manager:** Swift Package Manager

## Project Structure
ios/
├── HAL/
│   ├── App/
│   │   ├── HALApp.swift           // app entry point
│   │   └── AppState.swift         // top level app state
│   ├── Features/
│   │   ├── Auth/
│   │   │   ├── LoginView.swift
│   │   │   ├── RegisterView.swift
│   │   │   └── AuthViewModel.swift
│   │   ├── Devices/
│   │   │   ├── DeviceListView.swift
│   │   │   ├── DeviceDetailView.swift
│   │   │   ├── DeviceViewModel.swift
│   │   │   └── DeviceRowView.swift
│   │   ├── Pairing/
│   │   │   ├── PairingView.swift
│   │   │   ├── PairingViewModel.swift
│   │   │   └── WiFiPasswordView.swift
│   │   └── Diagnostics/
│   │       ├── DiagnosticsView.swift
│   │       └── DiagnosticsViewModel.swift
│   ├── Models/
│   │   ├── User.swift
│   │   ├── Device.swift
│   │   └── DiagnosticEvent.swift
│   ├── Network/
│   │   ├── HALClient.swift        // API client
│   │   ├── Endpoints.swift        // endpoint definitions
│   │   └── NetworkError.swift     // error types
│   ├── Storage/
│   │   └── KeychainHelper.swift   // JWT and sensitive storage
│   └── Utilities/
│       ├── Extensions/
│       └── Constants.swift
├── HALTests/
├── HALUITests/
└── README.md

## Architecture
- **Pattern:** MVVM
- ViewModels are ObservableObjects with @Published properties
- Views are purely declarative — no business logic in views
- Network layer is abstracted behind HALClient
- All sensitive data (JWT, server URL) stored in Keychain

## Data Models

### Device
```swift
struct Device: Identifiable, Codable {
    let id: String
    var name: String
    var currentState: SwitchState
    var desiredState: SwitchState
    var lastSeen: Date
}

enum SwitchState: String, Codable {
    case on
    case off
    case unknown
}
```

### DiagnosticEvent
```swift
struct DiagnosticEvent: Identifiable, Codable {
    let id: String
    let deviceID: String
    let resetReason: String
    let bootCount: Int
    let uptime: Int
    let createdAt: Date
}
```

## Network Layer

### HALClient
```swift
class HALClient {
    static let shared = HALClient()
    private var baseURL: String { Environment.current.baseURL }
    private var jwt: String { /* from Keychain */ }

    func getDevices() async throws -> [Device]
    func updateDesiredState(deviceID: String, state: SwitchState) async throws
    func getDiagnostics(deviceID: String) async throws -> [DiagnosticEvent]
    func login(username: String, password: String) async throws -> String
    func register(username: String, password: String) async throws
    func registerDevice(name: String) async throws -> String // returns API key
}
```

- Uses async/await throughout — no Combine
- URLSession for all networking
- JWT attached to all requests via Authorization header
- JWT stored in Keychain, never UserDefaults

## Pairing Flow
The pairing flow is the most complex feature in the app:

User taps "Add Device"
App calls POST /devices → server returns API key
User instructed to hold pairing button on HAL device for 3 seconds
App uses NEHotspotConfiguration to prompt user to join HAL AP
App POSTs to http://192.168.4.1/configure:
{
"ssid": "<current network SSID>",
"password": "<user entered>",
"serverUrl": "<stored server URL>",
"apiKey": "<from step 2>"
}
App uses NEHotspotConfiguration to rejoin home network
Device appears active in device list


### Key APIs
```swift
import NetworkExtension

// Join HAL access point
let config = NEHotspotConfiguration(ssid: "HAL-Device")
NEHotspotConfigurationManager.shared.apply(config) { error in }

// Get current SSID (pre-fill WiFi field)
import SystemConfiguration.CaptiveNetwork
// requires Access WiFi Information entitlement

// Rejoin home network
NEHotspotConfigurationManager.shared.removeConfiguration(forSSID: "HAL-Device")
```

### Required Entitlements
- `com.apple.developer.networking.HotspotConfiguration`
- `com.apple.developer.networking.wifi-info`

### Pairing UX Notes
- Pre-fill SSID from current network — user only enters password
- Password field has show/hide toggle
- "Go to WiFi Settings" button deep links to Settings for easy
  password copy: `UIApplication.shared.open(URL(string: "App-Prefs:WIFI")!)`
- Clear progress indication throughout pairing steps
- Auto timeout with helpful error if device AP not found
- LED color described in UI so user knows what to look for

## Storage
- **JWT** — Keychain
- **User preferences** — UserDefaults (non sensitive only)
- Never store sensitive data in UserDefaults
- Server URL is hardcoded per environment — see Environment enum

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

## Error Handling
- All network errors surfaced to user with actionable messages
- No silent failures
- Pairing flow has explicit error states for each step:
  - Device AP not found
  - Wrong WiFi password
  - Device failed to connect to home network
  - Server unreachable

## HomeKit Integration (Future)
- Nice to have, not in initial scope
- Will use brutella/hap on the server side — server acts as HomeKit bridge
- App will not implement HomeKit directly — native Home app handles that
- Note in backlog, do not implement until core features complete

## Developer Background
The developer is a senior iOS engineer with 12+ years of experience
in Swift/SwiftUI/UIKit. This is familiar territory. Claude should:
- Assume strong Swift and SwiftUI knowledge
- Focus on HAL specific implementation details
- Feel free to suggest idiomatic Swift without over-explaining basics
- Still explain reasoning behind architectural suggestions
- Flag any iOS version compatibility concerns

## Teach Mode Notes
- Developer is most experienced here — less hand holding needed on
  Swift basics
- Still explain reasoning behind every suggestion
- NEHotspotConfiguration and CaptiveNetwork APIs are potentially
  less familiar — explain these in detail
- HomeKit/HAP worth explaining when that feature is tackled
- Always ask clarifying questions before starting any task

## General Rules
- Ask clarifying questions before starting any task
- Explain reasoning behind every suggestion
- Suggest tests where appropriate — UI tests for pairing flow,
  unit tests for ViewModels and network layer
- Prefer native APIs — URLSession over Alamofire,
  SwiftUI over UIKit where possible
- Never refactor outside scope of current task
- Flag but do not implement improvements noticed outside current task
- Perform thorough code review when asked before commits