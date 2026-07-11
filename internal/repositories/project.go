package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
)

type ProjectRepository interface {
	List(ctx context.Context) ([]models.ProjectConfig, error)
	Get(ctx context.Context, id string) (*models.ProjectConfig, error)
	Create(ctx context.Context, p *models.ProjectConfig) error
	Delete(ctx context.Context, id string) error
}

type EnvRepository interface {
	GetVars(ctx context.Context, projectID string) (map[string]string, error)
	SetVar(ctx context.Context, projectID, key, value string) error
}

type ProjectSQLiteRepository struct {
	db           *sql.DB
	environments EnvironmentRepository
}

func NewProjectSQLiteRepository(db *sql.DB, envRepo EnvironmentRepository) *ProjectSQLiteRepository {
	return &ProjectSQLiteRepository{db: db, environments: envRepo}
}

func (r *ProjectSQLiteRepository) List(_ context.Context) ([]models.ProjectConfig, error) {
	rows, err := r.db.Query(`SELECT id, COALESCE(workspace_id, ''), COALESCE(team_id,''), name, COALESCE(description,''), created_at, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects []models.ProjectConfig
	for rows.Next() {
		var p models.ProjectConfig
		if err := rows.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *ProjectSQLiteRepository) Get(_ context.Context, id string) (*models.ProjectConfig, error) {
	row := r.db.QueryRow(`SELECT id, COALESCE(workspace_id, ''), COALESCE(team_id,''), name, COALESCE(description,''), created_at, updated_at FROM projects WHERE id = ?`, id)
	var p models.ProjectConfig
	err := row.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectSQLiteRepository) Create(ctx context.Context, p *models.ProjectConfig) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.Exec(
		`INSERT INTO projects (id, workspace_id, team_id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.WorkspaceID, p.TeamID, p.Name, p.Description, p.CreatedAt, p.UpdatedAt,
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

func (r *ProjectSQLiteRepository) Delete(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}

type EnvSQLiteRepository struct {
	db    *sql.DB
	vault Vault
}

func NewEnvSQLiteRepository(db *sql.DB, vault Vault) *EnvSQLiteRepository {
	return &EnvSQLiteRepository{db: db, vault: vault}
}

func (r *EnvSQLiteRepository) GetVars(_ context.Context, projectID string) (map[string]string, error) {
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

func (r *EnvSQLiteRepository) SetVar(_ context.Context, projectID, key, plaintextValue string) error {
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
