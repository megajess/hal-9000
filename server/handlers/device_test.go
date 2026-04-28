package handlers

import (
	"encoding/json"
	"fmt"
	"hal/middleware"
	"hal/models"
	"hal/store"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandleRegisterDevice(t *testing.T) {
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	body := strings.NewReader(`{ "name" : "Test Device" }`)
	req := requestWithUser(http.MethodPost, "/devices", body, "2-444-66666")
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	deviceHandler.HandleRegisterDevice(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d instead", w.Code)
	}

	var resp struct {
		models.Device
		APIKey string `json:"api_key"`
	}

	json.NewDecoder(w.Body).Decode(&resp)

	if resp.APIKey == "" {
		t.Error("expected a non-empty API key")
	}
}

func TestHandleRegisterDevice_missingNameValue(t *testing.T) {
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	body := strings.NewReader(`{ "name" : "" }`)
	req := requestWithUser(http.MethodPost, "/devices", body, "2-444-66666")
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	deviceHandler.HandleRegisterDevice(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d instead", w.Code)
	}
}

func TestHandleRegisterDevice_missingNameKey(t *testing.T) {
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	body := strings.NewReader(`{ }`)
	req := requestWithUser(http.MethodPost, "/devices", body, "2-444-66666")
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	deviceHandler.HandleRegisterDevice(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d instead", w.Code)
	}
}

func TestHandleRegisterDevice_missingUserID(t *testing.T) {
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	body := strings.NewReader(`{ "name" : "Test Device" }`)
	req := requestWithUser(http.MethodPost, "/devices", body, "")
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	deviceHandler.HandleRegisterDevice(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d instead", w.Code)
	}
}

func TestHandleDeviceList(t *testing.T) {
	testUserAID := "2-444-66666"

	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	device := createTestDevice(store, testUserAID)

	var resp []models.Device

	req := requestWithUser(http.MethodGet, "/devices", nil, testUserAID)
	w := httptest.NewRecorder()

	deviceHandler.HandleDeviceList(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected a 200, but got %d instead", w.Code)
	}

	json.NewDecoder(w.Body).Decode(&resp)

	if len(resp) != 1 {
		t.Errorf("expected 1 device, got %d instead", len(resp))
	}

	if len(resp) > 0 && resp[0].ID != device.ID {
		t.Error("device list returned does not match test device")
	}
}

func TestHandleDeviceList_emptyList(t *testing.T) {
	testUserAID := "2-444-66666"
	testUserBID := "1223334444"

	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	createTestDevice(store, testUserAID)

	var resp []models.Device

	req := requestWithUser(http.MethodGet, "/devices", nil, testUserBID)
	w := httptest.NewRecorder()

	deviceHandler.HandleDeviceList(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected a 200, but got %d instead", w.Code)
	}

	json.NewDecoder(w.Body).Decode(&resp)

	if len(resp) != 0 {
		t.Errorf("expected empty device list, got %d devices instead", len(resp))
	}
}

func TestHandleGetDevice(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	device := createTestDevice(store, testUserID)

	req := requestWithUser(http.MethodGet, "/devices", nil, testUserID)
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleGetDevice(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected a 200, got %d instead", w.Code)
	}
}

func TestHandleGetDevice_invalidUser(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	device := createTestDevice(store, testUserID)

	req := requestWithUser(http.MethodGet, "/devices", nil, "1223334444")
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleGetDevice(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected a 404, got %d instead", w.Code)
	}
}

func TestHandleGetDevice_missingDeviceID(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	createTestDevice(store, testUserID)

	req := requestWithUser(http.MethodGet, "/devices", nil, testUserID)

	w := httptest.NewRecorder()

	deviceHandler.HandleGetDevice(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected a 404, got %d instead", w.Code)
	}
}

func TestHandleUpdateDeviceName(t *testing.T) {
	testUserID := "2-444-66666"
	updatedName := "Updated Device Name"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	device := createTestDevice(store, testUserID)

	var resp models.Device

	// Update device request

	body := strings.NewReader(fmt.Sprintf(`{ "name" : "%s" }`, updatedName))
	req := requestWithUser(http.MethodPut, "/devices", body, testUserID)
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleUpdateDeviceName(w, req)

	json.NewDecoder(w.Body).Decode(&resp)

	if w.Code != http.StatusOK {
		t.Errorf("expected a 200, got %d instead", w.Code)
	}

	if resp.Name != updatedName {
		t.Errorf("expected name of %s, got %s instead", updatedName, resp.Name)
	}

	// Get device request

	body = strings.NewReader(fmt.Sprintf(`{ "name" : "%s" }`, updatedName))
	req = requestWithUser(http.MethodGet, "/devices", body, testUserID)
	req.SetPathValue("id", device.ID)

	w = httptest.NewRecorder()

	deviceHandler.HandleGetDevice(w, req)

	json.NewDecoder(w.Body).Decode(&resp)

	if w.Code != http.StatusOK {
		t.Errorf("expected a 200, got %d instead", w.Code)
	}

	if resp.Name != updatedName {
		t.Errorf("expected name of %s, got %s instead", updatedName, resp.Name)
	}
}

func TestHandleUpdateDeviceName_invalidUser(t *testing.T) {
	testUserID := "2-444-66666"
	updatedName := "Updated Device Name"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	device := createTestDevice(store, testUserID)

	var resp models.Device

	body := strings.NewReader(fmt.Sprintf(`{ "name" : "%s" }`, updatedName))
	req := requestWithUser(http.MethodPost, "/devices", body, "1223334444")
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleUpdateDeviceName(w, req)

	json.NewDecoder(w.Body).Decode(&resp)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected a 404, got %d instead", w.Code)
	}
}

func TestHandleUpdateDeviceName_missingNameKey(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	device := createTestDevice(store, testUserID)

	var resp models.Device

	body := strings.NewReader(`{  }`)
	req := requestWithUser(http.MethodPost, "/devices", body, testUserID)
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleUpdateDeviceName(w, req)

	json.NewDecoder(w.Body).Decode(&resp)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected a 400, got %d instead", w.Code)
	}
}

func TestHandleUpdateDeviceName_missingNameValue(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	device := createTestDevice(store, testUserID)

	var resp models.Device

	body := strings.NewReader(`{ "name" : "" }`)
	req := requestWithUser(http.MethodPost, "/devices", body, testUserID)
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleUpdateDeviceName(w, req)

	json.NewDecoder(w.Body).Decode(&resp)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected a 400, got %d instead", w.Code)
	}
}

func TestHandleUpdateDeviceState_on(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	device := createTestDevice(store, testUserID)

	var resp models.Device

	body := strings.NewReader(`{ "desired_state" : "on" }`)
	req := requestWithUser(http.MethodPut, "/devices", body, testUserID)
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleUpdateDeviceState(w, req)

	json.NewDecoder(w.Body).Decode(&resp)

	if w.Code != http.StatusOK {
		t.Errorf("expected a 200, got %d instead", w.Code)
	}
}

func TestHandleUpdateDeviceState_off(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	device := createTestDevice(store, testUserID)

	var resp models.Device

	body := strings.NewReader(`{ "desired_state" : "off" }`)
	req := requestWithUser(http.MethodPut, "/devices", body, testUserID)
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleUpdateDeviceState(w, req)

	json.NewDecoder(w.Body).Decode(&resp)

	if w.Code != http.StatusOK {
		t.Errorf("expected a 200, got %d instead", w.Code)
	}
}

func TestHandleUpdateDeviceState_invalidState(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	device := createTestDevice(store, testUserID)

	var resp models.Device

	body := strings.NewReader(`{ "desired_state" : "invalid state" }`)
	req := requestWithUser(http.MethodPut, "/devices", body, testUserID)
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleUpdateDeviceState(w, req)

	json.NewDecoder(w.Body).Decode(&resp)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected a 400, got %d instead", w.Code)
	}
}

func TestHandleUpdateDeviceState_invalidUser(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	device := createTestDevice(store, testUserID)

	var resp models.Device

	body := strings.NewReader(`{ "desired_state" : "on" }`)
	req := requestWithUser(http.MethodPut, "/devices", body, "1223334444")
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleUpdateDeviceState(w, req)

	json.NewDecoder(w.Body).Decode(&resp)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected a 404, got %d instead", w.Code)
	}
}

func TestHandleUpdateDeviceState_noStoredDevice(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	deviceHandler := createDeviceHandler(store)

	var resp models.Device

	body := strings.NewReader(`{ "desired_state" : "on" }`)
	req := requestWithUser(http.MethodPut, "/devices", body, testUserID)
	req.SetPathValue("id", "1223334444")

	w := httptest.NewRecorder()

	deviceHandler.HandleUpdateDeviceState(w, req)

	json.NewDecoder(w.Body).Decode(&resp)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected a 404, got %d instead", w.Code)
	}
}

func TestHandleDeleteDevice(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	device := createTestDevice(store, testUserID)

	deviceHandler := createDeviceHandler(store)

	// Delete request

	req := requestWithUser(http.MethodDelete, "/devices", nil, testUserID)
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleDeleteDevice(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected a 204, got %d instead", w.Code)
	}

	// Get device request

	req = requestWithUser(http.MethodDelete, "/devices", nil, testUserID)
	req.SetPathValue("id", device.ID)

	w = httptest.NewRecorder()

	deviceHandler.HandleGetDevice(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected a 404, got %d instead", w.Code)
	}
}

func TestHandleDeleteDevice_deviceNotFound(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()

	createTestDevice(store, testUserID)

	deviceHandler := createDeviceHandler(store)

	req := requestWithUser(http.MethodDelete, "/devices", nil, testUserID)
	req.SetPathValue("id", testUserID)

	w := httptest.NewRecorder()

	deviceHandler.HandleDeleteDevice(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected a 404, got %d instead", w.Code)
	}
}

func TestHandleDeleteDevice_invalidUserID(t *testing.T) {
	testUserID := "2-444-66666"
	store := createTestStore()
	device := createTestDevice(store, testUserID)

	deviceHandler := createDeviceHandler(store)

	req := requestWithUser(http.MethodDelete, "/devices", nil, "1223334444")
	req.SetPathValue("id", device.ID)

	w := httptest.NewRecorder()

	deviceHandler.HandleDeleteDevice(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected a 404, got %d instead", w.Code)
	}
}

// MARK: Private helper functions

func createDeviceHandler(s *store.MemoryStore) *DeviceHandler {
	return NewDeviceHandler(s)
}

func requestWithUser(method, path string, body io.Reader, userID string) *http.Request {
	req := httptest.NewRequest(method, path, body)
	ctx := middleware.ContextWithUserIDForTesting(req.Context(), userID)

	return req.WithContext(ctx)
}

func createTestDevice(s *store.MemoryStore, userID string) models.Device {
	device := models.Device{
		ID:           "test-device-id",
		UserID:       userID,
		Name:         "Test Device",
		APIKey:       "test-api-key",
		CurrentState: "off",
		DesiredState: "on",
		CreatedAt:    time.Now(),
	}
	s.CreateDevice(device)
	return device
}
