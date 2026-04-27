package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hal/middleware"
	"hal/models"
	"hal/store"
	"net/http"
	"time"
)

type DeviceHandler struct {
	store store.Store
}

func NewDeviceHandler(s store.Store) *DeviceHandler {
	return &DeviceHandler{store: s}
}

// HandleRegisterDevice handles POST /devices
// Creates a new device and returns it with the generated API key.
func (h *DeviceHandler) HandleRegisterDevice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)

		return
	}

	apiKey, err := generateAPIKey()

	if err != nil {
		http.Error(w, "failed to generate API key", http.StatusInternalServerError)

		return
	}

	deviceId, err := generateID()

	if err != nil {
		http.Error(w, "failed to generate device id", http.StatusInternalServerError)

		return
	}

	device := models.Device{
		ID:           deviceId,
		Name:         req.Name,
		APIKey:       apiKey,
		CurrentState: "off",
		DesiredState: "on",
		CreatedAt:    time.Now(),
	}

	if err := h.store.CreateDevice(device); err != nil {
		http.Error(w, "failed to create device", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(device); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

// HandlePoll handles GET /poll
// Called by the device every 1-2 seconds. Reads current state from query
// params, updates the store, then responds with 0 (no change) or 1 (toggle).
func (h *DeviceHandler) HandlePoll(w http.ResponseWriter, r *http.Request) {
	device, ok := middleware.DeviceFromContext(r.Context())

	if !ok {
		http.Error(w, "missing API key", http.StatusUnauthorized)

		return
	}

	// Convert query param (0/1) to state string ("off"/"on")
	reportedState := "off"

	if r.URL.Query().Get("state") == "1" {
		reportedState = "on"
	}

	if err := h.store.UpdateDeviceState(device.ID, reportedState, time.Now()); err != nil {
		http.Error(w, "failed to update state", http.StatusInternalServerError)
		return
	}

	// If current matches desired, no action needed.
	// If they differ, tell the device to toggle.
	if reportedState == device.DesiredState {
		w.Write([]byte("0"))
	} else {
		w.Write([]byte("1"))
	}
}

func (h *DeviceHandler) HandleUpdateState(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")

	if apiKey == "" {
		http.Error(w, "missing API key", http.StatusUnauthorized)
		return
	}

	device, err := h.store.GetDeviceByAPIKey(apiKey)

	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert query param (0/1) to state string ("off"/"on")
	desiredState := "off"

	if r.URL.Query().Get("state") == "1" {
		desiredState = "on"
	}

	if err := h.store.UpdateDeviceDesiredState(device.ID, desiredState); err != nil {
		http.Error(w, "error updating desired state", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func generateAPIKey() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func generateID() (string, error) {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
