package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

type OAuthRepository interface {
	ListProviders(ctx context.Context) ([]models.OAuthProviderConfig, error)
	GetProvider(ctx context.Context, idOrName string) (*models.OAuthProviderConfig, error)
	SaveProvider(ctx context.Context, p *models.OAuthProviderConfig) error
	GetUserTOTPSecret(ctx context.Context, userID string) (secret string, recoveryCodes []string, err error)
	UpdateUserTOTP(ctx context.Context, userID string, enabled bool, secret string, recoveryCodes []string) error
}

type OAuthSQLiteRepository struct {
	db *sql.DB
}

func NewOAuthSQLiteRepository(db *sql.DB) *OAuthSQLiteRepository {
	return &OAuthSQLiteRepository{db: db}
}

func (r *OAuthSQLiteRepository) ListProviders(ctx context.Context) ([]models.OAuthProviderConfig, error) {
	query := `SELECT id, provider_name, enabled, COALESCE(client_id, ''), COALESCE(client_secret, ''), COALESCE(redirect_uri, ''), COALESCE(base_url, ''), COALESCE(tenant, ''), created_at, updated_at FROM oauth_providers ORDER BY provider_name ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list oauth providers: %w", err)
	}
	defer rows.Close()
	var providers []models.OAuthProviderConfig
	for rows.Next() {
		var p models.OAuthProviderConfig
		if err := rows.Scan(&p.ID, &p.ProviderName, &p.Enabled, &p.ClientID, &p.ClientSecret, &p.RedirectURI, &p.BaseURL, &p.Tenant, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed scanning oauth provider: %w", err)
		}
		providers = append(providers, p)
	}
	return providers, nil
}

func (r *OAuthSQLiteRepository) GetProvider(ctx context.Context, idOrName string) (*models.OAuthProviderConfig, error) {
	query := `SELECT id, provider_name, enabled, COALESCE(client_id, ''), COALESCE(client_secret, ''), COALESCE(redirect_uri, ''), COALESCE(base_url, ''), COALESCE(tenant, ''), created_at, updated_at FROM oauth_providers WHERE id = ? OR provider_name = ?`
	row := r.db.QueryRowContext(ctx, query, idOrName, idOrName)
	var p models.OAuthProviderConfig
	if err := row.Scan(&p.ID, &p.ProviderName, &p.Enabled, &p.ClientID, &p.ClientSecret, &p.RedirectURI, &p.BaseURL, &p.Tenant, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("Provider", idOrName)
		}
		return nil, fmt.Errorf("failed to get oauth provider: %w", err)
	}
	return &p, nil
}

func (r *OAuthSQLiteRepository) SaveProvider(ctx context.Context, p *models.OAuthProviderConfig) error {
	now := time.Now().UTC()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now
	query := `INSERT INTO oauth_providers (
		id, provider_name, enabled, client_id, client_secret, redirect_uri, base_url, tenant, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		provider_name = excluded.provider_name,
		enabled = excluded.enabled,
		client_id = excluded.client_id,
		client_secret = excluded.client_secret,
		redirect_uri = excluded.redirect_uri,
		base_url = excluded.base_url,
		tenant = excluded.tenant,
		updated_at = excluded.updated_at`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.ProviderName, p.Enabled, p.ClientID, p.ClientSecret, p.RedirectURI, p.BaseURL, p.Tenant, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save oauth provider: %w", err)
	}
	return nil
}

func (r *OAuthSQLiteRepository) GetUserTOTPSecret(ctx context.Context, userID string) (string, []string, error) {
	var secret string
	var recovery string
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(totp_secret, ''), COALESCE(recovery_codes, '') FROM users WHERE id = ?`, userID).Scan(&secret, &recovery)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get totp secret: %w", err)
	}
	var codes []string
	if recovery != "" {
		for _, part := range strings.Split(recovery, ",") {
			if part = strings.TrimSpace(part); part != "" {
				codes = append(codes, part)
			}
		}
	}
	return secret, codes, nil
}

func (r *OAuthSQLiteRepository) UpdateUserTOTP(ctx context.Context, userID string, enabled bool, secret string, recoveryCodes []string) error {
	recoveryStr := strings.Join(recoveryCodes, ",")
	_, err := r.db.ExecContext(ctx, `UPDATE users SET totp_enabled = ?, totp_secret = ?, recovery_codes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, enabled, secret, recoveryStr, userID)
	if err != nil {
		return fmt.Errorf("failed to update user totp: %w", err)
	}
	return nil
}
