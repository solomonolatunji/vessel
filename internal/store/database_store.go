package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// CreateDatabase inserts a new managed database instance record and encrypts its password.
func (s *Store) CreateDatabase(db *types.DatabaseConfig) error {
	if db.ID == "" {
		db.ID = uuid.NewString()
	}
	now := time.Now()
	db.CreatedAt = now
	db.UpdatedAt = now

	encryptedPassword, err := s.vault.Encrypt(db.Password)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`INSERT INTO databases (
		id, project_id, environment_id, name, engine, version, port, username, encrypted_password, database_name, volume_path, container_id, status, internal_dns, external_dns, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		db.ID, db.ProjectID, db.EnvironmentID, db.Name, db.Engine, db.Version, db.Port, db.Username, encryptedPassword, db.DatabaseName, db.VolumePath, db.ContainerID, db.Status, db.InternalDNS, db.ExternalDNS, db.CreatedAt, db.UpdatedAt)
	return err
}

// GetDatabase retrieves a single managed database by ID and decrypts its password.
func (s *Store) GetDatabase(id string) (*types.DatabaseConfig, error) {
	var db types.DatabaseConfig
	var encryptedPassword string

	err := s.db.QueryRow(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, engine, version, port, username, encrypted_password, database_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at
		FROM databases WHERE id = ?`, id).Scan(
		&db.ID, &db.ProjectID, &db.EnvironmentID, &db.Name, &db.Engine, &db.Version, &db.Port, &db.Username, &encryptedPassword, &db.DatabaseName, &db.VolumePath, &db.ContainerID, &db.Status, &db.InternalDNS, &db.ExternalDNS, &db.CreatedAt, &db.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	plainPassword, err := s.vault.Decrypt(encryptedPassword)
	if err == nil {
		db.Password = plainPassword
	}
	return &db, nil
}

// ListDatabases retrieves all registered managed databases and decrypts their passwords.
func (s *Store) ListDatabases() ([]types.DatabaseConfig, error) {
	rows, err := s.db.Query(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, engine, version, port, username, encrypted_password, database_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM databases ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []types.DatabaseConfig
	for rows.Next() {
		var db types.DatabaseConfig
		var encryptedPassword string
		if err := rows.Scan(&db.ID, &db.ProjectID, &db.EnvironmentID, &db.Name, &db.Engine, &db.Version, &db.Port, &db.Username, &encryptedPassword, &db.DatabaseName, &db.VolumePath, &db.ContainerID, &db.Status, &db.InternalDNS, &db.ExternalDNS, &db.CreatedAt, &db.UpdatedAt); err != nil {
			return nil, err
		}
		if plainPassword, err := s.vault.Decrypt(encryptedPassword); err == nil {
			db.Password = plainPassword
		}
		databases = append(databases, db)
	}
	return databases, nil
}

// ListDatabasesByProject retrieves all managed databases linked to a specific project identifier.
func (s *Store) ListDatabasesByProject(projectID string) ([]types.DatabaseConfig, error) {
	rows, err := s.db.Query(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, engine, version, port, username, encrypted_password, database_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM databases WHERE project_id = ? ORDER BY created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []types.DatabaseConfig
	for rows.Next() {
		var db types.DatabaseConfig
		var encryptedPassword string
		if err := rows.Scan(&db.ID, &db.ProjectID, &db.EnvironmentID, &db.Name, &db.Engine, &db.Version, &db.Port, &db.Username, &encryptedPassword, &db.DatabaseName, &db.VolumePath, &db.ContainerID, &db.Status, &db.InternalDNS, &db.ExternalDNS, &db.CreatedAt, &db.UpdatedAt); err != nil {
			return nil, err
		}
		if plainPassword, err := s.vault.Decrypt(encryptedPassword); err == nil {
			db.Password = plainPassword
		}
		databases = append(databases, db)
	}
	return databases, nil
}

// ListDatabasesByEnvironment retrieves all managed databases linked to a specific environment identifier.
func (s *Store) ListDatabasesByEnvironment(environmentID string) ([]types.DatabaseConfig, error) {
	rows, err := s.db.Query(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, engine, version, port, username, encrypted_password, database_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM databases WHERE environment_id = ? ORDER BY created_at ASC`, environmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []types.DatabaseConfig
	for rows.Next() {
		var db types.DatabaseConfig
		var encryptedPassword string
		if err := rows.Scan(&db.ID, &db.ProjectID, &db.EnvironmentID, &db.Name, &db.Engine, &db.Version, &db.Port, &db.Username, &encryptedPassword, &db.DatabaseName, &db.VolumePath, &db.ContainerID, &db.Status, &db.InternalDNS, &db.ExternalDNS, &db.CreatedAt, &db.UpdatedAt); err != nil {
			return nil, err
		}
		if plainPassword, err := s.vault.Decrypt(encryptedPassword); err == nil {
			db.Password = plainPassword
		}
		databases = append(databases, db)
	}
	return databases, nil
}

// DeleteDatabase removes a managed database configuration record from SQLite.
func (s *Store) DeleteDatabase(id string) error {
	_, err := s.db.Exec(`DELETE FROM databases WHERE id = ?`, id)
	return err
}

// UpdateDatabaseStatus records a state transition and active container ID for a database instance.
func (s *Store) UpdateDatabaseStatus(id string, status string, containerID string) error {
	_, err := s.db.Exec(`UPDATE databases SET status = ?, container_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, containerID, id)
	return err
}
