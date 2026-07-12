package engine

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
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

func (d *Deployer) Deploy(ctx context.Context, project *models.ProjectConfig, sourceDir string, logWriter io.Writer) (string, error) {
	apps, err := d.store.ListAppServicesByProject(project.ID)
	if err == nil && len(apps) > 0 {
		return d.DeployAppService(ctx, apps[0], sourceDir, logWriter)
	}
	syntheticApp := &models.AppService{
		ID:           project.ID,
		ProjectID:    project.ID,
		Name:         project.Name,
		InternalPort: 3000,
	}
	return d.DeployAppService(ctx, syntheticApp, sourceDir, logWriter)
}

func (d *Deployer) DeployAppService(ctx context.Context, app *models.AppService, sourceDir string, logWriter io.Writer) (string, error) {
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🚀 [Deployer] Starting deployment for service: %s (ID: %s)\n", app.Name, app.ID)
	}

	if os.Getenv("DEPLOY_DRY_RUN") == "true" {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "🚀 [Deployer] Dry-run mode is enabled. Skipping actual build and run steps.\n")
		}
		newContainerName := fmt.Sprintf("%s-dryrun", utils.NormalizeContainerName(app.ID))
		return newContainerName, nil
	}

	if err := d.prepareServerlessCode(app, sourceDir, logWriter); err != nil {
		return "", err
	}

	imageTag, err := d.buildImage(ctx, app, sourceDir, logWriter)
	if err != nil {
		return "", err
	}

	envSlice, err := d.prepareEnvironmentVariables(app, logWriter)
	if err != nil {
		return "", err
	}

	newContainerName := fmt.Sprintf("%s-%s", utils.NormalizeContainerName(app.ID), uuid.New().String()[:8])
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🔄 [Deployer] Rolling out container %s with %d encrypted environment variables...\n", newContainerName, len(envSlice))
	}

	if err := d.startContainer(ctx, app, newContainerName, imageTag, envSlice); err != nil {
		return "", err
	}

	if err := d.verifyHealthCheck(ctx, app, newContainerName, logWriter); err != nil {
		return "", err
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "🎉 [Deployer] Health check passed! Container is ready.\n")
		fmt.Fprintf(logWriter, "🎉 [Deployer] Deployment successful! Container ID: %s\n", newContainerName)
	}

	d.scheduleCleanup(app, newContainerName, logWriter)
	return newContainerName, nil
}

func (d *Deployer) prepareServerlessCode(app *models.AppService, sourceDir string, logWriter io.Writer) error {
	if app.BuildEngine != string(StrategyServerless) {
		return nil
	}

	code, err := d.store.GetServerlessFunctionCode(app.ID)
	if err != nil {
		return fmt.Errorf("could not retrieve serverless code: %w", err)
	}

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		return fmt.Errorf("could not create source directory: %w", err)
	}

	var filename string
	switch code.Runtime {
	case "nodejs":
		filename = "index.js"
	case "python":
		filename = "main.py"
	case "go":
		filename = "main.go"
	default:
		filename = "main.txt"
	}

	filePath := filepath.Join(sourceDir, filename)
	if err := os.WriteFile(filePath, []byte(code.CodeContent), 0644); err != nil {
		return fmt.Errorf("could not write serverless code to file: %w", err)
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "📝 [Deployer] Wrote serverless function code to %s\n", filePath)
	}
	return nil
}

func (d *Deployer) buildImage(ctx context.Context, app *models.AppService, sourceDir string, logWriter io.Writer) (string, error) {
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
	return imageTag, nil
}

func (d *Deployer) prepareEnvironmentVariables(app *models.AppService, logWriter io.Writer) ([]string, error) {
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
	return envSlice, nil
}

func (d *Deployer) startContainer(ctx context.Context, app *models.AppService, containerName, imageTag string, envSlice []string) error {
	port := app.InternalPort
	if port <= 0 {
		port = 3000
	}
	memMB := 512
	cpuReq := 0.5
	_, err := d.containerManager.CreateAndStart(
		ctx,
		containerName,
		imageTag,
		app.ID,
		app.Domain,
		port,
		envSlice,
		memMB,
		cpuReq,
		app.HealthCheckPath,
	)
	if err != nil {
		return fmt.Errorf("container rollout failed: %w", err)
	}
	return nil
}

func (d *Deployer) verifyHealthCheck(ctx context.Context, app *models.AppService, containerName string, logWriter io.Writer) error {
	healthy := d.waitForHealthyContainer(ctx, containerName, app.HealthCheckPath, app.InternalPort)
	if !healthy {
		_ = d.containerManager.StopAndRemove(ctx, containerName)
		if logWriter != nil {
			fmt.Fprintf(logWriter, "❌ [Deployer] Health check failed. Rolling back to previous version.\n")
		}
		return fmt.Errorf("health check failed, deployment aborted")
	}
	return nil
}

func (d *Deployer) scheduleCleanup(app *models.AppService, newContainerName string, logWriter io.Writer) {
	prefix := utils.NormalizeContainerName(app.ID)
	go func() {
		time.Sleep(10 * time.Second)
		if logWriter != nil {
			fmt.Fprintf(logWriter, "🧹 [Deployer] Cleaning up old orphaned containers...\n")
		}
		_ = d.containerManager.CleanupOrphanedContainers(context.Background(), prefix, newContainerName)
	}()
}

func (d *Deployer) Stop(ctx context.Context, containerID string) error {
	stopTimeout := 10
	return d.containerManager.dockerClient.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &stopTimeout})
}

func (d *Deployer) Remove(ctx context.Context, containerID string) error {
	err := d.containerManager.dockerClient.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
	if err != nil && !client.IsErrNotFound(err) {
		return err
	}
	return nil
}

func (d *Deployer) waitForHealthyContainer(ctx context.Context, containerName string, healthCheckPath string, internalPort int) bool {
	for i := 0; i < 30; i++ {
		time.Sleep(2 * time.Second)
		inspect, err := d.containerManager.Inspect(ctx, containerName)
		if err == nil {
			if !inspect.State.Running {
				if inspect.State.Status == "exited" {
					break
				}
				continue
			}
			if healthCheckPath != "" {
				var containerIP string
				if net, ok := inspect.NetworkSettings.Networks[utils.GetRuntimeNetwork()]; ok {
					containerIP = net.IPAddress
				}
				if containerIP != "" {
					port := internalPort
					if port <= 0 {
						port = 3000
					}
					resp, err := http.Get(fmt.Sprintf("http://%s:%d%s", containerIP, port, healthCheckPath))
					if err == nil {
						resp.Body.Close()
						if resp.StatusCode >= 200 && resp.StatusCode < 400 {
							return true
						}
					}
				}
			} else {
				return true
			}
		}
	}
	return false
}
