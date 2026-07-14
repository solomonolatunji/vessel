package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

type ServiceVarRepository interface {
	Create(ctx context.Context, v *models.Variable) error
	Update(ctx context.Context, v *models.Variable) error
	GetByID(ctx context.Context, id string) (*models.Variable, error)
	ListByService(ctx context.Context, serviceID string) ([]*models.Variable, error)
	Delete(ctx context.Context, id string) error
}

type ServiceVarSQLiteRepository struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewServiceVarSQLiteRepository(db *sql.DB) *ServiceVarSQLiteRepository {
	return &ServiceVarSQLiteRepository{db: sqlx.NewDb(db, "sqlite")}
}

func (r *ServiceVarSQLiteRepository) Create(_ context.Context, v *models.Variable) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v.ID == "" {
		v.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	if v.CreatedAt.IsZero() {
		v.CreatedAt = now
	}
	v.UpdatedAt = now

	_, err := r.db.Exec(`INSERT INTO service_vars (id, service_id, environment_id, key, value, is_secret, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(service_id, key, environment_id) DO UPDATE SET value = excluded.value, is_secret = excluded.is_secret, updated_at = excluded.updated_at`,
		v.ID, v.ServiceID, v.EnvironmentID, v.Key, v.Value, v.IsSecret, v.CreatedAt, v.UpdatedAt)
	return err
}

func (r *ServiceVarSQLiteRepository) Update(_ context.Context, v *models.Variable) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	v.UpdatedAt = time.Now().UTC()

	_, err := r.db.Exec(`UPDATE service_vars SET key = ?, value = ?, is_secret = ?, updated_at = ? WHERE id = ?`,
		v.Key, v.Value, v.IsSecret, v.UpdatedAt, v.ID)
	return err
}

func (r *ServiceVarSQLiteRepository) GetByID(_ context.Context, id string) (*models.Variable, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var v models.Variable

	err := r.db.Get(&v, `SELECT id, service_id, COALESCE(environment_id, '') AS environment_id, key, value, is_secret, created_at, updated_at FROM service_vars WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Entity", id)
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *ServiceVarSQLiteRepository) ListByService(_ context.Context, serviceID string) ([]*models.Variable, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var list []*models.Variable
	err := r.db.Select(&list, `SELECT id, service_id, COALESCE(environment_id, '') AS environment_id, key, value, is_secret, created_at, updated_at FROM service_vars WHERE service_id = ? ORDER BY key ASC`, serviceID)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*models.Variable, 0)
	}
	return list, nil
}

func (r *ServiceVarSQLiteRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM service_vars WHERE id = ?`, id)
	return err
}
