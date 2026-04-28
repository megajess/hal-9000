package store

import (
	"hal/models"
	"sync"
	"time"
)

var _ Store = (*MemoryStore)(nil)

type MemoryStore struct {
	mu sync.RWMutex

	devices      map[string]models.Device // [devideID : Device]
	devicesByKey map[string]string        // [apiKey : deviceID]

	users           map[string]models.User // [userID : User]
	usersByUsername map[string]string      // [username : userID]

	userRefreshTokens map[string]refreshEntry // [token : refreshEntry]
}

type refreshEntry struct {
	userID    string
	expiresAt time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		devices:           make(map[string]models.Device),
		devicesByKey:      make(map[string]string),
		users:             make(map[string]models.User),
		usersByUsername:   make(map[string]string),
		userRefreshTokens: make(map[string]refreshEntry),
	}
}

func (s *MemoryStore) CreateDevice(device models.Device) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.devices[device.ID] = device
	s.devicesByKey[device.APIKey] = device.ID

	return nil
}

func (s *MemoryStore) GetDeviceByAPIKey(apiKey string) (models.Device, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	id, ok := s.devicesByKey[apiKey]

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

func (s *MemoryStore) CreateUser(user models.User) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	if _, ok := s.usersByUsername[user.Username]; ok {
		return ErrUsernameTaken
	}

	s.users[user.ID] = user
	s.usersByUsername[user.Username] = user.ID

	return nil
}

func (s *MemoryStore) GetUserByUsername(username string) (models.User, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	id, ok := s.usersByUsername[username]

	if !ok {
		return models.User{}, ErrUserNotFound
	}

	return s.users[id], nil
}

func (s *MemoryStore) GetUserByID(ID string) (models.User, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	user, ok := s.users[ID]

	if !ok {
		return models.User{}, ErrUserNotFound
	}

	return user, nil
}

func (s *MemoryStore) StoreRefreshToken(token string, userID string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	refreshEntry := refreshEntry{
		userID:    userID,
		expiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	s.userRefreshTokens[token] = refreshEntry

	return nil
}

func (s *MemoryStore) GetUserIDByRefreshToken(token string) (string, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	refreshEntry, ok := s.userRefreshTokens[token]

	if !ok || time.Now().After(refreshEntry.expiresAt) {
		return "", ErrRefreshTokenNotFound
	}

	return refreshEntry.userID, nil
}

func (s *MemoryStore) DeleteRefreshToken(token string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	_, ok := s.userRefreshTokens[token]

	if !ok {
		return ErrRefreshTokenNotFound
	}

	delete(s.userRefreshTokens, token)

	return nil
}

func (s *MemoryStore) GetDeviceByID(deviceID string) (models.Device, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	device, ok := s.devices[deviceID]

	if !ok {
		return models.Device{}, ErrDeviceNotFound
	}

	return device, nil
}

func (s *MemoryStore) GetDevicesByUserID(userID string) ([]models.Device, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	devices := []models.Device{}

	for _, device := range s.devices {
		if device.UserID == userID {
			devices = append(devices, device)
		}
	}

	return devices, nil
}

func (s *MemoryStore) UpdateDevice(device models.Device) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	_, ok := s.devices[device.ID]

	if !ok {
		return ErrDeviceNotFound
	}

	s.devices[device.ID] = device

	return nil
}

func (s *MemoryStore) DeleteDevice(deviceID string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	device, ok := s.devices[deviceID]

	if !ok {
		return ErrDeviceNotFound
	}

	delete(s.devices, deviceID)
	delete(s.devicesByKey, device.APIKey)

	return nil
}
