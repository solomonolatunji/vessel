package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v3"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

type ComposeDeployer struct {
	dockerClient *client.Client
}

func NewComposeDeployer(dockerClient *client.Client) *ComposeDeployer {
	return &ComposeDeployer{dockerClient: dockerClient}
}

func (cd *ComposeDeployer) Deploy(ctx context.Context, composePath string, projectID string) ([]*models.AppService, error) {
	compose, err := cd.parseComposeFile(composePath)
	if err != nil {
		return nil, err
	}

	results := make([]*models.AppService, 0, len(compose.Services))

	for name, svc := range compose.Services {
		app, err := cd.deployService(ctx, name, svc, projectID)
		if err != nil {
			return results, fmt.Errorf("service '%s': %w", name, err)
		}
		results = append(results, app)
	}

	return results, nil
}

func (cd *ComposeDeployer) parseComposeFile(path string) (*models.UserComposeFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read compose file: %w", err)
	}

	var compose models.UserComposeFile
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, fmt.Errorf("failed to parse compose file: %w", err)
	}

	if len(compose.Services) == 0 {
		return nil, fmt.Errorf("no services found in compose file")
	}

	return &compose, nil
}

func (cd *ComposeDeployer) deployService(ctx context.Context, name string, svc models.UserComposeService, projectID string) (*models.AppService, error) {
	port := extractPort(svc.Ports)
	containerName := fmt.Sprintf("vessl-comp-%s-%s", projectID[:8], name)
	aliasName := fmt.Sprintf("vessl-%s", name)

	app := &models.AppService{
		ID:           fmt.Sprintf("comp-%s-%s", projectID[:8], name),
		ProjectID:    projectID,
		Name:         name,
		InternalPort: port,
		Status:       models.AppServiceStatusRunning,
		ImageRef:     svc.Image,
	}

	if svc.Image == "" {
		return app, nil
	}

	envVars := buildEnvSlice(svc.Environment)
	if err := cd.startContainer(ctx, svc, containerName, aliasName, name, envVars); err != nil {
		return nil, err
	}

	return app, nil
}

func (cd *ComposeDeployer) startContainer(ctx context.Context, svc models.UserComposeService, containerName, aliasName, name string, envVars []string) error {
	networkName := utils.GetRuntimeNetwork()
	_ = utils.EnsureVesslNetwork(ctx, cd.dockerClient)

	containerCfg := &container.Config{
		Image: svc.Image,
		Env:   envVars,
	}

	hostCfg := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
	}
	applyVolumes(svc.Volumes, hostCfg)

	netCfg := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {
				Aliases: []string{containerName, aliasName, name},
			},
		},
	}

	created, err := cd.dockerClient.ContainerCreate(ctx, containerCfg, hostCfg, netCfg, nil, containerName)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	if err := cd.dockerClient.ContainerStart(ctx, created.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func extractPort(ports []string) int {
	if len(ports) == 0 {
		return 3000
	}
	parts := strings.Split(ports[0], ":")
	target := parts[0]
	if len(parts) >= 2 {
		target = parts[1]
	}
	var p int
	for _, c := range target {
		if c < '0' || c > '9' {
			break
		}
		p = p*10 + int(c-'0')
	}
	if p <= 0 {
		return 3000
	}
	return p
}

func applyVolumes(volumes []string, hostCfg *container.HostConfig) {
	for _, vol := range volumes {
		parts := strings.SplitN(vol, ":", 2)
		if len(parts) != 2 {
			continue
		}
		absPath, err := filepath.Abs(parts[0])
		if err != nil {
			continue
		}
		hostCfg.Mounts = append(hostCfg.Mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: absPath,
			Target: parts[1],
		})
	}
}

func buildEnvSlice(env map[string]string) []string {
	var result []string
	for k, v := range env {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}
