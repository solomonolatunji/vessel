package repositories_test

import (
	"context"
	"database/sql"
	"testing"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/repositories"

	_ "modernc.org/sqlite"
)

type testVault struct{}

func (testVault) Encrypt(plaintext string) (string, error) {
	return "encrypted:" + plaintext, nil
}

func (testVault) Decrypt(ciphertext string) (string, error) {
	return ciphertext, nil
}

func TestDatabaseRepositoryListAfterMigrations(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	repo := repositories.NewDatabaseRepo(db, testVault{})
	database := &models.Database{
		ProjectID:          "project-1",
		EnvironmentID:      "environment-1",
		Name:               "postgres",
		Engine:             models.DatabaseEnginePostgres,
		Version:            "16",
		Port:               5432,
		Username:           "postgres",
		Password:           "password",
		DatabaseName:       "app",
		VolumePath:         "/tmp/postgres",
		Status:             models.DatabaseStatusCreated,
		InternalDNS:        "postgres.internal",
		ExternalDNS:        "",
		CustomArgs:         "",
		LogicalReplication: true,
	}

	if err := repo.Create(context.Background(), database); err != nil {
		t.Fatalf("create database: %v", err)
	}

	databases, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("list databases: %v", err)
	}
	if len(databases) != 1 {
		t.Fatalf("expected 1 database, got %d", len(databases))
	}
	if databases[0].ProjectID != "project-1" {
		t.Fatalf("expected project id, got %q", databases[0].ProjectID)
	}
	if !databases[0].LogicalReplication {
		t.Fatal("expected logical replication to round trip")
	}
}
