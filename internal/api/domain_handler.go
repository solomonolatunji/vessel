package api

import (
	"encoding/json"
	"net/http"

	"github.com/solomonolatunji/vessel/internal/types"
)

// handleListDomains lists all custom domain configurations attached to a project.
func (s *Server) handleListDomains(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	domains, err := s.store.ListDomains(projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if domains == nil {
		domains = []types.DomainConfig{}
	}
	writeJSON(w, http.StatusOK, domains)
}

// handleAddDomain registers a new custom domain or sslip.io rule for a project and updates Caddy routing instantly.
func (s *Server) handleAddDomain(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	var d types.DomainConfig
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid domain configuration payload")
		return
	}

	d.ProjectID = projectID
	if d.DomainName == "" {
		writeError(w, http.StatusBadRequest, "domain_name is required")
		return
	}

	if err := s.store.AddDomain(&d); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = s.proxyManager.Reload(r.Context())
	writeJSON(w, http.StatusCreated, d)
}

// handleDeleteDomain removes a custom domain rule and reloads the Caddy reverse proxy table.
func (s *Server) handleDeleteDomain(w http.ResponseWriter, r *http.Request) {
	domainID := r.PathValue("id")
	if domainID == "" {
		writeError(w, http.StatusBadRequest, "missing domain id parameter")
		return
	}

	if err := s.store.DeleteDomain(domainID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = s.proxyManager.Reload(r.Context())
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
