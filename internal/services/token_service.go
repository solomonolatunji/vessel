package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"vessl.dev/vessl/internal/models"
)

type TokenService struct {
	secretKey []byte
}

func NewTokenService() *TokenService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = os.Getenv("VESSL_JWT_SECRET")
	}
	if secret == "" {
		secret = "vessl-super-secret-jwt-signing-key-change-in-prod"
	}
	return &TokenService{
		secretKey: []byte(secret),
	}
}

func (ts *TokenService) GenerateToken(u *models.User) (string, error) {
	if u == nil {
		return "", errors.New("user cannot be nil when generating token")
	}
	claims := jwt.MapClaims{
		"sub":         u.ID,
		"email":       u.Email,
		"role":        u.Role,
		"totpEnabled": u.TOTPEnabled,
		"exp":         time.Now().Add(72 * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
		"iss":         "vessl-auth",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secretKey)
}

func (ts *TokenService) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return ts.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims or signature")
	}
	return claims, nil
}

func (ts *TokenService) GeneratePasswordResetToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
		"iss":   "vessl-password-reset",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secretKey)
}

func (ts *TokenService) ValidatePasswordResetToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return ts.secretKey, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims or signature")
	}

	if iss, ok := claims["iss"].(string); !ok || iss != "vessl-password-reset" {
		return "", errors.New("invalid token issuer")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", errors.New("invalid token claims")
	}
	return email, nil
}
