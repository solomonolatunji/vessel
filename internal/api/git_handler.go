package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/solomonolatunji/vessel/internal/types"
)

func (s *Server) handleConnectGitProvider(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req types.GitConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	gp, err := s.gitService.SaveProvider(claims.UserID, &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, gp)
}

func (s *Server) handleGetGitProvidersStatus(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	status, err := s.gitService.GetConnectedProviders(claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, status)
}

func (s *Server) handleDisconnectGitProvider(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	provider := r.PathValue("provider")
	if provider == "" {
		writeError(w, http.StatusBadRequest, "missing provider parameter")
		return
	}

	if err := s.gitService.DisconnectProvider(claims.UserID, provider); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "disconnected"})
}

func (s *Server) handleListGitRepositories(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	provider := r.URL.Query().Get("provider")
	if provider == "" {
		writeError(w, http.StatusBadRequest, "missing provider query parameter (e.g. ?provider=github)")
		return
	}

	repos, err := s.gitService.ListRepositories(r.Context(), claims.UserID, provider)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, repos)
}

func (s *Server) handleGitWebhook(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "missing projectId parameter")
		return
	}

	project, err := s.store.GetProject(projectID)
	if err != nil || project == nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status":  "accepted",
		"message": fmt.Sprintf("triggering background build & deployment for %s", project.Name),
	})

	go func() {
		ctx := context.Background()
		sourceDir := filepath.Join("data", "builds", project.ID)
		if s.gitService != nil {
			if err := s.gitService.CloneOrPullRepository(ctx, project, sourceDir, nil); err != nil {
				log.Printf("❌ [GitWebhook] Git clone/pull failed for project %s (%s): %v", project.Name, project.ID, err)
				return
			}
		}
		if s.deployer != nil {
			containerID, err := s.deployer.Deploy(ctx, project, sourceDir, nil)
			if err != nil {
				log.Printf("❌ [GitWebhook] Deployment failed for project %s (%s): %v", project.Name, project.ID, err)
				return
			}
			log.Printf("✅ [GitWebhook] Successfully rolled out container %s for project %s (%s)", containerID, project.Name, project.ID)
		}
		if s.proxyManager != nil {
			_ = s.proxyManager.Reload(ctx)
		}
	}()
}
