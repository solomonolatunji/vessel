package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"vessel.dev/vessel/internal/api"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
)

func TestServiceSettingsAndProjectSettings(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "vessel-test-settings-*")
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
	regPayload := []byte(`{"email":"settingsuser@vessel.dev","password":"securepassword123","role":"developer"}`)
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

	// Step 1: Create Project
	projectPayload := map[string]interface{}{
		"name":          "Xdrive Settings Project",
		"repositoryUrl": "https://github.com/solomonolatunji/xdrive-web.git",
		"branch":        "main",
	}
	projBytes, _ := json.Marshal(projectPayload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewReader(projBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+tokenStr)
	createRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 created project, got %d: %s", createRec.Code, createRec.Body.String())
	}

	var project types.ProjectConfig
	json.Unmarshal(createRec.Body.Bytes(), &project)
	if project.ID == "" {
		t.Fatalf("created project had empty ID")
	}

	// Step 2: Get default environment
	listEnvReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/projects/%s/environments", project.ID), nil)
	listEnvReq.Header.Set("Authorization", "Bearer "+tokenStr)
	listEnvRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(listEnvRec, listEnvReq)

	if listEnvRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing environments, got %d", listEnvRec.Code)
	}

	var envs []*types.EnvironmentConfig
	json.Unmarshal(listEnvRec.Body.Bytes(), &envs)
	if len(envs) == 0 {
		t.Fatalf("expected at least 1 default environment")
	}
	env := envs[0]

	// Step 3: Create App Service (Backend API)
	appPayload := map[string]interface{}{
		"name":          "Backend API",
		"repositoryUrl": "https://github.com/solomonolatunji/xdrive-backend.git",
		"branch":        "main",
		"rootDirectory": "/src",
		"buildCommand":  "npm run build",
		"startCommand":  "npm run start:prod",
		"containerPort": 3000,
		"replicas":      2,
		"cpuRequest":    0.5,
		"memoryLimitMb": 512,
		"restartPolicy": "always",
	}
	appBytes, _ := json.Marshal(appPayload)
	createAppReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/environments/%s/apps", env.ID), bytes.NewReader(appBytes))
	createAppReq.Header.Set("Content-Type", "application/json")
	createAppReq.Header.Set("Authorization", "Bearer "+tokenStr)
	createAppRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(createAppRec, createAppReq)

	if createAppRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 created app service, got %d: %s", createAppRec.Code, createAppRec.Body.String())
	}
	var appService types.AppServiceConfig
	json.Unmarshal(createAppRec.Body.Bytes(), &appService)

	// Step 4: Test Trigger Service Deployment
	deployReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/services/%s/deploy", appService.ID), nil)
	deployReq.Header.Set("Authorization", "Bearer "+tokenStr)
	deployRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(deployRec, deployReq)

	if deployRec.Code != http.StatusAccepted {
		t.Fatalf("expected 202 accepted triggering deploy, got %d: %s", deployRec.Code, deployRec.Body.String())
	}
	var dep types.DeploymentRecord
	json.Unmarshal(deployRec.Body.Bytes(), &dep)
	if dep.Status != "BUILDING" {
		t.Fatalf("expected deployment status BUILDING, got %s", dep.Status)
	}

	// Step 5: List Service Deployments & Check Logs
	listDepsReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/services/%s/deployments", appService.ID), nil)
	listDepsReq.Header.Set("Authorization", "Bearer "+tokenStr)
	listDepsRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(listDepsRec, listDepsReq)

	if listDepsRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing deployments, got %d", listDepsRec.Code)
	}
	var deps []*types.DeploymentRecord
	json.Unmarshal(listDepsRec.Body.Bytes(), &deps)
	if len(deps) == 0 {
		t.Fatalf("expected at least 1 deployment record")
	}

	logsReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/deployments/%s/logs", dep.ID), nil)
	logsReq.Header.Set("Authorization", "Bearer "+tokenStr)
	logsRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(logsRec, logsReq)
	if logsRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK getting deployment logs, got %d", logsRec.Code)
	}

	// Step 6: Test Rollback Deployment
	rollbackReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/deployments/%s/rollback", dep.ID), nil)
	rollbackReq.Header.Set("Authorization", "Bearer "+tokenStr)
	rollbackRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rollbackRec, rollbackReq)
	if rollbackRec.Code != http.StatusAccepted {
		t.Fatalf("expected 202 accepted on rollback, got %d", rollbackRec.Code)
	}

	// Step 7: Test Service Metrics
	metricsReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/services/%s/metrics", appService.ID), nil)
	metricsReq.Header.Set("Authorization", "Bearer "+tokenStr)
	metricsRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(metricsRec, metricsReq)
	if metricsRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK getting service metrics, got %d", metricsRec.Code)
	}

	// Step 8: Test Service Variables & Bulk RAW Editor
	varPayload := map[string]interface{}{
		"key":           "PORT",
		"value":         "3000",
		"isSecret":      false,
		"environmentId": env.ID,
	}
	varBytes, _ := json.Marshal(varPayload)
	setVarReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/services/%s/variables", appService.ID), bytes.NewReader(varBytes))
	setVarReq.Header.Set("Content-Type", "application/json")
	setVarReq.Header.Set("Authorization", "Bearer "+tokenStr)
	setVarRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(setVarRec, setVarReq)
	if setVarRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 created service variable, got %d: %s", setVarRec.Code, setVarRec.Body.String())
	}

	bulkPayload := map[string]interface{}{
		"environmentId": env.ID,
		"variables": []map[string]interface{}{
			{"key": "DATABASE_URL", "value": "postgres://user:pass@localhost:5432/db", "isSecret": true},
			{"key": "REDIS_HOST", "value": "redis.internal", "isSecret": false},
		},
	}
	bulkBytes, _ := json.Marshal(bulkPayload)
	bulkReq := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/services/%s/variables/bulk", appService.ID), bytes.NewReader(bulkBytes))
	bulkReq.Header.Set("Content-Type", "application/json")
	bulkReq.Header.Set("Authorization", "Bearer "+tokenStr)
	bulkRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(bulkRec, bulkReq)
	if bulkRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on bulk set variables, got %d: %s", bulkRec.Code, bulkRec.Body.String())
	}

	// Step 9: Test Project Settings (Billing, Webhooks, Tokens, Members)
	billingReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/projects/%s/billing", project.ID), nil)
	billingReq.Header.Set("Authorization", "Bearer "+tokenStr)
	billingRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(billingRec, billingReq)
	if billingRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on billing, got %d", billingRec.Code)
	}

	webhookPayload := map[string]interface{}{
		"url":    "https://discord.com/api/webhooks/test",
		"events": []string{"deploy.success", "deploy.failure"},
	}
	webhookBytes, _ := json.Marshal(webhookPayload)
	whReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/projects/%s/webhooks", project.ID), bytes.NewReader(webhookBytes))
	whReq.Header.Set("Content-Type", "application/json")
	whReq.Header.Set("Authorization", "Bearer "+tokenStr)
	whRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(whRec, whReq)
	if whRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 created webhook, got %d: %s", whRec.Code, whRec.Body.String())
	}

	tokPayload := map[string]interface{}{
		"name":        "CI Deployment Token",
		"permissions": []string{"deploy", "read"},
	}
	tokBytes, _ := json.Marshal(tokPayload)
	tokReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/projects/%s/tokens", project.ID), bytes.NewReader(tokBytes))
	tokReq.Header.Set("Content-Type", "application/json")
	tokReq.Header.Set("Authorization", "Bearer "+tokenStr)
	tokRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(tokRec, tokReq)
	if tokRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 created token, got %d: %s", tokRec.Code, tokRec.Body.String())
	}

	memberPayload := map[string]interface{}{
		"email": "teammate@vessel.dev",
		"role":  "admin",
	}
	memberBytes, _ := json.Marshal(memberPayload)
	memReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/projects/%s/members", project.ID), bytes.NewReader(memberBytes))
	memReq.Header.Set("Content-Type", "application/json")
	memReq.Header.Set("Authorization", "Bearer "+tokenStr)
	memRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(memRec, memReq)
	if memRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 created project member invite, got %d: %s", memRec.Code, memRec.Body.String())
	}
}
