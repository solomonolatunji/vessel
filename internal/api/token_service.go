package api

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/solomonolatunji/vessel/internal/types"
)

// UserClaims defines the structured JWT claims embedded in authentication tokens.
type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// TokenService manages the cryptographic signing and validation of JSON Web Tokens.
type TokenService struct {
	secret []byte
}

// NewTokenService initializes a TokenService using the VESSEL_JWT_SECRET environment variable or a default signing key.
func NewTokenService() *TokenService {
	secret := os.Getenv("VESSEL_JWT_SECRET")
	if secret == "" {
		secret = "vessel-default-insecure-dev-secret-key-256bit"
	}
	return &TokenService{
		secret: []byte(secret),
	}
}

// GenerateToken creates a signed JWT string valid for 72 hours for an authenticated user session.
func (ts *TokenService) GenerateToken(u *types.User) (string, error) {
	claims := UserClaims{
		UserID: u.ID,
		Email:  u.Email,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "vessel-control-plane",
			Subject:   u.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secret)
}

// ValidateToken parses and verifies a JWT string, returning the extracted UserClaims if valid.
func (ts *TokenService) ValidateToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected token signing algorithm")
		}
		return ts.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired authorization token")
	}
	return claims, nil
}
