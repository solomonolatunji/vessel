package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"vessl.dev/vessl/internal/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"vessl.dev/vessl/internal/models"
)

type ProjectSettingsRepository interface {
	CreateToken(ctx context.Context, t *models.ProjectToken, fullToken string) error
	ListTokensByProject(ctx context.Context, projectID string) ([]*models.ProjectToken, error)
	DeleteToken(ctx context.Context, id, projectID string) error
	GetTokenByHash(ctx context.Context, tokenHash string) (*models.ProjectToken, error)
	UpdateTokenLastUsed(ctx context.Context, id string) error
	AddMember(ctx context.Context, m *models.ProjectMember) error
	GetMember(ctx context.Context, projectID, userID string) (*models.ProjectMember, error)
	ListMembers(ctx context.Context, projectID string) ([]*models.ProjectMember, error)
	RemoveMember(ctx context.Context, id, projectID string) error
	AcceptAllInvitesForUser(ctx context.Context, userID string) error
}

type ProjectSettingsRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewProjectSettingsRepo(db *sql.DB) *ProjectSettingsRepo {
	return &ProjectSettingsRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *ProjectSettingsRepo) deleteByIDAndProject(ctx context.Context, table, id, projectID, entityName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ? AND project_id = ?", table)
	res, err := r.db.ExecContext(ctx, query, id, projectID)
	if err != nil {
		return fmt.Errorf("delete %s: %w", entityName, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("%s not found", entityName)
	}
	return nil
}

func (r *ProjectSettingsRepo) CreateToken(ctx context.Context, t *models.ProjectToken, fullToken string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	scopesStr := strings.Join(t.Scopes, ",")
	ipStr := strings.Join(t.IPAllowlist, ",")
	var expiresAtVal interface{}
	if t.ExpiresAt != nil {
		expiresAtVal = t.ExpiresAt.Format(time.RFC3339)
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO project_tokens (id, project_id, environment_id, name, token_prefix, token_hash, scopes, ip_allowlist, expires_at, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.ProjectID, t.EnvironmentID, t.Name, t.TokenPrefix, fullToken, scopesStr, ipStr, expiresAtVal, t.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create token: %w", err)
	}
	return nil
}

func (r *ProjectSettingsRepo) ListTokensByProject(ctx context.Context, projectID string) ([]*models.ProjectToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, environment_id, name, token_prefix, scopes, ip_allowlist, expires_at, created_at
		 FROM project_tokens WHERE project_id = ? ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list tokens: %w", err)
	}
	defer rows.Close()
	var out []*models.ProjectToken
	for rows.Next() {
		var t models.ProjectToken
		var scopesStr, ipStr string
		var expiresAtStr sql.NullString
		var createdAtStr string
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.EnvironmentID, &t.Name, &t.TokenPrefix, &scopesStr, &ipStr, &expiresAtStr, &createdAtStr); err != nil {
			return nil, fmt.Errorf("scan token: %w", err)
		}
		if scopesStr != "" {
			t.Scopes = strings.Split(scopesStr, ",")
		} else {
			t.Scopes = []string{}
		}
		if ipStr != "" {
			t.IPAllowlist = strings.Split(ipStr, ",")
		} else {
			t.IPAllowlist = []string{}
		}
		if expiresAtStr.Valid && expiresAtStr.String != "" {
			parsed, _ := time.Parse(time.RFC3339, expiresAtStr.String)
			t.ExpiresAt = &parsed
		}
		parsedCreated, _ := time.Parse(time.RFC3339, createdAtStr)
		t.CreatedAt = parsedCreated
		out = append(out, &t)
	}
	return out, rows.Err()
}

func (r *ProjectSettingsRepo) DeleteToken(ctx context.Context, id, projectID string) error {
	return r.deleteByIDAndProject(ctx, "project_tokens", id, projectID, "token")
}

func (r *ProjectSettingsRepo) GetTokenByHash(ctx context.Context, tokenHash string) (*models.ProjectToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var t models.ProjectToken
	var scopesStr, ipStr string
	var expiresAtStr sql.NullString
	var createdAtStr string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, project_id, environment_id, name, token_prefix, scopes, ip_allowlist, expires_at, created_at
		 FROM project_tokens WHERE token_hash = ?`, tokenHash).
		Scan(&t.ID, &t.ProjectID, &t.EnvironmentID, &t.Name, &t.TokenPrefix, &scopesStr, &ipStr, &expiresAtStr, &createdAtStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("Token", tokenHash)
		}
		return nil, fmt.Errorf("get token by hash: %w", err)
	}
	if scopesStr != "" {
		t.Scopes = strings.Split(scopesStr, ",")
	}
	if ipStr != "" {
		t.IPAllowlist = strings.Split(ipStr, ",")
	}
	if expiresAtStr.Valid && expiresAtStr.String != "" {
		parsed, _ := time.Parse(time.RFC3339, expiresAtStr.String)
		t.ExpiresAt = &parsed
	}
	parsedCreated, _ := time.Parse(time.RFC3339, createdAtStr)
	t.CreatedAt = parsedCreated
	return &t, nil
}

func (r *ProjectSettingsRepo) UpdateTokenLastUsed(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE project_tokens SET last_used_at = ? WHERE id = ?`, time.Now().Format(time.RFC3339), id)
	return err
}

func (r *ProjectSettingsRepo) AddMember(ctx context.Context, m *models.ProjectMember) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	m.InvitedAt = now
	if m.Status == "" {
		m.Status = "pending"
	}
	if m.Permission == "" {
		m.Permission = "Can Edit"
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO project_members (id, project_id, user_id, email, permission, status, invited_at, accepted_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(project_id, user_id) DO UPDATE SET
		 permission = excluded.permission,
		 status = excluded.status`,
		m.ID, m.ProjectID, m.UserID, m.Email, m.Permission, m.Status, m.InvitedAt, m.AcceptedAt)
	if err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (r *ProjectSettingsRepo) GetMember(ctx context.Context, projectID, userID string) (*models.ProjectMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	row := r.db.QueryRowContext(ctx,
		`SELECT id, project_id, user_id, email, permission, status, invited_at, accepted_at
		 FROM project_members WHERE project_id = ? AND user_id = ?`, projectID, userID)
	var m models.ProjectMember
	var invitedAt, acceptedAt sql.NullString
	if err := row.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Email, &m.Permission, &m.Status, &invitedAt, &acceptedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get project member: %w", err)
	}
	if invitedAt.Valid {
		m.InvitedAt, _ = time.Parse(time.RFC3339, invitedAt.String)
	}
	if acceptedAt.Valid {
		m.AcceptedAt, _ = time.Parse(time.RFC3339, acceptedAt.String)
	}
	return &m, nil
}

func (r *ProjectSettingsRepo) ListMembers(ctx context.Context, projectID string) ([]*models.ProjectMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, user_id, email, permission, status, invited_at, accepted_at
		 FROM project_members WHERE project_id = ? ORDER BY invited_at ASC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	defer rows.Close()
	var out []*models.ProjectMember
	for rows.Next() {
		var m models.ProjectMember
		var acceptedAt sql.NullTime
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Email, &m.Permission, &m.Status, &m.InvitedAt, &acceptedAt); err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}
		if acceptedAt.Valid {
			m.AcceptedAt = acceptedAt.Time
		}
		out = append(out, &m)
	}
	return out, rows.Err()
}

func (r *ProjectSettingsRepo) RemoveMember(ctx context.Context, id, projectID string) error {
	return r.deleteByIDAndProject(ctx, "project_members", id, projectID, "member")
}

func (r *ProjectSettingsRepo) AcceptAllInvitesForUser(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().Format(time.RFC3339)
	_, err := r.db.ExecContext(ctx, "UPDATE project_members SET status = 'accepted', accepted_at = ? WHERE user_id = ? AND status = 'pending'", now, userID)
	return err
}
