package models

type DeviceUser struct {
	DeviceID string `db:"device_id" json:"device_id"`
	UserID   string `db:"user_id"   json:"user_id"`
	Username string `db:"username"  json:"username"`
	Role     string `db:"role"      json:"role"`
}
