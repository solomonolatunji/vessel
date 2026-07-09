package store

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// initProjectWebhooksTable initializes the project_webhooks table.
func (s *Store) initProjectWebhooksTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS project_webhooks (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		url TEXT NOT NULL,
		event_types TEXT NOT NULL,
		include_pr_environments BOOLEAN DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`
	_, err := s.db.Exec(query)
	return err
}

// CreateProjectWebhook registers a new webhook notification endpoint (`Project Settings` -> `Webhooks`).
func (s *Store) CreateProjectWebhook(w *types.ProjectWebhook) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now

	eventTypesStr := strings.Join(w.EventTypes, ",")
	query := `INSERT INTO project_webhooks (id, project_id, url, event_types, include_pr_environments, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, w.ID, w.ProjectID, w.URL, eventTypesStr, w.IncludePREnvironments, w.CreatedAt, w.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create project webhook: %w", err)
	}
	return nil
}

// ListProjectWebhooks returns all webhooks registered for a project (`Project Settings` -> `Webhooks`).
func (s *Store) ListProjectWebhooks(projectID string) ([]*types.ProjectWebhook, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, project_id, url, event_types, include_pr_environments, created_at, updated_at
		FROM project_webhooks WHERE project_id = ? ORDER BY created_at DESC`

	rows, err := s.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list project webhooks: %w", err)
	}
	defer rows.Close()

	var webhooks []*types.ProjectWebhook
	for rows.Next() {
		var w types.ProjectWebhook
		var eventsStr string
		var includePr int
		if err := rows.Scan(&w.ID, &w.ProjectID, &w.URL, &eventsStr, &includePr, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project webhook row: %w", err)
		}
		if eventsStr != "" {
			w.EventTypes = strings.Split(eventsStr, ",")
		} else {
			w.EventTypes = []string{}
		}
		w.IncludePREnvironments = includePr == 1
		webhooks = append(webhooks, &w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return webhooks, nil
}

// DeleteProjectWebhook removes a registered webhook from a project.
func (s *Store) DeleteProjectWebhook(id, projectID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec(`DELETE FROM project_webhooks WHERE id = ? AND project_id = ?`, id, projectID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("project webhook not found")
	}
	return nil
}
