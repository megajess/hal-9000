package main

import (
	"fmt"
	"hal/handlers"
	"hal/store"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("HAL_PORT")
	if port == "" {
		port = "8080"
	}

	s := store.NewMemoryStore()

	deviceHandler := handlers.NewDeviceHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /poll", deviceHandler.HandlePoll)
	mux.HandleFunc("POST /devices", deviceHandler.HandleRegisterDevice)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("HAL server listening on %s", addr)
	// TODO: Configure ReadTimeout, WriteTimeout, and IdleTimeout on an
	// http.Server before deploying. http.ListenAndServe with no timeouts
	// is vulnerable to slowloris-style attacks.
	log.Fatal(http.ListenAndServe(addr, mux))
}
