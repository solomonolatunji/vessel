package models

import "time"

type User struct {
	ID            string    `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	Name          string    `json:"name" db:"name"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	Role          string    `json:"role" db:"role"`
	TOTPEnabled   bool      `json:"totpEnabled" db:"totp_enabled"`
	OAuthProvider string    `json:"oauthProvider,omitempty" db:"oauth_provider"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type UserClaims struct {
	UserID      string `json:"sub"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	TOTPEnabled bool   `json:"totpEnabled"`
}

type PersonalAccessToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	TokenHash string    `json:"-" db:"token_hash"`
	Prefix    string    `json:"prefix" db:"prefix"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type UpdateProfileRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CreatePATRequest struct {
	Name string `json:"name"`
}

type CreatePATResponse struct {
	Token *PersonalAccessToken `json:"token"`
	Plain string               `json:"plain"`
}
