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
	db *sql.DB
	mu sync.Mutex
}

func NewJobSQLiteRepository(db *sql.DB) *JobSQLiteRepository {
	return &JobSQLiteRepository{db: db}
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
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		j.ID, j.ProjectID, j.Name, j.Schedule, j.Command, j.Status, j.LastRunAt, j.LastOutput, j.CreatedAt, j.UpdatedAt)
	return err
}

func (r *JobSQLiteRepository) GetByID(_ context.Context, id string) (*models.Job, error) {
	var j models.Job
	var lastRunAt sql.NullTime
	err := r.db.QueryRow(`SELECT id, project_id, name, schedule, command, status, last_run_at, COALESCE(last_output, ''), created_at, updated_at
		FROM jobs WHERE id = ?`, id).Scan(
		&j.ID, &j.ProjectID, &j.Name, &j.Schedule, &j.Command, &j.Status, &lastRunAt, &j.LastOutput, &j.CreatedAt, &j.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Entity", id)
	}
	if err != nil {
		return nil, err
	}
	if lastRunAt.Valid {
		j.LastRunAt = &lastRunAt.Time
	}
	return &j, nil
}

func (r *JobSQLiteRepository) ListAll(_ context.Context) ([]models.Job, error) {
	rows, err := r.db.Query(`SELECT id, project_id, name, schedule, command, status, last_run_at, COALESCE(last_output, ''), created_at, updated_at FROM jobs ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []models.Job
	for rows.Next() {
		var j models.Job
		var lastRunAt sql.NullTime
		if err := rows.Scan(&j.ID, &j.ProjectID, &j.Name, &j.Schedule, &j.Command, &j.Status, &lastRunAt, &j.LastOutput, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		if lastRunAt.Valid {
			j.LastRunAt = &lastRunAt.Time
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

func (r *JobSQLiteRepository) ListByProject(_ context.Context, projectID string) ([]models.Job, error) {
	rows, err := r.db.Query(`SELECT id, project_id, name, schedule, command, status, last_run_at, COALESCE(last_output, ''), created_at, updated_at
		FROM jobs WHERE project_id = ? ORDER BY created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []models.Job
	for rows.Next() {
		var j models.Job
		var lastRunAt sql.NullTime
		if err := rows.Scan(&j.ID, &j.ProjectID, &j.Name, &j.Schedule, &j.Command, &j.Status, &lastRunAt, &j.LastOutput, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		if lastRunAt.Valid {
			j.LastRunAt = &lastRunAt.Time
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
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
