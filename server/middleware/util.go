package middleware

import (
	"context"
	"hal/models"
)

type contextKey string

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)

	return id, ok
}

func DeviceFromContext(ctx context.Context) (models.Device, bool) {
	device, ok := ctx.Value(deviceKey).(models.Device)
	return device, ok
}
