package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

type EnvironmentRepository interface {
	Get(ctx context.Context, id string) (*models.EnvironmentConfig, error)
	ListByProject(ctx context.Context, projectID string) ([]models.EnvironmentConfig, error)
	Create(ctx context.Context, env *models.EnvironmentConfig) error
	Delete(ctx context.Context, id string) error
}

type DomainRepository interface {
	ListByProject(ctx context.Context, projectID string) ([]models.DomainConfig, error)
	ListAll(ctx context.Context) ([]models.DomainConfig, error)
	Create(ctx context.Context, d *models.DomainConfig) error
	Delete(ctx context.Context, id string) error
}

type EnvironmentSQLiteRepository struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewEnvironmentSQLiteRepository(db *sql.DB) *EnvironmentSQLiteRepository {
	return &EnvironmentSQLiteRepository{db: db}
}

func (r *EnvironmentSQLiteRepository) Get(_ context.Context, id string) (*models.EnvironmentConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	row := r.db.QueryRow(
		`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE id = ?`, id,
	)
	var env models.EnvironmentConfig
	var isDefault int
	err := row.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Environment", id)
	}
	if err != nil {
		return nil, err
	}
	env.IsDefault = isDefault == 1
	return &env, nil
}

func (r *EnvironmentSQLiteRepository) ListByProject(_ context.Context, projectID string) ([]models.EnvironmentConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rows, err := r.db.Query(
		`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE project_id = ? ORDER BY is_default DESC, created_at ASC`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}
	defer rows.Close()
	var envs []models.EnvironmentConfig
	for rows.Next() {
		var env models.EnvironmentConfig
		var isDefault int
		if err := rows.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan environment row: %w", err)
		}
		env.IsDefault = isDefault == 1
		envs = append(envs, env)
	}
	return envs, rows.Err()
}

func (r *EnvironmentSQLiteRepository) Create(_ context.Context, env *models.EnvironmentConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if env.ID == "" {
		env.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	env.CreatedAt = now
	env.UpdatedAt = now
	_, err := r.db.Exec(
		`INSERT INTO environments (id, project_id, name, is_default, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		env.ID, env.ProjectID, env.Name, env.IsDefault, env.CreatedAt, env.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}
	return nil
}

func (r *EnvironmentSQLiteRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM environments WHERE id = ?`, id)
	return err
}

type DomainSQLiteRepository struct {
	db *sql.DB
}

func NewDomainSQLiteRepository(db *sql.DB) *DomainSQLiteRepository {
	return &DomainSQLiteRepository{db: db}
}

func (r *DomainSQLiteRepository) ListByProject(_ context.Context, projectID string) ([]models.DomainConfig, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at FROM domains WHERE project_id = ? ORDER BY domain_name ASC`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var domains []models.DomainConfig
	for rows.Next() {
		var d models.DomainConfig
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.DomainName, &d.RedirectTo, &d.SSLCertStatus, &d.PathPrefix, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, rows.Err()
}

func (r *DomainSQLiteRepository) ListAll(ctx context.Context) ([]models.DomainConfig, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at FROM domains ORDER BY domain_name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var domains []models.DomainConfig
	for rows.Next() {
		var d models.DomainConfig
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.DomainName, &d.RedirectTo, &d.SSLCertStatus, &d.PathPrefix, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, rows.Err()
}

func (r *DomainSQLiteRepository) Create(_ context.Context, d *models.DomainConfig) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	_, err := r.db.Exec(
		`INSERT INTO domains (id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ProjectID, d.DomainName, d.RedirectTo, d.SSLCertStatus, d.PathPrefix, d.CreatedAt, d.UpdatedAt,
	)
	return err
}

func (r *DomainSQLiteRepository) Delete(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM domains WHERE id = ?`, id)
	return err
}
