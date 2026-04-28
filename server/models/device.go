package models

import "time"

type Device struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	APIKey       string    `json:"-"`
	CurrentState string    `json:"current_state"`
	DesiredState string    `json:"desired_state"`
	LastSeen     time.Time `json:"last_seen"`
	CreatedAt    time.Time `json:"created_at"`
}
