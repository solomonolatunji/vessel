package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"vessel.dev/vessel/internal/api"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
)

func TestAutomatedBackupsEndpoints(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "vessel_backup_test")
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

	// 1. Register admin user & get token
	registerPayload := []byte(`{"email":"backupadmin@vessel.dev","password":"securepassword123","role":"admin"}`)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerPayload))
	reqReg.Header.Set("Content-Type", "application/json")
	recReg := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recReg, reqReg)

	var registerResp map[string]any
	_ = json.NewDecoder(recReg.Body).Decode(&registerResp)
	tokenStr, _ := registerResp["token"].(string)

	// 2. Create a Project
	project := &types.ProjectConfig{Name: "Backup Test Project"}
	_ = s.CreateProject(project)

	// 3. Create a Database instance in the store to backup
	dbInstance := &types.DatabaseConfig{
		ID:        "test-pg-backup",
		ProjectID: project.ID,
		Name:      "PostgreSQL Prod",
		Engine:    "postgresql",
		Version:   "16",
		Port:      5432,
		Status:    "running",
	}
	if err := s.CreateDatabase(dbInstance); err != nil {
		t.Fatalf("failed to create dummy database: %v", err)
	}

	// 4. Create an S3 destination
	s3DestReq := types.S3Destination{
		ProjectID:       project.ID,
		Name:            "Offsite MinIO",
		Endpoint:        "s3.eu-west-1.amazonaws.com",
		Bucket:          "vessel-backups-prod",
		Region:          "eu-west-1",
		AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
		SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}
	payload, _ := json.Marshal(s3DestReq)
	reqS3 := httptest.NewRequest(http.MethodPost, "/api/s3-destinations", bytes.NewReader(payload))
	reqS3.Header.Set("Content-Type", "application/json")
	reqS3.Header.Set("Authorization", "Bearer "+tokenStr)
	recS3 := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recS3, reqS3)

	if recS3.Code != http.StatusCreated {
		t.Fatalf("expected status %d on create s3 dest, got %d. Body: %s", http.StatusCreated, recS3.Code, recS3.Body.String())
	}
	var createdS3 types.S3Destination
	_ = json.NewDecoder(recS3.Body).Decode(&createdS3)

	// Verify listing S3 destinations
	reqListS3 := httptest.NewRequest(http.MethodGet, "/api/s3-destinations?projectId="+project.ID, nil)
	reqListS3.Header.Set("Authorization", "Bearer "+tokenStr)
	recListS3 := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recListS3, reqListS3)
	if recListS3.Code != http.StatusOK {
		t.Fatalf("expected status %d on list s3 dests, got %d", http.StatusOK, recListS3.Code)
	}

	// 5. Create a Backup Schedule targeting the PostgreSQL instance and S3 destination
	backupReq := types.BackupConfig{
		ProjectID:       project.ID,
		DatabaseID:      dbInstance.ID,
		S3DestinationID: createdS3.ID,
		Name:            "Daily PG Dump",
		Schedule:        "0 2 * * *",
		RetentionDays:   14,
	}
	payload, _ = json.Marshal(backupReq)
	reqBackup := httptest.NewRequest(http.MethodPost, "/api/backups", bytes.NewReader(payload))
	reqBackup.Header.Set("Content-Type", "application/json")
	reqBackup.Header.Set("Authorization", "Bearer "+tokenStr)
	recBackup := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recBackup, reqBackup)

	if recBackup.Code != http.StatusCreated {
		t.Fatalf("expected status %d on create backup, got %d. Body: %s", http.StatusCreated, recBackup.Code, recBackup.Body.String())
	}
	var createdBackup types.BackupConfig
	_ = json.NewDecoder(recBackup.Body).Decode(&createdBackup)
	if createdBackup.ID == "" {
		t.Fatal("expected backup schedule ID to be generated")
	}

	// 6. Trigger Backup execution immediately
	reqTrig := httptest.NewRequest(http.MethodPost, "/api/backups/"+createdBackup.ID+"/trigger", nil)
	reqTrig.Header.Set("Authorization", "Bearer "+tokenStr)
	recTrig := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recTrig, reqTrig)

	if recTrig.Code != http.StatusAccepted {
		t.Fatalf("expected status %d on trigger backup, got %d. Body: %s", http.StatusAccepted, recTrig.Code, recTrig.Body.String())
	}

	var recResult types.BackupRecord
	_ = json.NewDecoder(recTrig.Body).Decode(&recResult)
	if recResult.Status != "completed" {
		t.Fatalf("expected backup status 'completed' in simulated local run, got '%s'. Logs: %s", recResult.Status, recResult.Logs)
	}
	if recResult.S3URL == "" {
		t.Fatalf("expected S3 URL to be generated after trigger, got empty string")
	}

	// 7. List backup execution records
	reqRecs := httptest.NewRequest(http.MethodGet, "/api/backups/"+createdBackup.ID+"/records", nil)
	reqRecs.Header.Set("Authorization", "Bearer "+tokenStr)
	recRecs := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recRecs, reqRecs)

	if recRecs.Code != http.StatusOK {
		t.Fatalf("expected status %d on list records, got %d", http.StatusOK, recRecs.Code)
	}

	var records []*types.BackupRecord
	_ = json.NewDecoder(recRecs.Body).Decode(&records)
	if len(records) == 0 {
		t.Fatal("expected at least 1 backup record")
	}

	// 8. Test retention policy pruning simulation
	// Insert an expired record from 30 days ago
	expiredRec := &types.BackupRecord{
		BackupConfigID: createdBackup.ID,
		ProjectID:      project.ID,
		DatabaseID:     dbInstance.ID,
		Status:         "completed",
		FilePath:       filepath.Join(tempDir, "old_expired_backup.sql"),
		StartedAt:      time.Now().Add(-35 * 24 * time.Hour).UTC().Format(time.RFC3339),
		CompletedAt:    time.Now().Add(-35 * 24 * time.Hour).UTC().Format(time.RFC3339),
	}
	_ = os.WriteFile(expiredRec.FilePath, []byte("old dump data"), 0600)
	_ = s.CreateBackupRecord(expiredRec)

	// Trigger again to verify retention cleanup kicks in
	reqTrig2 := httptest.NewRequest(http.MethodPost, "/api/backups/"+createdBackup.ID+"/trigger", nil)
	reqTrig2.Header.Set("Authorization", "Bearer "+tokenStr)
	recTrig2 := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recTrig2, reqTrig2)

	// Check if the old file was removed
	if _, err := os.Stat(expiredRec.FilePath); !os.IsNotExist(err) {
		t.Errorf("expected old backup file %s to be pruned by retention policy", expiredRec.FilePath)
	}

	// 9. Delete Backup Schedule and S3 destination
	reqDelBackup := httptest.NewRequest(http.MethodDelete, "/api/backups/"+createdBackup.ID, nil)
	reqDelBackup.Header.Set("Authorization", "Bearer "+tokenStr)
	recDelBackup := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recDelBackup, reqDelBackup)
	if recDelBackup.Code != http.StatusNoContent {
		t.Fatalf("expected status %d on delete backup, got %d", http.StatusNoContent, recDelBackup.Code)
	}

	reqDelS3 := httptest.NewRequest(http.MethodDelete, "/api/s3-destinations/"+createdS3.ID+"?projectId="+project.ID, nil)
	reqDelS3.Header.Set("Authorization", "Bearer "+tokenStr)
	recDelS3 := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recDelS3, reqDelS3)
	if recDelS3.Code != http.StatusNoContent {
		t.Fatalf("expected status %d on delete s3 dest, got %d", http.StatusNoContent, recDelS3.Code)
	}
}
