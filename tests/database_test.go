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
	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

func TestManagedDatabasesAndStorageEndpoints(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "vessel_db_storage_test")
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

	dbPayload := []byte(`{"name":"prod-pg","engine":"postgres","version":"16","username":"myuser","password":"mypassword","databaseName":"mydb"}`)
	reqCreateDB := httptest.NewRequest(http.MethodPost, "/api/databases", bytes.NewReader(dbPayload))
	reqCreateDB.Header.Set("Content-Type", "application/json")
	reqCreateDB.Header.Set("Authorization", "Bearer "+tokenStr)
	recCreateDB := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recCreateDB, reqCreateDB)

	if recCreateDB.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for POST /api/databases, got %d. Body: %s", recCreateDB.Code, recCreateDB.Body.String())
	}

	var createdDB types.DatabaseConfig
	if err := json.NewDecoder(recCreateDB.Body).Decode(&createdDB); err != nil {
		t.Fatalf("failed to decode created database response: %v", err)
	}

	if createdDB.Port != 5432 || createdDB.Engine != "postgres" {
		t.Errorf("expected port 5432 and engine postgres, got port %d and engine %s", createdDB.Port, createdDB.Engine)
	}

	reqGetDB := httptest.NewRequest(http.MethodGet, "/api/databases/"+createdDB.ID, nil)
	reqGetDB.Header.Set("Authorization", "Bearer "+tokenStr)
	recGetDB := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recGetDB, reqGetDB)

	if recGetDB.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GET /api/databases/%s, got %d", createdDB.ID, recGetDB.Code)
	}

	var fetchedDB types.DatabaseConfig
	_ = json.NewDecoder(recGetDB.Body).Decode(&fetchedDB)
	if fetchedDB.Password != "mypassword" {
		t.Errorf("expected decrypted password to be mypassword, got %s", fetchedDB.Password)
	}

	storagePayload := []byte(`{"name":"app-minio","accessKey":"minioadmin","secretKey":"miniosecret123","bucketName":"media"}`)
	reqCreateStorage := httptest.NewRequest(http.MethodPost, "/api/storage", bytes.NewReader(storagePayload))
	reqCreateStorage.Header.Set("Content-Type", "application/json")
	reqCreateStorage.Header.Set("Authorization", "Bearer "+tokenStr)
	recCreateStorage := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recCreateStorage, reqCreateStorage)

	if recCreateStorage.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for POST /api/storage, got %d. Body: %s", recCreateStorage.Code, recCreateStorage.Body.String())
	}

	var createdStorage types.StorageConfig
	_ = json.NewDecoder(recCreateStorage.Body).Decode(&createdStorage)
	if createdStorage.APIPort != 9000 || createdStorage.ConsolePort != 9001 {
		t.Errorf("expected API port 9000 and console 9001, got %d and %d", createdStorage.APIPort, createdStorage.ConsolePort)
	}
}
