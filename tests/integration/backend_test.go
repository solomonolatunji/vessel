package integration_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/docker/docker/client"
	_ "modernc.org/sqlite"

	vesslhttp "vessl.dev/vessl/internal/http"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/utils"
)

func TestVesslBackendInitialization(t *testing.T) {
	dataDir := t.TempDir()
	vlt, err := utils.NewVault(dataDir)
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}

	dbPath := filepath.Join(dataDir, "vessl.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}
	defer db.Close()

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	t.Setenv("VESSL_JWT_SECRET", "testsecret")

	dockerClient, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	server, err := vesslhttp.NewServer(db, vlt, nil, nil, dockerClient, dataDir)
	if err != nil {
		t.Fatalf("failed to initialize server: %v", err)
	}

	if server == nil {
		t.Fatalf("expected server to be non-nil")
	}

	if err := db.PingContext(context.Background()); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}
}
