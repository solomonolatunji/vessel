package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	List(ctx context.Context, limit, offset int) ([]models.AuditLog, error)
}

type AuditLogRepo struct {
	db *sql.DB
}

func NewAuditLogRepo(db *sql.DB) *AuditLogRepo {
	return &AuditLogRepo{db: db}
}

func (r *AuditLogRepo) Create(ctx context.Context, log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, user_id, action, resource, details, ip_address)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, log.ID, log.UserID, log.Action, log.Resource, log.Details, log.IPAddress)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	return nil
}

func (r *AuditLogRepo) List(ctx context.Context, limit, offset int) ([]models.AuditLog, error) {
	query := `
		SELECT id, user_id, action, resource, details, ip_address, created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var details, ipAddress sql.NullString
		if err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.Resource, &details, &ipAddress, &log.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		if details.Valid {
			log.Details = details.String
		}
		if ipAddress.Valid {
			log.IPAddress = ipAddress.String
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return logs, nil
}
