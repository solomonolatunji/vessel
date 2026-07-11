package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"vessel.dev/vessel/internal/models"
)

type TeamEmailSettingsRepository interface {
	GetByTeamID(ctx context.Context, teamID string) (*models.TeamEmailSettings, error)
	Save(ctx context.Context, settings *models.TeamEmailSettings) error
}

type TeamEmailSettingsSQLiteRepository struct {
	db    *sql.DB
	vault Vault
}

func NewTeamEmailSettingsSQLiteRepository(db *sql.DB, v Vault) *TeamEmailSettingsSQLiteRepository {
	return &TeamEmailSettingsSQLiteRepository{
		db:    db,
		vault: v,
	}
}

const teamEmailSettingsColumns = `id, team_id, smtp_host, smtp_port, smtp_user, encrypted_smtp_password, smtp_from_name, smtp_from_address, encrypted_resend_api_key, use_resend, created_at, updated_at`

func scanTeamEmailSettings(scanner interface{ Scan(dest ...any) error }, s *models.TeamEmailSettings, v Vault) error {
	var encryptedSMTPPassword, encryptedResendAPIKey string
	var createdAt, updatedAt string
	err := scanner.Scan(
		&s.ID, &s.TeamID, &s.SMTPHost, &s.SMTPPort, &s.SMTPUser, &encryptedSMTPPassword, &s.SMTPFromName, &s.SMTPFromAddress, &encryptedResendAPIKey, &s.UseResend, &createdAt, &updatedAt,
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

func (r *TeamEmailSettingsSQLiteRepository) GetByTeamID(ctx context.Context, teamID string) (*models.TeamEmailSettings, error) {
	query := fmt.Sprintf(`SELECT %s FROM team_email_settings WHERE team_id = ? LIMIT 1`, teamEmailSettingsColumns)
	row := r.db.QueryRowContext(ctx, query, teamID)

	var s models.TeamEmailSettings
	err := scanTeamEmailSettings(row, &s, r.vault)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // No settings configured yet
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team email settings: %w", err)
	}

	return &s, nil
}

func (r *TeamEmailSettingsSQLiteRepository) Save(ctx context.Context, s *models.TeamEmailSettings) error {
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
		INSERT INTO team_email_settings (%s)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(team_id) DO UPDATE SET
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
		s.ID, s.TeamID, s.SMTPHost, s.SMTPPort, s.SMTPUser, encryptedSMTPPassword, s.SMTPFromName, s.SMTPFromAddress, encryptedResendAPIKey, s.UseResend, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to save team email settings: %w", err)
	}

	return nil
}
