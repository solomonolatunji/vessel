package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
)

func (r *WorkspaceSQLiteRepository) CreateTrustedDomain(ctx context.Context, d *models.TrustedDomain) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now().UTC()
	}
	if d.Role == "" {
		d.Role = "developer"
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO workspace_trusted_domains (id, workspace_id, domain, role, created_at) VALUES (?, ?, ?, ?, ?)`,
		d.ID, d.WorkspaceID, d.Domain, d.Role, d.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create trusted domain: %w", err)
	}
	return nil
}

func (r *WorkspaceSQLiteRepository) ListTrustedDomains(ctx context.Context, workspaceID string) ([]*models.TrustedDomain, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx, `SELECT id, workspace_id, domain, role, created_at FROM workspace_trusted_domains WHERE workspace_id = ? ORDER BY created_at DESC`, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("list trusted domains: %w", err)
	}
	defer rows.Close()
	var list []*models.TrustedDomain
	for rows.Next() {
		var d models.TrustedDomain
		var createdAtStr string
		if err := rows.Scan(&d.ID, &d.WorkspaceID, &d.Domain, &d.Role, &createdAtStr); err != nil {
			return nil, err
		}
		d.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &d)
	}
	return list, nil
}

func (r *WorkspaceSQLiteRepository) DeleteTrustedDomain(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `DELETE FROM workspace_trusted_domains WHERE id = ?`, id)
	return err
}

func (r *WorkspaceSQLiteRepository) CreateSSHKey(ctx context.Context, key *models.SSHKey) error {
	if key.ID == "" {
		key.ID = uuid.NewString()
	}
	if key.CreatedAt.IsZero() {
		key.CreatedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO workspace_ssh_keys (id, workspace_id, name, public_key, created_at) VALUES (?, ?, ?, ?, ?)`,
		key.ID, key.WorkspaceID, key.Name, key.PublicKey, key.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create ssh key: %w", err)
	}
	return nil
}

func (r *WorkspaceSQLiteRepository) ListSSHKeys(ctx context.Context, workspaceID string) ([]*models.SSHKey, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx, `SELECT id, workspace_id, name, public_key, created_at FROM workspace_ssh_keys WHERE workspace_id = ? ORDER BY created_at DESC`, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("list ssh keys: %w", err)
	}
	defer rows.Close()
	var list []*models.SSHKey
	for rows.Next() {
		var k models.SSHKey
		var createdAtStr string
		if err := rows.Scan(&k.ID, &k.WorkspaceID, &k.Name, &k.PublicKey, &createdAtStr); err != nil {
			return nil, err
		}
		k.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &k)
	}
	return list, nil
}

func (r *WorkspaceSQLiteRepository) DeleteSSHKey(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `DELETE FROM workspace_ssh_keys WHERE id = ?`, id)
	return err
}

func (r *WorkspaceSQLiteRepository) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	if log.ID == "" {
		log.ID = uuid.NewString()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO workspace_audit_logs (id, workspace_id, project_id, environment_id, action, actor, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		log.ID, log.WorkspaceID, log.ProjectID, log.EnvironmentID, log.Action, log.Actor, log.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

func (r *WorkspaceSQLiteRepository) ListAuditLogs(ctx context.Context, workspaceID string, limit, offset int) ([]*models.AuditLog, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM workspace_audit_logs WHERE workspace_id = ?`, workspaceID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `SELECT id, workspace_id, COALESCE(project_id, ''), COALESCE(environment_id, ''), action, actor, created_at FROM workspace_audit_logs WHERE workspace_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, workspaceID, limit, offset)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	if rows == nil {
		return []*models.AuditLog{}, 0, nil
	}
	defer rows.Close()
	var list []*models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var createdAtStr string
		if err := rows.Scan(&log.ID, &log.WorkspaceID, &log.ProjectID, &log.EnvironmentID, &log.Action, &log.Actor, &createdAtStr); err != nil {
			return nil, 0, err
		}
		log.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &log)
	}
	return list, total, nil
}
