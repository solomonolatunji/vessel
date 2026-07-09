package orchestrator

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/client"
	"vessel.dev/vessel/internal/project"
	"vessel.dev/vessel/internal/service"
	"vessel.dev/vessel/internal/utils"
)

type Deployer struct {
	builder          Builder
	containerManager *ContainerManager
	store            DeployerStore
	EnvProvider      func(projectID string) (map[string]string, error)
}

func NewDeployer(dockerClient *client.Client, s DeployerStore) *Deployer {
	return &Deployer{
		builder:          NewBuilder(dockerClient),
		containerManager: NewContainerManager(dockerClient, s),
		store:            s,
	}
}

func (d *Deployer) Deploy(ctx context.Context, project *project.ProjectConfig, sourceDir string, logWriter io.Writer) (string, error) {
	apps, err := d.store.ListAppServicesByProject(project.ID)
	if err == nil && len(apps) > 0 {
		return d.DeployAppService(ctx, apps[0], sourceDir, logWriter)
	}

	syntheticApp := &service.AppService{
		ID:           project.ID,
		ProjectID:    project.ID,
		Name:         project.Name,
		InternalPort: 3000,
	}
	return d.DeployAppService(ctx, syntheticApp, sourceDir, logWriter)
}

func (d *Deployer) DeployAppService(ctx context.Context, app *service.AppService, sourceDir string, logWriter io.Writer) (string, error) {
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🚀 [Deployer] Starting deployment for service: %s (ID: %s)\n", app.Name, app.ID)
	}

	buildOpts := BuildOptions{
		ProjectID: app.ProjectID,
		ServiceID: app.ID,
		SourceDir: sourceDir,
		LogWriter: logWriter,
		AppConfig: app,
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
	memMB := 512
	cpuReq := 0.5

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
