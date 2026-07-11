package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"vessl.dev/vessl/internal/utils"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
)

type GitRepository interface {
	SaveProvider(ctx context.Context, gp *models.GitProviderConfig) error
	GetProvider(ctx context.Context, userID, provider string) (*models.GitProviderConfig, error)
	GetAnyProviderByType(ctx context.Context, provider string) (*models.GitProviderConfig, error)
	ListProvidersByUser(ctx context.Context, userID string) ([]*models.GitProviderConfig, error)
	DeleteProvider(ctx context.Context, userID, provider string) error
}

type GitSQLiteRepository struct {
	db    *sql.DB
	vault Vault
}

func NewGitSQLiteRepository(db *sql.DB, vault Vault) *GitSQLiteRepository {
	return &GitSQLiteRepository{db: db, vault: vault}
}

func (r *GitSQLiteRepository) SaveProvider(_ context.Context, gp *models.GitProviderConfig) error {
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
	_, err = r.db.Exec(
		`INSERT INTO user_git_providers (id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, provider) DO UPDATE SET encrypted_access_token = excluded.encrypted_access_token, account_name = excluded.account_name, updated_at = excluded.updated_at`,
		gp.ID, gp.UserID, gp.Provider, encryptedToken, gp.AccountName, gp.CreatedAt, gp.UpdatedAt,
	)
	return err
}

func (r *GitSQLiteRepository) GetProvider(_ context.Context, userID, provider string) (*models.GitProviderConfig, error) {
	if userID == "" {
		return r.GetAnyProviderByType(context.Background(), provider)
	}
	row := r.db.QueryRow(`SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE user_id = ? AND provider = ?`, userID, provider)
	var gp models.GitProviderConfig
	var encryptedToken string
	err := row.Scan(&gp.ID, &gp.UserID, &gp.Provider, &encryptedToken, &gp.AccountName, &gp.CreatedAt, &gp.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	decryptedToken, err := r.vault.Decrypt(encryptedToken)
	if err != nil {
		return nil, err
	}
	gp.AccessToken = decryptedToken
	return &gp, nil
}

func (r *GitSQLiteRepository) GetAnyProviderByType(_ context.Context, provider string) (*models.GitProviderConfig, error) {
	row := r.db.QueryRow(`SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE provider = ? LIMIT 1`, provider)
	var gp models.GitProviderConfig
	var encryptedToken string
	err := row.Scan(&gp.ID, &gp.UserID, &gp.Provider, &encryptedToken, &gp.AccountName, &gp.CreatedAt, &gp.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("AnyProvider", provider)
		}
		return nil, err
	}
	decryptedToken, err := r.vault.Decrypt(encryptedToken)
	if err != nil {
		return nil, err
	}
	gp.AccessToken = decryptedToken
	return &gp, nil
}

func (r *GitSQLiteRepository) ListProvidersByUser(_ context.Context, userID string) ([]*models.GitProviderConfig, error) {
	rows, err := r.db.Query(`SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.GitProviderConfig
	for rows.Next() {
		var gp models.GitProviderConfig
		var encryptedToken string
		if err := rows.Scan(&gp.ID, &gp.UserID, &gp.Provider, &encryptedToken, &gp.AccountName, &gp.CreatedAt, &gp.UpdatedAt); err != nil {
			return nil, err
		}
		decryptedToken, err := r.vault.Decrypt(encryptedToken)
		if err != nil {
			return nil, err
		}
		gp.AccessToken = decryptedToken
		list = append(list, &gp)
	}
	return list, nil
}

func (r *GitSQLiteRepository) DeleteProvider(_ context.Context, userID, provider string) error {
	_, err := r.db.Exec(`DELETE FROM user_git_providers WHERE user_id = ? AND provider = ?`, userID, provider)
	return err
}
