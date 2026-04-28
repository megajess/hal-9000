package store

import (
	"errors"
	"hal/models"
	"time"
)

var ErrDeviceNotFound = errors.New("device not found")
var ErrUserNotFound = errors.New("user not found")
var ErrRefreshTokenNotFound = errors.New("refresh token not found")
var ErrUsernameTaken = errors.New("username already taken")

type Store interface {
	CreateDevice(device models.Device) error
	GetDeviceByAPIKey(apiKey string) (models.Device, error)
	// TODO: timeOfUpdate is a testing hook that leaks into this interface.
	// The idiomatic Go approach is to inject a clock via a func field on the
	// concrete struct (e.g. `Now func() time.Time` on MemoryStore) so the
	// interface signature stays clean.
	UpdateDeviceState(deviceID string, reportedState string, timeOfUpdate ...time.Time) error

	CreateUser(user models.User) error
	GetUserByUsername(username string) (models.User, error)
	GetUserByID(ID string) (models.User, error)

	StoreRefreshToken(token string, userID string) error
	GetUserIDByRefreshToken(token string) (string, error)
	DeleteRefreshToken(token string) error

	GetDeviceByID(deviceID string) (models.Device, error)
	GetDevicesByUserID(userID string) ([]models.Device, error)
	UpdateDevice(device models.Device) error
	DeleteDevice(deviceID string) error
}
