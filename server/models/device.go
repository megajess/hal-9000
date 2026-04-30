package models

import "time"

type Device struct {
	ID           string    `json:"id"             db:"id"`
	Name         string    `json:"name"           db:"name"`
	APIKey       string    `json:"-"              db:"api_key"`
	CurrentState string    `json:"current_state"  db:"current_state"`
	DesiredState string    `json:"desired_state"  db:"desired_state"`
	LastSeen     time.Time `json:"last_seen"      db:"last_seen"`
	CreatedAt    time.Time `json:"created_at"     db:"created_at"`
}
