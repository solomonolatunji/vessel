package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/utils"
)

type PRPreviewService struct {
	repo        repositories.PRPreviewRepository
	appService  *AppService
	gitService  *GitService
	deployer    *engine.Deployer
	workerHub   *engine.WorkerHub
	projectRepo repositories.ProjectRepository
}

func NewPRPreviewService(
	repo repositories.PRPreviewRepository,
	appService *AppService,
	gitService *GitService,
	deployer *engine.Deployer,
	workerHub *engine.WorkerHub,
	projectRepo repositories.ProjectRepository,
) *PRPreviewService {
	return &PRPreviewService{
		repo:        repo,
		appService:  appService,
		gitService:  gitService,
		deployer:    deployer,
		workerHub:   workerHub,
		projectRepo: projectRepo,
	}
}

func (s *PRPreviewService) ListByApp(ctx context.Context, appID string) ([]*models.PRPreview, error) {
	return s.repo.GetByApp(ctx, appID)
}

type DeployPRPreviewOpts struct {
	AppID      string
	PRNumber   int
	CommitHash string
	Branch     string
}

func (s *PRPreviewService) DeployPRPreview(ctx context.Context, opts DeployPRPreviewOpts) (*models.PRPreview, error) {
	app, err := s.appService.GetAppService(ctx, opts.AppID)
	if err != nil || app == nil {
		return nil, utils.NewNotFoundError("AppService", opts.AppID)
	}
	previewDomain := fmt.Sprintf("pr-%d.%s", opts.PRNumber, app.Domain)
	if app.Domain == "" {
		magicDomain := os.Getenv("CODEDOCK_MAGIC_DOMAIN")
		if magicDomain == "" {
			magicDomain = "sslip.io"
		}
		previewDomain = fmt.Sprintf("pr-%d.%s.%s", opts.PRNumber, app.Name, magicDomain)
	}
	preview := &models.PRPreview{
		ID:            uuid.NewString(),
		ServiceID:     app.ID,
		ProjectID:     app.ProjectID,
		PRNumber:      opts.PRNumber,
		Branch:        opts.Branch,
		CommitHash:    opts.CommitHash,
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
		clonedApp.ID = preview.ID
		clonedApp.Branch = opts.Branch
		if err := s.gitService.CloneOrPullAppRepository(bgCtx, &clonedApp, sourceDir, nil); err != nil {
			slog.Warn("PR preview clone failed", "branch", opts.Branch, "err", err)
			preview.Status = "FAILED"
			_ = s.repo.Update(bgCtx, preview)
			return
		}
		clonedApp.Domain = previewDomain
		clonedApp.Name = fmt.Sprintf("%s-pr-%d", app.Name, opts.PRNumber)
		containerID, deployErr := s.deployer.DeployAppService(bgCtx, &clonedApp, sourceDir, nil)
		if deployErr != nil {
			slog.Warn("PR preview deploy failed", "err", deployErr)
			preview.Status = "FAILED"
			_ = s.repo.Update(bgCtx, preview)
			return
		}
		preview.ContainerID = containerID
		preview.Status = "READY"
		_ = s.repo.Update(bgCtx, preview)
		s.updateCommitStatus(bgCtx, app, opts.CommitHash, previewDomain)
	}()
	return preview, nil
}

func (s *PRPreviewService) DestroyPRPreview(ctx context.Context, appID string, prNumber int) error {
	previews, err := s.repo.GetByAppAndPR(ctx, appID, prNumber)
	if err != nil {
		return err
	}
	for _, p := range previews {
		_ = s.deployer.StopAppService(ctx, p.ID)
		sourceDir := filepath.Join(utils.GetDataDir(), "builds", "pr-previews", p.ID)
		_ = os.RemoveAll(sourceDir)
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
		"context":     "codedock/pr-preview",
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
