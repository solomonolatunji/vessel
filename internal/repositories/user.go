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
	"github.com/jmoiron/sqlx"

	"vessl.dev/vessl/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]models.User, int, error)
	CountUsers(ctx context.Context) (int, error)
	UpdateUser(ctx context.Context, u *models.User) error
	DeleteUser(ctx context.Context, id string) error
	CreatePAT(ctx context.Context, pat *models.PersonalAccessToken) error
	ListPATs(ctx context.Context, userID string) ([]*models.PersonalAccessToken, error)
	DeletePAT(ctx context.Context, id, userID string) error
}

type UserRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *UserRepo) CountUsers(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users")
	return count, err
}

func (r *UserRepo) CreateUser(ctx context.Context, u *models.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO users (id, email, name, password_hash, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, u.ID, u.Email, u.Name, u.PasswordHash, u.Role, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var u models.User
	err := r.db.GetContext(ctx, &u, `SELECT id, email, name, password_hash, role, created_at, updated_at
		FROM users WHERE email = ?`, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("User", email)
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var u models.User
	err := r.db.GetContext(ctx, &u, `SELECT id, email, name, password_hash, role, created_at, updated_at
		FROM users WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("User", id)
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) ListUsers(ctx context.Context, limit, offset int) ([]models.User, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var total int
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM users`); err != nil {
		return nil, 0, err
	}

	var users []models.User
	err := r.db.SelectContext(ctx, &users, `SELECT id, email, name, password_hash, role, created_at, updated_at FROM users ORDER BY created_at ASC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if users == nil {
		users = make([]models.User, 0)
	}
	return users, total, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, u *models.User) error {
	u.UpdatedAt = time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE users SET email = ?, name = ?, password_hash = ?, role = ?, updated_at = ? WHERE id = ?`,
		u.Email, u.Name, u.PasswordHash, u.Role, u.UpdatedAt, u.ID)
	return err
}

func (r *UserRepo) CreatePAT(ctx context.Context, pat *models.PersonalAccessToken) error {
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
	_, err := r.db.ExecContext(ctx, `INSERT INTO personal_access_tokens (id, user_id, name, token_hash, prefix, access_level, project_scope, allowed_projects, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		pat.ID, pat.UserID, pat.Name, pat.TokenHash, pat.Prefix, pat.AccessLevel, pat.ProjectScope, pat.AllowedProjects, pat.ExpiresAt.Format(time.RFC3339), pat.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to create personal access token: %w", err)
	}
	return nil
}

func (r *UserRepo) ListPATs(ctx context.Context, userID string) ([]*models.PersonalAccessToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, name, prefix, access_level, project_scope, allowed_projects, expires_at, created_at FROM personal_access_tokens WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list personal access tokens: %w", err)
	}
	defer rows.Close()

	var list []*models.PersonalAccessToken
	for rows.Next() {
		var (
			id, uid, name, prefix, accessLevel, projectScope, createdAt string
			allowedProjects, expiresAt                                  *string
		)
		if err := rows.Scan(&id, &uid, &name, &prefix, &accessLevel, &projectScope, &allowedProjects, &expiresAt, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan pat: %w", err)
		}

		cat, _ := time.Parse(time.RFC3339, createdAt)
		var eat time.Time
		if expiresAt != nil {
			parsed, err := time.Parse(time.RFC3339, *expiresAt)
			if err == nil {
				eat = parsed
			}
		}

		list = append(list, &models.PersonalAccessToken{
			ID:              id,
			UserID:          uid,
			Name:            name,
			Prefix:          prefix,
			AccessLevel:     accessLevel,
			ProjectScope:    projectScope,
			AllowedProjects: allowedProjects,
			ExpiresAt:       eat,
			CreatedAt:       cat,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *UserRepo) DeletePAT(ctx context.Context, id, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.ExecContext(ctx, `DELETE FROM personal_access_tokens WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete personal access token: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return utils.NewNotFoundError("PersonalAccessToken", id)
	}
	return nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return utils.NewNotFoundError("User", id)
	}
	return nil
}
