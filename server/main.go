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

	if os.Getenv("HAL_ENV") == "development" {
		mux.Handle("GET /", http.FileServer(http.Dir("static")))
		log.Println("Test harness available at http://localhost:" + port)
	}

	mux.Handle("GET /poll", apiKeyMiddleware.Require(http.HandlerFunc(deviceHandler.HandlePoll)))

	// Device handlers
	mux.Handle("GET /devices", authMiddleware.Require(http.HandlerFunc(deviceHandler.HandleDeviceList)))
	mux.Handle("POST /devices", authMiddleware.Require(http.HandlerFunc(deviceHandler.HandleRegisterDevice)))
	mux.Handle("GET /devices/{id}", authMiddleware.Require(http.HandlerFunc(deviceHandler.HandleGetDevice)))
	mux.Handle("PUT /devices/{id}/name", authMiddleware.Require(http.HandlerFunc(deviceHandler.HandleUpdateDeviceName)))
	mux.Handle("PUT /devices/{id}/state", authMiddleware.Require(http.HandlerFunc(deviceHandler.HandleUpdateDeviceState)))
	mux.Handle("DELETE /devices/{id}", authMiddleware.Require(http.HandlerFunc(deviceHandler.HandleDeleteDevice)))

	// Auth handlers
	mux.HandleFunc("POST /auth/register", authHandler.HandleRegistration)
	mux.HandleFunc("POST /auth/login", authHandler.HandleLogin)
	mux.HandleFunc("POST /auth/refresh", authHandler.HandleRefresh)
	mux.HandleFunc("POST /auth/logout", authHandler.HandleLogout)

	addr := fmt.Sprintf(":%s", port)

	log.Printf("HAL server listening on %s", addr)
	// TODO: Configure ReadTimeout, WriteTimeout, and IdleTimeout on an
	// http.Server before deploying. http.ListenAndServe with no timeouts
	// is vulnerable to slowloris-style attacks.
	log.Fatal(http.ListenAndServe(addr, mux))
}
