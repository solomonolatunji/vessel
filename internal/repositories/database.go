package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"
	"vessl.dev/vessl/internal/utils"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
)

type DatabaseRepository interface {
	Create(ctx context.Context, db *models.Database) error
	GetByID(ctx context.Context, id string) (*models.Database, error)
	List(ctx context.Context) ([]*models.Database, error)
	ListByProject(ctx context.Context, projectID string) ([]*models.Database, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, db *models.Database) error
}

type DatabaseSQLiteRepository struct {
	db    *sql.DB
	mu    sync.Mutex
	vault Vault
}

func NewDatabaseSQLiteRepository(db *sql.DB, vault Vault) *DatabaseSQLiteRepository {
	return &DatabaseSQLiteRepository{db: db, vault: vault}
}

func (r *DatabaseSQLiteRepository) Create(_ context.Context, db *models.Database) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if db.ID == "" {
		db.ID = uuid.NewString()
	}
	now := time.Now()
	db.CreatedAt = now
	db.UpdatedAt = now
	encryptedPassword, err := r.vault.Encrypt(db.Password)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`INSERT INTO databases (
		id, project_id, environment_id, name, engine, version, port, username, encrypted_password, database_name, volume_path, container_id, status, internal_dns, external_dns, custom_args, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		db.ID, db.ProjectID, db.EnvironmentID, db.Name, db.Engine, db.Version, db.Port, db.Username, encryptedPassword, db.DatabaseName, db.VolumePath, db.ContainerID, db.Status, db.InternalDNS, db.ExternalDNS, db.CustomArgs, db.CreatedAt, db.UpdatedAt)
	return err
}

const listDatabaseQuery = `SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, engine, version, port, username, encrypted_password, database_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), COALESCE(custom_args, ''), created_at, updated_at FROM databases`

func scanDatabase(scanner interface {
	Scan(dest ...any) error
}, d *models.Database, encryptedPassword *string,
) error {
	return scanner.Scan(
		&d.ID, &d.ProjectID, &d.EnvironmentID, &d.Name, &d.Engine, &d.Version,
		&d.Port, &d.Username, encryptedPassword, &d.DatabaseName, &d.VolumePath,
		&d.ContainerID, &d.Status, &d.InternalDNS, &d.ExternalDNS, &d.CustomArgs, &d.CreatedAt, &d.UpdatedAt,
	)
}

func (r *DatabaseSQLiteRepository) decryptPassword(encrypted string, d *models.Database) {
	if plain, err := r.vault.Decrypt(encrypted); err == nil {
		d.Password = plain
	}
}

func (r *DatabaseSQLiteRepository) GetByID(_ context.Context, id string) (*models.Database, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	row := r.db.QueryRow(listDatabaseQuery+` WHERE id = ?`, id)
	var d models.Database
	var encryptedPassword string
	if err := scanDatabase(row, &d, &encryptedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("Entity", id)
		}
		return nil, err
	}
	r.decryptPassword(encryptedPassword, &d)
	return &d, nil
}

func (r *DatabaseSQLiteRepository) List(_ context.Context) ([]*models.Database, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.Query(listDatabaseQuery + ` ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.Database
	for rows.Next() {
		var d models.Database
		var encryptedPassword string
		if err := scanDatabase(rows, &d, &encryptedPassword); err != nil {
			return nil, err
		}
		r.decryptPassword(encryptedPassword, &d)
		list = append(list, &d)
	}
	return list, nil
}

func (r *DatabaseSQLiteRepository) ListByProject(_ context.Context, projectID string) ([]*models.Database, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.Query(listDatabaseQuery+` WHERE project_id = ? ORDER BY created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.Database
	for rows.Next() {
		var d models.Database
		var encryptedPassword string
		if err := scanDatabase(rows, &d, &encryptedPassword); err != nil {
			return nil, err
		}
		r.decryptPassword(encryptedPassword, &d)
		list = append(list, &d)
	}
	return list, nil
}

func (r *DatabaseSQLiteRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM databases WHERE id = ?`, id)
	return err
}

func (r *DatabaseSQLiteRepository) Update(_ context.Context, db *models.Database) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	db.UpdatedAt = time.Now()
	encryptedPassword, err := r.vault.Encrypt(db.Password)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`UPDATE databases SET project_id = ?, environment_id = ?, name = ?, engine = ?, version = ?, port = ?, username = ?, encrypted_password = ?, database_name = ?, volume_path = ?, container_id = ?, status = ?, internal_dns = ?, external_dns = ?, custom_args = ?, updated_at = ? WHERE id = ?`,
		db.ProjectID, db.EnvironmentID, db.Name, db.Engine, db.Version, db.Port, db.Username, encryptedPassword, db.DatabaseName, db.VolumePath, db.ContainerID, db.Status, db.InternalDNS, db.ExternalDNS, db.CustomArgs, db.UpdatedAt, db.ID)
	return err
}
