package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
	"vessl.dev/vessl/internal/utils"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
)

type S3DestinationRepository interface {
	CreateS3Destination(ctx context.Context, dest *models.S3Destination) error
	ListS3Destinations(ctx context.Context, projectID string) ([]*models.S3Destination, error)
	GetS3Destination(ctx context.Context, id string) (*models.S3Destination, error)
	DeleteS3Destination(ctx context.Context, id, projectID string) error
}

type S3DestinationSQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewS3DestinationSQLiteRepository(db *sql.DB) *S3DestinationSQLiteRepository {
	return &S3DestinationSQLiteRepository{db: db}
}

func (r *S3DestinationSQLiteRepository) CreateS3Destination(_ context.Context, dest *models.S3Destination) error {
	if dest.ID == "" {
		dest.ID = uuid.New().String()
	}
	if dest.CreatedAt == "" {
		dest.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`INSERT INTO s3_destinations (id, project_id, name, endpoint, bucket, region, access_key_id, secret_access_key, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		dest.ID, dest.ProjectID, dest.Name, dest.Endpoint, dest.Bucket, dest.Region, dest.AccessKeyID, dest.SecretAccessKey, dest.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create s3 destination: %w", err)
	}
	return nil
}

func (r *S3DestinationSQLiteRepository) ListS3Destinations(_ context.Context, projectID string) ([]*models.S3Destination, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.Query(`SELECT id, project_id, name, endpoint, bucket, region, access_key_id, secret_access_key, created_at
		FROM s3_destinations WHERE project_id = ? ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list s3 destinations: %w", err)
	}
	defer rows.Close()
	var list []*models.S3Destination
	for rows.Next() {
		var dest models.S3Destination
		if err := rows.Scan(&dest.ID, &dest.ProjectID, &dest.Name, &dest.Endpoint, &dest.Bucket, &dest.Region, &dest.AccessKeyID, &dest.SecretAccessKey, &dest.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &dest)
	}
	return list, nil
}

func (r *S3DestinationSQLiteRepository) GetS3Destination(_ context.Context, id string) (*models.S3Destination, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	row := r.db.QueryRow(`SELECT id, project_id, name, endpoint, bucket, region, access_key_id, secret_access_key, created_at
		FROM s3_destinations WHERE id = ?`, id)
	var dest models.S3Destination
	err := row.Scan(&dest.ID, &dest.ProjectID, &dest.Name, &dest.Endpoint, &dest.Bucket, &dest.Region, &dest.AccessKeyID, &dest.SecretAccessKey, &dest.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("S3Destination", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get s3 destination %s: %w", id, err)
	}
	return &dest, nil
}

func (r *S3DestinationSQLiteRepository) DeleteS3Destination(_ context.Context, id, projectID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.Exec("DELETE FROM s3_destinations WHERE id = ? AND project_id = ?", id, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete s3 destination: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("s3 destination not found or unauthorized")
	}
	return nil
}
