package repositories

import (
	"context"
	"database/sql"
	"errors"

	"vessl.dev/vessl/internal/models"
)

type PRPreviewRepository interface {
	Create(ctx context.Context, p *models.PRPreview) error
	GetByAppAndPR(ctx context.Context, appID string, prNumber int) ([]*models.PRPreview, error)
	Update(ctx context.Context, p *models.PRPreview) error
	Delete(ctx context.Context, id string) error
}

type prPreviewRepo struct {
	db *sql.DB
}

func NewPRPreviewRepository(db *sql.DB) PRPreviewRepository {
	return &prPreviewRepo{db: db}
}

func (r *prPreviewRepo) Create(ctx context.Context, p *models.PRPreview) error {
	q := `INSERT INTO pr_previews (id, service_id, project_id, pr_number, branch, commit_hash, status, preview_domain, container_id, created_at, updated_at)
	      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, q, p.ID, p.ServiceID, p.ProjectID, p.PRNumber, p.Branch, p.CommitHash, p.Status, p.PreviewDomain, p.ContainerID, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *prPreviewRepo) GetByAppAndPR(ctx context.Context, appID string, prNumber int) ([]*models.PRPreview, error) {
	q := `SELECT id, service_id, project_id, pr_number, branch, commit_hash, status, preview_domain, container_id, created_at, updated_at
	      FROM pr_previews WHERE service_id = ? AND pr_number = ?`
	rows, err := r.db.QueryContext(ctx, q, appID, prNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var previews []*models.PRPreview
	for rows.Next() {
		p := &models.PRPreview{}
		if err := rows.Scan(&p.ID, &p.ServiceID, &p.ProjectID, &p.PRNumber, &p.Branch, &p.CommitHash, &p.Status, &p.PreviewDomain, &p.ContainerID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		previews = append(previews, p)
	}
	return previews, nil
}

func (r *prPreviewRepo) Update(ctx context.Context, p *models.PRPreview) error {
	q := `UPDATE pr_previews SET status = ?, preview_domain = ?, container_id = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, p.Status, p.PreviewDomain, p.ContainerID, p.UpdatedAt, p.ID)
	return err
}

func (r *prPreviewRepo) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	q := `DELETE FROM pr_previews WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}
