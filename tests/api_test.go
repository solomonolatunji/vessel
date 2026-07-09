package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/solomonolatunji/vessel/internal/api"
	"github.com/solomonolatunji/vessel/internal/proxy"
	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

func TestProjectHandlerAndSslipFallback(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "vessel_api_test")
	_ = os.RemoveAll(tempDir)
	_ = os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "vessel.db")
	s, err := store.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to init store: %v", err)
	}
	defer s.Close()

	proxyCfg := proxy.NewCaddyConfig(tempDir, "admin@test.local")
	proxyMgr := proxy.NewProxyManager(proxyCfg, s, nil)

	srv := api.NewServer(s, nil, proxyMgr, nil)

	// Register user & obtain auth token
	regPayload := []byte(`{"email":"admin@vessel.dev","password":"password123","role":"admin"}`)
	regReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regPayload))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(regRec, regReq)

	var tokenStr string
	for _, cookie := range regRec.Result().Cookies() {
		if cookie.Name == "vessel_token" {
			tokenStr = cookie.Value
			break
		}
	}

	payload := []byte(`{"name":"My Node Service","internal_port":3000}`)
	req := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	if tokenStr != "" {
		req.Header.Set("Authorization", "Bearer "+tokenStr)
	}
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d. Body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var created types.ProjectConfig
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if created.Name != "My Node Service" {
		t.Errorf("expected name 'My Node Service', got '%s'", created.Name)
	}

	services, err := s.ListAppServicesByProject(created.ID)
	if err != nil || len(services) == 0 {
		t.Fatalf("expected auto-created service, got %d (err: %v)", len(services), err)
	}

	if !strings.HasSuffix(services[0].Domain, ".sslip.io") {
		t.Errorf("expected auto-generated sslip.io fallback domain on service, got '%s'", services[0].Domain)
	}
}
