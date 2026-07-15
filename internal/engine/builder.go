package engine

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/docker/docker/client"
	"golang.org/x/sync/semaphore"

	"vessl.dev/vessl/internal/models"
)

type BuildStrategy string

const (
	StrategyDockerfile BuildStrategy = "dockerfile"
	StrategyRailpack   BuildStrategy = "railpack"
	StrategyNixpacks   BuildStrategy = "nixpacks"
	StrategyBuildpacks BuildStrategy = "buildpacks"
	StrategyServerless BuildStrategy = "serverless"
)

type BuildOptions struct {
	ProjectID      string                `json:"projectId"`
	ServiceID      string                `json:"serviceId,omitempty"`
	SourceDir      string                `json:"sourceDir"`
	DockerfilePath string                `json:"dockerfilePath,omitempty"`
	LogWriter      io.Writer             `json:"-"`
	ProjectConfig  *models.ProjectConfig `json:"projectConfig,omitempty"`
	AppConfig      *models.AppService    `json:"appConfig,omitempty"`
	EnvVars        map[string]string     `json:"envVars,omitempty"`
}

type Builder interface {
	Build(ctx context.Context, opts BuildOptions) (string, error)
}

type EngineBuilder struct {
	dockerClient      *client.Client
	dockerfileBuilder *DockerfileBuilder
	railpackBuilder   *RailpackBuilder
	sem               *semaphore.Weighted
}

func defaultConcurrentBuilds() int64 {
	if s := os.Getenv("VESSL_MAX_CONCURRENT_BUILDS"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil && v > 0 {
			return v
		}
	}
	return 2
}

func NewBuilder(dockerClient *client.Client) *EngineBuilder {
	return &EngineBuilder{
		dockerClient:      dockerClient,
		dockerfileBuilder: NewDockerfileBuilder(dockerClient),
		railpackBuilder:   NewRailpackBuilder(dockerClient),
		sem:               semaphore.NewWeighted(defaultConcurrentBuilds()),
	}
}

func (b *EngineBuilder) Build(ctx context.Context, opts BuildOptions) (string, error) {
	if err := b.sem.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("failed to acquire build slot: %w", err)
	}
	defer b.sem.Release(1)

	strategy := b.DetectStrategy(opts.SourceDir, opts.DockerfilePath, opts.AppConfig)
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "🚀 [Builder] Detected build strategy: %s\n", strategy)
	}
	switch strategy {
	case StrategyDockerfile:
		imageTag, err := b.dockerfileBuilder.Build(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("dockerfile build failed: %w", err)
		}
		if opts.AppConfig != nil && opts.AppConfig.StaticOutput != "" {
			imageTag, err = b.wrapStaticOutput(ctx, imageTag, opts.AppConfig.StaticOutput, opts.LogWriter)
		}
		return imageTag, err
	case StrategyRailpack, StrategyNixpacks, StrategyBuildpacks:
		imageTag, err := b.railpackBuilder.Build(ctx, opts, string(strategy))
		if err != nil {
			return "", fmt.Errorf("%s build failed: %w", strategy, err)
		}
		if opts.AppConfig != nil && opts.AppConfig.StaticOutput != "" {
			imageTag, err = b.wrapStaticOutput(ctx, imageTag, opts.AppConfig.StaticOutput, opts.LogWriter)
		}
		return imageTag, err
	default:
		return "", fmt.Errorf("unsupported build strategy: %s", strategy)
	}
}

func (b *EngineBuilder) DetectStrategy(sourceDir, dockerfilePath string, app *models.AppService) BuildStrategy {
	if app != nil && app.BuildEngine == models.BuildEngineServerless {
		return StrategyServerless
	}

	if dockerfilePath != "" {
		if _, err := os.Stat(filepath.Join(sourceDir, dockerfilePath)); err == nil {
			return StrategyDockerfile
		}
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "Dockerfile")); err == nil {
		return StrategyDockerfile
	}

	if app != nil && app.BuildEngine != "" && app.BuildEngine != models.BuildEngineAuto {
		return BuildStrategy(app.BuildEngine)
	}

	return StrategyRailpack
}

func (b *EngineBuilder) wrapStaticOutput(ctx context.Context, baseImageTag string, staticOutputDir string, logWriter io.Writer) (string, error) {
	if logWriter != nil {
		fmt.Fprintf(logWriter, "📦 [StaticOutput] Extracting %s and wrapping in NGINX alpine...\n", staticOutputDir)
	}

	tmpDir, err := os.MkdirTemp("", "vessl-static-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	sourcePath := staticOutputDir
	if !filepath.IsAbs(sourcePath) {
		sourcePath = filepath.Join("/app", sourcePath)
	}

	dockerfileContent := fmt.Sprintf(`
FROM %s AS builder
FROM nginx:alpine
COPY --from=builder %s /usr/share/nginx/html
EXPOSE 80
`, baseImageTag, sourcePath)

	if err := os.WriteFile(filepath.Join(tmpDir, "Dockerfile"), []byte(dockerfileContent), 0644); err != nil {
		return "", err
	}

	opts := BuildOptions{
		ProjectID:      baseImageTag + "-static",
		SourceDir:      tmpDir,
		DockerfilePath: "Dockerfile",
		LogWriter:      logWriter,
	}

	return b.dockerfileBuilder.Build(ctx, opts)
}
