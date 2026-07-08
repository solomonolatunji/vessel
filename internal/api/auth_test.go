package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

func TestAuthEndpointsAndRBAC(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "vessel_auth_test")
	_ = os.RemoveAll(tempDir)
	_ = os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "vessel.db")
	s, err := store.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to init store: %v", err)
	}
	defer s.Close()

	srv := NewServer(s, nil, nil, nil)

	registerPayload := []byte(`{"email":"solomon@vessel.dev","password":"securepassword123","role":"admin"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerPayload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on register, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	var registerResp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&registerResp); err != nil {
		t.Fatalf("failed to decode register response: %v", err)
	}

	tokenStr, ok := registerResp["token"].(string)
	if !ok || tokenStr == "" {
		t.Fatalf("expected valid token string, got: %v", registerResp["token"])
	}

	meReqNoAuth := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	meRecNoAuth := httptest.NewRecorder()
	srv.Handler().ServeHTTP(meRecNoAuth, meReqNoAuth)

	if meRecNoAuth.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for unauthenticated /api/auth/me, got %d", meRecNoAuth.Code)
	}

	meReqAuth := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	meReqAuth.Header.Set("Authorization", "Bearer "+tokenStr)
	meRecAuth := httptest.NewRecorder()
	srv.Handler().ServeHTTP(meRecAuth, meReqAuth)

	if meRecAuth.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for authenticated /api/auth/me, got %d. Body: %s", meRecAuth.Code, meRecAuth.Body.String())
	}

	var user types.User
	if err := json.NewDecoder(meRecAuth.Body).Decode(&user); err != nil {
		t.Fatalf("failed to decode user profile: %v", err)
	}

	if user.Email != "solomon@vessel.dev" || user.Role != "admin" {
		t.Errorf("expected solomon@vessel.dev [admin], got %s [%s]", user.Email, user.Role)
	}
	if user.PasswordHash != "" {
		t.Errorf("expected PasswordHash to be stripped in API response, got: %s", user.PasswordHash)
	}
}
