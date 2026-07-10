package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/vault"
)

type VercelRepository struct {
	db    *sql.DB
	vault *vault.Vault
}

func NewVercelRepository(db *sql.DB, v *vault.Vault) *VercelRepository {
	return &VercelRepository{db: db, vault: v}
}

func (r *VercelRepository) SaveAccount(ctx context.Context, account *models.UserVercelAccount) error {
	if account.ID == "" {
		account.ID = uuid.NewString()
	}
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	encryptedToken, err := r.vault.Encrypt(account.AccessToken)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO user_vercel_accounts (id, user_id, encrypted_access_token, team_id, account_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, team_id) DO UPDATE SET
			encrypted_access_token = excluded.encrypted_access_token,
			account_name = excluded.account_name,
			updated_at = excluded.updated_at
	`

	var teamID sql.NullString
	if account.TeamID != nil && *account.TeamID != "" {
		teamID = sql.NullString{String: *account.TeamID, Valid: true}
	} else {
		teamID = sql.NullString{Valid: false}
	}

	_, err = r.db.ExecContext(ctx, query,
		account.ID, account.UserID, encryptedToken, teamID, account.AccountName, account.CreatedAt, account.UpdatedAt,
	)
	return err
}

func (r *VercelRepository) GetAccount(ctx context.Context, userID string, teamID *string) (*models.UserVercelAccount, error) {
	var account models.UserVercelAccount
	var encryptedToken string
	var sqlTeamID sql.NullString

	query := `
		SELECT id, user_id, encrypted_access_token, team_id, account_name, created_at, updated_at
		FROM user_vercel_accounts
		WHERE user_id = ? AND team_id IS ?
	`

	var tID interface{} = nil
	if teamID != nil && *teamID != "" {
		tID = *teamID
	}

	err := r.db.QueryRowContext(ctx, query, userID, tID).Scan(
		&account.ID, &account.UserID, &encryptedToken, &sqlTeamID, &account.AccountName, &account.CreatedAt, &account.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if sqlTeamID.Valid {
		account.TeamID = &sqlTeamID.String
	}

	decryptedToken, err := r.vault.Decrypt(encryptedToken)
	if err != nil {
		return nil, err
	}
	account.AccessToken = decryptedToken

	return &account, nil
}

func (r *VercelRepository) GetAccountsForUser(ctx context.Context, userID string) ([]*models.UserVercelAccount, error) {
	query := `
		SELECT id, user_id, encrypted_access_token, team_id, account_name, created_at, updated_at
		FROM user_vercel_accounts
		WHERE user_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*models.UserVercelAccount
	for rows.Next() {
		var account models.UserVercelAccount
		var encryptedToken string
		var sqlTeamID sql.NullString

		if err := rows.Scan(&account.ID, &account.UserID, &encryptedToken, &sqlTeamID, &account.AccountName, &account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, err
		}

		if sqlTeamID.Valid {
			account.TeamID = &sqlTeamID.String
		}

		decryptedToken, err := r.vault.Decrypt(encryptedToken)
		if err != nil {
			return nil, err
		}
		account.AccessToken = decryptedToken
		accounts = append(accounts, &account)
	}
	return accounts, nil
}
