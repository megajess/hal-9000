package handlers

import (
	"encoding/json"
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

	deviceId := generateID()

	userID, ok := middleware.UserIDFromContext(r.Context())

	if userID == "" || !ok {
		http.Error(w, "could not get userID from context", http.StatusUnauthorized)

		return
	}

	// NOTE: Desired state is set to "on" during regristration to provide a physical
	// verification that regisration has succeeded
	device := models.Device{
		ID:           deviceId,
		Name:         req.Name,
		APIKey:       apiKey,
		CurrentState: "off",
		DesiredState: "on",
		CreatedAt:    time.Now(),
	}

	if err := h.store.CreateDevice(device, userID); err != nil {
		http.Error(w, "failed to create device", http.StatusInternalServerError)
		return
	}

	resp := struct {
		models.Device
		APIKey string `json:"api_key"`
	}{
		Device: device,
		APIKey: device.APIKey,
	}

	writeJSONResponse(w, http.StatusCreated, resp)
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

	if err := h.store.UpdateDeviceState(device.ID, reportedState); err != nil {
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

func (h *DeviceHandler) HandleDeviceList(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())

	if userID == "" || !ok {
		http.Error(w, "could not get userID from context", http.StatusUnauthorized)

		return
	}

	devices, err := h.store.GetDevicesByUserID(userID)

	if err != nil {
		http.Error(w, "error getting devices", http.StatusInternalServerError)

		return
	}

	writeJSONResponse(w, http.StatusOK, devices)
}

func (h *DeviceHandler) HandleGetDevice(w http.ResponseWriter, r *http.Request) {
	device, ok := h.getOwnedDevice(w, r)

	if !ok {
		return
	}

	writeJSONResponse(w, http.StatusOK, device)
}

func (h *DeviceHandler) HandleUpdateDeviceName(w http.ResponseWriter, r *http.Request) {
	device, ok := h.getOwnedDevice(w, r)

	if !ok {
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "could not decode name from request", http.StatusBadRequest)

		return
	}

	newDeviceName := req.Name

	// TODO: Should there be more checks for a valid name?
	if newDeviceName == "" {
		http.Error(w, "must provide a valid name for device", http.StatusBadRequest)

		return
	}

	device.Name = newDeviceName

	err := h.store.UpdateDevice(device)

	if err != nil {
		http.Error(w, "error updating device", http.StatusInternalServerError)

		return
	}

	writeJSONResponse(w, http.StatusOK, device)
}

func (h *DeviceHandler) HandleUpdateDeviceState(w http.ResponseWriter, r *http.Request) {
	device, ok := h.getOwnedDevice(w, r)

	if !ok {
		return
	}

	var req struct {
		DesiredState string `json:"desired_state"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "could not decode desired state from request", http.StatusBadRequest)

		return
	}

	newDeviceDesiredState := req.DesiredState

	if newDeviceDesiredState != "on" && newDeviceDesiredState != "off" {
		http.Error(w, `only desired states of "on", or "off" are accepted`, http.StatusBadRequest)

		return
	}

	device.DesiredState = newDeviceDesiredState

	err := h.store.UpdateDevice(device)

	if err != nil {
		http.Error(w, "error updating device", http.StatusInternalServerError)

		return
	}

	writeJSONResponse(w, http.StatusOK, device)
}

func (h *DeviceHandler) HandleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	device, ok := h.getOwnedDevice(w, r)

	if !ok {
		return
	}

	if err := h.store.DeleteDevice(device.ID); err != nil {
		http.Error(w, "error deleting device", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MARK: Private methods

func (h *DeviceHandler) getOwnedDevice(w http.ResponseWriter, r *http.Request) (models.Device, bool) {
	userID, ok := middleware.UserIDFromContext(r.Context())

	if userID == "" || !ok {
		http.Error(w, "could not get userID from context", http.StatusUnauthorized)

		return models.Device{}, false
	}

	deviceID := r.PathValue("id")
	device, err := h.store.GetDeviceByID(deviceID)

	if err != nil {
		http.Error(w, "device not found", http.StatusNotFound)

		return models.Device{}, false
	}

	users, err := h.store.GetUsersByDeviceID(deviceID)

	if err != nil {
		http.Error(w, "device not found", http.StatusNotFound)

		return models.Device{}, false
	}

	for _, du := range users {
		if du.UserID == userID {
			return device, true
		}
	}

	http.Error(w, "device not found", http.StatusNotFound)

	return models.Device{}, false
}
