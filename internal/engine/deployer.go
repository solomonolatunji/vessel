package engine

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/containerd/errdefs"
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
	EnvInterpolator  func(projectID string) (map[string]map[string]string, error)
}

func NewDeployer(dockerClient *client.Client, s DeployerStore) *Deployer {
	return &Deployer{
		builder:          NewBuilder(dockerClient),
		containerManager: NewContainerManager(dockerClient, s),
		store:            s,
	}
}

func (d *Deployer) DeployAppService(ctx context.Context, app *models.AppService, sourceDir string, logWriter io.Writer) (string, error) {
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🚀 [Deployer] Starting deployment for service: %s (ID: %s)\n", app.Name, app.ID)
	}

	if utils.IsDryRun() {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "🚀 [Deployer] Dry-run mode is enabled. Skipping actual build and run steps.\n")
		}
		newContainerName := fmt.Sprintf("%s-dryrun", utils.NormalizeContainerName(app.ID))
		return newContainerName, nil
	}

	if err := d.prepareServerlessCode(app, sourceDir, logWriter); err != nil {
		return "", err
	}

	envVarsMap, err := d.getEnvironmentVariables(app, logWriter)
	if err != nil {
		return "", err
	}

	imageTag, err := d.buildImage(ctx, BuildImageOpts{
		App:        app,
		SourceDir:  sourceDir,
		EnvVarsMap: envVarsMap,
		LogWriter:  logWriter,
	})
	if err != nil {
		return "", err
	}

	envSlice := make([]string, 0, len(envVarsMap))
	for k, v := range envVarsMap {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}

	newContainerName := fmt.Sprintf("%s-%s", utils.NormalizeContainerName(app.ID), uuid.New().String()[:8])
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🔄 [Deployer] Rolling out container %s with %d encrypted environment variables...\n", newContainerName, len(envSlice))
	}

	startedNames, err := d.startContainer(ctx, StartContainerOpts{
		App:           app,
		ContainerName: newContainerName,
		ImageTag:      imageTag,
		EnvSlice:      envSlice,
	})
	if err != nil {
		return "", err
	}

	if err := d.verifyHealthCheck(ctx, app, startedNames[0], logWriter); err != nil {
		return "", err
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "🎉 [Deployer] Health check passed! Container is ready.\n")
		fmt.Fprintf(logWriter, "🎉 [Deployer] Deployment successful! Replicas started: %d\n", len(startedNames))
	}

	d.scheduleCleanup(app, startedNames, logWriter)
	return startedNames[0], nil
}

func (d *Deployer) prepareServerlessCode(app *models.AppService, sourceDir string, logWriter io.Writer) error {
	if app.BuildEngine != models.BuildEngineServerless {
		return nil
	}

	code, err := d.store.GetServerlessFunctionCode(app.ID)
	if err != nil {
		return fmt.Errorf("could not retrieve serverless code: %w", err)
	}

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		return fmt.Errorf("could not create source directory: %w", err)
	}

	if code.Runtime == "nodejs" {
		handlerPath := filepath.Join(sourceDir, "handler.js")
		if err := os.WriteFile(handlerPath, []byte(code.CodeContent), 0644); err != nil {
			return err
		}

		pkgJSON := `{
  "name": "vessl-function",
  "version": "1.0.0",
  "main": "server.js",
  "scripts": {
    "start": "node server.js"
  },
  "dependencies": {
    "express": "^4.18.2"
  }
}`
		if err := os.WriteFile(filepath.Join(sourceDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
			return err
		}

		serverJS := `const express = require('express');
const app = express();
app.use(express.json());

let handler;
try {
  handler = require('./handler.js');
} catch (err) {
  console.error("Failed to load handler.js:", err);
  process.exit(1);
}

app.all('*', async (req, res) => {
  try {
    if (typeof handler === 'function') {
      await handler(req, res);
    } else if (handler.default && typeof handler.default === 'function') {
      await handler.default(req, res);
    } else {
      res.status(500).send("No valid function exported in handler.js");
    }
  } catch (err) {
    console.error(err);
    res.status(500).send("Internal Server Error");
  }
});

const port = process.env.PORT || 3000;
app.listen(port, () => console.log('Function listening on port', port));`
		if err := os.WriteFile(filepath.Join(sourceDir, "server.js"), []byte(serverJS), 0644); err != nil {
			return err
		}

		if logWriter != nil {
			fmt.Fprintf(logWriter, "📝 [Deployer] Wrapped NodeJS serverless function with Express\n")
		}
		return nil
	}

	if code.Runtime == "python" {
		handlerPath := filepath.Join(sourceDir, "handler.py")
		if err := os.WriteFile(handlerPath, []byte(code.CodeContent), 0644); err != nil {
			return err
		}

		reqs := "Flask==2.3.2\nWerkzeug==2.3.4"
		if err := os.WriteFile(filepath.Join(sourceDir, "requirements.txt"), []byte(reqs), 0644); err != nil {
			return err
		}

		serverPy := `import os
from flask import Flask, request
import handler

app = Flask(__name__)

@app.route('/', defaults={'path': ''}, methods=['GET', 'POST', 'PUT', 'DELETE', 'PATCH'])
@app.route('/<path:path>', methods=['GET', 'POST', 'PUT', 'DELETE', 'PATCH'])
def catch_all(path):
    if hasattr(handler, 'main'):
        return handler.main(request)
    elif hasattr(handler, 'handler'):
        return handler.handler(request)
    else:
        return "No valid function exported in handler.py", 500

if __name__ == '__main__':
    port = int(os.environ.get('PORT', 3000))
    app.run(host='0.0.0.0', port=port)
`
		if err := os.WriteFile(filepath.Join(sourceDir, "main.py"), []byte(serverPy), 0644); err != nil {
			return err
		}

		if logWriter != nil {
			fmt.Fprintf(logWriter, "📝 [Deployer] Wrapped Python serverless function with Flask\n")
		}
		return nil
	}

	if code.Runtime == "go" {
		handlerPath := filepath.Join(sourceDir, "handler.go")
		if err := os.WriteFile(handlerPath, []byte(code.CodeContent), 0644); err != nil {
			return err
		}

		serverGo := `package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", Handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Printf("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
`
		if err := os.WriteFile(filepath.Join(sourceDir, "main.go"), []byte(serverGo), 0644); err != nil {
			return err
		}

		goMod := `module vessl-function

go 1.23
`
		if err := os.WriteFile(filepath.Join(sourceDir, "go.mod"), []byte(goMod), 0644); err != nil {
			return err
		}

		if logWriter != nil {
			fmt.Fprintf(logWriter, "📝 [Deployer] Wrapped Go serverless function with net/http\n")
		}
		return nil
	}

	var filename string
	switch code.Runtime {
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

type BuildImageOpts struct {
	App        *models.AppService
	SourceDir  string
	EnvVarsMap map[string]string
	LogWriter  io.Writer
}

func (d *Deployer) buildImage(ctx context.Context, opts BuildImageOpts) (string, error) {
	buildOpts := BuildOptions{
		ProjectID: opts.App.ProjectID,
		ServiceID: opts.App.ID,
		SourceDir: opts.SourceDir,
		LogWriter: opts.LogWriter,
		AppConfig: opts.App,
		EnvVars:   opts.EnvVarsMap,
	}
	imageTag, err := d.builder.Build(ctx, buildOpts)
	if err != nil {
		return "", fmt.Errorf("build phase failed: %w", err)
	}
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "✅ [Deployer] Successfully built OCI image: %s\n", imageTag)
	}
	return imageTag, nil
}

type StartContainerOpts struct {
	App           *models.AppService
	ContainerName string
	ImageTag      string
	EnvSlice      []string
}

func (d *Deployer) startContainer(ctx context.Context, opts StartContainerOpts) ([]string, error) {
	port := opts.App.InternalPort
	if port <= 0 {
		port = defaultAppPort()
	}
	if opts.App.StaticOutput != "" {
		port = 80 // NGINX alpine default port
	}

	replicas := opts.App.Replicas
	if replicas <= 0 {
		replicas = 1
	}

	var startedNames []string

	for i := 0; i < replicas; i++ {
		containerName := opts.ContainerName
		if replicas > 1 {
			containerName = fmt.Sprintf("%s-%d", opts.ContainerName, i)
		}

		containerOpts := ContainerRunOptions{
			Name:            containerName,
			ImageTag:        opts.ImageTag,
			ServiceID:       opts.App.ID,
			Domain:          opts.App.Domain,
			InternalPort:    port,
			RuntimeMode:     opts.App.RuntimeMode,
			Envs:            opts.EnvSlice,
			MemoryLimitMB:   defaultMemoryMB(),
			CPURequest:      defaultCPURequest(),
			HealthCheckPath: opts.App.HealthCheckPath,
		}

		_, err := d.containerManager.CreateAndStart(ctx, containerOpts)
		if err != nil {
			return startedNames, fmt.Errorf("container rollout failed for replica %d: %w", i, err)
		}
		startedNames = append(startedNames, containerName)
	}

	return startedNames, nil
}

func (d *Deployer) scheduleCleanup(app *models.AppService, newContainerNames []string, logWriter io.Writer) {
	prefix := utils.NormalizeContainerName(app.ID)
	go func() {
		time.Sleep(10 * time.Second)
		if logWriter != nil {
			fmt.Fprintf(logWriter, "🧹 [Deployer] Cleaning up old orphaned containers...\n")
		}
		_ = d.containerManager.CleanupOrphanedContainers(context.Background(), prefix, newContainerNames)
	}()
}

func (d *Deployer) Stop(ctx context.Context, containerID string) error {
	stopTimeout := 10
	return d.containerManager.dockerClient.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &stopTimeout})
}

func (d *Deployer) StopAppService(ctx context.Context, appID string) error {
	prefix := utils.NormalizeContainerName(appID)
	return d.containerManager.CleanupOrphanedContainers(ctx, prefix, []string{})
}

func (d *Deployer) RestartAppService(ctx context.Context, appID string) error {
	prefix := utils.NormalizeContainerName(appID)
	containers, err := d.containerManager.dockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	timeout := 10
	restarted := 0
	for _, ctn := range containers {
		for _, name := range ctn.Names {
			if strings.HasPrefix(name, "/"+prefix+"-") {
				if err := d.containerManager.dockerClient.ContainerRestart(ctx, ctn.ID, container.StopOptions{Timeout: &timeout}); err != nil {
					return fmt.Errorf("failed to restart container %s: %w", ctn.ID, err)
				}
				restarted++
				break
			}
		}
	}
	if restarted == 0 {
		return fmt.Errorf("no containers found for app")
	}
	return nil
}

func (d *Deployer) Remove(ctx context.Context, containerID string) error {
	err := d.containerManager.dockerClient.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
	if err != nil && !errdefs.IsNotFound(err) {
		return err
	}
	return nil
}

func (d *Deployer) DeployImage(ctx context.Context, app *models.AppService, logWriter io.Writer) (string, error) {
	if app.ImageRef == "" {
		return "", fmt.Errorf("image ref is empty")
	}

	if utils.IsDryRun() {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "🚀 [Deployer] Dry-run mode is enabled. Skipping image deploy.\n")
		}
		newContainerName := fmt.Sprintf("%s-dryrun", utils.NormalizeContainerName(app.ID))
		return newContainerName, nil
	}

	port := app.InternalPort
	if port <= 0 {
		port = defaultAppPort()
	}

	containerName := fmt.Sprintf("%s-%s", utils.NormalizeContainerName(app.ID), uuid.New().String()[:8])

	startedNames, err := d.startContainer(ctx, StartContainerOpts{
		App:           app,
		ContainerName: containerName,
		ImageTag:      app.ImageRef,
		EnvSlice:      nil,
	})
	if err != nil {
		return "", err
	}

	if err := d.verifyHealthCheck(ctx, app, startedNames[0], logWriter); err != nil {
		return "", err
	}

	d.scheduleCleanup(app, startedNames, logWriter)
	return startedNames[0], nil
}
