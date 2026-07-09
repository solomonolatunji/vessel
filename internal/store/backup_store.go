package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// CreateBackupConfig inserts a new automated backup schedule into the store.
func (s *Store) CreateBackupConfig(cfg *types.BackupConfig) error {
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

	s.mu.Lock()
	defer s.mu.Unlock()

	query := `INSERT INTO backup_configs (id, project_id, database_id, storage_id, s3_destination_id, name, schedule, retention_days, status, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, cfg.ID, cfg.ProjectID, cfg.DatabaseID, cfg.StorageID, cfg.S3DestinationID, cfg.Name, cfg.Schedule, cfg.RetentionDays, cfg.Status, cfg.CreatedAt, cfg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create backup config: %w", err)
	}
	return nil
}

// GetBackupConfig retrieves a backup config by ID.
func (s *Store) GetBackupConfig(id string) (*types.BackupConfig, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, project_id, database_id, storage_id, s3_destination_id, name, schedule, retention_days, status, created_at, updated_at
	          FROM backup_configs WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var cfg types.BackupConfig
	err := row.Scan(&cfg.ID, &cfg.ProjectID, &cfg.DatabaseID, &cfg.StorageID, &cfg.S3DestinationID, &cfg.Name, &cfg.Schedule, &cfg.RetentionDays, &cfg.Status, &cfg.CreatedAt, &cfg.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get backup config %s: %w", id, err)
	}
	return &cfg, nil
}

// ListBackupConfigs retrieves all backup configs for a project.
func (s *Store) ListBackupConfigs(projectID string) ([]*types.BackupConfig, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, project_id, database_id, storage_id, s3_destination_id, name, schedule, retention_days, status, created_at, updated_at
	          FROM backup_configs WHERE project_id = ? ORDER BY created_at DESC`
	rows, err := s.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup configs: %w", err)
	}
	defer rows.Close()

	var list []*types.BackupConfig
	for rows.Next() {
		var cfg types.BackupConfig
		if err := rows.Scan(&cfg.ID, &cfg.ProjectID, &cfg.DatabaseID, &cfg.StorageID, &cfg.S3DestinationID, &cfg.Name, &cfg.Schedule, &cfg.RetentionDays, &cfg.Status, &cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &cfg)
	}
	return list, nil
}

// ListAllActiveBackupConfigs returns all active scheduled backups across all projects.
func (s *Store) ListAllActiveBackupConfigs() ([]*types.BackupConfig, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, project_id, database_id, storage_id, s3_destination_id, name, schedule, retention_days, status, created_at, updated_at
	          FROM backup_configs WHERE status = 'active'`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active backup configs: %w", err)
	}
	defer rows.Close()

	var list []*types.BackupConfig
	for rows.Next() {
		var cfg types.BackupConfig
		if err := rows.Scan(&cfg.ID, &cfg.ProjectID, &cfg.DatabaseID, &cfg.StorageID, &cfg.S3DestinationID, &cfg.Name, &cfg.Schedule, &cfg.RetentionDays, &cfg.Status, &cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &cfg)
	}
	return list, nil
}

// DeleteBackupConfig removes a backup schedule from the database.
func (s *Store) DeleteBackupConfig(id, projectID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec("DELETE FROM backup_configs WHERE id = ? AND project_id = ?", id, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete backup config: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("backup config not found or unauthorized")
	}
	return nil
}

// CreateBackupRecord stores a new execution log entry when a backup starts.
func (s *Store) CreateBackupRecord(rec *types.BackupRecord) error {
	if rec.ID == "" {
		rec.ID = uuid.New().String()
	}
	if rec.StartedAt == "" {
		rec.StartedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if rec.Status == "" {
		rec.Status = "running"
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	query := `INSERT INTO backup_records (id, backup_config_id, project_id, database_id, status, file_path, file_size_bytes, s3_url, logs, started_at, completed_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, rec.ID, rec.BackupConfigID, rec.ProjectID, rec.DatabaseID, rec.Status, rec.FilePath, rec.FileSizeBytes, rec.S3URL, rec.Logs, rec.StartedAt, rec.CompletedAt)
	if err != nil {
		return fmt.Errorf("failed to create backup record: %w", err)
	}
	return nil
}

// UpdateBackupRecord updates the status, file details, and completion timestamp of a backup run.
func (s *Store) UpdateBackupRecord(id, status, filePath, s3URL, logs string, fileSizeBytes int64, completedAt string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `UPDATE backup_records SET status = ?, file_path = ?, s3_url = ?, logs = ?, file_size_bytes = ?, completed_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, status, filePath, s3URL, logs, fileSizeBytes, completedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update backup record %s: %w", id, err)
	}
	return nil
}

// ListBackupRecords returns all backup run execution logs for a backup schedule.
func (s *Store) ListBackupRecords(backupConfigID string) ([]*types.BackupRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, backup_config_id, project_id, database_id, status, file_path, file_size_bytes, s3_url, logs, started_at, completed_at
	          FROM backup_records WHERE backup_config_id = ? ORDER BY started_at DESC`
	rows, err := s.db.Query(query, backupConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup records: %w", err)
	}
	defer rows.Close()

	var list []*types.BackupRecord
	for rows.Next() {
		var rec types.BackupRecord
		if err := rows.Scan(&rec.ID, &rec.BackupConfigID, &rec.ProjectID, &rec.DatabaseID, &rec.Status, &rec.FilePath, &rec.FileSizeBytes, &rec.S3URL, &rec.Logs, &rec.StartedAt, &rec.CompletedAt); err != nil {
			return nil, err
		}
		list = append(list, &rec)
	}
	return list, nil
}

// GetBackupRecord retrieves a backup run record by ID.
func (s *Store) GetBackupRecord(id string) (*types.BackupRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, backup_config_id, project_id, database_id, status, file_path, file_size_bytes, s3_url, logs, started_at, completed_at
	          FROM backup_records WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var rec types.BackupRecord
	err := row.Scan(&rec.ID, &rec.BackupConfigID, &rec.ProjectID, &rec.DatabaseID, &rec.Status, &rec.FilePath, &rec.FileSizeBytes, &rec.S3URL, &rec.Logs, &rec.StartedAt, &rec.CompletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get backup record %s: %w", id, err)
	}
	return &rec, nil
}

// CreateS3Destination registers a new offsite S3/MinIO destination.
func (s *Store) CreateS3Destination(dest *types.S3Destination) error {
	if dest.ID == "" {
		dest.ID = uuid.New().String()
	}
	if dest.CreatedAt == "" {
		dest.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	query := `INSERT INTO s3_destinations (id, project_id, name, endpoint, bucket, region, access_key_id, secret_access_key, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, dest.ID, dest.ProjectID, dest.Name, dest.Endpoint, dest.Bucket, dest.Region, dest.AccessKeyID, dest.SecretAccessKey, dest.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create s3 destination: %w", err)
	}
	return nil
}

// GetS3Destination retrieves an S3 destination by ID.
func (s *Store) GetS3Destination(id string) (*types.S3Destination, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, project_id, name, endpoint, bucket, region, access_key_id, secret_access_key, created_at
	          FROM s3_destinations WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var dest types.S3Destination
	err := row.Scan(&dest.ID, &dest.ProjectID, &dest.Name, &dest.Endpoint, &dest.Bucket, &dest.Region, &dest.AccessKeyID, &dest.SecretAccessKey, &dest.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get s3 destination %s: %w", id, err)
	}
	return &dest, nil
}

// ListS3Destinations lists all registered S3 backup targets for a project.
func (s *Store) ListS3Destinations(projectID string) ([]*types.S3Destination, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, project_id, name, endpoint, bucket, region, access_key_id, secret_access_key, created_at
	          FROM s3_destinations WHERE project_id = ? ORDER BY created_at DESC`
	rows, err := s.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list s3 destinations: %w", err)
	}
	defer rows.Close()

	var list []*types.S3Destination
	for rows.Next() {
		var dest types.S3Destination
		if err := rows.Scan(&dest.ID, &dest.ProjectID, &dest.Name, &dest.Endpoint, &dest.Bucket, &dest.Region, &dest.AccessKeyID, &dest.SecretAccessKey, &dest.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &dest)
	}
	return list, nil
}

// DeleteS3Destination removes an S3 destination target.
func (s *Store) DeleteS3Destination(id, projectID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec("DELETE FROM s3_destinations WHERE id = ? AND project_id = ?", id, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete s3 destination: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("s3 destination not found or unauthorized")
	}
	return nil
}
