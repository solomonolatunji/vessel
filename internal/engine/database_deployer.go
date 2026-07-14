package engine

import (
	"context"
	"fmt"
	"io"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"vessl.dev/vessl/internal/models"

	"vessl.dev/vessl/internal/utils"
)

type DatabaseDeployer struct {
	dockerClient *client.Client
	store        DatabaseDeployerStore
}

func NewDatabaseDeployer(dockerClient *client.Client, s DatabaseDeployerStore) *DatabaseDeployer {
	return &DatabaseDeployer{
		dockerClient: dockerClient,
		store:        s,
	}
}

func (d *DatabaseDeployer) SpinUp(ctx context.Context, dbConfig *models.Database) (string, error) {
	if d.dockerClient == nil {
		return "", fmt.Errorf("docker daemon connection is not available")
	}
	containerName := utils.NormalizeContainerName(fmt.Sprintf("vessl-db-%s", dbConfig.Name))
	_ = d.dockerClient.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})
	var imageName string
	var envVars []string
	var cmd []string
	var containerMountPath string
	tmplMgr, err := NewTemplateManager()
	if err != nil {
		return "", fmt.Errorf("failed to initialize template manager: %w", err)
	}

	composeFile, err := tmplMgr.GetTemplate(strings.ToLower(dbConfig.Engine))
	if err != nil {
		return "", fmt.Errorf("unsupported database engine %s: %w", dbConfig.Engine, err)
	}

	tmplService, exists := composeFile.Services[strings.ToLower(dbConfig.Engine)]
	if !exists {
		for _, s := range composeFile.Services {
			tmplService = s
			break
		}
	}

	imageName = tmplService.Image
	if dbConfig.Version != "" && !strings.Contains(imageName, ":") {
		imageName = imageName + ":" + dbConfig.Version
	} else if dbConfig.Version != "" {
		parts := strings.Split(imageName, ":")
		imageName = parts[0] + ":" + dbConfig.Version
	}

	for _, ev := range tmplService.Environment {
		resolved := strings.ReplaceAll(ev, "${db.password}", dbConfig.Password)
		resolved = strings.ReplaceAll(resolved, "${db.username}", dbConfig.Username)
		resolved = strings.ReplaceAll(resolved, "${db.database_name}", dbConfig.DatabaseName)
		envVars = append(envVars, resolved)
	}

	for i := 0; i < len(tmplService.Command); i++ {
		c := tmplService.Command[i]
		if c == "--requirepass" && dbConfig.Password == "" {
			i++
			continue
		}
		resolved := strings.ReplaceAll(c, "${db.password}", dbConfig.Password)
		resolved = strings.ReplaceAll(resolved, "${db.username}", dbConfig.Username)
		if resolved != "" {
			cmd = append(cmd, resolved)
		}
	}

	if dbConfig.CustomArgs != "" {
		args := strings.Fields(dbConfig.CustomArgs)
		cmd = append(cmd, args...)
	}

	if len(tmplService.Volumes) > 0 {
		parts := strings.Split(tmplService.Volumes[0], ":")
		if len(parts) >= 2 {
			containerMountPath = parts[1]
		}
	} else {
		containerMountPath = "/data"
	}
	pullResp, err := d.dockerClient.ImagePull(ctx, imageName, dockertypes.ImagePullOptions{})
	if err == nil {
		_, _ = io.Copy(io.Discard, pullResp)
		_ = pullResp.Close()
	}
	volumeName := fmt.Sprintf("vessl-db-data-%s", dbConfig.ID)

	if err := utils.EnsureVesslNetwork(ctx, d.dockerClient); err != nil {
		return "", fmt.Errorf("failed to ensure Docker network: %w", err)
	}
	containerCfg := &container.Config{
		Image: imageName,
		Env:   envVars,
		Cmd:   cmd,
	}
	hostCfg := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "always"},
		Resources: container.Resources{
			Memory:   utils.MegaBytesToBytes(utils.DefaultDBMemoryMB()),
			NanoCPUs: utils.CPURequestToNanoCPUs(utils.DefaultDBCPURequest()),
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volumeName,
				Target: containerMountPath,
			},
		},
	}
	netCfg := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			utils.GetRuntimeNetwork(): {
				Aliases: []string{containerName, dbConfig.Name},
			},
		},
	}
	if d.store != nil {
		settings, _ := d.store.GetServerSettings()
		if settings != nil {
			ApplyCustomDNS(hostCfg, settings.CustomDNSResolvers)
		}
	}
	created, err := d.dockerClient.ContainerCreate(ctx, containerCfg, hostCfg, netCfg, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("failed to create database container: %w", err)
	}
	if err := d.dockerClient.ContainerStart(ctx, created.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start database container: %w", err)
	}
	internalDNS := fmt.Sprintf("%s:%d", containerName, dbConfig.Port)
	_ = d.store.UpdateDatabaseStatus(dbConfig.ID, "running", created.ID)
	dbConfig.ContainerID = created.ID
	dbConfig.Status = "running"
	dbConfig.InternalDNS = internalDNS
	return created.ID, nil
}

func (d *DatabaseDeployer) Stop(ctx context.Context, dbID string) error {
	if d.dockerClient == nil {
		return fmt.Errorf("docker daemon connection is not available")
	}
	dbConfig, err := d.store.GetDatabase(dbID)
	if err != nil || dbConfig == nil {
		return utils.NewNotFoundError("Database", dbID)
	}
	containerName := utils.NormalizeContainerName(fmt.Sprintf("vessl-db-%s", dbConfig.Name))
	_ = d.dockerClient.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})
	return d.store.UpdateDatabaseStatus(dbID, "stopped", "")
}
