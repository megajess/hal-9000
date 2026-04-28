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

// ContextWithUserIDForTesting returns a copy of ctx with the given userID injected.
// This is intended for use in tests only — in production the userID is
// set by the AuthMiddleware.Require handler.
func ContextWithUserIDForTesting(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
