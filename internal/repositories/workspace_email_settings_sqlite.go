package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"vessl.dev/vessl/internal/models"
)

type WorkspaceEmailSettingsRepository interface {
	GetByWorkspaceID(ctx context.Context, workspaceID string) (*models.WorkspaceEmailSettings, error)
	Save(ctx context.Context, settings *models.WorkspaceEmailSettings) error
}

type WorkspaceEmailSettingsSQLiteRepository struct {
	db    *sqlx.DB
	vault Vault
}

func NewWorkspaceEmailSettingsSQLiteRepository(db *sql.DB, v Vault) *WorkspaceEmailSettingsSQLiteRepository {
	return &WorkspaceEmailSettingsSQLiteRepository{
		db:    sqlx.NewDb(db, "sqlite"),
		vault: v,
	}
}

const teamEmailSettingsColumns = `id, workspace_id, smtp_host, smtp_port, smtp_user, encrypted_smtp_password, smtp_from_name, smtp_from_address, encrypted_resend_api_key, use_resend, created_at, updated_at`

func scanWorkspaceEmailSettings(scanner interface{ Scan(dest ...any) error }, s *models.WorkspaceEmailSettings, v Vault) error {
	var encryptedSMTPPassword, encryptedResendAPIKey string
	var createdAt, updatedAt string
	err := scanner.Scan(
		&s.ID, &s.WorkspaceID, &s.SMTPHost, &s.SMTPPort, &s.SMTPUser, &encryptedSMTPPassword, &s.SMTPFromName, &s.SMTPFromAddress, &encryptedResendAPIKey, &s.UseResend, &createdAt, &updatedAt,
	)
	if err != nil {
		return err
	}

	if encryptedSMTPPassword != "" {
		decrypted, err := v.Decrypt(encryptedSMTPPassword)
		if err == nil {
			s.SMTPPassword = decrypted
		}
	}
	if encryptedResendAPIKey != "" {
		decrypted, err := v.Decrypt(encryptedResendAPIKey)
		if err == nil {
			s.ResendAPIKey = decrypted
		}
	}

	s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	s.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return nil
}

func (r *WorkspaceEmailSettingsSQLiteRepository) GetByWorkspaceID(ctx context.Context, workspaceID string) (*models.WorkspaceEmailSettings, error) {
	query := fmt.Sprintf(`SELECT %s FROM workspace_email_settings WHERE workspace_id = ? LIMIT 1`, teamEmailSettingsColumns)
	row := r.db.QueryRowContext(ctx, query, workspaceID)

	var s models.WorkspaceEmailSettings
	err := scanWorkspaceEmailSettings(row, &s, r.vault)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team email settings: %w", err)
	}

	return &s, nil
}

func (r *WorkspaceEmailSettingsSQLiteRepository) Save(ctx context.Context, s *models.WorkspaceEmailSettings) error {
	now := time.Now().UTC().Format(time.RFC3339)

	var encryptedSMTPPassword, encryptedResendAPIKey string
	var err error

	if s.SMTPPassword != "" {
		encryptedSMTPPassword, err = r.vault.Encrypt(s.SMTPPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt smtp password: %w", err)
		}
	}
	if s.ResendAPIKey != "" {
		encryptedResendAPIKey, err = r.vault.Encrypt(s.ResendAPIKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt resend api key: %w", err)
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO workspace_email_settings (%s)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(workspace_id) DO UPDATE SET
			smtp_host = excluded.smtp_host,
			smtp_port = excluded.smtp_port,
			smtp_user = excluded.smtp_user,
			encrypted_smtp_password = excluded.encrypted_smtp_password,
			smtp_from_name = excluded.smtp_from_name,
			smtp_from_address = excluded.smtp_from_address,
			encrypted_resend_api_key = excluded.encrypted_resend_api_key,
			use_resend = excluded.use_resend,
			updated_at = excluded.updated_at
	`, teamEmailSettingsColumns)

	_, err = r.db.ExecContext(ctx, query,
		s.ID, s.WorkspaceID, s.SMTPHost, s.SMTPPort, s.SMTPUser, encryptedSMTPPassword, s.SMTPFromName, s.SMTPFromAddress, encryptedResendAPIKey, s.UseResend, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to save team email settings: %w", err)
	}

	return nil
}
