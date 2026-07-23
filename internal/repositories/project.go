package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock.dev/codedock/internal/models"
)

type ProjectRepository interface {
	List(ctx context.Context, limit, offset int) ([]models.ProjectConfig, int, error)
	Get(ctx context.Context, id string) (*models.ProjectConfig, error)
	Create(ctx context.Context, p *models.ProjectConfig) error
	CreateWithMember(ctx context.Context, p *models.ProjectConfig, userID, role string) error
	Delete(ctx context.Context, id string) error
}

type EnvRepository interface {
	GetVars(ctx context.Context, projectID string) (map[string]string, error)
	SetVar(ctx context.Context, projectID, key, value string) error
}

type ProjectRepo struct {
	db           *sqlx.DB
	environments EnvironmentRepository
}

func NewProjectRepo(db *sql.DB, envRepo EnvironmentRepository) *ProjectRepo {
	return &ProjectRepo{db: sqlx.NewDb(db, "sqlite"), environments: envRepo}
}

func (r *ProjectRepo) List(_ context.Context, limit, offset int) ([]models.ProjectConfig, int, error) {
	var total int
	var err error
	var projects []models.ProjectConfig

	if err = r.db.Get(&total, `SELECT COUNT(*) FROM projects`); err != nil {
		return nil, 0, err
	}
	err = r.db.Select(&projects, `SELECT id, name, COALESCE(description,'') AS description, created_at, updated_at FROM projects ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)

	if err != nil {
		return nil, 0, err
	}
	if projects == nil {
		projects = make([]models.ProjectConfig, 0)
	}
	return projects, total, nil
}

func (r *ProjectRepo) Get(_ context.Context, id string) (*models.ProjectConfig, error) {
	var p models.ProjectConfig
	err := r.db.Get(&p, `SELECT id, name, COALESCE(description,'') AS description, created_at, updated_at FROM projects WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepo) Create(ctx context.Context, p *models.ProjectConfig) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.Exec(
		`INSERT INTO projects (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Description, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return err
	}
	defaultEnv := &models.EnvironmentConfig{
		ProjectID: p.ID,
		Name:      "production",
		IsDefault: true,
	}
	return r.environments.Create(ctx, defaultEnv)
}

func (r *ProjectRepo) CreateWithMember(ctx context.Context, p *models.ProjectConfig, userID, role string) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO projects (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Description, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO project_members (project_id, user_id, role, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		p.ID, userID, role, models.MemberStatusAccepted, now, now,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	defaultEnv := &models.EnvironmentConfig{
		ProjectID: p.ID,
		Name:      "production",
		IsDefault: true,
	}
	return r.environments.Create(ctx, defaultEnv)
}

func (r *ProjectRepo) Delete(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}

type EnvRepo struct {
	db    *sqlx.DB
	vault Vault
}

func NewEnvRepo(db *sql.DB, vault Vault) *EnvRepo {
	return &EnvRepo{db: sqlx.NewDb(db, "sqlite"), vault: vault}
}

func (r *EnvRepo) GetVars(_ context.Context, projectID string) (map[string]string, error) {
	rows, err := r.db.Query(`SELECT key, encrypted_value FROM env_vars WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	envs := make(map[string]string)
	for rows.Next() {
		var key, encrypted string
		if err := rows.Scan(&key, &encrypted); err != nil {
			return nil, err
		}
		plaintext, err := r.vault.Decrypt(encrypted)
		if err != nil {
			continue
		}
		envs[key] = plaintext
	}
	return envs, rows.Err()
}

func (r *EnvRepo) SetVar(_ context.Context, projectID, key, plaintextValue string) error {
	encrypted, err := r.vault.Encrypt(plaintextValue)
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = r.db.Exec(
		`INSERT INTO env_vars (id, project_id, key, encrypted_value, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(project_id, key) DO UPDATE SET encrypted_value = excluded.encrypted_value, updated_at = excluded.updated_at`,
		uuid.NewString(), projectID, key, encrypted, now, now,
	)
	return err
}
