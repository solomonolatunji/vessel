package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"
	"vessl.dev/vessl/internal/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"vessl.dev/vessl/internal/models"
)

type JobRepository interface {
	Create(ctx context.Context, j *models.Job) error
	GetByID(ctx context.Context, id string) (*models.Job, error)
	ListAll(ctx context.Context) ([]models.Job, error)
	ListByProject(ctx context.Context, projectID string) ([]models.Job, error)
	Update(ctx context.Context, j *models.Job) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id, status string, lastRunAt *time.Time, output string) error
}

type JobSQLiteRepository struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewJobSQLiteRepository(db *sql.DB) *JobSQLiteRepository {
	return &JobSQLiteRepository{db: sqlx.NewDb(db, "sqlite")}
}

func (r *JobSQLiteRepository) Create(_ context.Context, j *models.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if j.ID == "" {
		j.ID = uuid.NewString()
	}
	now := time.Now()
	j.CreatedAt = now
	j.UpdatedAt = now
	if j.Status == "" {
		j.Status = "active"
	}
	_, err := r.db.Exec(`INSERT INTO jobs (
		id, project_id, name, schedule, command, status, last_run_at, last_output, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		j.ID, j.ProjectID, j.Name, j.Schedule, j.Command, j.Status, j.LastRunAt, j.LastOutput, j.CreatedAt, j.UpdatedAt)
	return err
}

func (r *JobSQLiteRepository) GetByID(_ context.Context, id string) (*models.Job, error) {
	var j models.Job
	err := r.db.Get(&j, `SELECT id, project_id, name, schedule, command, status, last_run_at, COALESCE(last_output, '') AS last_output, created_at, updated_at
		FROM jobs WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Entity", id)
	}
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *JobSQLiteRepository) ListAll(_ context.Context) ([]models.Job, error) {
	var jobs []models.Job
	err := r.db.Select(&jobs, `SELECT id, project_id, name, schedule, command, status, last_run_at, COALESCE(last_output, '') AS last_output, created_at, updated_at FROM jobs ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	if jobs == nil {
		jobs = make([]models.Job, 0)
	}
	return jobs, nil
}

func (r *JobSQLiteRepository) ListByProject(_ context.Context, projectID string) ([]models.Job, error) {
	var jobs []models.Job
	err := r.db.Select(&jobs, `SELECT id, project_id, name, schedule, command, status, last_run_at, COALESCE(last_output, '') AS last_output, created_at, updated_at
		FROM jobs WHERE project_id = ? ORDER BY created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	if jobs == nil {
		jobs = make([]models.Job, 0)
	}
	return jobs, nil
}

func (r *JobSQLiteRepository) Update(_ context.Context, j *models.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	j.UpdatedAt = time.Now()
	_, err := r.db.Exec(`UPDATE jobs SET name = ?, schedule = ?, command = ?, status = ?, last_run_at = ?, last_output = ?, updated_at = ? WHERE id = ?`,
		j.Name, j.Schedule, j.Command, j.Status, j.LastRunAt, j.LastOutput, j.UpdatedAt, j.ID)
	return err
}

func (r *JobSQLiteRepository) Delete(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM jobs WHERE id = ?`, id)
	return err
}

func (r *JobSQLiteRepository) UpdateStatus(_ context.Context, id, status string, lastRunAt *time.Time, output string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	_, err := r.db.Exec(`UPDATE jobs SET status = ?, last_run_at = ?, last_output = ?, updated_at = ? WHERE id = ?`,
		status, lastRunAt, output, now, id)
	return err
}
