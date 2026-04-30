package store

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRefreshTokenExpiration(t *testing.T) {
	store := NewMemoryStore()
	token := uuid.New().String()
	userID := "2-444-66666"

	store.userRefreshTokens[token] = refreshEntry{
		userID:    userID,
		expiresAt: time.Now().Add(-1 * time.Hour),
	}

	if _, err := store.GetUserIDByRefreshToken(token); err != ErrRefreshTokenExpired {
		t.Fatalf("get userID by refresh token returned an unexpected errpr: %v", err)
	}
}
