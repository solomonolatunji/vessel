package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// AddDomain registers a new custom domain routing rule for a project.
func (s *Store) AddDomain(d *types.DomainConfig) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now

	_, err := s.db.Exec(`INSERT INTO domains (id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ProjectID, d.DomainName, d.RedirectTo, d.SSLCertStatus, d.PathPrefix, d.CreatedAt, d.UpdatedAt)
	return err
}

// ListDomains returns all custom domain configurations attached to the specified project ID.
func (s *Store) ListDomains(projectID string) ([]types.DomainConfig, error) {
	rows, err := s.db.Query(`SELECT id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at
		FROM domains WHERE project_id = ? ORDER BY domain_name ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []types.DomainConfig
	for rows.Next() {
		var d types.DomainConfig
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.DomainName, &d.RedirectTo, &d.SSLCertStatus, &d.PathPrefix, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, nil
}
