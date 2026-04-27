package store_test

import (
	"hal/models"
	"hal/store"
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	deviceID         = "some-device-id"
	deviceName       = "some-name"
	apiKey           = "2-444-66666"
	userID           = "1-22-333-4444"
	username         = "Billiam"
	userPasswordHash = "qwerty"
)

func createTestStore() *store.MemoryStore {
	return store.NewMemoryStore()
}

func createDevice() models.Device {
	return models.Device{
		ID:           deviceID,
		Name:         deviceName,
		APIKey:       apiKey,
		CurrentState: "off",
		DesiredState: "on",
		CreatedAt:    time.Now(),
	}
}

func createTestuser() models.User {
	return models.User{
		ID:           userID,
		Username:     username,
		PasswordHash: userPasswordHash,
		CreatedAt:    time.Now(),
	}
}

func TestCreateDevice(t *testing.T) {
	store := createTestStore()
	deviceModel := createDevice()

	err := store.CreateDevice(deviceModel)

	if err != nil {
		t.Fatalf("CreateDevice returned unexpected error: %v", err)
	}

	device, err := store.GetDeviceByAPIKey(apiKey)

	if err != nil {
		t.Fatalf("Could not find device with API key %q", apiKey)
	}

	if device.Name != deviceName {
		t.Errorf("Expected name %q, got %q instead", deviceName, device.Name)
	}

	if device.ID != deviceID {
		t.Errorf("Expected ID %q, got %q instead", deviceID, device.ID)
	}

	if device.APIKey != apiKey {
		t.Errorf("Expected API key %q, got %q instead", apiKey, device.APIKey)
	}
}

func TestGetDeviceByAPIKey_NotFound(t *testing.T) {
	store := createTestStore()

	_, err := store.GetDeviceByAPIKey("unknown-api-key")

	if err == nil {
		t.Fatal("Expected an error for unknown API key, got nil instead")
	}
}

func TestUpdateDeviceState(t *testing.T) {
	store := createTestStore()
	device := createDevice()

	if err := store.CreateDevice(device); err != nil {
		t.Fatalf("CreateDevice returned unexpected error: %v", err)
	}

	if err := store.UpdateDeviceState(device.ID, "on"); err != nil {
		t.Fatalf("Update device state returned an unexpected error: %v", err)
	}

	updatedDevice, err := store.GetDeviceByAPIKey(device.APIKey)

	if err != nil {
		t.Fatalf("Get device by API key returned an unexpected error: %v", err)
	}

	if updatedDevice.CurrentState != "on" {
		t.Errorf("Expected device state to be %q, got %q instead!", "on", updatedDevice.CurrentState)
	}

	if updatedDevice.LastSeen.IsZero() {
		t.Error("Expected updated device to have non-zero value for last seen!")
	}
}

func TestUpdateDeviceState_NotFound(t *testing.T) {
	store := createTestStore()
	device := createDevice()

	if err := store.CreateDevice(device); err != nil {
		t.Fatalf("CreateDevice returned unexpected error: %v", err)
	}

	if err := store.UpdateDeviceState("unknown-device-id", "on"); err == nil {
		t.Fatal("Expected an error for unknown device id, but got nil instead!")
	}
}

func TestCreateUser(t *testing.T) {
	store := createTestStore()
	user := createTestuser()

	if err := store.CreateUser(user); err != nil {
		t.Fatalf("CreateUser returned unexpected error: %v", err)
	}
}

func TestCreateUser_FindByID(t *testing.T) {
	store := createTestStore()
	user := createTestuser()

	if err := store.CreateUser(user); err != nil {
		t.Fatalf("CreateUser returned unexpected error: %v", err)
	}

	userByID, err := store.GetUserByID(user.ID)

	if err != nil {
		t.Fatalf("Finding user by ID returned unexpected error: %v", err)
	}

	if userByID != user {
		t.Fatal("Fetching newly saved user by ID returned incorrect user")
	}
}

func TestCreateUser_FindByUsername(t *testing.T) {
	store := createTestStore()
	user := createTestuser()

	if err := store.CreateUser(user); err != nil {
		t.Fatalf("CreateUser returned unexpected error: %v", err)
	}

	userByName, err := store.GetUserByUsername(user.Username)

	if err != nil {
		t.Fatalf("Finding user by username returned unexpected error: %v", err)
	}

	if userByName != user {
		t.Fatal("Fetching newly saved user by username returned incorrect user")
	}
}

func TestStoreRefreshToken(t *testing.T) {
	store := createTestStore()
	user := createTestuser()

	if err := store.CreateUser(user); err != nil {
		t.Fatalf("CreateUser returned unexpected error: %v", err)
	}

	refreshToken := uuid.New().String()

	if err := store.StoreRefreshToken(refreshToken, user.ID); err != nil {
		t.Fatalf("store refresh token returned and unexpected error: %v", err)
	}

	userID, err := store.GetUserIDByRefreshToken(refreshToken)

	if err != nil {
		t.Fatalf("get userId by refresh token returned and unexpected error: %v", err)
	}

	if userID != user.ID {
		t.Fatal("get userId by refresh token returned an incorrect userID")
	}

	if err := store.DeleteRefreshToken(refreshToken); err != nil {
		t.Fatalf("delete refresh toke returned an unexpected error: %v", err)
	}

	if _, err = store.GetUserIDByRefreshToken(refreshToken); err == nil {
		t.Fatal("was able to find userID by refresh token, after refresh token deletion")
	}
}

func TestRefreshTokenRotation(t *testing.T) {
	store := createTestStore()
	user := createTestuser()

	tokenA := uuid.New().String()
	tokenB := uuid.New().String()

	store.StoreRefreshToken(tokenA, user.ID)
	store.StoreRefreshToken(tokenB, user.ID)

	_, err := store.GetUserIDByRefreshToken(tokenA)

	if err == nil {
		t.Fatal("expected tokenA to be invalid after rotation")
	}
}
