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
)

type WorkspaceRepository interface {
	Create(ctx context.Context, ws *models.Workspace) error
	Get(ctx context.Context, id string) (*models.Workspace, error)
	List(ctx context.Context, ownerID string) ([]*models.Workspace, error)
	Update(ctx context.Context, ws *models.Workspace) error
	Delete(ctx context.Context, id, ownerID string) error
	CreateTrustedDomain(ctx context.Context, d *models.TrustedDomain) error
	ListTrustedDomains(ctx context.Context, teamID string) ([]*models.TrustedDomain, error)
	DeleteTrustedDomain(ctx context.Context, id string) error
	CreateSSHKey(ctx context.Context, key *models.SSHKey) error
	ListSSHKeys(ctx context.Context, teamID string) ([]*models.SSHKey, error)
	DeleteSSHKey(ctx context.Context, id string) error
	CreateAuditLog(ctx context.Context, log *models.AuditLog) error
	ListAuditLogs(ctx context.Context, teamID string, limit int) ([]*models.AuditLog, error)
}

type WorkspaceSQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewWorkspaceSQLiteRepository(db *sql.DB) *WorkspaceSQLiteRepository {
	return &WorkspaceSQLiteRepository{db: db}
}

func (r *WorkspaceSQLiteRepository) Migrate(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS workspaces (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			avatar_url TEXT NOT NULL DEFAULT '',
			preferred_region TEXT NOT NULL DEFAULT 'local',
			owner_id TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS workspace_trusted_domains (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			domain TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'developer',
			created_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS workspace_ssh_keys (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			name TEXT NOT NULL,
			public_key TEXT NOT NULL,
			created_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS workspace_audit_logs (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			project_id TEXT DEFAULT '',
			environment_id TEXT DEFAULT '',
			action TEXT NOT NULL,
			actor TEXT NOT NULL,
			created_at TEXT NOT NULL
		);
	`)
	return err
}

func (r *WorkspaceSQLiteRepository) Create(ctx context.Context, ws *models.Workspace) error {
	if ws.ID == "" {
		ws.ID = uuid.NewString()
	}
	if ws.CreatedAt.IsZero() {
		ws.CreatedAt = time.Now().UTC()
	}
	if ws.UpdatedAt.IsZero() {
		ws.UpdatedAt = time.Now().UTC()
	}
	if ws.PreferredRegion == "" {
		ws.PreferredRegion = "local"
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO workspaces (id, name, avatar_url, preferred_region, owner_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		ws.ID, ws.Name, ws.AvatarURL, ws.PreferredRegion, ws.OwnerID, ws.CreatedAt.Format(time.RFC3339), ws.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create workspace: %w", err)
	}
	return nil
}

func (r *WorkspaceSQLiteRepository) Get(ctx context.Context, id string) (*models.Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var ws models.Workspace
	var createdStr, updatedStr string
	err := r.db.QueryRowContext(ctx, `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM workspaces WHERE id = ?`, id).
		Scan(&ws.ID, &ws.Name, &ws.AvatarURL, &ws.PreferredRegion, &ws.OwnerID, &createdStr, &updatedStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace: %w", err)
	}
	ws.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	ws.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &ws, nil
}

func (r *WorkspaceSQLiteRepository) List(ctx context.Context, ownerID string) ([]*models.Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM workspaces WHERE owner_id = ? ORDER BY created_at DESC`, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	defer rows.Close()
	var list []*models.Workspace
	for rows.Next() {
		var ws models.Workspace
		var createdStr, updatedStr string
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.AvatarURL, &ws.PreferredRegion, &ws.OwnerID, &createdStr, &updatedStr); err != nil {
			return nil, err
		}
		ws.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		ws.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		list = append(list, &ws)
	}
	return list, nil
}

func (r *WorkspaceSQLiteRepository) Update(ctx context.Context, ws *models.Workspace) error {
	ws.UpdatedAt = time.Now().UTC()
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE workspaces SET name = ?, avatar_url = ?, preferred_region = ?, updated_at = ? WHERE id = ?`,
		ws.Name, ws.AvatarURL, ws.PreferredRegion, ws.UpdatedAt.Format(time.RFC3339), ws.ID)
	if err != nil {
		return fmt.Errorf("update workspace: %w", err)
	}
	return nil
}

func (r *WorkspaceSQLiteRepository) Delete(ctx context.Context, id, ownerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var count int
	_ = r.db.QueryRowContext(ctx, "SELECT count(*) FROM workspaces WHERE owner_id = ?", ownerID).Scan(&count)
	if count <= 1 {
		return errors.New("cannot delete your last workspace. To delete your account, visit Account Settings")
	}
	_, err := r.db.ExecContext(ctx, "DELETE FROM workspaces WHERE id = ? AND owner_id = ?", id, ownerID)
	return err
}
