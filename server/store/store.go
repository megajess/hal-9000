package store

import (
	"hal/models"
	"time"
)

type Store interface {
	CreateDevice(device models.Device) error
	GetDeviceByAPIKey(apiKey string) (models.Device, error)
	// TODO: timeOfUpdate is a testing hook that leaks into this interface.
	// The idiomatic Go approach is to inject a clock via a func field on the
	// concrete struct (e.g. `Now func() time.Time` on MemoryStore) so the
	// interface signature stays clean.
	UpdateDeviceState(deviceID string, reportedState string, timeOfUpdate ...time.Time) error
}
