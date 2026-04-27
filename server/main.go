package main

import (
	"fmt"
	"hal/handlers"
	"hal/middleware"
	"hal/store"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv("HAL_PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("HAL_JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("HAL_JWT_SECRET environment variable is required")
	}

	s := store.NewMemoryStore()

	deviceHandler := handlers.NewDeviceHandler(s)
	authHandler := handlers.NewAuthHandler(s, jwtSecret)

	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)
	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(s)

	mux := http.NewServeMux()

	// Device handlers
	mux.Handle("GET /poll", apiKeyMiddleware.Require(http.HandlerFunc(deviceHandler.HandlePoll)))
	mux.Handle("POST /devices", authMiddleware.Require(http.HandlerFunc(deviceHandler.HandleRegisterDevice)))
	mux.Handle("PATCH /devices/state", authMiddleware.Require(http.HandlerFunc(deviceHandler.HandleUpdateState)))

	// Auth handlers
	mux.HandleFunc("POST /auth/register", authHandler.HandleRegistration)
	mux.HandleFunc("POST /auth/login", authHandler.HandleLogin)
	mux.HandleFunc("POST /auth/refresh", authHandler.HandleRefresh)

	addr := fmt.Sprintf(":%s", port)

	log.Printf("HAL server listening on %s", addr)
	// TODO: Configure ReadTimeout, WriteTimeout, and IdleTimeout on an
	// http.Server before deploying. http.ListenAndServe with no timeouts
	// is vulnerable to slowloris-style attacks.
	log.Fatal(http.ListenAndServe(addr, mux))
}
