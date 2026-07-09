package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"vessel.dev/vessel/internal/api"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
)

func TestScheduledJobsEndpoints(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "vessel_job_test")
	_ = os.RemoveAll(tempDir)
	_ = os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "vessel.db")
	s, err := store.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to init store: %v", err)
	}
	defer s.Close()

	srv := api.NewServer(s, nil, nil, nil)

	registerPayload := []byte(`{"email":"admin@vessel.dev","password":"securepassword123","role":"admin"}`)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerPayload))
	reqReg.Header.Set("Content-Type", "application/json")
	recReg := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recReg, reqReg)

	var registerResp map[string]any
	_ = json.NewDecoder(recReg.Body).Decode(&registerResp)
	tokenStr, _ := registerResp["token"].(string)

	project := &types.ProjectConfig{
		Name: "Job Test Project",
	}
	if err := s.CreateProject(project); err != nil {
		t.Fatalf("failed to create dummy project: %v", err)
	}

	jobReq := types.JobConfig{
		ProjectID: project.ID,
		Name:      "Daily Cleanup",
		Schedule:  "0 0 * * *",
		Command:   "php artisan schedule:run",
	}
	payload, _ := json.Marshal(jobReq)
	req := httptest.NewRequest(http.MethodPost, "/api/jobs", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d. Body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var createdJob types.JobConfig
	if err := json.NewDecoder(rec.Body).Decode(&createdJob); err != nil {
		t.Fatalf("failed to decode created job: %v", err)
	}
	if createdJob.ID == "" {
		t.Fatal("expected job ID to be generated")
	}

	reqList := httptest.NewRequest(http.MethodGet, "/api/jobs?projectId="+project.ID, nil)
	reqList.Header.Set("Authorization", "Bearer "+tokenStr)
	recList := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recList, reqList)

	if recList.Code != http.StatusOK {
		t.Fatalf("expected status %d on list jobs, got %d", http.StatusOK, recList.Code)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/api/jobs/"+createdJob.ID, nil)
	reqGet.Header.Set("Authorization", "Bearer "+tokenStr)
	recGet := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Fatalf("expected status %d on get job, got %d", http.StatusOK, recGet.Code)
	}

	reqDel := httptest.NewRequest(http.MethodDelete, "/api/jobs/"+createdJob.ID, nil)
	reqDel.Header.Set("Authorization", "Bearer "+tokenStr)
	recDel := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recDel, reqDel)

	if recDel.Code != http.StatusOK {
		t.Fatalf("expected status %d on delete job, got %d", http.StatusOK, recDel.Code)
	}
}
