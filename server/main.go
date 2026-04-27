package main

import (
	"fmt"
	"hal/handlers"
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

	mux := http.NewServeMux()

	// Device handlers
	mux.HandleFunc("GET /poll", deviceHandler.HandlePoll)
	mux.HandleFunc("POST /devices", deviceHandler.HandleRegisterDevice)
	mux.HandleFunc("PATCH /devices/state", deviceHandler.HandleUpdateState)

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
