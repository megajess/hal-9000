package models

import "time"

type Device struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	// TODO: APIKey is included in all JSON responses. Correct for registration,
	// but will expose keys in GET /devices list responses added in Phase 3.
	// Fix: use json:"-" here and include the key only in the registration response type.
	APIKey       string    `json:"api_key"`
	CurrentState string    `json:"current_state"`
	DesiredState string    `json:"desired_state"`
	LastSeen     time.Time `json:"last_seen"`
	CreatedAt    time.Time `json:"created_at"`
}
