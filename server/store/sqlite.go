package store

import (
	"database/sql"
	"errors"
	"hal/models"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

var _ Store = (*SQLiteStore)(nil)

type SQLiteStore struct {
	db  *sqlx.DB
	Now func() time.Time
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sqlx.Open("sqlite3", dbPath+"?_foreign_keys=on")

	if err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db, Now: time.Now}, nil
}

func RunMigrations(dbPath string) error {
	m, err := migrate.New("file://migrations", "sqlite3://"+dbPath)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (s *SQLiteStore) CreateDevice(device models.Device, ownerUserID string) error {
	tx, err := s.db.Beginx()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = s.db.NamedExec(`INSERT INTO devices (
		id,
		name,
		api_key,
		current_state,
		desired_state,
		last_seen,
		created_at
	) VALUES (
	 	:id,
		:name,
		:api_key,
		:current_state,
		:desired_state,
		:last_seen,
		:created_at
	)`, device)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO device_users (device_id, user_id, role) VALUES (?, ?, ?)`,
		device.ID, ownerUserID, "owner",
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLiteStore) GetDeviceByAPIKey(apiKey string) (models.Device, error) {
	var device models.Device

	err := s.db.Get(&device, "SELECT * FROM devices WHERE api_key = ?", apiKey)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Device{}, ErrDeviceNotFound
	}

	return device, err
}

func (s *SQLiteStore) UpdateDeviceState(deviceID string, reportedState string) error {
	results, err := s.db.Exec(`UPDATE devices SET current_state = ?, last_seen = ?`, reportedState, s.Now())

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrDeviceNotFound
	}

	return nil
}

func (s *SQLiteStore) CreateUser(user models.User) error {
	_, err := s.db.NamedExec(`INSERT INTO users (
		id,
		username,
		password_hash,
		created_at
	) VALUES (
		:id,
		:username,
		:password_hash,
		:created_at
	)`, user)

	var sqliteErr sqlite3.Error

	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return ErrUsernameTaken
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStore) GetUserByUsername(username string) (models.User, error) {
	var user models.User

	err := s.db.Get(&user, "SELECT * FROM users WHERE username = ?", username)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}

	return user, err
}

func (s *SQLiteStore) GetUserByID(ID string) (models.User, error) {
	var user models.User

	err := s.db.Get(&user, "SELECT * FROM users WHERE id = ?", ID)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}

	return user, err
}

func (s *SQLiteStore) StoreRefreshToken(token string, userID string) error {
	_, err := s.db.Exec(`INSERT INTO refresh_tokens
		(token, user_id, expires_at)
		VALUES (?, ?, ?)`, token, userID, s.Now())

	return err
}

func (s *SQLiteStore) GetUserIDByRefreshToken(token string) (string, error) {
	var entry struct {
		userID    string    `db:"user_id"`
		ExpiresAt time.Time `db:"expires_at"`
	}

	err := s.db.Get(&entry, "SELECT * FROM refresh_tokens WHERE token = ?", token)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrRefreshTokenNotFound
	}

	if err != nil {
		return "", err
	}

	if s.Now().After(entry.ExpiresAt) {
		return "", ErrRefreshTokenExpired
	}

	return entry.userID, nil
}

func (s *SQLiteStore) DeleteRefreshToken(token string) error {
	results, err := s.db.Exec("DELETE FROM refresh_tokens WHERE token = ?", token)

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}

func (s *SQLiteStore) GetDeviceByID(deviceID string) (models.Device, error) {
	var device models.Device

	err := s.db.Get(&device, "SELECT * FROM devices WHERE id = ?", deviceID)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Device{}, ErrDeviceNotFound
	}

	return device, err
}

func (s *SQLiteStore) GetDevicesByUserID(userID string) ([]models.Device, error) {
	var devices []models.Device

	err := s.db.Select(&devices, `
	SELECT d.* FROM devices d
	JOIN device_users du ON d.id = du.device_id
	WHERE du.user_id = ?`, userID)

	return devices, err
}

func (s *SQLiteStore) UpdateDevice(device models.Device) error {
	results, err := s.db.NamedExec(`UPDATE devices SET
		id = :id,
		name = :name,
		api_key = :api_key,
		current_state = :current_state,
		desired_state = :desired_state,
		last_seen = :last_seen,
		created_at = :created_at`, device)

	rows, err := results.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrDeviceNotFound
	}

	return nil
}

func (s *SQLiteStore) DeleteDevice(deviceID string) error {
	results, err := s.db.Exec("DELETE FROM devices WHERE id = ?", deviceID)

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrDeviceNotFound
	}

	return nil
}

func (s *SQLiteStore) AddUserToDevice(deviceID, userID, role string) error {
	_, err := s.db.Exec(`INSERT INTO device_users (device_id, user_id, role) VALUES (?, ?, ?)`, deviceID, userID, role)

	return err
}

func (s *SQLiteStore) GetUsersByDeviceID(deviceID string) ([]models.DeviceUser, error) {
	var deviceUsers []models.DeviceUser

	err := s.db.Select(&deviceUsers, `SELECT du.*, u.username
		FROM device_users
		INNER JOIN users u ON u.id = du.user_id
		WHERE du.device_id = ?`, deviceID)

	return deviceUsers, err
}

func (s *SQLiteStore) RemoveUserFromDevice(deviceID, userID string) error {
	results, err := s.db.Exec(`DELETE FROM device_users WHERE device_id = ? AND user_id = ?`, deviceID, userID)

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotAssociated
	}

	return nil
}
