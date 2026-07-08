package orchestrator

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/solomonolatunji/vessel/internal/types"
)

// BuildStrategy identifies the container compilation mechanism chosen for the project.
type BuildStrategy string

const (
	// StrategyDockerfile indicates building directly from an existing Dockerfile in the repository.
	StrategyDockerfile BuildStrategy = "dockerfile"
	// StrategyRailpack indicates zero-configuration build auto-detection using Railpack or Nixpacks.
	StrategyRailpack BuildStrategy = "railpack"
)

// BuildOptions contains options required by builders to generate an OCI image.
type BuildOptions struct {
	ProjectID      string              `json:"projectId"`
	SourceDir      string              `json:"sourceDir"`
	DockerfilePath string              `json:"dockerfilePath,omitempty"`
	LogWriter      io.Writer           `json:"-"`
	ProjectConfig  *types.ProjectConfig `json:"projectConfig,omitempty"`
}

// Builder defines the interface required for any container image builder strategy.
type Builder interface {
	Build(ctx context.Context, opts BuildOptions) (string, error)
}

// EngineBuilder coordinates container image builds and delegates to specific strategy implementations.
type EngineBuilder struct {
	dockerClient       *client.Client
	dockerfileBuilder  *DockerfileBuilder
	railpackBuilder    *RailpackBuilder
}

// NewBuilder instantiates a new EngineBuilder wired to the provided Docker client without global state.
func NewBuilder(dockerClient *client.Client) *EngineBuilder {
	return &EngineBuilder{
		dockerClient:      dockerClient,
		dockerfileBuilder: NewDockerfileBuilder(dockerClient),
		railpackBuilder:   NewRailpackBuilder(dockerClient),
	}
}

// Build inspects the source directory and dispatches to the appropriate builder strategy.
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

// DetectStrategy examines the project filesystem to determine if a Dockerfile or Railpack should be used.
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
