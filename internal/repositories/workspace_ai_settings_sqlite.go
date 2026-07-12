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

type WorkspaceAISettingsRepository interface {
	Get(ctx context.Context, workspaceID string) (*models.WorkspaceAISettings, error)
	Save(ctx context.Context, settings *models.WorkspaceAISettings) error
}

type WorkspaceAISettingsSQLiteRepository struct {
	db    *sql.DB
	vault Vault
}

func NewWorkspaceAISettingsSQLiteRepository(db *sql.DB, vault Vault) *WorkspaceAISettingsSQLiteRepository {
	return &WorkspaceAISettingsSQLiteRepository{db: db, vault: vault}
}

func (r *WorkspaceAISettingsSQLiteRepository) Get(ctx context.Context, workspaceID string) (*models.WorkspaceAISettings, error) {
	query := `SELECT id, team_id, provider, encrypted_api_key, created_at, updated_at FROM team_ai_settings WHERE team_id = ?`
	row := r.db.QueryRowContext(ctx, query, workspaceID)

	var s models.WorkspaceAISettings
	var encryptedKey string
	if err := row.Scan(&s.ID, &s.WorkspaceID, &s.Provider, &encryptedKey, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("WorkspaceAISettings", workspaceID)
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

func (r *WorkspaceAISettingsSQLiteRepository) Save(ctx context.Context, settings *models.WorkspaceAISettings) error {
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

	if _, err := r.db.ExecContext(ctx, query, settings.ID, settings.WorkspaceID, settings.Provider, encryptedKey, settings.CreatedAt, settings.UpdatedAt); err != nil {
		return fmt.Errorf("failed to save team AI settings: %w", err)
	}
	return nil
}
