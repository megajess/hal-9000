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
	altUserID        = "1223334444"
	altUsername      = "Bobbert"
	userPasswordHash = "qwerty"
)

func createTestStore() *store.MemoryStore {
	return store.NewMemoryStore()
}

func createTestDevice() models.Device {
	return createDevice(deviceID, deviceName, nil, apiKey, "off", "on")
}

func createTestDeviceForUser(user models.User) models.Device {
	return createDevice(deviceID, deviceName, &user.ID, apiKey, "off", "on")
}

func createDevice(id string, name string, userID *string, apiKey string, currentState string, desiredState string) models.Device {
	var _userID string

	if userID != nil {
		_userID = *userID
	}

	return models.Device{
		ID:           id,
		UserID:       _userID,
		Name:         name,
		APIKey:       apiKey,
		CurrentState: currentState,
		DesiredState: desiredState,
		CreatedAt:    time.Now(),
	}
}

func createTestUser() models.User {
	return createUser(userID, username, userPasswordHash)
}

func createAlternateTestUser() models.User {
	return createUser(altUserID, altUsername, userPasswordHash)
}

func createUser(id string, username string, passwordHash string) models.User {
	return models.User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
}

func TestCreateDevice(t *testing.T) {
	store := createTestStore()
	deviceModel := createTestDevice()

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
	device := createTestDevice()

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
	device := createTestDevice()

	if err := store.CreateDevice(device); err != nil {
		t.Fatalf("CreateDevice returned unexpected error: %v", err)
	}

	if err := store.UpdateDeviceState("unknown-device-id", "on"); err == nil {
		t.Fatal("Expected an error for unknown device id, but got nil instead!")
	}
}

func TestCreateUser(t *testing.T) {
	store := createTestStore()
	user := createTestUser()

	if err := store.CreateUser(user); err != nil {
		t.Fatalf("CreateUser returned unexpected error: %v", err)
	}
}

func TestCreateUser_duplicateUsername(t *testing.T) {
	store := createTestStore()
	user := createTestUser()

	if err := store.CreateUser(user); err != nil {
		t.Fatalf("CreateUser returned unexpected error: %v", err)
	}

	if err := store.CreateUser(user); err == nil {
		t.Fatal("did not get error when creating duplicated user")
	}
}

func TestCreateUser_FindByID(t *testing.T) {
	store := createTestStore()
	user := createTestUser()

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
	user := createTestUser()

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
	user := createTestUser()

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
	user := createTestUser()

	tokenA := uuid.New().String()
	tokenB := uuid.New().String()

	store.StoreRefreshToken(tokenA, user.ID)
	store.DeleteRefreshToken(tokenA)
	store.StoreRefreshToken(tokenB, user.ID)

	_, err := store.GetUserIDByRefreshToken(tokenA)

	if err == nil {
		t.Fatal("expected tokenA to be invalid after rotation")
	}
}

func TestGetDeviceByID(t *testing.T) {
	store := createTestStore()
	device := createTestDevice()

	err := store.CreateDevice(device)

	if err != nil {
		t.Fatalf("CreateDevice returned unexpected error: %v", err)
	}

	fetchedDevice, err := store.GetDeviceByID(device.ID)

	if err != nil {
		t.Fatal("got an unexpected error getting device by id")
	}

	if fetchedDevice.ID != deviceID {
		t.Fatal("got incorrect device when getting device by id")
	}
}

func TestGetDevicesbyUserID(t *testing.T) {
	store := createTestStore()
	user := createTestUser()
	device := createTestDeviceForUser(user)

	err := store.CreateDevice(device)

	if err != nil {
		t.Fatalf("CreateDevice returned unexpected error: %v", err)
	}

	fetchedDevices, err := store.GetDevicesByUserID(user.ID)

	if err != nil {
		t.Fatal("got an unexpected error fetching devices by userID")
	}

	if len(fetchedDevices) != 1 {
		t.Fatalf("expected 1 device for user, but got %d instead", len(fetchedDevices))
	}

	if len(fetchedDevices) > 0 && fetchedDevices[0].ID != device.ID {
		t.Fatal("Got incorrect device when fetching by userID")
	}
}

func TestUpdateDevice(t *testing.T) {
	updatedName := "UPDATED"
	store := createTestStore()
	device := createTestDevice()

	err := store.CreateDevice(device)

	if err != nil {
		t.Fatalf("CreateDevice returned unexpected error: %v", err)
	}

	updatedDevice := device

	updatedDevice.Name = updatedName

	if err := store.UpdateDevice(updatedDevice); err != nil {
		t.Fatal("got unexpected error updating device")
	}

	fetchedDevice, err := store.GetDeviceByID(device.ID)

	if err != nil {
		t.Fatal("got unexpected error fetching updated device")
	}

	if fetchedDevice.ID != device.ID {
		t.Fatal("fetched device does not match original device")
	}

	if fetchedDevice.Name != updatedName {
		t.Fatal("device did not update")
	}
}

func TestDeleteDevice(t *testing.T) {
	store := createTestStore()
	device := createTestDevice()

	err := store.CreateDevice(device)

	if err != nil {
		t.Fatalf("CreateDevice returned unexpected error: %v", err)
	}

	err = store.DeleteDevice(device.ID)

	if err != nil {
		t.Fatal("got unexpected error deleting device")
	}

	_, err = store.GetDeviceByID(device.ID)

	if err == nil {
		t.Fatal("device did not delete as expected")
	}
}
