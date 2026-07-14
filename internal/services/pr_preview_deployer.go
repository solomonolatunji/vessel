package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/utils"
)

type PRPreviewService struct {
	repo       repositories.PRPreviewRepository
	appService *AppService
	gitService *GitService
	deployer   *engine.Deployer
}

func NewPRPreviewService(
	repo repositories.PRPreviewRepository,
	appService *AppService,
	gitService *GitService,
	deployer *engine.Deployer,
) *PRPreviewService {
	return &PRPreviewService{
		repo:       repo,
		appService: appService,
		gitService: gitService,
		deployer:   deployer,
	}
}

func (s *PRPreviewService) DeployPRPreview(ctx context.Context, appID string, prNumber int, commitHash, branch string) (*models.PRPreview, error) {
	app, err := s.appService.GetAppService(ctx, appID)
	if err != nil || app == nil {
		return nil, errors.New("app service not found")
	}
	previewDomain := fmt.Sprintf("pr-%d.%s", prNumber, app.Domain)
	if app.Domain == "" {
		previewDomain = fmt.Sprintf("pr-%d.%s.sslip.io", prNumber, app.Name)
	}
	preview := &models.PRPreview{
		ID:            uuid.NewString(),
		ServiceID:     app.ID,
		ProjectID:     app.ProjectID,
		PRNumber:      prNumber,
		Branch:        branch,
		CommitHash:    commitHash,
		Status:        "BUILDING",
		PreviewDomain: previewDomain,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, preview); err != nil {
		return nil, err
	}
	go func() {
		bgCtx := context.Background()
		sourceDir := filepath.Join(utils.GetDataDir(), "builds", "pr-previews", preview.ID)
		clonedApp := *app
		clonedApp.Branch = branch
		if err := s.gitService.CloneOrPullAppRepository(bgCtx, &clonedApp, sourceDir, nil); err != nil {
			log.Printf("[PRPreview] failed to clone PR branch %s: %v", branch, err)
			preview.Status = "FAILED"
			_ = s.repo.Update(bgCtx, preview)
			return
		}
		clonedApp.Domain = previewDomain
		clonedApp.Name = fmt.Sprintf("%s-pr-%d", app.Name, prNumber)
		containerID, deployErr := s.deployer.DeployAppService(bgCtx, &clonedApp, sourceDir, nil)
		if deployErr != nil {
			log.Printf("[PRPreview] failed to deploy: %v", deployErr)
			preview.Status = "FAILED"
			_ = s.repo.Update(bgCtx, preview)
			return
		}
		preview.ContainerID = containerID
		preview.Status = "READY"
		_ = s.repo.Update(bgCtx, preview)
		s.updateCommitStatus(bgCtx, app, commitHash, previewDomain)
	}()
	return preview, nil
}

func (s *PRPreviewService) DestroyPRPreview(ctx context.Context, appID string, prNumber int) error {
	previews, err := s.repo.GetByAppAndPR(ctx, appID, prNumber)
	if err != nil {
		return err
	}
	for _, p := range previews {
		if p.ContainerID != "" {
			_ = s.deployer.Stop(ctx, p.ContainerID)
			_ = s.deployer.Remove(ctx, p.ContainerID)
		}
		_ = s.repo.Delete(ctx, p.ID)
	}
	return nil
}

func (s *PRPreviewService) updateCommitStatus(ctx context.Context, app *models.AppService, commitHash, previewDomain string) {
	if app.RepositoryURL == "" {
		return
	}
	repoParts := strings.Split(strings.TrimSuffix(app.RepositoryURL, ".git"), "/")
	if len(repoParts) < 2 {
		return
	}
	owner := repoParts[len(repoParts)-2]
	repo := repoParts[len(repoParts)-1]
	token := ""
	if strings.Contains(app.RepositoryURL, "@") {
		parts := strings.Split(app.RepositoryURL, "@")
		authParts := strings.Split(parts[0], "://")
		if len(authParts) == 2 {
			creds := strings.Split(authParts[1], ":")
			if len(creds) == 2 {
				token = creds[1]
			} else {
				token = creds[0]
			}
		}
	}
	if token == "" || !strings.Contains(app.RepositoryURL, "github.com") {
		return
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/statuses/%s", owner, repo, commitHash)
	payload := map[string]string{
		"state":       "success",
		"target_url":  "https://" + previewDomain,
		"description": "PR Preview is ready",
		"context":     "vessl/pr-preview",
	}
	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}
