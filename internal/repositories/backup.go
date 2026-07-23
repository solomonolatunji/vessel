package repositories

import (
	"codedock.dev/codedock/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"

	"codedock.dev/codedock/internal/models"
)

type BackupRepository interface {
	CreateConfig(ctx context.Context, cfg *models.BackupConfig) error
	UpdateConfig(ctx context.Context, cfg *models.BackupConfig) error
	GetConfigByID(ctx context.Context, id string) (*models.BackupConfig, error)
	ListConfigs(ctx context.Context) ([]*models.BackupConfig, error)
	ListAllActiveConfigs(ctx context.Context) ([]*models.BackupConfig, error)
	DeleteConfig(ctx context.Context, id string) error
	CreateRecord(ctx context.Context, rec *models.BackupRecord) error
	GetRecordByID(ctx context.Context, id string) (*models.BackupRecord, error)
	ListRecordsByConfig(ctx context.Context, backupConfigID string) ([]*models.BackupRecord, error)
	UpdateRecord(ctx context.Context, rec *models.BackupRecord) error
	DeleteRecord(ctx context.Context, id string) error
}

type BackupRepo struct {
	db    *sqlx.DB
	mu    sync.Mutex
	vault Vault
}

func NewBackupRepo(db *sql.DB, v Vault) *BackupRepo {
	return &BackupRepo{db: sqlx.NewDb(db, "sqlite"), vault: v}
}

func (r *BackupRepo) EnsureTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS backup_configs (
			id TEXT PRIMARY KEY,
			database_id TEXT,
			storage_id TEXT,
			s3_destination_id TEXT,
			name TEXT NOT NULL,
			description TEXT,
			db_user TEXT,
			db_password TEXT,
			backup_enabled INTEGER DEFAULT 1,
			s3_enabled INTEGER DEFAULT 0,
			disable_local INTEGER DEFAULT 0,
			schedule TEXT NOT NULL,
			timezone TEXT DEFAULT 'UTC',
			timeout INTEGER DEFAULT 3600,
			retention_days INTEGER DEFAULT 7,
			max_backups INTEGER DEFAULT 0,
			max_storage_gb INTEGER DEFAULT 0,
			status TEXT DEFAULT 'active',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS backup_records (
			id TEXT PRIMARY KEY,
			backup_config_id TEXT NOT NULL,
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
			description TEXT DEFAULT '',
			provider TEXT DEFAULT 's3',
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

func (r *BackupRepo) CreateConfig(ctx context.Context, cfg *models.BackupConfig) error {
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
	if cfg.Timeout == 0 {
		cfg.Timeout = 3600
	}
	if cfg.RetentionDays < 0 {
		cfg.RetentionDays = 0
	}
	if cfg.DbPassword != "" && r.vault != nil {
		enc, err := r.vault.Encrypt(cfg.DbPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt db password: %w", err)
		}
		cfg.DbPassword = enc
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO backup_configs (id, database_id, s3_destination_id, name, description, db_user, db_password, backup_enabled, s3_enabled, disable_local, schedule, timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		cfg.ID, cfg.DatabaseID, cfg.S3DestinationID, cfg.Name, cfg.Description, cfg.DbUser, cfg.DbPassword, cfg.BackupEnabled, cfg.S3Enabled, cfg.DisableLocal, cfg.Schedule, cfg.Timezone, cfg.Timeout, cfg.RetentionDays, cfg.MaxBackups, cfg.MaxStorageGB, cfg.Status, cfg.CreatedAt, cfg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create backup config: %w", err)
	}
	cfg.DbPassword = "********"
	return nil
}

func (r *BackupRepo) GetConfigByID(ctx context.Context, id string) (*models.BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg models.BackupConfig
	err := r.db.GetContext(ctx, &cfg, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(s3_destination_id, '') as s3_destination_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, COALESCE(db_password, '') as db_password, backup_enabled, s3_enabled, disable_local, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at
		FROM backup_configs WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Config", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get backup config %s: %w", id, err)
	}
	if cfg.DbPassword != "" && r.vault != nil {
		dec, err := r.vault.Decrypt(cfg.DbPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt db password: %w", err)
		}
		cfg.DbPassword = dec
	}
	return &cfg, nil
}

func (r *BackupRepo) UpdateConfig(ctx context.Context, cfg *models.BackupConfig) error {
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if cfg.DbPassword != "" && cfg.DbPassword != "********" && r.vault != nil {
		enc, err := r.vault.Encrypt(cfg.DbPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt db password: %w", err)
		}
		cfg.DbPassword = enc
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if cfg.DbPassword == "********" || cfg.DbPassword == "" {
		res, err := r.db.ExecContext(ctx, `UPDATE backup_configs SET database_id=?, s3_destination_id=?, name=?, description=?, db_user=?, backup_enabled=?, s3_enabled=?, disable_local=?, schedule=?, timezone=?, timeout=?, retention_days=?, max_backups=?, max_storage_gb=?, updated_at=? WHERE id=?`,
			cfg.DatabaseID, cfg.S3DestinationID, cfg.Name, cfg.Description, cfg.DbUser, cfg.BackupEnabled, cfg.S3Enabled, cfg.DisableLocal, cfg.Schedule, cfg.Timezone, cfg.Timeout, cfg.RetentionDays, cfg.MaxBackups, cfg.MaxStorageGB, cfg.UpdatedAt, cfg.ID)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}
		if affected == 0 {
			return utils.NewNotFoundError("BackupConfig", cfg.ID)
		}
		return nil
	}

	res, err := r.db.ExecContext(ctx, `UPDATE backup_configs SET database_id=?, s3_destination_id=?, name=?, description=?, db_user=?, db_password=?, backup_enabled=?, s3_enabled=?, disable_local=?, schedule=?, timezone=?, timeout=?, retention_days=?, max_backups=?, max_storage_gb=?, updated_at=? WHERE id=?`,
		cfg.DatabaseID, cfg.S3DestinationID, cfg.Name, cfg.Description, cfg.DbUser, cfg.DbPassword, cfg.BackupEnabled, cfg.S3Enabled, cfg.DisableLocal, cfg.Schedule, cfg.Timezone, cfg.Timeout, cfg.RetentionDays, cfg.MaxBackups, cfg.MaxStorageGB, cfg.UpdatedAt, cfg.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if affected == 0 {
		return utils.NewNotFoundError("BackupConfig", cfg.ID)
	}
	return nil
}

func (r *BackupRepo) ListConfigs(ctx context.Context) ([]*models.BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.BackupConfig
	err := r.db.SelectContext(ctx, &list, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(s3_destination_id, '') as s3_destination_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, COALESCE(db_password, '') as db_password, backup_enabled, s3_enabled, disable_local, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at
		FROM backup_configs ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup configs: %w", err)
	}
	if list == nil {
		list = make([]*models.BackupConfig, 0)
	}
	if r.vault != nil {
		for _, cfg := range list {
			if cfg.DbPassword != "" {
				dec, err := r.vault.Decrypt(cfg.DbPassword)
				if err != nil {
					return nil, fmt.Errorf("failed to decrypt db password: %w", err)
				}
				cfg.DbPassword = dec
			}
		}
	}
	return list, nil
}

func (r *BackupRepo) ListAllActiveConfigs(ctx context.Context) ([]*models.BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.BackupConfig
	err := r.db.SelectContext(ctx, &list, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(s3_destination_id, '') as s3_destination_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, COALESCE(db_password, '') as db_password, backup_enabled, s3_enabled, disable_local, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at
		FROM backup_configs WHERE status = 'active'`)
	if err != nil {
		return nil, fmt.Errorf("failed to list active backup configs: %w", err)
	}
	if list == nil {
		list = make([]*models.BackupConfig, 0)
	}
	if r.vault != nil {
		for _, cfg := range list {
			if cfg.DbPassword != "" {
				dec, err := r.vault.Decrypt(cfg.DbPassword)
				if err != nil {
					return nil, fmt.Errorf("failed to decrypt db password: %w", err)
				}
				cfg.DbPassword = dec
			}
		}
	}
	return list, nil
}

func (r *BackupRepo) DeleteConfig(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.ExecContext(ctx, "DELETE FROM backup_configs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete backup config: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return utils.NewNotFoundError("BackupConfig", id)
	}
	return nil
}

func (r *BackupRepo) CreateRecord(ctx context.Context, rec *models.BackupRecord) error {
	if rec.ID == "" {
		rec.ID = uuid.New().String()
	}
	if rec.StartedAt == "" {
		rec.StartedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if rec.Status == "" {
		rec.Status = models.BackupRecordStatusRunning
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO backup_records (id, backup_config_id, database_id, status, file_path, file_size_bytes, s3_url, logs, started_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rec.ID, rec.BackupConfigID, rec.DatabaseID, rec.Status, rec.FilePath, rec.FileSizeBytes, rec.S3URL, rec.Logs, rec.StartedAt, rec.CompletedAt)
	if err != nil {
		return fmt.Errorf("failed to create backup record: %w", err)
	}
	return nil
}

func (r *BackupRepo) ListRecordsByConfig(ctx context.Context, backupConfigID string) ([]*models.BackupRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.BackupRecord
	err := r.db.SelectContext(ctx, &list, `SELECT id, backup_config_id, COALESCE(database_id, '') as database_id, status, COALESCE(file_path, '') as file_path, file_size_bytes, COALESCE(s3_url, '') as s3_url, COALESCE(logs, '') as logs, started_at, COALESCE(completed_at, '') as completed_at
		FROM backup_records WHERE backup_config_id = ? ORDER BY started_at DESC`, backupConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup records: %w", err)
	}
	if list == nil {
		list = make([]*models.BackupRecord, 0)
	}
	return list, nil
}

func (r *BackupRepo) GetRecordByID(ctx context.Context, id string) (*models.BackupRecord, error) {
	var rec models.BackupRecord
	err := r.db.GetContext(ctx, &rec, `
		SELECT id, backup_config_id, COALESCE(database_id, '') as database_id, status, COALESCE(file_path, '') as file_path, file_size_bytes, COALESCE(s3_url, '') as s3_url, COALESCE(logs, '') as logs, started_at, COALESCE(completed_at, '') as completed_at
		FROM backup_records WHERE id = ?`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("Record", id)
		}
		return nil, err
	}
	return &rec, nil
}

func (r *BackupRepo) UpdateRecord(ctx context.Context, rec *models.BackupRecord) error {
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
		return utils.NewNotFoundError("BackupRecord", rec.ID)
	}
	return nil
}

func (r *BackupRepo) DeleteRecord(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, "DELETE FROM backup_records WHERE id=?", id)
	return err
}
