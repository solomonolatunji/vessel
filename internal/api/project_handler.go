package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
	"github.com/solomonolatunji/vessel/internal/utils"
)

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

// handleCreateProject parses project creation payloads and generates an initial service or auto-named project.
func (s *Server) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID                 string `json:"id"`
		TeamID             string `json:"teamId,omitempty"`
		Name               string `json:"name"`
		Description        string `json:"description,omitempty"`
		RepositoryURL      string `json:"repositoryUrl,omitempty"`
		RepositoryURLSnake string `json:"repository_url,omitempty"`
		Branch             string `json:"branch,omitempty"`
		InternalPort       int    `json:"internalPort,omitempty"`
		InternalPortSnake  int    `json:"internal_port,omitempty"`
		Domain             string `json:"domain,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid project configuration payload")
		return
	}

	if req.Name == "" {
		req.Name = fmt.Sprintf("project-%s", uuid.NewString()[:8])
	}

	p := &types.ProjectConfig{
		ID:          req.ID,
		TeamID:      req.TeamID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.store.CreateProject(p); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	port := req.InternalPort
	if port <= 0 {
		port = req.InternalPortSnake
	}
	if port <= 0 {
		port = 3000
	}

	repo := req.RepositoryURL
	if repo == "" {
		repo = req.RepositoryURLSnake
	}

	domain := req.Domain
	if domain == "" {
		domain = utils.GenerateSslipDomain(req.Name, "")
	}

	branch := req.Branch
	if branch == "" {
		branch = "main"
	}

	envs, _ := s.store.ListEnvironments(p.ID)
	envID := "env-prod"
	if len(envs) > 0 {
		envID = envs[0].ID
	}

	app := &types.AppServiceConfig{
		ProjectID:     p.ID,
		EnvironmentID: envID,
		Name:          req.Name,
		RepositoryURL: repo,
		Branch:        branch,
		InternalPort:  port,
		Domain:        domain,
	}
	_ = s.store.CreateAppService(app)

	_ = s.proxyManager.Reload(r.Context())
	writeJSON(w, http.StatusCreated, p)
}

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
	if s.gitService != nil {
		_ = s.gitService.CloneOrPullRepository(r.Context(), project, sourceDir, nil)
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
