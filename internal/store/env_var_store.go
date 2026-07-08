package store

import (
	"time"

	"github.com/google/uuid"
)

// SetEnvVar encrypts a plaintext environment variable value and stores it in SQLite.
func (s *Store) SetEnvVar(projectID, key, plaintextValue string) error {
	encrypted, err := s.vault.Encrypt(plaintextValue)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = s.db.Exec(`INSERT INTO env_vars (id, project_id, key, encrypted_value, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(project_id, key) DO UPDATE SET encrypted_value = excluded.encrypted_value, updated_at = excluded.updated_at`,
		uuid.NewString(), projectID, key, encrypted, now, now)
	return err
}

// GetEnvVars retrieves and decrypts all environment variables for a given project ID.
func (s *Store) GetEnvVars(projectID string) (map[string]string, error) {
	rows, err := s.db.Query(`SELECT key, encrypted_value FROM env_vars WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	envs := make(map[string]string)
	for rows.Next() {
		var key, encrypted string
		if err := rows.Scan(&key, &encrypted); err != nil {
			return nil, err
		}
		plaintext, err := s.vault.Decrypt(encrypted)
		if err != nil {
			continue
		}
		envs[key] = plaintext
	}
	return envs, nil
}
