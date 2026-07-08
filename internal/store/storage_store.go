package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// CreateStorage inserts a new managed storage cluster record and encrypts its secret key.
func (s *Store) CreateStorage(st *types.StorageConfig) error {
	if st.ID == "" {
		st.ID = uuid.NewString()
	}
	now := time.Now()
	st.CreatedAt = now
	st.UpdatedAt = now

	encryptedSecretKey, err := s.vault.Encrypt(st.SecretKey)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`INSERT INTO storage (
		id, project_id, environment_id, name, type, api_port, console_port, access_key, encrypted_secret_key, bucket_name, volume_path, container_id, status, internal_dns, external_dns, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		st.ID, st.ProjectID, st.EnvironmentID, st.Name, st.Type, st.APIPort, st.ConsolePort, st.AccessKey, encryptedSecretKey, st.BucketName, st.VolumePath, st.ContainerID, st.Status, st.InternalDNS, st.ExternalDNS, st.CreatedAt, st.UpdatedAt)
	return err
}

// GetStorage retrieves a single managed object storage instance by ID and decrypts its secret key.
func (s *Store) GetStorage(id string) (*types.StorageConfig, error) {
	var st types.StorageConfig
	var encryptedSecretKey string

	err := s.db.QueryRow(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, type, api_port, console_port, access_key, encrypted_secret_key, bucket_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at
		FROM storage WHERE id = ?`, id).Scan(
		&st.ID, &st.ProjectID, &st.EnvironmentID, &st.Name, &st.Type, &st.APIPort, &st.ConsolePort, &st.AccessKey, &encryptedSecretKey, &st.BucketName, &st.VolumePath, &st.ContainerID, &st.Status, &st.InternalDNS, &st.ExternalDNS, &st.CreatedAt, &st.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	plainSecretKey, err := s.vault.Decrypt(encryptedSecretKey)
	if err == nil {
		st.SecretKey = plainSecretKey
	}
	return &st, nil
}

// ListStorage retrieves all registered object storage instances and decrypts their secret keys.
func (s *Store) ListStorage() ([]types.StorageConfig, error) {
	rows, err := s.db.Query(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, type, api_port, console_port, access_key, encrypted_secret_key, bucket_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM storage ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var storages []types.StorageConfig
	for rows.Next() {
		var st types.StorageConfig
		var encryptedSecretKey string
		if err := rows.Scan(&st.ID, &st.ProjectID, &st.EnvironmentID, &st.Name, &st.Type, &st.APIPort, &st.ConsolePort, &st.AccessKey, &encryptedSecretKey, &st.BucketName, &st.VolumePath, &st.ContainerID, &st.Status, &st.InternalDNS, &st.ExternalDNS, &st.CreatedAt, &st.UpdatedAt); err != nil {
			return nil, err
		}
		if plainSecretKey, err := s.vault.Decrypt(encryptedSecretKey); err == nil {
			st.SecretKey = plainSecretKey
		}
		storages = append(storages, st)
	}
	return storages, nil
}

// ListStorageByProject retrieves all managed object storage instances linked to a specific project identifier.
func (s *Store) ListStorageByProject(projectID string) ([]types.StorageConfig, error) {
	rows, err := s.db.Query(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, type, api_port, console_port, access_key, encrypted_secret_key, bucket_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM storage WHERE project_id = ? ORDER BY created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var storages []types.StorageConfig
	for rows.Next() {
		var st types.StorageConfig
		var encryptedSecretKey string
		if err := rows.Scan(&st.ID, &st.ProjectID, &st.EnvironmentID, &st.Name, &st.Type, &st.APIPort, &st.ConsolePort, &st.AccessKey, &encryptedSecretKey, &st.BucketName, &st.VolumePath, &st.ContainerID, &st.Status, &st.InternalDNS, &st.ExternalDNS, &st.CreatedAt, &st.UpdatedAt); err != nil {
			return nil, err
		}
		if plainSecretKey, err := s.vault.Decrypt(encryptedSecretKey); err == nil {
			st.SecretKey = plainSecretKey
		}
		storages = append(storages, st)
	}
	return storages, nil
}

// ListStorageByEnvironment retrieves all managed object storage instances linked to a specific environment identifier.
func (s *Store) ListStorageByEnvironment(environmentID string) ([]types.StorageConfig, error) {
	rows, err := s.db.Query(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, type, api_port, console_port, access_key, encrypted_secret_key, bucket_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM storage WHERE environment_id = ? ORDER BY created_at ASC`, environmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var storages []types.StorageConfig
	for rows.Next() {
		var st types.StorageConfig
		var encryptedSecretKey string
		if err := rows.Scan(&st.ID, &st.ProjectID, &st.EnvironmentID, &st.Name, &st.Type, &st.APIPort, &st.ConsolePort, &st.AccessKey, &encryptedSecretKey, &st.BucketName, &st.VolumePath, &st.ContainerID, &st.Status, &st.InternalDNS, &st.ExternalDNS, &st.CreatedAt, &st.UpdatedAt); err != nil {
			return nil, err
		}
		if plainSecretKey, err := s.vault.Decrypt(encryptedSecretKey); err == nil {
			st.SecretKey = plainSecretKey
		}
		storages = append(storages, st)
	}
	return storages, nil
}

// DeleteStorage removes an object storage configuration record from SQLite.
func (s *Store) DeleteStorage(id string) error {
	_, err := s.db.Exec(`DELETE FROM storage WHERE id = ?`, id)
	return err
}

// UpdateStorageStatus updates the status and container identifier of an object storage service.
func (s *Store) UpdateStorageStatus(id string, status string, containerID string) error {
	_, err := s.db.Exec(`UPDATE storage SET status = ?, container_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, containerID, id)
	return err
}
