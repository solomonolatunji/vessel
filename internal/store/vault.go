package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Vault handles AES-256-GCM encryption and decryption for sensitive project secrets.
type Vault struct {
	key []byte
}

// NewVault initializes or loads the 256-bit encryption key stored in the persistent data directory.
func NewVault(dataDir string) (*Vault, error) {
	keyPath := filepath.Join(dataDir, ".vault_key")

	if keyData, err := os.ReadFile(keyPath); err == nil {
		if len(keyData) == 32 {
			return &Vault{key: keyData}, nil
		}
	}

	newKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, newKey); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(keyPath, newKey, 0600); err != nil {
		return nil, err
	}

	return &Vault{key: newKey}, nil
}

// Encrypt locks plaintext using AES-256-GCM and returns a base64-encoded ciphertext with prepended nonce.
func (v *Vault) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decodes a base64 ciphertext and unlocks it with AES-256-GCM.
func (v *Vault) Decrypt(encrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(v.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
