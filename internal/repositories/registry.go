package repositories

import (
	"context"
	"database/sql"
	"time"

	"codedock.run/codedock/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RegistryRepository interface {
	Create(ctx context.Context, registry *models.Registry) error
	ListByProject(ctx context.Context, projectID string) ([]*models.Registry, error)
	Get(ctx context.Context, id string) (*models.Registry, error)
	Delete(ctx context.Context, id string) error
}

type sqliteRegistryRepository struct {
	db *sqlx.DB
}

func NewRegistryRepository(db *sql.DB) RegistryRepository {
	return &sqliteRegistryRepository{db: sqlx.NewDb(db, "sqlite")}
}

func (r *sqliteRegistryRepository) Create(ctx context.Context, registry *models.Registry) error {
	if registry.ID == "" {
		registry.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	registry.CreatedAt = now
	registry.UpdatedAt = now

	query := `
		INSERT INTO registries (id, project_id, name, registry_url, username, password_token, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		registry.ID, registry.ProjectID, registry.Name, registry.RegistryURL,
		registry.Username, registry.PasswordToken, registry.CreatedAt, registry.UpdatedAt)
	return err
}

func (r *sqliteRegistryRepository) ListByProject(ctx context.Context, projectID string) ([]*models.Registry, error) {
	var registries []*models.Registry
	query := `SELECT id, project_id, name, registry_url, username, password_token, created_at, updated_at FROM registries WHERE project_id = ? ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &registries, query, projectID)
	if err != nil {
		return nil, err
	}
	if registries == nil {
		registries = make([]*models.Registry, 0)
	}
	return registries, nil
}

func (r *sqliteRegistryRepository) Get(ctx context.Context, id string) (*models.Registry, error) {
	var registry models.Registry
	query := `SELECT id, project_id, name, registry_url, username, password_token, created_at, updated_at FROM registries WHERE id = ?`
	err := r.db.GetContext(ctx, &registry, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &registry, err
}

func (r *sqliteRegistryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM registries WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
