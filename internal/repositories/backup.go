package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/models"
)

type BackupRepository interface {
	CreateConfig(ctx context.Context, cfg *models.BackupConfig) error
	GetConfigByID(ctx context.Context, id string) (*models.BackupConfig, error)
	ListConfigsByProject(ctx context.Context, projectID string) ([]*models.BackupConfig, error)
	ListAllActiveConfigs(ctx context.Context) ([]*models.BackupConfig, error)
	DeleteConfig(ctx context.Context, id, projectID string) error

	CreateRecord(ctx context.Context, rec *models.BackupRecord) error
	GetRecordByID(ctx context.Context, id string) (*models.BackupRecord, error)
	ListRecordsByConfig(ctx context.Context, backupConfigID string) ([]*models.BackupRecord, error)
	UpdateRecord(ctx context.Context, rec *models.BackupRecord) error
}

type BackupSQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewBackupSQLiteRepository(db *sql.DB) *BackupSQLiteRepository {
	return &BackupSQLiteRepository{db: db}
}

func (r *BackupSQLiteRepository) EnsureTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS backup_configs (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			database_id TEXT,
			storage_id TEXT,
			s3_destination_id TEXT,
			name TEXT NOT NULL,
			schedule TEXT NOT NULL,
			retention_days INTEGER DEFAULT 7,
			status TEXT DEFAULT 'active',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS backup_records (
			id TEXT PRIMARY KEY,
			backup_config_id TEXT NOT NULL,
			project_id TEXT NOT NULL,
			database_id TEXT,
			status TEXT DEFAULT 'running',
			file_path TEXT,
			file_size_bytes INTEGER DEFAULT 0,
			s3_url TEXT,
			logs TEXT,
			started_at TEXT NOT NULL,
			completed_at TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS s3_destinations (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			endpoint TEXT NOT NULL,
			bucket TEXT NOT NULL,
			region TEXT,
			access_key_id TEXT,
			secret_access_key TEXT,
			created_at TEXT NOT NULL
		)`,
	}
	for _, q := range queries {
		if _, err := r.db.Exec(q); err != nil {
			return fmt.Errorf("failed to create backup table: %w", err)
		}
	}
	return nil
}

func (r *BackupSQLiteRepository) CreateConfig(_ context.Context, cfg *models.BackupConfig) error {
	if cfg.ID == "" {
		cfg.ID = uuid.New().String()
	}
	if cfg.CreatedAt == "" {
		cfg.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	cfg.UpdatedAt = cfg.CreatedAt
	if cfg.Status == "" {
		cfg.Status = "active"
	}
	if cfg.RetentionDays <= 0 {
		cfg.RetentionDays = 7
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec(`INSERT INTO backup_configs (id, project_id, database_id, storage_id, s3_destination_id, name, schedule, retention_days, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		cfg.ID, cfg.ProjectID, cfg.DatabaseID, cfg.StorageID, cfg.S3DestinationID, cfg.Name, cfg.Schedule, cfg.RetentionDays, cfg.Status, cfg.CreatedAt, cfg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create backup config: %w", err)
	}
	return nil
}

func (r *BackupSQLiteRepository) GetConfigByID(_ context.Context, id string) (*models.BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	row := r.db.QueryRow(`SELECT id, project_id, database_id, storage_id, s3_destination_id, name, schedule, retention_days, status, created_at, updated_at
		FROM backup_configs WHERE id = ?`, id)

	var cfg models.BackupConfig
	err := row.Scan(&cfg.ID, &cfg.ProjectID, &cfg.DatabaseID, &cfg.StorageID, &cfg.S3DestinationID, &cfg.Name, &cfg.Schedule, &cfg.RetentionDays, &cfg.Status, &cfg.CreatedAt, &cfg.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get backup config %s: %w", id, err)
	}
	return &cfg, nil
}

func (r *BackupSQLiteRepository) ListConfigsByProject(_ context.Context, projectID string) ([]*models.BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query(`SELECT id, project_id, database_id, storage_id, s3_destination_id, name, schedule, retention_days, status, created_at, updated_at
		FROM backup_configs WHERE project_id = ? ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup configs: %w", err)
	}
	defer rows.Close()

	var list []*models.BackupConfig
	for rows.Next() {
		var cfg models.BackupConfig
		if err := rows.Scan(&cfg.ID, &cfg.ProjectID, &cfg.DatabaseID, &cfg.StorageID, &cfg.S3DestinationID, &cfg.Name, &cfg.Schedule, &cfg.RetentionDays, &cfg.Status, &cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &cfg)
	}
	return list, nil
}

func (r *BackupSQLiteRepository) ListAllActiveConfigs(_ context.Context) ([]*models.BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query(`SELECT id, project_id, database_id, storage_id, s3_destination_id, name, schedule, retention_days, status, created_at, updated_at
		FROM backup_configs WHERE status = 'active'`)
	if err != nil {
		return nil, fmt.Errorf("failed to list active backup configs: %w", err)
	}
	defer rows.Close()

	var list []*models.BackupConfig
	for rows.Next() {
		var cfg models.BackupConfig
		if err := rows.Scan(&cfg.ID, &cfg.ProjectID, &cfg.DatabaseID, &cfg.StorageID, &cfg.S3DestinationID, &cfg.Name, &cfg.Schedule, &cfg.RetentionDays, &cfg.Status, &cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &cfg)
	}
	return list, nil
}

func (r *BackupSQLiteRepository) DeleteConfig(_ context.Context, id, projectID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	res, err := r.db.Exec("DELETE FROM backup_configs WHERE id = ? AND project_id = ?", id, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete backup config: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("backup config not found or unauthorized")
	}
	return nil
}

func (r *BackupSQLiteRepository) CreateRecord(_ context.Context, rec *models.BackupRecord) error {
	if rec.ID == "" {
		rec.ID = uuid.New().String()
	}
	if rec.StartedAt == "" {
		rec.StartedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if rec.Status == "" {
		rec.Status = "running"
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec(`INSERT INTO backup_records (id, backup_config_id, project_id, database_id, status, file_path, file_size_bytes, s3_url, logs, started_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rec.ID, rec.BackupConfigID, rec.ProjectID, rec.DatabaseID, rec.Status, rec.FilePath, rec.FileSizeBytes, rec.S3URL, rec.Logs, rec.StartedAt, rec.CompletedAt)
	if err != nil {
		return fmt.Errorf("failed to create backup record: %w", err)
	}
	return nil
}

func (r *BackupSQLiteRepository) ListRecordsByConfig(_ context.Context, backupConfigID string) ([]*models.BackupRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query(`SELECT id, backup_config_id, project_id, database_id, status, file_path, file_size_bytes, s3_url, logs, started_at, completed_at
		FROM backup_records WHERE backup_config_id = ? ORDER BY started_at DESC`, backupConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup records: %w", err)
	}
	defer rows.Close()

	var list []*models.BackupRecord
	for rows.Next() {
		var rec models.BackupRecord
		if err := rows.Scan(&rec.ID, &rec.BackupConfigID, &rec.ProjectID, &rec.DatabaseID, &rec.Status, &rec.FilePath, &rec.FileSizeBytes, &rec.S3URL, &rec.Logs, &rec.StartedAt, &rec.CompletedAt); err != nil {
			return nil, err
		}
		list = append(list, &rec)
	}
	return list, nil
}

func (r *BackupSQLiteRepository) GetRecordByID(ctx context.Context, id string) (*models.BackupRecord, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, backup_config_id, project_id, COALESCE(database_id, ''), status, COALESCE(file_path, ''), file_size_bytes, COALESCE(s3_url, ''), COALESCE(logs, ''), created_at, COALESCE(completed_at, '')
		FROM backup_records WHERE id = ?`, id)

	var rec models.BackupRecord
	err := row.Scan(&rec.ID, &rec.BackupConfigID, &rec.ProjectID, &rec.DatabaseID, &rec.Status, &rec.FilePath, &rec.FileSizeBytes, &rec.S3URL, &rec.Logs, &rec.StartedAt, &rec.CompletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

func (r *BackupSQLiteRepository) UpdateRecord(ctx context.Context, rec *models.BackupRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	res, err := r.db.ExecContext(ctx, `
		UPDATE backup_records 
		SET status = ?, file_path = ?, s3_url = ?, logs = ?, file_size_bytes = ?, completed_at = ?
		WHERE id = ?`,
		rec.Status, rec.FilePath, rec.S3URL, rec.Logs, rec.FileSizeBytes, rec.CompletedAt, rec.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("backup record not found")
	}
	return nil
}
