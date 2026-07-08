package orchestrator

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/client"
	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
	"github.com/solomonolatunji/vessel/internal/utils"
)

// Deployer orchestrates full zero-downtime application builds, secret injection, and container switchover.
type Deployer struct {
	builder          Builder
	containerManager *ContainerManager
	store            *store.Store
}

// NewDeployer initializes a Deployer wired to the container build engine, Docker lifecycle manager, and store.
func NewDeployer(dockerClient *client.Client, s *store.Store) *Deployer {
	return &Deployer{
		builder:          NewBuilder(dockerClient),
		containerManager: NewContainerManager(dockerClient),
		store:            s,
	}
}

// Deploy executes the complete deployment sequence for a given project configuration.
func (d *Deployer) Deploy(ctx context.Context, project *types.ProjectConfig, sourceDir string, logWriter io.Writer) (string, error) {
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🚀 [Deployer] Starting deployment for project: %s (ID: %s)\n", project.Name, project.ID)
	}

	buildOpts := BuildOptions{
		ProjectID:      project.ID,
		SourceDir:      sourceDir,
		DockerfilePath: project.DockerfilePath,
		LogWriter:      logWriter,
		ProjectConfig:  project,
	}

	imageTag, err := d.builder.Build(ctx, buildOpts)
	if err != nil {
		return "", fmt.Errorf("build phase failed: %w", err)
	}
	if logWriter != nil {
		fmt.Fprintf(logWriter, "✅ [Deployer] Successfully built OCI image: %s\n", imageTag)
	}

	envVarsMap, err := d.store.GetEnvVars(project.ID)
	if err != nil && logWriter != nil {
		fmt.Fprintf(logWriter, "⚠️ [Deployer] Warning: could not load environment variables: %v\n", err)
	}

	var envSlice []string
	for k, v := range envVarsMap {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}

	containerName := utils.NormalizeContainerName(project.ID)
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🔄 [Deployer] Rolling out container %s with %d encrypted environment variables...\n", containerName, len(envSlice))
	}

	containerID, err := d.containerManager.CreateAndStart(
		ctx,
		containerName,
		imageTag,
		project.InternalPort,
		envSlice,
		project.MemoryLimitMB,
		project.CPURequest,
	)
	if err != nil {
		return "", fmt.Errorf("container rollout failed: %w", err)
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "🎉 [Deployer] Deployment successful! Container ID: %s\n", containerID[:12])
	}

	return containerID, nil
}
