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
	EnvProvider      func(projectID string) (map[string]string, error)
}

// NewDeployer initializes a Deployer wired to the container build engine, Docker lifecycle manager, and store.
func NewDeployer(dockerClient *client.Client, s *store.Store) *Deployer {
	return &Deployer{
		builder:          NewBuilder(dockerClient),
		containerManager: NewContainerManager(dockerClient),
		store:            s,
	}
}

// Deploy executes the complete deployment sequence for a given project configuration or its primary application service.
func (d *Deployer) Deploy(ctx context.Context, project *types.ProjectConfig, sourceDir string, logWriter io.Writer) (string, error) {
	apps, err := d.store.ListAppServicesByProject(project.ID)
	if err == nil && len(apps) > 0 {
		return d.DeployAppService(ctx, apps[0], sourceDir, logWriter)
	}

	syntheticApp := &types.AppServiceConfig{
		ID:            project.ID,
		ProjectID:     project.ID,
		Name:          project.Name,
		InternalPort:  3000,
		MemoryLimitMB: 512,
		CPURequest:    0.5,
	}
	return d.DeployAppService(ctx, syntheticApp, sourceDir, logWriter)
}

// DeployAppService executes the complete zero-downtime deployment sequence for a specific application service container.
func (d *Deployer) DeployAppService(ctx context.Context, app *types.AppServiceConfig, sourceDir string, logWriter io.Writer) (string, error) {
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🚀 [Deployer] Starting deployment for service: %s (ID: %s)\n", app.Name, app.ID)
	}

	buildOpts := BuildOptions{
		ProjectID:      app.ProjectID,
		ServiceID:      app.ID,
		SourceDir:      sourceDir,
		DockerfilePath: app.DockerfilePath,
		LogWriter:      logWriter,
		AppConfig:      app,
	}

	imageTag, err := d.builder.Build(ctx, buildOpts)
	if err != nil {
		return "", fmt.Errorf("build phase failed: %w", err)
	}
	if logWriter != nil {
		fmt.Fprintf(logWriter, "✅ [Deployer] Successfully built OCI image: %s\n", imageTag)
	}

	envVarsMap, err := d.store.GetEnvVars(app.ProjectID)
	if err != nil && logWriter != nil {
		fmt.Fprintf(logWriter, "⚠️ [Deployer] Warning: could not load shared project environment variables: %v\n", err)
	}
	if envVarsMap == nil {
		envVarsMap = make(map[string]string)
	}

	// Merge service-specific variables over shared project variables
	serviceVars, _ := d.store.ListServiceVariables(app.ID)
	for _, sv := range serviceVars {
		envVarsMap[sv.Key] = sv.Value
	}

	if d.EnvProvider != nil {
		if linkedEnvs, err := d.EnvProvider(app.ProjectID); err == nil {
			for k, v := range linkedEnvs {
				if _, exists := envVarsMap[k]; !exists {
					envVarsMap[k] = v
				}
			}
			if logWriter != nil && len(linkedEnvs) > 0 {
				fmt.Fprintf(logWriter, "🔗 [Deployer] Automatically linked %d service connection strings (DATABASE_URL, REDIS_URL, etc.)\n", len(linkedEnvs))
			}
		}
	}

	var envSlice []string
	for k, v := range envVarsMap {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}

	containerName := utils.NormalizeContainerName(app.ID)
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🔄 [Deployer] Rolling out container %s with %d encrypted environment variables...\n", containerName, len(envSlice))
	}

	port := app.InternalPort
	if port <= 0 {
		port = 3000
	}
	memMB := app.MemoryLimitMB
	if memMB <= 0 {
		memMB = 512
	}
	cpuReq := app.CPURequest
	if cpuReq <= 0 {
		cpuReq = 0.5
	}

	containerID, err := d.containerManager.CreateAndStart(
		ctx,
		containerName,
		imageTag,
		port,
		envSlice,
		memMB,
		cpuReq,
	)
	if err != nil {
		return "", fmt.Errorf("container rollout failed: %w", err)
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "🎉 [Deployer] Deployment successful! Container ID: %s\n", containerID[:12])
	}

	return containerID, nil
}
