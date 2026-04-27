package store

import (
	"hal/models"
	"sync"
	"time"
)

type MemoryStore struct {
	mu sync.RWMutex

	devices map[string]models.Device // [devideID : Device]
	byKey   map[string]string        // [apiKey : deviceID]

	users      map[string]models.User // [userID : User]
	byUsername map[string]string      // [username : userID]

	refreshTokens     map[string]string       // [userID : token]
	userRefreshTokens map[string]refreshEntry // [token : refreshEntry]
}

type refreshEntry struct {
	userID    string
	expiresAt time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		devices:           make(map[string]models.Device),
		byKey:             make(map[string]string),
		users:             make(map[string]models.User),
		byUsername:        make(map[string]string),
		refreshTokens:     make(map[string]string),
		userRefreshTokens: make(map[string]refreshEntry),
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

func (s *MemoryStore) UpdateDeviceDesiredState(deviceID string, desiredState string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	device, ok := s.devices[deviceID]

	if !ok {
		return ErrDeviceNotFound
	}

	device.DesiredState = desiredState
	s.devices[deviceID] = device

	return nil
}

func (s *MemoryStore) CreateUser(user models.User) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	if _, ok := s.byUsername[user.Username]; ok {
		return ErrUsernameTaken
	}

	s.users[user.ID] = user
	s.byUsername[user.Username] = user.ID

	return nil
}

func (s *MemoryStore) GetUserByUsername(username string) (models.User, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	id, ok := s.byUsername[username]

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

	if oldToken, ok := s.refreshTokens[userID]; ok {
		delete(s.userRefreshTokens, oldToken)
	}

	refreshEntry := refreshEntry{
		userID:    userID,
		expiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	s.refreshTokens[userID] = token
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

	refreshEntry := s.userRefreshTokens[token]

	delete(s.userRefreshTokens, token)
	delete(s.refreshTokens, refreshEntry.userID)

	return nil
}
