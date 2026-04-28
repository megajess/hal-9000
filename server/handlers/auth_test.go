package handlers

import (
	"encoding/json"
	"fmt"
	"hal/store"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleRegistration(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)

	body := strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d instead", w.Code)
	}
}

func TestHandleRegistration_missingUsername(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)
	body := strings.NewReader(`{ "password" : "P@ssword123" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d instead", w.Code)
	}
}

func TestHandleRegistration_emptyUsername(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)
	body := strings.NewReader(`{ "username" : "", "password" : "P@ssword123" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d instead", w.Code)
	}
}

func TestHandleRegistration_missingPassword(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)
	body := strings.NewReader(`{ "username" : "billiam" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d instead", w.Code)
	}
}

func TestHandleRegistration_emptyPassword(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)
	body := strings.NewReader(`{ "username" : "billiam", "password" : "" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d instead", w.Code)
	}
}

func TestPasswordTooLong(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)
	body := strings.NewReader(`{ "username" : "billiam", "password" : "a really long password that should be over 72 bytes, so lets add some more characters to get it to be big enough for this test, that is testing how the system handles passwords that are too long!" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d instead", w.Code)
	}
}

func TestDuplicateUsername(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)

	// First request

	body := strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d instead", w.Code)
	}

	// Second request with duplicate username

	body = strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req = httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w = httptest.NewRecorder()

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d instead", w.Code)
	}
}

func TestLogin(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)

	// Register request

	body := strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d instead", w.Code)
	}

	// Login request

	body = strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	w = httptest.NewRecorder()

	authHandler.HandleLogin(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d instead", w.Code)
	}

	var resp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	json.NewDecoder(w.Body).Decode(&resp)

	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}

	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestLogout(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)

	// Register request

	body := strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d instead", w.Code)
	}

	// Login request

	body = strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	w = httptest.NewRecorder()

	authHandler.HandleLogin(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d instead", w.Code)
	}

	var resp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	json.NewDecoder(w.Body).Decode(&resp)

	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}

	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}

	// Logout request

	body = strings.NewReader(fmt.Sprintf(`{ "refresh_token" : "%s" }`, resp.RefreshToken))
	req = httptest.NewRequest(http.MethodPost, "/auth/logout", body)
	w = httptest.NewRecorder()

	authHandler.HandleLogout(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d instead", w.Code)
	}

	// Refresh request

	body = strings.NewReader(fmt.Sprintf(`{ "refresh_token" : "%s" }`, resp.RefreshToken))
	req = httptest.NewRequest(http.MethodPost, "/auth/refresh", body)
	w = httptest.NewRecorder()

	authHandler.HandleRefresh(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d instead", w.Code)
	}
}

func TestLogin_wrongPassword(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)

	// Register request

	body := strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d instead", w.Code)
	}

	// Login request

	body = strings.NewReader(`{ "username" : "billiam", "password" : "NOPE" }`)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	w = httptest.NewRecorder()

	authHandler.HandleLogin(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d instead", w.Code)
	}
}

func TestLogin_unknownUser(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)

	body := strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", body)
	w := httptest.NewRecorder()

	authHandler.HandleLogin(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d instead", w.Code)
	}
}

func TestHandleRefresh(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)

	// Register request

	body := strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d instead", w.Code)
	}

	// Login request

	body = strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	w = httptest.NewRecorder()

	authHandler.HandleLogin(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d instead", w.Code)
	}

	var resp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	json.NewDecoder(w.Body).Decode(&resp)

	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}

	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}

	// Login request using refresh token

	refreshTokenBodyString := fmt.Sprintf(`{ "refresh_token" : "%s" }`, resp.RefreshToken)
	body = strings.NewReader(refreshTokenBodyString)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	w = httptest.NewRecorder()

	authHandler.HandleRefresh(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d instead", w.Code)
	}

	json.NewDecoder(w.Body).Decode(&resp)

	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}

	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestHandleRefresh_invalidToken(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)

	// Register request

	body := strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d instead", w.Code)
	}

	// Login request

	body = strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	w = httptest.NewRecorder()

	authHandler.HandleLogin(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d instead", w.Code)
	}

	var resp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	json.NewDecoder(w.Body).Decode(&resp)

	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}

	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}

	// Login request with invalid refresh token

	body = strings.NewReader(`{ "refresh_token" : "not a valid token" }`)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	w = httptest.NewRecorder()

	authHandler.HandleRefresh(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d instead", w.Code)
	}
}

func TestHandleRefresh_tokenRotation(t *testing.T) {
	store := createTestStore()
	authHandler := createAuthHandler(store)

	// Register request

	body := strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	authHandler.HandleRegistration(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d instead", w.Code)
	}

	// Login request

	body = strings.NewReader(`{ "username" : "billiam", "password" : "P@ssword123" }`)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	w = httptest.NewRecorder()

	authHandler.HandleLogin(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d instead", w.Code)
	}

	var resp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	json.NewDecoder(w.Body).Decode(&resp)

	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}

	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}

	refreshToken := resp.RefreshToken

	// Login request with refresh token

	refreshTokenBodyString := fmt.Sprintf(`{ "refresh_token" : "%s" }`, refreshToken)
	body = strings.NewReader(refreshTokenBodyString)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	w = httptest.NewRecorder()

	authHandler.HandleRefresh(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d instead", w.Code)
	}

	json.NewDecoder(w.Body).Decode(&resp)

	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}

	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}

	// Second login request with refresh token

	refreshTokenBodyString = fmt.Sprintf(`{ "refresh_token" : "%s" }`, refreshToken)
	body = strings.NewReader(refreshTokenBodyString)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	w = httptest.NewRecorder()

	authHandler.HandleRefresh(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d instead", w.Code)
	}
}

// MARK: Private helper functions

func createTestStore() *store.MemoryStore {
	return store.NewMemoryStore()
}

func createAuthHandler(store *store.MemoryStore) *AuthHandler {
	return NewAuthHandler(store, "shhh-its-a-secret")
}
