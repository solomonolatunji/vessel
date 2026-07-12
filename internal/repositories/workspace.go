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
	CreateAuditLog(ctx context.Context, log *models.AuditLog) error
	ListAuditLogs(ctx context.Context, workspaceID string, limit, offset int) ([]*models.AuditLog, int, error)
	ListWorkspacesByUser(ctx context.Context, userID string) ([]*models.Workspace, error)
	AddMember(ctx context.Context, member *models.WorkspaceMember) error
	RemoveMember(ctx context.Context, workspaceID, userID string) error
	GetMember(ctx context.Context, workspaceID, userID string) (*models.WorkspaceMember, error)
	ListMembers(ctx context.Context, workspaceID string) ([]*models.WorkspaceMember, error)
	CreateInvite(ctx context.Context, invite *models.WorkspaceInvite) error
	GetInviteByToken(ctx context.Context, token string) (*models.WorkspaceInvite, error)
	DeleteInvite(ctx context.Context, id string) error
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

func (r *WorkspaceSQLiteRepository) List(ctx context.Context, ownerID string, limit, offset int) ([]*models.Workspace, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM workspaces WHERE owner_id = ?`, ownerID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count workspaces: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM workspaces WHERE owner_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, ownerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list workspaces: %w", err)
	}
	defer rows.Close()
	var list []*models.Workspace
	for rows.Next() {
		var ws models.Workspace
		var createdStr, updatedStr string
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.AvatarURL, &ws.PreferredRegion, &ws.OwnerID, &createdStr, &updatedStr); err != nil {
			return nil, 0, err
		}
		ws.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		ws.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		list = append(list, &ws)
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

func (r *WorkspaceSQLiteRepository) CreateWorkspace(ctx context.Context, workspace *models.Workspace) error {
	if workspace.ID == "" {
		workspace.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	if workspace.CreatedAt.IsZero() {
		workspace.CreatedAt = now
	}
	workspace.UpdatedAt = workspace.CreatedAt
	r.mu.Lock()
	defer r.mu.Unlock()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO workspaces (id, name, avatar_url, preferred_region, owner_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		workspace.ID, workspace.Name, workspace.AvatarURL, workspace.PreferredRegion, workspace.OwnerID, workspace.CreatedAt.Format(time.RFC3339), workspace.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("insert workspace: %w", err)
	}
	ownerMember := &models.WorkspaceMember{
		ID:          uuid.NewString(),
		WorkspaceID: workspace.ID,
		UserID:      workspace.OwnerID,
		UserEmail:   "",
		Role:        "Owner",
		JoinedAt:    now,
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO workspace_members (id, workspace_id, user_id, user_email, role, joined_at) VALUES (?, ?, ?, ?, ?, ?)`,
		ownerMember.ID, ownerMember.WorkspaceID, ownerMember.UserID, ownerMember.UserEmail, ownerMember.Role, ownerMember.JoinedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("insert workspace owner: %w", err)
	}
	return tx.Commit()
}

func (r *WorkspaceSQLiteRepository) GetWorkspaceByID(ctx context.Context, id string) (*models.Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var t models.Workspace
	var createdStr, updatedStr string
	err := r.db.QueryRowContext(ctx, `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM workspaces WHERE id = ?`, id).
		Scan(&t.ID, &t.Name, &t.AvatarURL, &t.PreferredRegion, &t.OwnerID, &createdStr, &updatedStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Workspace", id)
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace %s: %w", id, err)
	}
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &t, nil
}

func (r *WorkspaceSQLiteRepository) ListWorkspacesByUser(ctx context.Context, userID string) ([]*models.Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := `SELECT t.id, t.name, t.avatar_url, t.preferred_region, t.owner_id, t.created_at, t.updated_at
	          FROM workspaces t
	          JOIN workspace_members m ON t.id = m.workspace_id
	          WHERE m.user_id = ? ORDER BY t.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list workspaces by user: %w", err)
	}
	defer rows.Close()
	var list []*models.Workspace
	for rows.Next() {
		var t models.Workspace
		var createdStr, updatedStr string
		if err := rows.Scan(&t.ID, &t.Name, &t.AvatarURL, &t.PreferredRegion, &t.OwnerID, &createdStr, &updatedStr); err != nil {
			return nil, err
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		list = append(list, &t)
	}
	return list, nil
}

func (r *WorkspaceSQLiteRepository) UpdateWorkspace(ctx context.Context, workspace *models.Workspace) error {
	workspace.UpdatedAt = time.Now().UTC()
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE workspaces SET name = ?, avatar_url = ?, preferred_region = ?, updated_at = ? WHERE id = ?`,
		workspace.Name, workspace.AvatarURL, workspace.PreferredRegion, workspace.UpdatedAt.Format(time.RFC3339), workspace.ID)
	return err
}

func (r *WorkspaceSQLiteRepository) DeleteWorkspace(ctx context.Context, id, ownerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.ExecContext(ctx, `DELETE FROM workspaces WHERE id = ? AND owner_id = ?`, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete workspace: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("workspace not found or unauthorized (must be owner)")
	}
	_, _ = r.db.ExecContext(ctx, `DELETE FROM workspace_members WHERE workspace_id = ?`, id)
	_, _ = r.db.ExecContext(ctx, `DELETE FROM workspace_invites WHERE workspace_id = ?`, id)
	return nil
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

func (r *WorkspaceSQLiteRepository) GetMember(ctx context.Context, workspaceID, userID string) (*models.WorkspaceMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var m models.WorkspaceMember
	var joinedStr string
	err := r.db.QueryRowContext(ctx, `SELECT id, workspace_id, user_id, user_email, role, joined_at FROM workspace_members WHERE workspace_id = ? AND user_id = ?`, workspaceID, userID).
		Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.UserEmail, &m.Role, &joinedStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace member: %w", err)
	}
	m.JoinedAt, _ = time.Parse(time.RFC3339, joinedStr)
	return &m, nil
}

func (r *WorkspaceSQLiteRepository) ListMembers(ctx context.Context, workspaceID string) ([]*models.WorkspaceMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx, `SELECT id, workspace_id, user_id, user_email, role, joined_at FROM workspace_members WHERE workspace_id = ? ORDER BY joined_at ASC`, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("list workspace members: %w", err)
	}
	defer rows.Close()
	var list []*models.WorkspaceMember
	for rows.Next() {
		var m models.WorkspaceMember
		var joinedStr string
		if err := rows.Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.UserEmail, &m.Role, &joinedStr); err != nil {
			return nil, err
		}
		m.JoinedAt, _ = time.Parse(time.RFC3339, joinedStr)
		list = append(list, &m)
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
	var expStr, createdStr string
	err := r.db.QueryRowContext(ctx, `SELECT id, workspace_id, email, role, token, invited_by, expires_at, created_at FROM workspace_invites WHERE token = ?`, token).
		Scan(&inv.ID, &inv.WorkspaceID, &inv.Email, &inv.Role, &inv.Token, &inv.InvitedBy, &expStr, &createdStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Invite", token)
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace invite: %w", err)
	}
	inv.ExpiresAt, _ = time.Parse(time.RFC3339, expStr)
	inv.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	return &inv, nil
}

func (r *WorkspaceSQLiteRepository) DeleteInvite(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `DELETE FROM workspace_invites WHERE id = ?`, id)
	return err
}
