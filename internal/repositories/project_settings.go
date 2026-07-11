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

	"vessl.dev/vessl/internal/models"
)

type ProjectSettingsRepository interface {
	CreateWebhook(ctx context.Context, w *models.Webhook) error
	ListWebhooksByProject(ctx context.Context, projectID string) ([]*models.Webhook, error)
	DeleteWebhook(ctx context.Context, id, projectID string) error
	CreateToken(ctx context.Context, t *models.ProjectToken, fullToken string) error
	ListTokensByProject(ctx context.Context, projectID string) ([]*models.ProjectToken, error)
	DeleteToken(ctx context.Context, id, projectID string) error
	GetTokenByHash(ctx context.Context, tokenHash string) (*models.ProjectToken, error)
	UpdateTokenLastUsed(ctx context.Context, id string) error
	AddMember(ctx context.Context, m *models.ProjectMember) error
	ListMembers(ctx context.Context, projectID string) ([]*models.ProjectMember, error)
	RemoveMember(ctx context.Context, id, projectID string) error
}

type ProjectSettingsSQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewProjectSettingsSQLiteRepository(db *sql.DB) *ProjectSettingsSQLiteRepository {
	return &ProjectSettingsSQLiteRepository{db: db}
}

func (r *ProjectSettingsSQLiteRepository) CreateWebhook(ctx context.Context, w *models.Webhook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now
	eventTypesStr := strings.Join(w.EventTypes, ",")
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO project_webhooks (id, project_id, url, event_types, include_pr_environments, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		w.ID, w.ProjectID, w.URL, eventTypesStr, w.IncludePREnvironments, w.CreatedAt, w.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}
	return nil
}

func (r *ProjectSettingsSQLiteRepository) ListWebhooksByProject(ctx context.Context, projectID string) ([]*models.Webhook, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, url, event_types, include_pr_environments, created_at, updated_at
		 FROM project_webhooks WHERE project_id = ? ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list webhooks: %w", err)
	}
	defer rows.Close()
	var out []*models.Webhook
	for rows.Next() {
		var w models.Webhook
		var eventsStr string
		var includePr int
		if err := rows.Scan(&w.ID, &w.ProjectID, &w.URL, &eventsStr, &includePr, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan webhook: %w", err)
		}
		if eventsStr != "" {
			w.EventTypes = strings.Split(eventsStr, ",")
		} else {
			w.EventTypes = []string{}
		}
		w.IncludePREnvironments = includePr == 1
		out = append(out, &w)
	}
	return out, rows.Err()
}

func (r *ProjectSettingsSQLiteRepository) deleteByIDAndProject(ctx context.Context, table, id, projectID, entityName string) error {
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

func (r *ProjectSettingsSQLiteRepository) DeleteWebhook(ctx context.Context, id, projectID string) error {
	return r.deleteByIDAndProject(ctx, "project_webhooks", id, projectID, "webhook")
}

func (r *ProjectSettingsSQLiteRepository) CreateToken(ctx context.Context, t *models.ProjectToken, fullToken string) error {
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

func (r *ProjectSettingsSQLiteRepository) ListTokensByProject(ctx context.Context, projectID string) ([]*models.ProjectToken, error) {
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

func (r *ProjectSettingsSQLiteRepository) DeleteToken(ctx context.Context, id, projectID string) error {
	return r.deleteByIDAndProject(ctx, "project_tokens", id, projectID, "token")
}

func (r *ProjectSettingsSQLiteRepository) GetTokenByHash(ctx context.Context, tokenHash string) (*models.ProjectToken, error) {
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

func (r *ProjectSettingsSQLiteRepository) UpdateTokenLastUsed(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE project_tokens SET last_used_at = ? WHERE id = ?`, time.Now().Format(time.RFC3339), id)
	return err
}

func (r *ProjectSettingsSQLiteRepository) AddMember(ctx context.Context, m *models.ProjectMember) error {
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
		 ON CONFLICT(project_id, email) DO UPDATE SET
		 permission = excluded.permission,
		 status = excluded.status`,
		m.ID, m.ProjectID, m.UserID, m.Email, m.Permission, m.Status, m.InvitedAt, m.AcceptedAt)
	if err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (r *ProjectSettingsSQLiteRepository) ListMembers(ctx context.Context, projectID string) ([]*models.ProjectMember, error) {
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

func (r *ProjectSettingsSQLiteRepository) RemoveMember(ctx context.Context, id, projectID string) error {
	return r.deleteByIDAndProject(ctx, "project_members", id, projectID, "member")
}
