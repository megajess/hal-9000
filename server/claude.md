# Project HAL - Server Context

## Overview
The HAL server is a Go REST API responsible for:
- Managing user accounts and authentication
- Managing device registration and state
- Acting as the communication bridge between devices and the iOS app/web UI
- Logging device diagnostics and health
- Serving the web UI as a static frontend

## Development Environment
- **Language:** Go
- **Editor:** VS Code
- **Deployment target:** Linux (DigitalOcean droplet or similar)
- **Deployment method:** Single binary, built with `go build`

## Project Structure
server/
├── main.go                 // entry point, server setup, route registration
├── go.mod
├── go.sum
├── .env                    // secrets, never committed
├── .env.example            // committed, documents required env vars
├── handlers/
│   ├── auth.go             // login, register, JWT issuance
│   ├── device.go           // device poll endpoint
│   ├── user.go             // user account management
│   └── diagnostic.go       // diagnostic event logging
├── middleware/
│   ├── auth.go             // JWT validation middleware
│   └── apikey.go           // device API key validation middleware
├── models/
│   ├── user.go             // User struct
│   ├── device.go           // Device struct
│   └── diagnostic.go       // DiagnosticEvent struct
├── store/
│   ├── store.go            // storage interface
│   ├── memory.go           // in memory implementation (development)
│   └── sqlite.go           // SQLite implementation (production)
├── static/
│   └── index.html          // web UI served by Go binary
└── README.md

## Data Models

### User
```go
type User struct {
    ID           string    
    Username     string    
    PasswordHash string    
    CreatedAt    time.Time 
}
```

### Device
```go
type Device struct {
    ID           string    
    UserID       string    // owner
    Name         string    // friendly name e.g. "Living Room"
    APIKey       string    // unique per device, used for authentication
    CurrentState string    // reported by device — "on" or "off"
    DesiredState string    // set by user via app or web UI
    LastSeen     time.Time // last successful poll
    CreatedAt    time.Time 
}
```

### DiagnosticEvent
```go
type DiagnosticEvent struct {
    ID          string    
    DeviceID    string    
    ResetReason string    
    BootCount   int       
    Uptime      int       
    CreatedAt   time.Time 
}
```

## API Endpoints

### Authentication (JWT)
POST /auth/register        // create user account
POST /auth/login           // returns JWT
POST /auth/refresh         // refresh JWT

### Devices (JWT required)
GET    /devices            // list all devices for authenticated user
POST   /devices            // register new device, returns API key
GET    /devices/:id        // get device state
PUT    /devices/:id        // update device desired state or name
DELETE /devices/:id        // remove device

### Device Poll (API key required)
GET /poll?state=1&boot=42&reason=0&uptime=3600
X-API-Key: unique-per-device-key
Response: 0 or 1 (plain text)

### Diagnostics (JWT required)
GET /devices/:id/diagnostics    // get diagnostic history for device

## Authentication Model

### Users — JWT
- User logs in with username and password
- Server issues a short lived JWT (15 minutes)
- Refresh token issued alongside JWT, longer lived (7 days)
- All user facing endpoints validate JWT via middleware
- Passwords hashed with bcrypt

### Devices — API Key
- Each device has a unique API key generated at registration
- API key sent in X-API-Key header on every poll
- Server looks up device by API key
- Device has no concept of users or accounts
- If a device is compromised, revoke only that device's API key

## Poll Endpoint Logic

Validate API key → identify device
Update device CurrentState from query params
Update device LastSeen timestamp
Check if reset reason indicates unexpected reboot
→ if so, log DiagnosticEvent
Compare CurrentState to DesiredState
→ if match: respond "0" (no change)
→ if differ: respond "1" (toggle)


## Storage Strategy
- **Development:** in memory store — fast, no setup required
- **Production:** SQLite — simple, single file, no separate DB server
- Storage interface abstracts the implementation:

```go
type Store interface {
    // Users
    CreateUser(user User) error
    GetUserByUsername(username string) (User, error)

    // Devices
    CreateDevice(device Device) error
    GetDeviceByAPIKey(apiKey string) (Device, error)
    GetDevicesByUserID(userID string) ([]Device, error)
    UpdateDeviceState(deviceID string, current string) error
    UpdateDeviceDesiredState(deviceID string, desired string) error
    DeleteDevice(deviceID string) error

    // Diagnostics
    CreateDiagnosticEvent(event DiagnosticEvent) error
    GetDiagnosticsByDeviceID(deviceID string) ([]DiagnosticEvent, error)
}
```

Swapping from memory to SQLite requires no changes outside the store package.

## Security Considerations
- HTTPS required in production — server is internet facing
- JWT secret stored in environment variable, never hardcoded
- Passwords hashed with bcrypt, never stored plain
- API keys generated with crypto/rand, not math/rand
- Rate limiting on auth endpoints to prevent brute force
- Device API keys never logged
- .env never committed — .env.example documents required variables

## Environment Variables
HAL_PORT=8080
HAL_JWT_SECRET=your-secret-here
HAL_ENV=development|production
HAL_DB_PATH=./hal.db          // SQLite only

## Web UI
- Served as static HTML from the Go binary
- Single page, no framework — vanilla JS
- Shows all devices for authenticated user
- Toggle desired state per device
- Shows device health/last seen
- Shows diagnostic history per device

## Developer Background
The developer is a senior iOS/Swift engineer with Go experience that
is slightly rusty. When working in this component:
- Remind of Go idioms where relevant — error handling patterns,
  interface usage, package structure
- Explain standard library packages being used
- Flag any concurrency concerns — the poll endpoint will be hit
  frequently by multiple devices simultaneously
- Prefer standard library over third party where possible
- The developer has used Express.js previously — analogies to that
  are helpful where appropriate

## Teach Mode Notes
- Explain Go patterns that may differ from Swift or Express.js
- Always explain middleware concepts when introducing them
- Flag any security implications of suggested approaches
- Concurrency is a particular area to explain carefully —
  Go's concurrency model may be unfamiliar
- Explain interface usage and why it benefits this codebase

## General Rules
- Ask clarifying questions before starting any task
- Explain reasoning behind every suggestion
- Suggest tests where appropriate and explain what they cover
- Prefer standard library — explain why if third party is needed
- Never refactor outside scope of current task
- Flag but do not implement improvements noticed outside current task
- Perform thorough code review when asked before commits