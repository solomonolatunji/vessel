package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/solomonolatunji/vessel/internal/types"
	"github.com/solomonolatunji/vessel/internal/utils"
)

// handleListProjects returns a JSON array of all registered applications across the platform.
func (s *Server) handleListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := s.store.ListProjects()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if projects == nil {
		projects = []types.ProjectConfig{}
	}
	writeJSON(w, http.StatusOK, projects)
}

// handleCreateProject parses project creation payloads and generates a default wildcard sslip.io domain when none is supplied.
func (s *Server) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	var p types.ProjectConfig
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid project configuration payload")
		return
	}

	if p.Name == "" {
		writeError(w, http.StatusBadRequest, "project name is required")
		return
	}

	if p.Domain == "" {
		p.Domain = utils.GenerateSslipDomain(p.Name, "")
	}
	if p.InternalPort <= 0 {
		p.InternalPort = 3000
	}

	if err := s.store.CreateProject(&p); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = s.proxyManager.Reload(r.Context())
	writeJSON(w, http.StatusCreated, p)
}

// handleGetProject retrieves the full details of a specific project by ID.
func (s *Server) handleGetProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	project, err := s.store.GetProject(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	writeJSON(w, http.StatusOK, project)
}

// handleDeleteProject removes a project record from SQLite and triggers a Caddy configuration reload.
func (s *Server) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	if err := s.store.DeleteProject(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = s.proxyManager.Reload(r.Context())
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// handleDeployProject triggers the multi-language build pipeline and container switchover for the target project.
func (s *Server) handleDeployProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	project, err := s.store.GetProject(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	sourceDir := filepath.Join("data", "builds", id)
	if s.gitService != nil && project.RepositoryURL != "" {
		if err := s.gitService.CloneOrPullRepository(r.Context(), project, sourceDir, nil); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("git checkout failed: %v", err))
			return
		}
	}

	containerID, err := s.deployer.Deploy(r.Context(), project, sourceDir, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("deployment rollout failed: %v", err))
		return
	}

	_ = s.proxyManager.Reload(r.Context())
	writeJSON(w, http.StatusOK, map[string]string{
		"status":       "deployed",
		"container_id": containerID,
	})
}
