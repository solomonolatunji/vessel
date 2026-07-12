package license

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	WorkspaceID string `json:"team_id"`
	Plan        string `json:"plan"` // e.g. "enterprise"
	MaxSeats    int    `json:"max_seats"`
	ExpiresAt   int64  `json:"exp"`
	jwt.RegisteredClaims
}

var (
	ErrInvalidLicense = errors.New("invalid license key")
	ErrExpiredLicense = errors.New("expired license key")
)

func GenerateLicense(privateKeyBase64, workspaceID, plan string, maxSeats int, expiry time.Time) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", fmt.Errorf("invalid private key encoding: %w", err)
	}

	if len(keyBytes) != ed25519.PrivateKeySize {
		return "", errors.New("invalid private key size")
	}

	privKey := ed25519.PrivateKey(keyBytes)

	claims := Claims{
		WorkspaceID: workspaceID,
		Plan:        plan,
		MaxSeats:    maxSeats,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			Issuer:    "vessl-cloud",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(privKey)
}

func VerifyLicense(publicKeyBase64, licenseKey string) (*Claims, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("invalid public key encoding: %w", err)
	}

	if len(keyBytes) != ed25519.PublicKeySize {
		return nil, errors.New("invalid public key size")
	}

	pubKey := ed25519.PublicKey(keyBytes)

	token, err := jwt.ParseWithClaims(licenseKey, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredLicense
		}
		return nil, ErrInvalidLicense
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidLicense
}
