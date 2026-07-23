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

type ScheduledTaskRepository interface {
	Create(ctx context.Context, j *models.ScheduledTask) error
	GetByID(ctx context.Context, id string) (*models.ScheduledTask, error)
	ListAll(ctx context.Context) ([]models.ScheduledTask, error)
	ListByProject(ctx context.Context, projectID string) ([]models.ScheduledTask, error)
	ListByService(ctx context.Context, serviceID string) ([]models.ScheduledTask, error)
	Update(ctx context.Context, j *models.ScheduledTask) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status models.ScheduledTaskStatus, lastRunAt *time.Time, output string) error
}

type ScheduledTaskRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewScheduledTaskRepo(db *sql.DB) *ScheduledTaskRepo {
	return &ScheduledTaskRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *ScheduledTaskRepo) Create(_ context.Context, j *models.ScheduledTask) error {
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
	_, err := r.db.Exec(`INSERT INTO scheduled_tasks (
		id, service_id, name, schedule, command, status, last_run_at, last_output, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		j.ID, j.ServiceID, j.Name, j.Schedule, j.Command, j.Status, j.LastRunAt, j.LastOutput, j.CreatedAt, j.UpdatedAt)
	return err
}

func (r *ScheduledTaskRepo) GetByID(_ context.Context, id string) (*models.ScheduledTask, error) {
	var j models.ScheduledTask
	err := r.db.Get(&j, `SELECT id, service_id, name, schedule, command, status, last_run_at, COALESCE(last_output, '') AS last_output, created_at, updated_at
		FROM scheduled_tasks WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Entity", id)
	}
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *ScheduledTaskRepo) ListAll(_ context.Context) ([]models.ScheduledTask, error) {
	var scheduledTasks []models.ScheduledTask
	err := r.db.Select(&scheduledTasks, `SELECT id, service_id, name, schedule, command, status, last_run_at, COALESCE(last_output, '') AS last_output, created_at, updated_at FROM scheduled_tasks ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	if scheduledTasks == nil {
		scheduledTasks = make([]models.ScheduledTask, 0)
	}
	return scheduledTasks, nil
}

func (r *ScheduledTaskRepo) ListByProject(_ context.Context, projectID string) ([]models.ScheduledTask, error) {
	var scheduledTasks []models.ScheduledTask
	var err error
	if projectID == "" {
		err = r.db.Select(&scheduledTasks, `SELECT id, service_id, name, schedule, command, status, last_run_at, COALESCE(last_output, '') AS last_output, created_at, updated_at FROM scheduled_tasks ORDER BY created_at ASC`)
	} else {
		err = r.db.Select(&scheduledTasks, `SELECT st.id, st.service_id, st.name, st.schedule, st.command, st.status, st.last_run_at, COALESCE(st.last_output, '') AS last_output, st.created_at, st.updated_at
		FROM scheduled_tasks st JOIN app_services a ON st.service_id = a.id WHERE a.project_id = ? ORDER BY st.created_at ASC`, projectID)
	}
	if err != nil {
		return nil, err
	}
	if scheduledTasks == nil {
		scheduledTasks = make([]models.ScheduledTask, 0)
	}
	return scheduledTasks, nil
}

func (r *ScheduledTaskRepo) ListByService(_ context.Context, serviceID string) ([]models.ScheduledTask, error) {
	var scheduledTasks []models.ScheduledTask
	err := r.db.Select(&scheduledTasks, `SELECT id, service_id, name, schedule, command, status, last_run_at, COALESCE(last_output, '') AS last_output, created_at, updated_at
		FROM scheduled_tasks WHERE service_id = ? ORDER BY created_at ASC`, serviceID)
	if err != nil {
		return nil, err
	}
	if scheduledTasks == nil {
		scheduledTasks = make([]models.ScheduledTask, 0)
	}
	return scheduledTasks, nil
}

func (r *ScheduledTaskRepo) Update(_ context.Context, j *models.ScheduledTask) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	j.UpdatedAt = time.Now()
	_, err := r.db.Exec(`UPDATE scheduled_tasks SET
		service_id = ?, name = ?, schedule = ?, command = ?, status = ?, last_run_at = ?, last_output = ?, updated_at = ?
		WHERE id = ?`,
		j.ServiceID, j.Name, j.Schedule, j.Command, j.Status, j.LastRunAt, j.LastOutput, j.UpdatedAt, j.ID)
	return err
}

func (r *ScheduledTaskRepo) Delete(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM scheduled_tasks WHERE id = ?`, id)
	return err
}

func (r *ScheduledTaskRepo) UpdateStatus(_ context.Context, id string, status models.ScheduledTaskStatus, lastRunAt *time.Time, output string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	_, err := r.db.Exec(`UPDATE scheduled_tasks SET status = ?, last_run_at = ?, last_output = ?, updated_at = ? WHERE id = ?`,
		status, lastRunAt, output, now, id)
	return err
}
