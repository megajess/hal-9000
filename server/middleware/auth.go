package middleware

import (
	"context"
	"fmt"
	"hal/models"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const userIDKey contextKey = "userID"

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: secret,
	}
}

func (m *AuthMiddleware) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		_, tokenString, ok := strings.Cut(authHeader, "Bearer ")

		if tokenString == "" || !ok {
			http.Error(w, "invalid authorization", http.StatusUnauthorized)

			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid authorization", http.StatusUnauthorized)

			return
		}

		claims, ok := token.Claims.(*models.Claims)

		if !ok {
			http.Error(w, "invalid authorization", http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
