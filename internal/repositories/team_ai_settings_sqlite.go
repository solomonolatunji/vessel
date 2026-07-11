package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

type TeamAISettingsRepository interface {
	Get(ctx context.Context, teamID string) (*models.TeamAISettings, error)
	Save(ctx context.Context, settings *models.TeamAISettings) error
}

type TeamAISettingsSQLiteRepository struct {
	db    *sql.DB
	vault Vault
}

func NewTeamAISettingsSQLiteRepository(db *sql.DB, vault Vault) *TeamAISettingsSQLiteRepository {
	return &TeamAISettingsSQLiteRepository{db: db, vault: vault}
}

func (r *TeamAISettingsSQLiteRepository) Get(ctx context.Context, teamID string) (*models.TeamAISettings, error) {
	query := `SELECT id, team_id, provider, encrypted_api_key, created_at, updated_at FROM team_ai_settings WHERE team_id = ?`
	row := r.db.QueryRowContext(ctx, query, teamID)

	var s models.TeamAISettings
	var encryptedKey string
	if err := row.Scan(&s.ID, &s.TeamID, &s.Provider, &encryptedKey, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("TeamAISettings", teamID)
		}
		return nil, fmt.Errorf("failed to get team AI settings: %w", err)
	}

	if key, err := r.vault.Decrypt(encryptedKey); err == nil {
		s.APIKey = key
	} else {
		s.APIKey = encryptedKey
	}

	return &s, nil
}

func (r *TeamAISettingsSQLiteRepository) Save(ctx context.Context, settings *models.TeamAISettings) error {
	query := `
		INSERT INTO team_ai_settings (id, team_id, provider, encrypted_api_key, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(team_id) DO UPDATE SET
			provider = excluded.provider,
			encrypted_api_key = excluded.encrypted_api_key,
			updated_at = CURRENT_TIMESTAMP
	`
	encryptedKey, err := r.vault.Encrypt(settings.APIKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt API key: %w", err)
	}

	if settings.CreatedAt.IsZero() {
		settings.CreatedAt = time.Now()
	}
	settings.UpdatedAt = time.Now()

	if _, err := r.db.ExecContext(ctx, query, settings.ID, settings.TeamID, settings.Provider, encryptedKey, settings.CreatedAt, settings.UpdatedAt); err != nil {
		return fmt.Errorf("failed to save team AI settings: %w", err)
	}
	return nil
}
