package api

import (
	"encoding/json"
	"net/http"

	"github.com/solomonolatunji/vessel/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type authPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// handleRegister creates a new user account, hashes the password with bcrypt, issues a JWT token, and sets a secure session cookie.
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload authPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid authentication payload format")
		return
	}

	if payload.Email == "" || payload.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required fields")
		return
	}

	existing, err := s.store.GetUserByEmail(payload.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existing != nil {
		writeError(w, http.StatusConflict, "user account with this email already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to encrypt user password")
		return
	}

	role := payload.Role
	if role == "" {
		role = "member"
	}

	user := &types.User{
		Email:        payload.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := s.store.CreateUser(user); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := s.tokenService.GenerateToken(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to issue authentication token")
		return
	}

	setAuthCookie(w, token)
	user.PasswordHash = ""
	writeJSON(w, http.StatusCreated, map[string]any{
		"token": token,
		"user":  user,
	})
}

// handleLogin validates credentials against stored bcrypt hashes, issues a JWT token, and sets a session cookie.
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload authPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid login credentials format")
		return
	}

	user, err := s.store.GetUserByEmail(payload.Email)
	if err != nil || user == nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := s.tokenService.GenerateToken(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to issue authentication token")
		return
	}

	setAuthCookie(w, token)
	user.PasswordHash = ""
	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  user,
	})
}

// handleGetCurrentUser returns the profile details of the currently authenticated user session.
func (s *Server) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	user, err := s.store.GetUserByID(claims.UserID)
	if err != nil || user == nil {
		writeError(w, http.StatusNotFound, "user account not found")
		return
	}

	user.PasswordHash = ""
	writeJSON(w, http.StatusOK, user)
}

// handleLogout clears the vessel_token HTTP-only session cookie from the client browser.
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "vessel_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	writeJSON(w, http.StatusOK, map[string]string{"status": "logged_out"})
}

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "vessel_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   72 * 3600,
	})
}
