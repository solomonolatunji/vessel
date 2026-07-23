package repositories

import (
	"codedock.dev/codedock/internal/utils"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock.dev/codedock/internal/models"
)

type GitRepository interface {
	SaveProvider(ctx context.Context, gp *models.GitProviderConfig) error
	GetProvider(ctx context.Context, userID, provider string) (*models.GitProviderConfig, error)
	GetAnyProviderByType(ctx context.Context, provider string) (*models.GitProviderConfig, error)
	ListProvidersByUser(ctx context.Context, userID string) ([]*models.GitProviderConfig, error)
	DeleteProvider(ctx context.Context, userID, provider string) error
}

type GitRepo struct {
	db    *sqlx.DB
	vault Vault
}

func NewGitRepo(db *sql.DB, vault Vault) *GitRepo {
	return &GitRepo{db: sqlx.NewDb(db, "sqlite"), vault: vault}
}

func (r *GitRepo) SaveProvider(ctx context.Context, gp *models.GitProviderConfig) error {
	if gp.ID == "" {
		gp.ID = uuid.NewString()
	}
	now := time.Now()
	gp.CreatedAt = now
	gp.UpdatedAt = now
	encryptedToken, err := r.vault.Encrypt(gp.AccessToken)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO user_git_providers (id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, provider) DO UPDATE SET encrypted_access_token = excluded.encrypted_access_token, account_name = excluded.account_name, updated_at = excluded.updated_at`,
		gp.ID, gp.UserID, gp.Provider, encryptedToken, gp.AccountName, gp.CreatedAt, gp.UpdatedAt,
	)
	return err
}

func (r *GitRepo) GetProvider(ctx context.Context, userID, provider string) (*models.GitProviderConfig, error) {
	if userID == "" {
		return r.GetAnyProviderByType(ctx, provider)
	}
	var gp models.GitProviderConfig
	err := r.db.GetContext(ctx, &gp, `SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE user_id = ? AND provider = ?`, userID, provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	decryptedToken, err := r.vault.Decrypt(gp.AccessToken)
	if err != nil {
		return nil, err
	}
	gp.AccessToken = decryptedToken
	return &gp, nil
}

func (r *GitRepo) GetAnyProviderByType(ctx context.Context, provider string) (*models.GitProviderConfig, error) {
	var gp models.GitProviderConfig
	err := r.db.GetContext(ctx, &gp, `SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE provider = ? LIMIT 1`, provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("AnyProvider", provider)
		}
		return nil, err
	}
	decryptedToken, err := r.vault.Decrypt(gp.AccessToken)
	if err != nil {
		return nil, err
	}
	gp.AccessToken = decryptedToken
	return &gp, nil
}

func (r *GitRepo) ListProvidersByUser(ctx context.Context, userID string) ([]*models.GitProviderConfig, error) {
	var list []*models.GitProviderConfig
	err := r.db.SelectContext(ctx, &list, `SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	for _, gp := range list {
		decryptedToken, err := r.vault.Decrypt(gp.AccessToken)
		if err != nil {
			return nil, err
		}
		gp.AccessToken = decryptedToken
	}
	return list, nil
}

func (r *GitRepo) DeleteProvider(ctx context.Context, userID, provider string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_git_providers WHERE user_id = ? AND provider = ?`, userID, provider)
	return err
}
