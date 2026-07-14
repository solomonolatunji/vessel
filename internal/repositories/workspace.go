package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

type WorkspaceRepository interface {
	Create(ctx context.Context, ws *models.Workspace) error
	Get(ctx context.Context, id string) (*models.Workspace, error)
	List(ctx context.Context, ownerID string, limit, offset int) ([]*models.Workspace, int, error)
	Update(ctx context.Context, ws *models.Workspace) error
	Delete(ctx context.Context, id, ownerID string) error
	CreateTrustedDomain(ctx context.Context, d *models.TrustedDomain) error
	ListTrustedDomains(ctx context.Context, workspaceID string) ([]*models.TrustedDomain, error)
	DeleteTrustedDomain(ctx context.Context, id string) error
	CreateSSHKey(ctx context.Context, key *models.SSHKey) error
	ListSSHKeys(ctx context.Context, workspaceID string) ([]*models.SSHKey, error)
	DeleteSSHKey(ctx context.Context, id string) error
	ListAuditLogs(ctx context.Context, workspaceID string, limit, offset int) ([]*models.AuditLog, int, error)
	ListWorkspacesByUser(ctx context.Context, userID string) ([]*models.Workspace, error)
	AddMember(ctx context.Context, member *models.WorkspaceMember) error
	RemoveMember(ctx context.Context, workspaceID, userID string) error
	ListMembers(ctx context.Context, workspaceID string) ([]*models.WorkspaceMember, error)
	CreateInvite(ctx context.Context, invite *models.WorkspaceInvite) error
	GetInviteByToken(ctx context.Context, token string) (*models.WorkspaceInvite, error)
	DeleteInvite(ctx context.Context, id string) error
}

type WorkspaceSQLiteRepository struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewWorkspaceSQLiteRepository(db *sql.DB) *WorkspaceSQLiteRepository {
	return &WorkspaceSQLiteRepository{db: sqlx.NewDb(db, "sqlite")}
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
			workspace_id TEXT NOT NULL,
			domain TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'developer',
			created_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS workspace_ssh_keys (
			id TEXT PRIMARY KEY,
			workspace_id TEXT NOT NULL,
			name TEXT NOT NULL,
			public_key TEXT NOT NULL,
			created_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS workspace_audit_logs (
			id TEXT PRIMARY KEY,
			workspace_id TEXT NOT NULL,
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
	err := r.db.GetContext(ctx, &ws, `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM workspaces WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace: %w", err)
	}
	return &ws, nil
}

func (r *WorkspaceSQLiteRepository) List(ctx context.Context, ownerID string, limit, offset int) ([]*models.Workspace, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var total int
	err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM workspaces WHERE owner_id = ?`, ownerID)
	if err != nil {
		return nil, 0, fmt.Errorf("count workspaces: %w", err)
	}

	var list []*models.Workspace
	err = r.db.SelectContext(ctx, &list, `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM workspaces WHERE owner_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, ownerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list workspaces: %w", err)
	}
	if list == nil {
		list = make([]*models.Workspace, 0)
	}
	return list, total, nil
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

func (r *WorkspaceSQLiteRepository) ListWorkspacesByUser(ctx context.Context, userID string) ([]*models.Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := `SELECT t.id, t.name, t.avatar_url, t.preferred_region, t.owner_id, t.created_at, t.updated_at
	          FROM workspaces t
	          JOIN workspace_members m ON t.id = m.workspace_id
	          WHERE m.user_id = ? ORDER BY t.created_at DESC`
	var list []*models.Workspace
	err := r.db.SelectContext(ctx, &list, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list workspaces by user: %w", err)
	}
	if list == nil {
		list = make([]*models.Workspace, 0)
	}
	return list, nil
}

func (r *WorkspaceSQLiteRepository) AddMember(ctx context.Context, member *models.WorkspaceMember) error {
	if member.ID == "" {
		member.ID = uuid.NewString()
	}
	if member.JoinedAt.IsZero() {
		member.JoinedAt = time.Now().UTC()
	}
	if member.Role == "" {
		member.Role = "Member"
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO workspace_members (id, workspace_id, user_id, user_email, role, joined_at) VALUES (?, ?, ?, ?, ?, ?)`,
		member.ID, member.WorkspaceID, member.UserID, member.UserEmail, member.Role, member.JoinedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("add workspace member: %w", err)
	}
	return nil
}

func (r *WorkspaceSQLiteRepository) RemoveMember(ctx context.Context, workspaceID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.ExecContext(ctx, `DELETE FROM workspace_members WHERE workspace_id = ? AND user_id = ? AND role != 'Owner'`, workspaceID, userID)
	if err != nil {
		return fmt.Errorf("remove workspace member: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("member not found or cannot remove workspace Owner")
	}
	return nil
}

func (r *WorkspaceSQLiteRepository) ListMembers(ctx context.Context, workspaceID string) ([]*models.WorkspaceMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.WorkspaceMember
	err := r.db.SelectContext(ctx, &list, `SELECT id, workspace_id, user_id, user_email, role, joined_at FROM workspace_members WHERE workspace_id = ? ORDER BY joined_at ASC`, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("list workspace members: %w", err)
	}
	if list == nil {
		list = make([]*models.WorkspaceMember, 0)
	}
	return list, nil
}

func (r *WorkspaceSQLiteRepository) CreateInvite(ctx context.Context, invite *models.WorkspaceInvite) error {
	if invite.ID == "" {
		invite.ID = uuid.NewString()
	}
	if invite.Token == "" {
		invite.Token = uuid.NewString()
	}
	now := time.Now().UTC()
	if invite.CreatedAt.IsZero() {
		invite.CreatedAt = now
	}
	if invite.ExpiresAt.IsZero() {
		invite.ExpiresAt = now.Add(7 * 24 * time.Hour)
	}
	if invite.Role == "" {
		invite.Role = "Member"
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO workspace_invites (id, workspace_id, email, role, token, invited_by, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		invite.ID, invite.WorkspaceID, invite.Email, invite.Role, invite.Token, invite.InvitedBy, invite.ExpiresAt.Format(time.RFC3339), invite.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create workspace invite: %w", err)
	}
	return nil
}

func (r *WorkspaceSQLiteRepository) GetInviteByToken(ctx context.Context, token string) (*models.WorkspaceInvite, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var inv models.WorkspaceInvite
	err := r.db.GetContext(ctx, &inv, `SELECT id, workspace_id, email, role, token, invited_by, expires_at, created_at FROM workspace_invites WHERE token = ?`, token)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Invite", token)
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace invite: %w", err)
	}
	return &inv, nil
}

func (r *WorkspaceSQLiteRepository) DeleteInvite(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `DELETE FROM workspace_invites WHERE id = ?`, id)
	return err
}
