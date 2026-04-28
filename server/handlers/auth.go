package handlers

import (
	"encoding/json"
	"hal/models"
	"hal/store"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store     store.Store
	jwtSecret string
}

func NewAuthHandler(s store.Store, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		store:     s,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Username == "" || req.Password == "") {
		http.Error(w, "username and password required for registration", http.StatusBadRequest)

		return
	}

	if len(req.Password) > 72 {
		http.Error(w, "password must be 72 characters or less", http.StatusBadRequest)

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "failed to generate password hash", http.StatusInternalServerError)

		return
	}

	user := models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := h.store.CreateUser(user); err != nil {
		if err == store.ErrUsernameTaken {
			http.Error(w, "username already exists", http.StatusConflict)

			return
		}

		http.Error(w, "error creating user", http.StatusInternalServerError)

		return
	}

	resp := struct {
		User models.User `json:"user"`
	}{
		User: user,
	}

	writeJSONResponse(w, http.StatusCreated, resp)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Username == "" || req.Password == "") {
		http.Error(w, "username and password required to login", http.StatusBadRequest)

		return
	}

	user, err := h.store.GetUserByUsername(req.Username)

	if err != nil {
		if err == store.ErrUserNotFound {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)

		return
	}

	accessTokenString, err := h.generateAccessTokenFor(user.ID)

	if err != nil {
		http.Error(w, "failed to create access token", http.StatusInternalServerError)

		return
	}

	refreshToken := h.generateRefreshToken()

	if err := h.store.StoreRefreshToken(refreshToken, user.ID); err != nil {
		http.Error(w, "failed to create refresh token", http.StatusInternalServerError)

		return
	}

	resp := struct {
		User         models.User `json:"user"`
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
	}{
		User:         user,
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
	}

	writeJSONResponse(w, http.StatusOK, resp)
}

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)

		return
	}

	userID, err := h.store.GetUserIDByRefreshToken(req.RefreshToken)

	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)

		return
	}

	accessTokenString, err := h.generateAccessTokenFor(userID)

	if err != nil {
		http.Error(w, "failed to create access token", http.StatusInternalServerError)

		return
	}

	refreshTokenString := h.generateRefreshToken()

	if err := h.store.StoreRefreshToken(refreshTokenString, userID); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	if err := h.store.DeleteRefreshToken(req.RefreshToken); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	resp := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	writeJSONResponse(w, http.StatusOK, resp)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
		h.store.DeleteRefreshToken(req.RefreshToken)
	}

	w.WriteHeader(http.StatusNoContent)
}

// MARK: Private methods

func (h *AuthHandler) generateAccessTokenFor(userID string) (string, error) {
	claims := models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(h.jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *AuthHandler) generateRefreshToken() string {
	return uuid.New().String()
}
