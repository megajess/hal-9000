package store

import (
	"errors"
	"hal/models"
	"sync"
	"time"
)

var ErrDeviceNotFound = errors.New("device not found")

type MemoryStore struct {
	mu      sync.RWMutex
	devices map[string]models.Device // keyed by device ID
	byKey   map[string]string        // api key → device ID
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		devices: make(map[string]models.Device),
		byKey:   make(map[string]string),
	}
}

func (s *MemoryStore) CreateDevice(device models.Device) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.devices[device.ID] = device
	s.byKey[device.APIKey] = device.ID

	return nil
}

func (s *MemoryStore) GetDeviceByAPIKey(apiKey string) (models.Device, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	id, ok := s.byKey[apiKey]

	if !ok {
		return models.Device{}, ErrDeviceNotFound
	}

	return s.devices[id], nil
}

func (s *MemoryStore) UpdateDeviceState(deviceID string, reportedState string, timeOfUpdate ...time.Time) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	lastSeen := time.Now()

	if len(timeOfUpdate) > 0 {
		lastSeen = timeOfUpdate[0]
	}

	device, ok := s.devices[deviceID]

	if !ok {
		return ErrDeviceNotFound
	}

	device.CurrentState = reportedState
	device.LastSeen = lastSeen
	s.devices[deviceID] = device

	return nil
}
