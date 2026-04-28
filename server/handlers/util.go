package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func writeJSONResponse(w http.ResponseWriter, status int, resp any) {
	var buffer bytes.Buffer

	if err := json.NewEncoder(&buffer).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buffer.Bytes())
}

func generateAPIKey() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func generateID() string {
	return uuid.New().String()
}
