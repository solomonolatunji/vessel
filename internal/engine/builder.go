package orchestrator

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"vessel.dev/vessel/internal/project"
	"vessel.dev/vessel/internal/service"
)

type BuildStrategy string

const (
	StrategyDockerfile BuildStrategy = "dockerfile"
	StrategyRailpack   BuildStrategy = "railpack"
)

type BuildOptions struct {
	ProjectID      string                 `json:"projectId"`
	ServiceID      string                 `json:"serviceId,omitempty"`
	SourceDir      string                 `json:"sourceDir"`
	DockerfilePath string                 `json:"dockerfilePath,omitempty"`
	LogWriter      io.Writer              `json:"-"`
	ProjectConfig  *project.ProjectConfig `json:"projectConfig,omitempty"`
	AppConfig      *service.AppService    `json:"appConfig,omitempty"`
}

type Builder interface {
	Build(ctx context.Context, opts BuildOptions) (string, error)
}

type EngineBuilder struct {
	dockerClient      *client.Client
	dockerfileBuilder *DockerfileBuilder
	railpackBuilder   *RailpackBuilder
}

func NewBuilder(dockerClient *client.Client) *EngineBuilder {
	return &EngineBuilder{
		dockerClient:      dockerClient,
		dockerfileBuilder: NewDockerfileBuilder(dockerClient),
		railpackBuilder:   NewRailpackBuilder(dockerClient),
	}
}

func (b *EngineBuilder) Build(ctx context.Context, opts BuildOptions) (string, error) {
	strategy := b.DetectStrategy(opts.SourceDir, opts.DockerfilePath)
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "🚀 [Builder] Detected build strategy: %s\n", strategy)
	}

	switch strategy {
	case StrategyDockerfile:
		imageTag, err := b.dockerfileBuilder.Build(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("dockerfile build failed: %w", err)
		}
		return imageTag, nil
	case StrategyRailpack:
		imageTag, err := b.railpackBuilder.Build(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("railpack/nixpacks build failed: %w", err)
		}
		return imageTag, nil
	default:
		return "", fmt.Errorf("unsupported build strategy: %s", strategy)
	}
}

func (b *EngineBuilder) DetectStrategy(sourceDir, dockerfilePath string) BuildStrategy {
	if dockerfilePath != "" {
		if _, err := os.Stat(filepath.Join(sourceDir, dockerfilePath)); err == nil {
			return StrategyDockerfile
		}
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "Dockerfile")); err == nil {
		return StrategyDockerfile
	}
	return StrategyRailpack
}
