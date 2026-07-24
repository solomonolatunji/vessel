package worker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/models"
	"github.com/docker/docker/client"
)

func (d *WorkerDaemon) processDeployment(ctx context.Context, commandID string, payload models.WorkerDeployAppPayload) error {
	// Initialize Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	// Create workspace
	workspace := filepath.Join(os.TempDir(), "codedock-worker", payload.AppID)
	_ = os.RemoveAll(workspace)
	if err := os.MkdirAll(workspace, 0755); err != nil {
		return err
	}

	// Clone repository
	if payload.GitRepoURL != "" {
		cloneURL := payload.GitRepoURL
		if payload.GitAuthToken != "" {
			cloneURL = strings.Replace(cloneURL, "https://", fmt.Sprintf("https://oauth2:%s@", payload.GitAuthToken), 1)
		}

		args := []string{"clone", "--depth", "1"}
		if payload.GitBranch != "" {
			args = append(args, "--branch", payload.GitBranch)
		}
		args = append(args, cloneURL, workspace)

		cmd := exec.CommandContext(ctx, "git", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git clone failed: %v\n%s", err, string(out))
		}
	}

	// Initialize Deployer
	store := NewWorkerLocalStore(payload)
	deployer := engine.NewDeployer(dockerClient, store)

	app := &models.AppService{
		ID:              payload.AppID,
		Name:            payload.AppID,
		Domain:          payload.Domain,
		RuntimeMode:     models.RuntimeMode(payload.RuntimeMode),
		HealthCheckPath: payload.HealthCheckPath,
		MemoryLimit:     payload.MemoryLimitMB,
		CPULimit:        float64(payload.CPURequest),
		Volumes:         nil,
		InternalPort:    8080,
	}
	if len(payload.Ports) > 0 {
		app.InternalPort = 80
	}

	if app.RuntimeMode == "" {
		app.RuntimeMode = models.RuntimeModeWeb
	}

	_, err = deployer.DeployAppService(ctx, app, workspace, os.Stdout)
	return err
}
