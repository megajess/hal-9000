# Project HAL - Server

Go REST API serving as the communication bridge between HAL devices
and the iOS app/web UI.

## Requirements
- Go 1.21 or later
- SQLite (production)

## Project Structure
server/
├── main.go
├── go.mod
├── go.sum
├── .env
├── .env.example
├── handlers/
│   ├── auth.go
│   ├── device.go
│   ├── user.go
│   └── diagnostic.go
├── middleware/
│   ├── auth.go
│   └── apikey.go
├── models/
│   ├── user.go
│   ├── device.go
│   └── diagnostic.go
├── store/
│   ├── store.go
│   ├── memory.go
│   └── sqlite.go
└── static/
└── index.html

## Getting Started

### Environment Variables
Copy .env.example to .env and fill in values:
```bash
cp .env.example .env
```

Required variables:
HAL_PORT=8080
HAL_JWT_SECRET=your-secret-here
HAL_ENV=development|production
HAL_DB_PATH=./hal.db

### Running Locally
```bash
go run main.go
```

### Building for Production
```bash
go build -o hal-server
```

### Deploying
```bash
# Build for Linux if cross compiling from Mac
GOOS=linux GOARCH=amd64 go build -o hal-server

# Copy to server
scp hal-server user@your-server:/path/to/hal

# Run
./hal-server
```

## API Reference

### Authentication
POST /auth/register
Body: { "username": "string", "password": "string" }
POST /auth/login
Body: { "username": "string", "password": "string" }
Returns: { "token": "jwt" }
POST /auth/refresh
Header: Authorization: Bearer <jwt>
Returns: { "token": "jwt" }

### Devices
All device endpoints require Authorization: Bearer <jwt> header
GET    /devices
Returns: [Device]
POST   /devices
Body: { "name": "string" }
Returns: { "device": Device, "apiKey": "string" }
GET    /devices/:id
Returns: Device
PUT    /devices/:id
Body: { "name": "string", "desiredState": "on|off" }
DELETE /devices/:id

### Device Poll
Requires X-API-Key header
GET /poll?state=1&boot=42&reason=0&uptime=3600
X-API-Key: unique-per-device-key
Returns: 0 or 1 (plain text)

### Diagnostics
Requires Authorization: Bearer <jwt> header
GET /devices/:id/diagnostics
Returns: [DiagnosticEvent]

## Authentication Model
- **Users** — JWT based, 15 minute expiry with 7 day refresh token
- **Devices** — unique per device API key sent in X-API-Key header
- Passwords hashed with bcrypt
- API keys generated with crypto/rand

## Development vs Production
- **Development** — in memory store, no SQLite setup required
- **Production** — SQLite, single file database

Switch via HAL_ENV environment variable.

## Security Notes
- HTTPS required in production
- JWT secret must be cryptographically random and kept secret
- API keys are never logged
- Rate limiting applied to auth endpoints
- .env is never committed — use .env.example as reference