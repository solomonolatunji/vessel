package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/models"
)

type S3DestinationRepository interface {
	CreateS3Destination(ctx context.Context, dest *models.S3Destination) error
	ListS3Destinations(ctx context.Context, projectID string) ([]*models.S3Destination, error)
	GetS3Destination(ctx context.Context, id string) (*models.S3Destination, error)
	DeleteS3Destination(ctx context.Context, id, projectID string) error
}

type S3DestinationSQLiteRepository struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewS3DestinationSQLiteRepository(db *sql.DB) *S3DestinationSQLiteRepository {
	return &S3DestinationSQLiteRepository{db: sqlx.NewDb(db, "sqlite")}
}

func (r *S3DestinationSQLiteRepository) CreateS3Destination(ctx context.Context, dest *models.S3Destination) error {
	if dest.ID == "" {
		dest.ID = uuid.New().String()
	}
	if dest.CreatedAt == "" {
		dest.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO s3_destinations (id, project_id, name, endpoint, bucket, region, access_key_id, secret_access_key, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		dest.ID, dest.ProjectID, dest.Name, dest.Endpoint, dest.Bucket, dest.Region, dest.AccessKeyID, dest.SecretAccessKey, dest.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create s3 destination: %w", err)
	}
	return nil
}

func (r *S3DestinationSQLiteRepository) ListS3Destinations(ctx context.Context, projectID string) ([]*models.S3Destination, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.S3Destination
	err := r.db.SelectContext(ctx, &list, `SELECT id, project_id, name, endpoint, bucket, COALESCE(region, '') as region, COALESCE(access_key_id, '') as access_key_id, COALESCE(secret_access_key, '') as secret_access_key, created_at
		FROM s3_destinations WHERE project_id = ? ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list s3 destinations: %w", err)
	}
	if list == nil {
		list = make([]*models.S3Destination, 0)
	}
	return list, nil
}

func (r *S3DestinationSQLiteRepository) GetS3Destination(ctx context.Context, id string) (*models.S3Destination, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var dest models.S3Destination
	err := r.db.GetContext(ctx, &dest, `SELECT id, project_id, name, endpoint, bucket, COALESCE(region, '') as region, COALESCE(access_key_id, '') as access_key_id, COALESCE(secret_access_key, '') as secret_access_key, created_at
		FROM s3_destinations WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("S3Destination", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get s3 destination %s: %w", id, err)
	}
	return &dest, nil
}

func (r *S3DestinationSQLiteRepository) DeleteS3Destination(ctx context.Context, id, projectID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.ExecContext(ctx, "DELETE FROM s3_destinations WHERE id = ? AND project_id = ?", id, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete s3 destination: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return utils.NewNotFoundError("S3Destination", id)
	}
	return nil
}
