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
		secret = os.Getenv("VESSEL_JWT_SECRET")
	}
	if secret == "" {
		secret = "vessel-super-secret-jwt-signing-key-change-in-prod"
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
		"iss":         "vessel-auth",
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
