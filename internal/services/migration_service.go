package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

const bundleManifestVersion = "1"

type BundleManifest struct {
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	Databases []string  `json:"databases"`
	HasSQLite bool      `json:"hasSqlite"`
}

type MigrationService struct {
	dbRepo  repositories.DatabaseRepository
	dataDir string
}

func NewMigrationService(dbRepo repositories.DatabaseRepository, dataDir string) *MigrationService {
	return &MigrationService{dbRepo: dbRepo, dataDir: dataDir}
}

func (s *MigrationService) Export(ctx context.Context, passphrase string) ([]byte, error) {
	files := make(map[string][]byte)

	sqliteData, err := s.dumpSQLite()
	if err != nil {
		return nil, fmt.Errorf("sqlite dump failed: %w", err)
	}
	files["vessl.db.sql"] = sqliteData

	dbs, err := s.dbRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}

	manifest := BundleManifest{
		Version:   bundleManifestVersion,
		CreatedAt: time.Now().UTC(),
		HasSQLite: true,
	}

	for _, db := range dbs {
		dump, err := s.dumpDatabase(ctx, db)
		if err != nil {
			continue
		}
		ext := dumpExtension(db.Engine)
		filename := fmt.Sprintf("databases/%s%s", db.Name, ext)
		files[filename] = dump
		manifest.Databases = append(manifest.Databases, db.Name)
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest: %w", err)
	}
	files["manifest.json"] = manifestData

	var tarBuf bytes.Buffer
	if err := engine.CreateTarGz(&tarBuf, files); err != nil {
		return nil, fmt.Errorf("failed to create bundle archive: %w", err)
	}

	var encBuf bytes.Buffer
	if err := engine.EncryptBundle(&tarBuf, &encBuf, passphrase); err != nil {
		return nil, fmt.Errorf("failed to encrypt bundle: %w", err)
	}

	return encBuf.Bytes(), nil
}

func (s *MigrationService) Import(ctx context.Context, bundleData []byte, passphrase string) (*BundleManifest, error) {
	var decBuf bytes.Buffer
	if err := engine.DecryptBundle(bytes.NewReader(bundleData), &decBuf, passphrase); err != nil {
		return nil, err
	}

	files, err := engine.ExtractTarGz(&decBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to extract bundle: %w", err)
	}

	manifestData, ok := files["manifest.json"]
	if !ok {
		return nil, fmt.Errorf("bundle is missing manifest.json")
	}
	var manifest BundleManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	if sqlData, ok := files["vessl.db.sql"]; ok {
		if err := s.restoreSQLite(sqlData); err != nil {
			return nil, fmt.Errorf("sqlite restore failed: %w", err)
		}
	}

	dbs, _ := s.dbRepo.List(ctx)
	dbsByName := make(map[string]*models.Database, len(dbs))
	for _, db := range dbs {
		dbsByName[db.Name] = db
	}

	for _, dbName := range manifest.Databases {
		db, found := dbsByName[dbName]
		if !found {
			continue
		}
		for _, ext := range []string{".sql", ".rdb", ".dump"} {
			key := fmt.Sprintf("databases/%s%s", dbName, ext)
			if data, ok := files[key]; ok {
				_ = s.restoreDatabase(ctx, db, data, ext)
				break
			}
		}
	}

	return &manifest, nil
}

func (s *MigrationService) dumpSQLite() ([]byte, error) {
	dbPath := filepath.Join(s.dataDir, "vessl.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("sqlite db not found at %s", dbPath)
	}
	out, err := exec.Command("sqlite3", dbPath, ".dump").Output()
	if err != nil {
		data, readErr := os.ReadFile(dbPath)
		if readErr != nil {
			return nil, fmt.Errorf("sqlite3 dump failed and file read failed: %w", err)
		}
		return data, nil
	}
	return out, nil
}

func (s *MigrationService) restoreSQLite(sqlData []byte) error {
	dbPath := filepath.Join(s.dataDir, "vessl.db")
	backupPath := dbPath + ".bak"
	_ = os.Rename(dbPath, backupPath)

	tmpSQL := dbPath + ".import.sql"
	if err := os.WriteFile(tmpSQL, sqlData, 0600); err != nil {
		_ = os.Rename(backupPath, dbPath)
		return err
	}
	defer os.Remove(tmpSQL)

	cmd := exec.Command("sqlite3", dbPath)
	cmd.Stdin = bytes.NewReader(sqlData)
	if err := cmd.Run(); err != nil {
		_ = os.Rename(backupPath, dbPath)
		return fmt.Errorf("sqlite3 restore failed: %w", err)
	}
	_ = os.Remove(backupPath)
	return nil
}

func (s *MigrationService) dumpDatabase(_ context.Context, db *models.Database) ([]byte, error) {
	containerName := db.InternalDNS
	if containerName == "" {
		return nil, fmt.Errorf("no internal DNS for database %s", db.Name)
	}

	switch db.Engine {
	case "postgres", "timescaledb":
		cmd := exec.Command("docker", "exec", containerName,
			"pg_dump", "-U", db.Username, db.DatabaseName)
		return cmd.Output()
	case "mysql", "mariadb":
		cmd := exec.Command("docker", "exec", containerName,
			"mysqldump", "-u", db.Username, fmt.Sprintf("-p%s", db.Password), db.DatabaseName)
		return cmd.Output()
	case "mongodb":
		cmd := exec.Command("docker", "exec", containerName,
			"mongodump", "--archive", "--authenticationDatabase=admin",
			fmt.Sprintf("--username=%s", db.Username), fmt.Sprintf("--password=%s", db.Password))
		return cmd.Output()
	case "clickhouse":
		cmd := exec.Command("docker", "exec", containerName,
			"clickhouse-client", "--user", db.Username, "--password", db.Password,
			"--query", fmt.Sprintf("SELECT * FROM %s FORMAT Native", db.DatabaseName))
		return cmd.Output()
	case "redis":
		cmd := exec.Command("docker", "exec", containerName,
			"redis-cli", "--rdb", "/tmp/dump.rdb")
		if err := cmd.Run(); err != nil {
			return nil, err
		}
		return exec.Command("docker", "exec", containerName, "cat", "/tmp/dump.rdb").Output()
	default:
		return nil, fmt.Errorf("unsupported engine for dump: %s", db.Engine)
	}
}

func (s *MigrationService) restoreDatabase(_ context.Context, db *models.Database, data []byte, ext string) error {
	containerName := db.InternalDNS
	if containerName == "" {
		return fmt.Errorf("no internal DNS for database %s", db.Name)
	}

	switch db.Engine {
	case "postgres", "timescaledb":
		cmd := exec.Command("docker", "exec", "-i", containerName,
			"psql", "-U", db.Username, db.DatabaseName)
		cmd.Stdin = bytes.NewReader(data)
		return cmd.Run()
	case "mysql", "mariadb":
		cmd := exec.Command("docker", "exec", "-i", containerName,
			"mysql", "-u", db.Username, fmt.Sprintf("-p%s", db.Password), db.DatabaseName)
		cmd.Stdin = bytes.NewReader(data)
		return cmd.Run()
	case "mongodb":
		cmd := exec.Command("docker", "exec", "-i", containerName,
			"mongorestore", "--archive", "--authenticationDatabase=admin",
			fmt.Sprintf("--username=%s", db.Username), fmt.Sprintf("--password=%s", db.Password))
		cmd.Stdin = bytes.NewReader(data)
		return cmd.Run()
	case "redis":
		if ext != ".rdb" {
			return fmt.Errorf("redis restore only supports .rdb format")
		}
		copyCmd := exec.Command("docker", "cp", "-", fmt.Sprintf("%s:/tmp/vessl-import.rdb", containerName))
		copyCmd.Stdin = bytes.NewReader(data)
		if err := copyCmd.Run(); err != nil {
			return err
		}
		return exec.Command("docker", "exec", containerName,
			"redis-cli", "DEBUG", "RELOAD").Run()
	default:
		return fmt.Errorf("unsupported engine for restore: %s", db.Engine)
	}
}

func dumpExtension(engine models.DatabaseEngine) string {
	switch engine {
	case models.DatabaseEngineRedis:
		return ".rdb"
	default:
		return ".sql"
	}
}
