package repositories

import (
	"codedock.dev/codedock/internal/utils"
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock.dev/codedock/internal/models"
)

type DatabaseRepository interface {
	Create(ctx context.Context, db *models.Database) error
	GetByID(ctx context.Context, id string) (*models.Database, error)
	List(ctx context.Context) ([]*models.Database, error)
	ListByProject(ctx context.Context, projectID string) ([]*models.Database, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, db *models.Database) error
}

type DatabaseRepo struct {
	db    *sqlx.DB
	mu    sync.Mutex
	vault Vault
}

func NewDatabaseRepo(db *sql.DB, vault Vault) *DatabaseRepo {
	return &DatabaseRepo{db: sqlx.NewDb(db, "sqlite"), vault: vault}
}

func (r *DatabaseRepo) Create(_ context.Context, db *models.Database) error {
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
		id, project_id, environment_id, name, engine, version, port, username, encrypted_password, database_name, volume_path, container_id, status, internal_dns, external_dns, custom_args, logical_replication, cpu_limit, memory_limit, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		db.ID, db.ProjectID, db.EnvironmentID, db.Name, db.Engine, db.Version, db.Port, db.Username, encryptedPassword, db.DatabaseName, db.VolumePath, db.ContainerID, db.Status, db.InternalDNS, db.ExternalDNS, db.CustomArgs, db.LogicalReplication, db.CPULimit, db.MemoryLimit, db.CreatedAt, db.UpdatedAt)
	return err
}

const listDatabaseQuery = `SELECT id, COALESCE(project_id, '') AS project_id, COALESCE(environment_id, '') AS environment_id, name, engine, version, port, username, encrypted_password, database_name, volume_path, COALESCE(container_id, '') AS container_id, status, COALESCE(internal_dns, '') AS internal_dns, COALESCE(external_dns, '') AS external_dns, COALESCE(custom_args, '') AS custom_args, COALESCE(logical_replication, 0) AS logical_replication, COALESCE(cpu_limit, 0) AS cpu_limit, COALESCE(memory_limit, 0) AS memory_limit, created_at, updated_at FROM databases`

func (r *DatabaseRepo) decryptPassword(encrypted string, d *models.Database) {
	if plain, err := r.vault.Decrypt(encrypted); err == nil {
		d.Password = plain
	}
}

func (r *DatabaseRepo) GetByID(_ context.Context, id string) (*models.Database, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var d models.Database
	err := r.db.Get(&d, listDatabaseQuery+` WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Entity", id)
	}
	if err != nil {
		return nil, err
	}
	r.decryptPassword(d.EncryptedPassword, &d)
	return &d, nil
}

func (r *DatabaseRepo) List(_ context.Context) ([]*models.Database, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.Database
	err := r.db.Select(&list, listDatabaseQuery+` ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*models.Database, 0)
	}
	for _, d := range list {
		r.decryptPassword(d.EncryptedPassword, d)
	}
	return list, nil
}

func (r *DatabaseRepo) ListByProject(_ context.Context, projectID string) ([]*models.Database, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.Database
	err := r.db.Select(&list, listDatabaseQuery+` WHERE project_id = ? ORDER BY created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*models.Database, 0)
	}
	for _, d := range list {
		r.decryptPassword(d.EncryptedPassword, d)
	}
	return list, nil
}

func (r *DatabaseRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM databases WHERE id = ?`, id)
	return err
}

func (r *DatabaseRepo) Update(_ context.Context, db *models.Database) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	db.UpdatedAt = time.Now()
	encryptedPassword, err := r.vault.Encrypt(db.Password)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`UPDATE databases SET project_id = ?, environment_id = ?, name = ?, engine = ?, version = ?, port = ?, username = ?, encrypted_password = ?, database_name = ?, volume_path = ?, container_id = ?, status = ?, internal_dns = ?, external_dns = ?, custom_args = ?, logical_replication = ?, cpu_limit = ?, memory_limit = ?, updated_at = ? WHERE id = ?`,
		db.ProjectID, db.EnvironmentID, db.Name, db.Engine, db.Version, db.Port, db.Username, encryptedPassword, db.DatabaseName, db.VolumePath, db.ContainerID, db.Status, db.InternalDNS, db.ExternalDNS, db.CustomArgs, db.LogicalReplication, db.CPULimit, db.MemoryLimit, db.UpdatedAt, db.ID)
	return err
}
