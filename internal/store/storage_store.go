package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// CreateStorage inserts a new object storage instance record and encrypts its secret key.
func (s *Store) CreateStorage(sc *types.StorageConfig) error {
	if sc.ID == "" {
		sc.ID = uuid.NewString()
	}
	now := time.Now()
	sc.CreatedAt = now
	sc.UpdatedAt = now

	encryptedSecret, err := s.vault.Encrypt(sc.SecretKey)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`INSERT INTO storage (
		id, name, type, api_port, console_port, access_key, encrypted_secret_key, bucket_name, volume_path, container_id, status, internal_dns, external_dns, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sc.ID, sc.Name, sc.Type, sc.APIPort, sc.ConsolePort, sc.AccessKey, encryptedSecret, sc.BucketName, sc.VolumePath, sc.ContainerID, sc.Status, sc.InternalDNS, sc.ExternalDNS, sc.CreatedAt, sc.UpdatedAt)
	return err
}

// GetStorage retrieves a single object storage configuration and decrypts its secret key.
func (s *Store) GetStorage(id string) (*types.StorageConfig, error) {
	var sc types.StorageConfig
	var encryptedSecret string

	err := s.db.QueryRow(`SELECT id, name, type, api_port, console_port, access_key, encrypted_secret_key, bucket_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at
		FROM storage WHERE id = ?`, id).Scan(
		&sc.ID, &sc.Name, &sc.Type, &sc.APIPort, &sc.ConsolePort, &sc.AccessKey, &encryptedSecret, &sc.BucketName, &sc.VolumePath, &sc.ContainerID, &sc.Status, &sc.InternalDNS, &sc.ExternalDNS, &sc.CreatedAt, &sc.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	plainSecret, err := s.vault.Decrypt(encryptedSecret)
	if err == nil {
		sc.SecretKey = plainSecret
	}
	return &sc, nil
}

// ListStorage retrieves all object storage configurations and decrypts their secret keys.
func (s *Store) ListStorage() ([]types.StorageConfig, error) {
	rows, err := s.db.Query(`SELECT id, name, type, api_port, console_port, access_key, encrypted_secret_key, bucket_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM storage ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var storages []types.StorageConfig
	for rows.Next() {
		var sc types.StorageConfig
		var encryptedSecret string
		if err := rows.Scan(&sc.ID, &sc.Name, &sc.Type, &sc.APIPort, &sc.ConsolePort, &sc.AccessKey, &encryptedSecret, &sc.BucketName, &sc.VolumePath, &sc.ContainerID, &sc.Status, &sc.InternalDNS, &sc.ExternalDNS, &sc.CreatedAt, &sc.UpdatedAt); err != nil {
			return nil, err
		}
		if plainSecret, err := s.vault.Decrypt(encryptedSecret); err == nil {
			sc.SecretKey = plainSecret
		}
		storages = append(storages, sc)
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
