package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
	"vessl.dev/vessl/internal/utils"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, u *models.User) error
	CreatePAT(ctx context.Context, pat *models.PersonalAccessToken) error
	ListPATs(ctx context.Context, userID string) ([]*models.PersonalAccessToken, error)
	DeletePAT(ctx context.Context, id, userID string) error
}

type UserSQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewUserSQLiteRepository(db *sql.DB) *UserSQLiteRepository {
	return &UserSQLiteRepository{db: db}
}

func (r *UserSQLiteRepository) CreateUser(ctx context.Context, u *models.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`, u.ID, u.Email, u.PasswordHash, u.Role, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *UserSQLiteRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var u models.User
	err := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash, role, created_at, updated_at
		FROM users WHERE email = ?`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("User", email)
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserSQLiteRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var u models.User
	err := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash, role, created_at, updated_at
		FROM users WHERE id = ?`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("User", id)
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserSQLiteRepository) ListUsers(ctx context.Context) ([]models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx, `SELECT id, email, password_hash, role, created_at, updated_at FROM users ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserSQLiteRepository) UpdateUser(ctx context.Context, u *models.User) error {
	u.UpdatedAt = time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE users SET email = ?, password_hash = ?, role = ?, updated_at = ? WHERE id = ?`,
		u.Email, u.PasswordHash, u.Role, u.UpdatedAt, u.ID)
	return err
}

func (r *UserSQLiteRepository) CreatePAT(ctx context.Context, pat *models.PersonalAccessToken) error {
	if pat.ID == "" {
		pat.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	if pat.CreatedAt.IsZero() {
		pat.CreatedAt = now
	}
	if pat.ExpiresAt.IsZero() {
		pat.ExpiresAt = now.Add(365 * 24 * time.Hour)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO personal_access_tokens (id, user_id, name, token_hash, prefix, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		pat.ID, pat.UserID, pat.Name, pat.TokenHash, pat.Prefix, pat.ExpiresAt.Format(time.RFC3339), pat.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to create personal access token: %w", err)
	}
	return nil
}

func (r *UserSQLiteRepository) ListPATs(ctx context.Context, userID string) ([]*models.PersonalAccessToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, name, prefix, expires_at, created_at FROM personal_access_tokens WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list personal access tokens: %w", err)
	}
	defer rows.Close()
	var list []*models.PersonalAccessToken
	for rows.Next() {
		var pat models.PersonalAccessToken
		var expStr, createdStr string
		if err := rows.Scan(&pat.ID, &pat.UserID, &pat.Name, &pat.Prefix, &expStr, &createdStr); err != nil {
			return nil, err
		}
		pat.ExpiresAt, _ = time.Parse(time.RFC3339, expStr)
		pat.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		list = append(list, &pat)
	}
	return list, nil
}

func (r *UserSQLiteRepository) DeletePAT(ctx context.Context, id, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.ExecContext(ctx, `DELETE FROM personal_access_tokens WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete personal access token: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("personal access token not found or unauthorized")
	}
	return nil
}
