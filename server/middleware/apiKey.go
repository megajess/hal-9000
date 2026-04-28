package middleware

import (
	"context"
	"hal/store"
	"net/http"
)

const deviceKey contextKey = "device"

type APIKeyMiddleware struct {
	store store.Store
}

func NewAPIKeyMiddleware(s store.Store) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		store: s,
	}
}

func (m *APIKeyMiddleware) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			http.Error(w, "invalid api key", http.StatusUnauthorized)

			return
		}

		device, err := m.store.GetDeviceByAPIKey(apiKey)

		if err != nil {
			http.Error(w, "invalid api key", http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(r.Context(), deviceKey, device)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
