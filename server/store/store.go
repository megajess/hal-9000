package store

import (
	"errors"
	"hal/models"
)

var ErrDeviceNotFound = errors.New("device not found")
var ErrUserNotFound = errors.New("user not found")
var ErrRefreshTokenNotFound = errors.New("refresh token not found")
var ErrUsernameTaken = errors.New("username already taken")
var ErrRefreshTokenExpired = errors.New("refresh token expired")
var ErrUserNotAssociated = errors.New("user not associated with device")

type Store interface {
	CreateDevice(device models.Device, ownerUserID string) error
	GetDeviceByAPIKey(apiKey string) (models.Device, error)
	GetDeviceByID(deviceID string) (models.Device, error)
	GetDevicesByUserID(userID string) ([]models.Device, error)
	UpdateDeviceState(deviceID string, reportedState string) error
	UpdateDevice(device models.Device) error
	DeleteDevice(deviceID string) error

	CreateUser(user models.User) error
	GetUserByUsername(username string) (models.User, error)
	GetUserByID(ID string) (models.User, error)

	AddUserToDevice(deviceID, userID, role string) error
	RemoveUserFromDevice(deviceID, userID string) error
	GetUsersByDeviceID(deviceID string) ([]models.DeviceUser, error)

	StoreRefreshToken(token string, userID string) error
	GetUserIDByRefreshToken(token string) (string, error)
	DeleteRefreshToken(token string) error
}
