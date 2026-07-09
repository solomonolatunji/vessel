package project

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/environment"
)

// SQLiteRepository implements Repository against a SQLite database.
type SQLiteRepository struct {
	db           *sql.DB
	environments environment.Repository
}

// NewSQLiteRepository constructs a SQLiteRepository backed by the given db and environment repository.
func NewSQLiteRepository(db *sql.DB, envRepo environment.Repository) *SQLiteRepository {
	return &SQLiteRepository{db: db, environments: envRepo}
}

// List retrieves all ProjectConfig records ordered by creation date descending.
func (r *SQLiteRepository) List(_ context.Context) ([]ProjectConfig, error) {
	rows, err := r.db.Query(`SELECT id, COALESCE(workspace_id, ''), COALESCE(team_id,''), name, COALESCE(description,''), created_at, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []ProjectConfig
	for rows.Next() {
		var p ProjectConfig
		if err := rows.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// Get retrieves a single ProjectConfig by its ID.
func (r *SQLiteRepository) Get(_ context.Context, id string) (*ProjectConfig, error) {
	row := r.db.QueryRow(`SELECT id, COALESCE(workspace_id, ''), COALESCE(team_id,''), name, COALESCE(description,''), created_at, updated_at FROM projects WHERE id = ?`, id)
	var p ProjectConfig
	err := row.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Create inserts a new project and creates its default production environment.
func (r *SQLiteRepository) Create(ctx context.Context, p *ProjectConfig) error {
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

	defaultEnv := &environment.Config{
		ProjectID: p.ID,
		Name:      "production",
		IsDefault: true,
	}
	return r.environments.Create(ctx, defaultEnv)
}

// Delete removes a project record by ID.
func (r *SQLiteRepository) Delete(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}
