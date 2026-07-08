package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

	srv := NewServer(s, nil, proxyMgr, nil)

	payload := []byte(`{"name":"My Node Service","internal_port":3000}`)
	req := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
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

	if !strings.HasSuffix(created.Domain, ".sslip.io") {
		t.Errorf("expected auto-generated sslip.io fallback domain, got '%s'", created.Domain)
	}
}
