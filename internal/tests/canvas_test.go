package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"vessel.dev/vessel/internal/api"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
)

func TestRailwayCanvasAndEnvironments(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "vessel-test-canvas-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	st, err := store.NewStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer st.Close()

	proxyCfg := proxy.NewCaddyConfig(tempDir, "admin@test.local")
	proxyMgr := proxy.NewProxyManager(proxyCfg, st, nil)
	srv := api.NewServer(st, nil, proxyMgr, nil)

	// Register user and obtain auth token
	regPayload := []byte(`{"email":"canvasuser@vessel.dev","password":"securepassword123","role":"developer"}`)
	regReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regPayload))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(regRec, regReq)

	if regRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 created on register, got %d: %s", regRec.Code, regRec.Body.String())
	}

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

	// Step 1: Create a Project Workspace (e.g. Xdrive) with an initial repository
	projectPayload := map[string]interface{}{
		"name":          "Xdrive",
		"repositoryUrl": "https://github.com/solomonolatunji/xdrive-web.git",
		"branch":        "main",
		"internalPort":  3000,
	}
	body, _ := json.Marshal(projectPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 created on project creation, got %d: %s", w.Code, w.Body.String())
	}

	var project types.ProjectConfig
	json.Unmarshal(w.Body.Bytes(), &project)
	if project.ID == "" {
		t.Fatalf("expected non-empty project ID")
	}

	// Step 2: Verify default "production" environment and initial app service were auto-created
	envsReq := httptest.NewRequest(http.MethodGet, "/api/projects/"+project.ID+"/environments", nil)
	envsReq.Header.Set("Authorization", "Bearer "+tokenStr)
	w = httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, envsReq)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on list environments, got %d: %s", w.Code, w.Body.String())
	}

	var envs []*types.EnvironmentConfig
	json.Unmarshal(w.Body.Bytes(), &envs)
	if len(envs) != 1 || envs[0].Name != "production" || !envs[0].IsDefault {
		t.Fatalf("expected 1 default production environment, got %+v", envs)
	}
	prodEnvID := envs[0].ID

	// Step 3: Create an additional "staging" environment inside the project workspace
	stagingPayload := map[string]interface{}{
		"name": "staging",
	}
	stagingBody, _ := json.Marshal(stagingPayload)
	createEnvReq := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID+"/environments", bytes.NewBuffer(stagingBody))
	createEnvReq.Header.Set("Content-Type", "application/json")
	createEnvReq.Header.Set("Authorization", "Bearer "+tokenStr)
	w = httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, createEnvReq)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 when creating staging environment, got %d: %s", w.Code, w.Body.String())
	}

	// Step 4: Register two more Git apps inside production (e.g. wallet-bot and recovery)
	appPayload := map[string]interface{}{
		"name":          "wallet-bot",
		"repositoryUrl": "https://github.com/solomonolatunji/wallet-bot.git",
		"branch":        "main",
		"internalPort":  8000,
	}
	appBody, _ := json.Marshal(appPayload)
	createAppReq := httptest.NewRequest(http.MethodPost, "/api/environments/"+prodEnvID+"/apps", bytes.NewBuffer(appBody))
	createAppReq.Header.Set("Content-Type", "application/json")
	createAppReq.Header.Set("Authorization", "Bearer "+tokenStr)
	w = httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, createAppReq)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 on create app service, got %d: %s", w.Code, w.Body.String())
	}

	// Step 5: Query Project Canvas Summary (`GET /api/projects/{id}/summary`)
	summaryReq := httptest.NewRequest(http.MethodGet, "/api/projects/"+project.ID+"/summary", nil)
	summaryReq.Header.Set("Authorization", "Bearer "+tokenStr)
	w = httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, summaryReq)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on project canvas summary, got %d: %s", w.Code, w.Body.String())
	}

	var summary types.ProjectCanvasSummary
	json.Unmarshal(w.Body.Bytes(), &summary)
	if summary.AppsCount != 2 { // initial Xdrive app + wallet-bot
		t.Fatalf("expected 2 apps in project canvas summary, got %d", summary.AppsCount)
	}
	if summary.EnvironmentsCount != 2 { // production + staging
		t.Fatalf("expected 2 environments in project canvas summary, got %d", summary.EnvironmentsCount)
	}

	// Step 6: Query Environment Canvas (`GET /api/environments/{id}/canvas`)
	canvasReq := httptest.NewRequest(http.MethodGet, "/api/environments/"+prodEnvID+"/canvas", nil)
	canvasReq.Header.Set("Authorization", "Bearer "+tokenStr)
	w = httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, canvasReq)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on environment canvas, got %d: %s", w.Code, w.Body.String())
	}

	var canvas types.EnvironmentCanvas
	json.Unmarshal(w.Body.Bytes(), &canvas)
	if len(canvas.Apps) != 2 {
		t.Fatalf("expected 2 apps in production environment canvas, got %d", len(canvas.Apps))
	}
	if canvas.Environment.Name != "production" {
		t.Fatalf("expected environment name 'production', got '%s'", canvas.Environment.Name)
	}
}
