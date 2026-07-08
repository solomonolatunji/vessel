package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/solomonolatunji/vessel/internal/api"
	"github.com/solomonolatunji/vessel/internal/proxy"
	"github.com/solomonolatunji/vessel/internal/store"
)

func TestGitEndpointsAndWebhooks(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "vessel_git_test")
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
	regPayload := []byte(`{"email":"gituser@vessel.dev","password":"securepassword123","role":"developer"}`)
	regReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regPayload))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(regRec, regReq)

	if regRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d: %s", regRec.Code, regRec.Body.String())
	}

	// Extract cookie token
	var tokenStr string
	for _, cookie := range regRec.Result().Cookies() {
		if cookie.Name == "vessel_token" {
			tokenStr = cookie.Value
			break
		}
	}
	if tokenStr == "" {
		t.Fatalf("missing vessel_token cookie from register response")
	}

	// 1. Check Git provider status before connection
	statusReq := httptest.NewRequest(http.MethodGet, "/api/git/status", nil)
	statusReq.Header.Set("Authorization", "Bearer "+tokenStr)
	statusRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(statusRec, statusReq)

	if statusRec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on get status, got %d: %s", statusRec.Code, statusRec.Body.String())
	}

	// 2. Connect GitHub provider token
	connectPayload := []byte(`{"provider":"github","accessToken":"ghp_mocktoken123456789","accountName":"octocat"}`)
	connectReq := httptest.NewRequest(http.MethodPost, "/api/git/connect", bytes.NewReader(connectPayload))
	connectReq.Header.Set("Authorization", "Bearer "+tokenStr)
	connectReq.Header.Set("Content-Type", "application/json")
	connectRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(connectRec, connectReq)

	if connectRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 on git connect, got %d: %s", connectRec.Code, connectRec.Body.String())
	}

	// 3. Check status again (should show connected: true for github)
	statusReq2 := httptest.NewRequest(http.MethodGet, "/api/git/status", nil)
	statusReq2.Header.Set("Authorization", "Bearer "+tokenStr)
	statusRec2 := httptest.NewRecorder()
	srv.Handler().ServeHTTP(statusRec2, statusReq2)

	if statusRec2.Code != http.StatusOK {
		t.Fatalf("expected status 200 on get status, got %d: %s", statusRec2.Code, statusRec2.Body.String())
	}

	var providerStatuses []map[string]any
	if err := json.NewDecoder(statusRec2.Body).Decode(&providerStatuses); err != nil {
		t.Fatalf("failed to decode provider statuses: %v", err)
	}

	foundGitHub := false
	for _, ps := range providerStatuses {
		if ps["provider"] == "github" && ps["connected"] == true && ps["accountName"] == "octocat" {
			foundGitHub = true
		}
	}
	if !foundGitHub {
		t.Fatalf("expected connected github status, got %v", providerStatuses)
	}

	// 4. Create a project and trigger Git webhook push notification
	projPayload := []byte(`{"name":"Git Webhook App","repositoryUrl":"https://github.com/octocat/hello-world.git","branch":"main"}`)
	projReq := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewReader(projPayload))
	projReq.Header.Set("Authorization", "Bearer "+tokenStr)
	projReq.Header.Set("Content-Type", "application/json")
	projRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(projRec, projReq)

	if projRec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 on project creation, got %d: %s", projRec.Code, projRec.Body.String())
	}

	var createdProj map[string]any
	_ = json.NewDecoder(projRec.Body).Decode(&createdProj)
	projID := createdProj["id"].(string)

	webhookReq := httptest.NewRequest(http.MethodPost, "/api/webhooks/git/"+projID, bytes.NewReader([]byte(`{"ref":"refs/heads/main"}`)))
	webhookReq.Header.Set("Content-Type", "application/json")
	webhookRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(webhookRec, webhookReq)

	if webhookRec.Code != http.StatusAccepted {
		t.Fatalf("expected status 202 accepted on git webhook trigger, got %d: %s", webhookRec.Code, webhookRec.Body.String())
	}

	// 5. Disconnect GitHub provider
	disconnectReq := httptest.NewRequest(http.MethodDelete, "/api/git/connect/github", nil)
	disconnectReq.Header.Set("Authorization", "Bearer "+tokenStr)
	disconnectRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(disconnectRec, disconnectReq)

	if disconnectRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on git disconnect, got %d: %s", disconnectRec.Code, disconnectRec.Body.String())
	}
}
