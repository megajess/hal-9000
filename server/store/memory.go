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

	deviceUsers map[string][]models.DeviceUser // [deviceID : []DeviceUser]

	Now func() time.Time
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
		deviceUsers:       make(map[string][]models.DeviceUser),
		Now:               time.Now,
	}
}

func (s *MemoryStore) CreateDevice(device models.Device, ownerUserID string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	user, ok := s.users[ownerUserID]

	if !ok {
		return ErrUserNotFound
	}

	s.devices[device.ID] = device
	s.devicesByKey[device.APIKey] = device.ID
	s.deviceUsers[device.ID] = []models.DeviceUser{
		{
			DeviceID: device.ID,
			UserID:   ownerUserID,
			Username: s.users[user.Username].Username,
			Role:     "owner",
		},
	}

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

func (s *MemoryStore) UpdateDeviceState(deviceID string, reportedState string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	lastSeen := s.Now()

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
	if !ok {
		return "", ErrRefreshTokenNotFound
	}

	if s.Now().After(refreshEntry.expiresAt) {
		return "", ErrRefreshTokenExpired
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

	for deviceID, users := range s.deviceUsers {
		for _, du := range users {
			if du.UserID == userID {
				if device, ok := s.devices[deviceID]; ok {
					devices = append(devices, device)
				}
				break
			}
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
	delete(s.deviceUsers, deviceID)

	return nil
}

func (s *MemoryStore) AddUserToDevice(deviceID, userID, role string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	if _, ok := s.devices[deviceID]; !ok {
		return ErrDeviceNotFound
	}

	user, ok := s.users[userID]

	if !ok {
		return ErrUserNotFound
	}

	for _, du := range s.deviceUsers[deviceID] {
		if du.UserID == userID {
			return nil
		}
	}

	s.deviceUsers[deviceID] = append(s.deviceUsers[deviceID], models.DeviceUser{
		DeviceID: deviceID,
		UserID:   userID,
		Username: user.Username,
		Role:     role,
	})

	return nil
}

func (s *MemoryStore) GetUsersByDeviceID(deviceID string) ([]models.DeviceUser, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	if _, ok := s.devices[deviceID]; !ok {
		return nil, ErrDeviceNotFound
	}

	return append([]models.DeviceUser(nil), s.deviceUsers[deviceID]...), nil
}

func (s *MemoryStore) RemoveUserFromDevice(deviceID, userID string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	users, ok := s.deviceUsers[deviceID]

	if !ok {
		return ErrDeviceNotFound
	}

	for i, du := range users {
		if du.UserID == userID {
			s.deviceUsers[deviceID] = append(users[:i], users[i+1:]...)
			return nil
		}
	}

	return ErrUserNotFound
}
