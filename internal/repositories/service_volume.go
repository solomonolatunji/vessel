package repositories

import (
	"context"
	"database/sql"

	"codedock.dev/codedock/internal/models"
	"github.com/jmoiron/sqlx"
)

type ServiceVolumeRepository interface {
	Create(ctx context.Context, volume *models.ServiceVolume) error
	GetByID(ctx context.Context, id string) (*models.ServiceVolume, error)
	ListByService(ctx context.Context, serviceID string) ([]models.ServiceVolume, error)
	Delete(ctx context.Context, id string) error
}

type serviceVolumeRepo struct {
	db *sqlx.DB
}

func NewServiceVolumeRepo(db *sql.DB) ServiceVolumeRepository {
	return &serviceVolumeRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *serviceVolumeRepo) Create(ctx context.Context, volume *models.ServiceVolume) error {
	query := `
		INSERT INTO service_volumes (id, service_id, host_path, container_path, created_at)
		VALUES (:id, :service_id, :host_path, :container_path, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, volume)
	return err
}

func (r *serviceVolumeRepo) GetByID(ctx context.Context, id string) (*models.ServiceVolume, error) {
	var volume models.ServiceVolume
	query := `SELECT * FROM service_volumes WHERE id = ?`
	err := r.db.GetContext(ctx, &volume, query, id)
	return &volume, err
}

func (r *serviceVolumeRepo) ListByService(ctx context.Context, serviceID string) ([]models.ServiceVolume, error) {
	var volumes []models.ServiceVolume
	query := `SELECT * FROM service_volumes WHERE service_id = ? ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &volumes, query, serviceID)
	return volumes, err
}

func (r *serviceVolumeRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM service_volumes WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
