package repos

import (
	"context"
	"database/sql"
	"time"

	"vessel.dev/vessel/internal/models"
)

// AuthRepo defines the data-access interface for cloud user authentication.
type AuthRepo interface {
	CreateUser(ctx context.Context, user *models.CloudUser) error
	GetUserByEmail(ctx context.Context, email string) (*models.CloudUser, error)
	GetUserByID(ctx context.Context, id string) (*models.CloudUser, error)
	GetUserByVerifyToken(ctx context.Context, token string) (*models.CloudUser, error)
	SaveOTP(ctx context.Context, userID, code string, expiresAt time.Time) error
	ClearOTP(ctx context.Context, userID string) error
	UpdatePassword(ctx context.Context, userID, hash string) error
	MarkEmailVerified(ctx context.Context, userID string) error
	SaveVerifyToken(ctx context.Context, userID, token string, expiresAt time.Time) error
}

type authRepo struct {
	db *sql.DB
}

// NewAuthRepo creates an AuthRepo backed by a *sql.DB.
func NewAuthRepo(db *sql.DB) AuthRepo {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(ctx context.Context, user *models.CloudUser) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO cloud_users (id, email, full_name, password_hash, role, email_verified, verify_token, verify_token_expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID, user.Email, user.FullName, user.PasswordHash, user.Role, user.EmailVerified,
		nullString(user.VerifyToken), nullTime(user.VerifyTokenExpiresAt),
	)
	return err
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*models.CloudUser, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, full_name, password_hash, role, email_verified, verified_at,
			verify_token, verify_token_expires_at, otp_code, otp_expires_at, created_at
		FROM cloud_users WHERE email = $1`, email)
	return scanUser(row)
}

func (r *authRepo) GetUserByID(ctx context.Context, id string) (*models.CloudUser, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, full_name, password_hash, role, email_verified, verified_at,
			verify_token, verify_token_expires_at, otp_code, otp_expires_at, created_at
		FROM cloud_users WHERE id = $1`, id)
	return scanUser(row)
}

func (r *authRepo) GetUserByVerifyToken(ctx context.Context, token string) (*models.CloudUser, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, full_name, password_hash, role, email_verified, verified_at,
			verify_token, verify_token_expires_at, otp_code, otp_expires_at, created_at
		FROM cloud_users WHERE verify_token = $1`, token)
	return scanUser(row)
}

func (r *authRepo) SaveOTP(ctx context.Context, userID, code string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE cloud_users SET otp_code = $1, otp_expires_at = $2 WHERE id = $3`,
		code, expiresAt, userID)
	return err
}

func (r *authRepo) ClearOTP(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE cloud_users SET otp_code = NULL, otp_expires_at = NULL WHERE id = $1`,
		userID)
	return err
}

func (r *authRepo) UpdatePassword(ctx context.Context, userID, hash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE cloud_users SET password_hash = $1 WHERE id = $2`,
		hash, userID)
	return err
}

func (r *authRepo) MarkEmailVerified(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE cloud_users SET email_verified = TRUE, verified_at = NOW(), verify_token = NULL, verify_token_expires_at = NULL WHERE id = $1`,
		userID)
	return err
}

func (r *authRepo) SaveVerifyToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE cloud_users SET verify_token = $1, verify_token_expires_at = $2 WHERE id = $3`,
		token, expiresAt, userID)
	return err
}

// scanUser scans a single SQL row into a CloudUser. Returns nil, nil on ErrNoRows.
func scanUser(row *sql.Row) (*models.CloudUser, error) {
	var u models.CloudUser
	var verifiedAt sql.NullTime
	var verifyToken sql.NullString
	var verifyTokenExpiresAt sql.NullTime
	var otpCode sql.NullString
	var otpExpiresAt sql.NullTime

	err := row.Scan(
		&u.ID, &u.Email, &u.FullName, &u.PasswordHash, &u.Role, &u.EmailVerified,
		&verifiedAt, &verifyToken, &verifyTokenExpiresAt, &otpCode, &otpExpiresAt, &u.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if verifiedAt.Valid {
		u.VerifiedAt = &verifiedAt.Time
	}
	if verifyToken.Valid {
		u.VerifyToken = verifyToken.String
	}
	if verifyTokenExpiresAt.Valid {
		u.VerifyTokenExpiresAt = &verifyTokenExpiresAt.Time
	}
	if otpCode.Valid {
		u.OTPCode = otpCode.String
	}
	if otpExpiresAt.Valid {
		u.OTPExpiresAt = &otpExpiresAt.Time
	}

	return &u, nil
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}
