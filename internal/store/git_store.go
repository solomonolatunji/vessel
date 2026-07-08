package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// SaveGitProvider encrypts and stores a user's Git platform credentials.
func (s *Store) SaveGitProvider(gp *types.GitProviderConfig) error {
	if gp.ID == "" {
		gp.ID = uuid.NewString()
	}
	now := time.Now()
	gp.CreatedAt = now
	gp.UpdatedAt = now

	encryptedToken, err := s.vault.Encrypt(gp.AccessToken)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`INSERT INTO user_git_providers (id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, provider) DO UPDATE SET encrypted_access_token = excluded.encrypted_access_token, account_name = excluded.account_name, updated_at = excluded.updated_at`,
		gp.ID, gp.UserID, gp.Provider, encryptedToken, gp.AccountName, gp.CreatedAt, gp.UpdatedAt,
	)
	return err
}

// GetGitProvider retrieves and decrypts a user's stored Git access token for a given provider (github or gitlab).
func (s *Store) GetGitProvider(userID, provider string) (*types.GitProviderConfig, error) {
	if userID == "" {
		return s.GetAnyGitProviderByProvider(provider)
	}

	row := s.db.QueryRow(`SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE user_id = ? AND provider = ?`, userID, provider)

	var gp types.GitProviderConfig
	var encryptedToken string
	err := row.Scan(&gp.ID, &gp.UserID, &gp.Provider, &encryptedToken, &gp.AccountName, &gp.CreatedAt, &gp.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	decryptedToken, err := s.vault.Decrypt(encryptedToken)
	if err != nil {
		return nil, err
	}
	gp.AccessToken = decryptedToken

	return &gp, nil
}

// GetAnyGitProviderByProvider retrieves the first available encrypted Git access token for a provider when cloning system-wide.
func (s *Store) GetAnyGitProviderByProvider(provider string) (*types.GitProviderConfig, error) {
	row := s.db.QueryRow(`SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE provider = ? LIMIT 1`, provider)

	var gp types.GitProviderConfig
	var encryptedToken string
	err := row.Scan(&gp.ID, &gp.UserID, &gp.Provider, &encryptedToken, &gp.AccountName, &gp.CreatedAt, &gp.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	decryptedToken, err := s.vault.Decrypt(encryptedToken)
	if err != nil {
		return nil, err
	}
	gp.AccessToken = decryptedToken

	return &gp, nil
}

// DeleteGitProvider removes a stored Git provider connection for a user.
func (s *Store) DeleteGitProvider(userID, provider string) error {
	_, err := s.db.Exec(`DELETE FROM user_git_providers WHERE user_id = ? AND provider = ?`, userID, provider)
	return err
}
