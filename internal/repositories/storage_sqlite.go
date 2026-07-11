package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/utils"
)

type StorageRepository interface {
	Create(ctx context.Context, s *models.Storage) error
	GetByID(ctx context.Context, id string) (*models.Storage, error)
	List(ctx context.Context) ([]*models.Storage, error)
	ListByProject(ctx context.Context, projectID string) ([]*models.Storage, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, s *models.Storage) error
}

type StorageSQLiteRepository struct {
	db    *sql.DB
	mu    sync.Mutex
	vault Vault
}

func NewStorageSQLiteRepository(db *sql.DB, vault Vault) *StorageSQLiteRepository {
	return &StorageSQLiteRepository{db: db, vault: vault}
}

const listStorageQuery = `SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, type, api_port, console_port, access_key, encrypted_secret_key, bucket_name, COALESCE(volume_path, ''), COALESCE(container_id, ''), COALESCE(status, 'stopped'), COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM storage`

func scanStorage(scanner interface {
	Scan(dest ...any) error
}, s *models.Storage, encryptedSecretKey *string,
) error {
	return scanner.Scan(
		&s.ID, &s.ProjectID, &s.EnvironmentID, &s.Name, &s.Type,
		&s.APIPort, &s.ConsolePort, &s.AccessKey, encryptedSecretKey,
		&s.BucketName, &s.VolumePath, &s.ContainerID, &s.Status,
		&s.InternalDNS, &s.ExternalDNS, &s.CreatedAt, &s.UpdatedAt,
	)
}

func (r *StorageSQLiteRepository) decryptSecretKey(encrypted string, s *models.Storage) {
	if plain, err := r.vault.Decrypt(encrypted); err == nil {
		s.SecretKey = plain
	}
}

func (r *StorageSQLiteRepository) Create(_ context.Context, s *models.Storage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	encryptedSecretKey, err := r.vault.Encrypt(s.SecretKey)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`INSERT INTO storage (
		id, project_id, environment_id, name, type, api_port, console_port,
		access_key, encrypted_secret_key, bucket_name, volume_path,
		container_id, status, internal_dns, external_dns, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.ID, s.ProjectID, s.EnvironmentID, s.Name, s.Type,
		s.APIPort, s.ConsolePort, s.AccessKey, encryptedSecretKey,
		s.BucketName, s.VolumePath, s.ContainerID, s.Status,
		s.InternalDNS, s.ExternalDNS, s.CreatedAt, s.UpdatedAt,
	)
	return err
}

func (r *StorageSQLiteRepository) GetByID(_ context.Context, id string) (*models.Storage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	row := r.db.QueryRow(listStorageQuery+` WHERE id = ?`, id)
	var s models.Storage
	var encryptedSecretKey string
	if err := scanStorage(row, &s, &encryptedSecretKey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("Entity", id)
		}
		return nil, err
	}
	r.decryptSecretKey(encryptedSecretKey, &s)
	return &s, nil
}

func (r *StorageSQLiteRepository) List(_ context.Context) ([]*models.Storage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.Query(listStorageQuery + ` ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.Storage
	for rows.Next() {
		var s models.Storage
		var encryptedSecretKey string
		if err := scanStorage(rows, &s, &encryptedSecretKey); err != nil {
			return nil, err
		}
		r.decryptSecretKey(encryptedSecretKey, &s)
		list = append(list, &s)
	}
	return list, nil
}

func (r *StorageSQLiteRepository) ListByProject(_ context.Context, projectID string) ([]*models.Storage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.Query(listStorageQuery+` WHERE project_id = ? ORDER BY created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.Storage
	for rows.Next() {
		var s models.Storage
		var encryptedSecretKey string
		if err := scanStorage(rows, &s, &encryptedSecretKey); err != nil {
			return nil, err
		}
		r.decryptSecretKey(encryptedSecretKey, &s)
		list = append(list, &s)
	}
	return list, nil
}

func (r *StorageSQLiteRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM storage WHERE id = ?`, id)
	return err
}

func (r *StorageSQLiteRepository) Update(_ context.Context, s *models.Storage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s.UpdatedAt = time.Now()
	encryptedSecretKey, err := r.vault.Encrypt(s.SecretKey)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`UPDATE storage SET project_id = ?, environment_id = ?, name = ?, type = ?, api_port = ?, console_port = ?, access_key = ?, encrypted_secret_key = ?, bucket_name = ?, volume_path = ?, container_id = ?, status = ?, internal_dns = ?, external_dns = ?, updated_at = ? WHERE id = ?`,
		s.ProjectID, s.EnvironmentID, s.Name, s.Type, s.APIPort, s.ConsolePort, s.AccessKey, encryptedSecretKey, s.BucketName, s.VolumePath, s.ContainerID, s.Status, s.InternalDNS, s.ExternalDNS, s.UpdatedAt, s.ID)
	return err
}
