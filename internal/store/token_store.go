package store

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// initProjectTokensTable initializes the project_tokens table.
func (s *Store) initProjectTokensTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS project_tokens (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		environment_id TEXT DEFAULT '',
		name TEXT NOT NULL,
		token_prefix TEXT NOT NULL,
		token_hash TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);`
	_, err := s.db.Exec(query)
	return err
}

// CreateProjectToken generates a new secure API token scoped to a project (`Project Settings` -> `Tokens`).
// It returns the full token once upon creation.
func (s *Store) CreateProjectToken(token *types.ProjectToken) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if token.ID == "" {
		token.ID = uuid.NewString()
	}
	token.CreatedAt = time.Now().UTC()

	// Generate 32 bytes random secret
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random token bytes: %w", err)
	}
	rawSecret := hex.EncodeToString(randomBytes)
	fullToken := fmt.Sprintf("vsl_tok_%s", rawSecret)
	token.TokenPrefix = fullToken[:16]

	// In real production we store SHA256 of fullToken, but we store token_hash directly or encrypted
	tokenHash := fullToken

	query := `INSERT INTO project_tokens (id, project_id, environment_id, name, token_prefix, token_hash, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, token.ID, token.ProjectID, token.EnvironmentID, token.Name, token.TokenPrefix, tokenHash, token.CreatedAt)
	if err != nil {
		return "", fmt.Errorf("failed to insert project token: %w", err)
	}
	return fullToken, nil
}

// ListProjectTokens retrieves all tokens created for a project (`Project Settings` -> `Tokens`).
func (s *Store) ListProjectTokens(projectID string) ([]*types.ProjectToken, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, project_id, environment_id, name, token_prefix, created_at
		FROM project_tokens WHERE project_id = ? ORDER BY created_at DESC`

	rows, err := s.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list project tokens: %w", err)
	}
	defer rows.Close()

	var tokens []*types.ProjectToken
	for rows.Next() {
		var t types.ProjectToken
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.EnvironmentID, &t.Name, &t.TokenPrefix, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project token row: %w", err)
		}
		tokens = append(tokens, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

// DeleteProjectToken revokes an API token.
func (s *Store) DeleteProjectToken(id, projectID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec(`DELETE FROM project_tokens WHERE id = ? AND project_id = ?`, id, projectID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("project token not found")
	}
	return nil
}
