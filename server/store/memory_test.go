package store_test

import (
	"hal/models"
	"hal/store"
	"testing"
	"time"
)

const (
	id     = "some-device-id"
	name   = "some-name"
	apiKey = "2-444-66666"
)

func createTestStore() *store.MemoryStore {
	return store.NewMemoryStore()
}

func createDevice() models.Device {
	return models.Device{
		ID:           id,
		Name:         name,
		APIKey:       apiKey,
		CurrentState: "off",
		DesiredState: "on",
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

	if device.Name != name {
		t.Errorf("Expected name %q, got %q instead", name, device.Name)
	}

	if device.ID != id {
		t.Errorf("Expected ID %q, got %q instead", id, device.ID)
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
