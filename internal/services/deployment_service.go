package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type DeploymentService struct {
	repo         repositories.DeploymentRepository
	appRepo      repositories.AppServiceRepository
	projectRepo  repositories.ProjectRepository
	deployer     *engine.Deployer
	gitService   *GitService
	statsMonitor *engine.StatsMonitor
}

func NewDeploymentService(
	r repositories.DeploymentRepository,
	ar repositories.AppServiceRepository,
	pr repositories.ProjectRepository,
	d *engine.Deployer,
	gs *GitService,
	sm *engine.StatsMonitor,
) *DeploymentService {
	return &DeploymentService{
		repo:         r,
		appRepo:      ar,
		projectRepo:  pr,
		deployer:     d,
		gitService:   gs,
		statsMonitor: sm,
	}
}

func (s *DeploymentService) CreateDeployment(ctx context.Context, d *models.Deployment) (*models.Deployment, error) {
	if d == nil || d.ServiceID == "" {
		return nil, errors.New("valid deployment with serviceId required")
	}
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	if d.Status == "" {
		d.Status = models.DeploymentStatusPending
	}
	now := time.Now()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	d.UpdatedAt = now
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DeploymentService) GetDeployment(ctx context.Context, id string) (*models.Deployment, error) {
	if id == "" {
		return nil, errors.New("deployment id required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *DeploymentService) ListByService(ctx context.Context, serviceID string, limit, offset int) ([]*models.Deployment, int, error) {
	if serviceID == "" {
		return nil, 0, errors.New("service id required")
	}
	return s.repo.ListByService(ctx, serviceID, limit, offset)
}

func (s *DeploymentService) UpdateDeployment(ctx context.Context, d *models.Deployment) error {
	if d == nil || d.ID == "" {
		return errors.New("valid deployment required for update")
	}
	d.UpdatedAt = time.Now()
	return s.repo.Update(ctx, d)
}

type DeployStatusOpts struct {
	ID          string
	Status      models.DeploymentStatus
	BuildLogs   string
	ContainerID string
}

func (s *DeploymentService) UpdateStatus(ctx context.Context, opts DeployStatusOpts) error {
	if opts.ID == "" {
		return errors.New("deployment id required")
	}
	return s.repo.UpdateStatus(ctx, opts.ID, opts.Status, opts.BuildLogs, opts.ContainerID)
}

func (s *DeploymentService) ExecuteDeploymentAsync(d *models.Deployment) {
	go func() {
		bgCtx := context.Background()
		if s.deployer == nil || s.appRepo == nil || s.gitService == nil {
			_ = s.UpdateStatus(bgCtx, DeployStatusOpts{ID: d.ID, Status: models.DeploymentStatusFailed, BuildLogs: "Deployment dependencies missing\n", ContainerID: ""})
			return
		}

		app, err := s.appRepo.GetByID(bgCtx, d.ServiceID)
		if err != nil {
			_ = s.UpdateStatus(bgCtx, DeployStatusOpts{ID: d.ID, Status: models.DeploymentStatusFailed, BuildLogs: fmt.Sprintf("Failed to get app service: %v\n", err), ContainerID: ""})
			return
		}

		if app.ImageRef != "" {
			d.Status = models.DeploymentStatusPulling
			_ = s.repo.Update(bgCtx, d)

			containerID, err := s.deployer.DeployImage(bgCtx, app, nil)
			if err != nil {
				_ = s.UpdateStatus(bgCtx, DeployStatusOpts{ID: d.ID, Status: models.DeploymentStatusFailed, BuildLogs: fmt.Sprintf("Image deploy failed: %v\n", err), ContainerID: ""})
				return
			}

			_ = s.UpdateStatus(bgCtx, DeployStatusOpts{ID: d.ID, Status: models.DeploymentStatusReady, BuildLogs: "Deployment succeeded.\n", ContainerID: containerID})
			app.ContainerID = containerID
			_ = s.appRepo.Update(bgCtx, app)
			return
		}

		sourceDir := fmt.Sprintf("data/builds/%s/%s", app.ID, d.ID)

		d.Status = models.DeploymentStatusCloning
		_ = s.repo.Update(bgCtx, d)

		if err := s.gitService.CloneOrPullAppRepository(bgCtx, app, sourceDir, nil); err != nil {
			_ = s.UpdateStatus(bgCtx, DeployStatusOpts{ID: d.ID, Status: models.DeploymentStatusFailed, BuildLogs: fmt.Sprintf("Git clone failed: %v\n", err), ContainerID: ""})
			return
		}

		app.Icon = detectAppIcon(sourceDir)
		_ = s.appRepo.Update(bgCtx, app)

		d.Status = models.DeploymentStatusBuilding
		_ = s.repo.Update(bgCtx, d)

		containerID, err := s.deployer.DeployAppService(bgCtx, app, sourceDir, nil)
		if err != nil {
			_ = s.UpdateStatus(bgCtx, DeployStatusOpts{ID: d.ID, Status: models.DeploymentStatusFailed, BuildLogs: fmt.Sprintf("Deployment failed: %v\n", err), ContainerID: ""})
			return
		}

		_ = s.UpdateStatus(bgCtx, DeployStatusOpts{ID: d.ID, Status: models.DeploymentStatusReady, BuildLogs: "Deployment succeeded.\n", ContainerID: containerID})

		app.ContainerID = containerID
		_ = s.appRepo.Update(bgCtx, app)
	}()
}

func (s *DeploymentService) DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error) {
	if s.deployer == nil || s.appRepo == nil {
		return "", errors.New("deployer or app repo not available")
	}
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return "", err
	}
	containerID, err := s.deployer.DeployAppService(ctx, app, sourceDir, logWriter)
	if err == nil && containerID != "" {
		app.ContainerID = containerID
		_ = s.appRepo.Update(ctx, app)
	}
	return containerID, err
}

func (s *DeploymentService) GetMetrics(ctx context.Context, appID string) (*engine.ContainerHealth, error) {
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return nil, err
	}
	if app.ContainerID == "" {
		return &engine.ContainerHealth{Status: engine.ContainerHealthStatusNotDeployed}, nil
	}
	if s.statsMonitor == nil {
		return nil, errors.New("stats monitor not available")
	}
	return s.statsMonitor.GetHealth(ctx, app.ContainerID)
}

func detectAppIcon(sourceDir string) string {
	// Check package.json dependencies first if it exists
	if b, err := os.ReadFile(filepath.Join(sourceDir, "package.json")); err == nil {
		s := strings.ToLower(string(b))
		if strings.Contains(s, "\"next\"") {
			return "nextjs"
		}
		if strings.Contains(s, "\"nuxt\"") {
			return "nuxt"
		}
		if strings.Contains(s, "\"@sveltejs/kit\"") {
			return "sveltekit"
		}
		if strings.Contains(s, "\"@solidjs/start\"") {
			return "solidstart"
		}
		if strings.Contains(s, "\"@remix-run") {
			return "remix"
		}
		if strings.Contains(s, "\"astro\"") {
			return "astro"
		}
		if strings.Contains(s, "\"@tanstack/start\"") {
			return "tanstack"
		}
		if strings.Contains(s, "\"expo\"") {
			return "expo"
		}
		if strings.Contains(s, "\"react-native\"") {
			return "react-native"
		}
		if strings.Contains(s, "\"vite\"") {
			return "vite"
		}
		if strings.Contains(s, "\"vue\"") {
			return "vue"
		}
		if strings.Contains(s, "\"svelte\"") {
			return "svelte"
		}
		if strings.Contains(s, "\"@angular/core\"") {
			return "angular"
		}
		if strings.Contains(s, "\"@builder.io/qwik\"") {
			return "qwik"
		}
		if strings.Contains(s, "\"gatsby\"") {
			return "gatsby"
		}
		if strings.Contains(s, "\"@redwoodjs/") {
			return "redwoodjs"
		}
		if strings.Contains(s, "\"electron\"") {
			return "electron"
		}
		if strings.Contains(s, "\"@tauri-apps/") {
			return "tauri"
		}
		if strings.Contains(s, "\"hono\"") {
			return "hono"
		}
		if strings.Contains(s, "\"elysia\"") {
			return "elysia"
		}
		if strings.Contains(s, "\"@nestjs/core\"") {
			return "nestjs"
		}
		if strings.Contains(s, "\"fastify\"") {
			return "fastify"
		}
		if strings.Contains(s, "\"express\"") {
			return "express"
		}
		if strings.Contains(s, "\"koa\"") {
			return "koa"
		}
		if strings.Contains(s, "\"@adonisjs/") {
			return "adonisjs"
		}
		if strings.Contains(s, "\"@strapi/") {
			return "strapi"
		}
		if strings.Contains(s, "\"payload\"") {
			return "payload"
		}
		if strings.Contains(s, "\"@trpc/") {
			return "trpc"
		}
		if strings.Contains(s, "\"graphql\"") {
			return "graphql"
		}
		if strings.Contains(s, "\"react\"") {
			return "react"
		}
	}

	// Fallback to file-based detection
	files := map[string]string{
		"next.config.js": "nextjs", "next.config.mjs": "nextjs",
		"vite.config.ts": "vite", "vite.config.js": "vite",
		"astro.config.mjs": "astro", "nuxt.config.ts": "nuxt",
		"nest-cli.json": "nestjs", "svelte.config.js": "sveltekit",
		"remix.config.js": "remix", "angular.json": "angular",
		"pom.xml": "java", "build.gradle": "java",
		"manage.py": "django", "Gemfile": "ruby",
		"composer.json": "php", "Cargo.toml": "rust",
		"go.mod": "golang", "artisan": "laravel",
	}

	for file, icon := range files {
		if _, err := os.Stat(filepath.Join(sourceDir, file)); err == nil {
			return icon
		}
	}

	// Specific text content checks
	if b, err := os.ReadFile(filepath.Join(sourceDir, "requirements.txt")); err == nil {
		s := strings.ToLower(string(b))
		if strings.Contains(s, "fastapi") {
			return "fastapi"
		}
		if strings.Contains(s, "flask") {
			return "flask"
		}
		if strings.Contains(s, "django") {
			return "django"
		}
		return "python"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "package.json")); err == nil {
		return "nodejs"
	}

	return "git"
}
