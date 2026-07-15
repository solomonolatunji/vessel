package engine

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
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

	composeFile, err := tmplMgr.GetTemplate(strings.ToLower(string(dbConfig.Engine)))
	if err != nil {
		return "", fmt.Errorf("unsupported database engine %s: %w", dbConfig.Engine, err)
	}

	tmplService, exists := composeFile.Services[strings.ToLower(string(dbConfig.Engine))]
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

	if dbConfig.LogicalReplication && (strings.Contains(imageName, "postgres") || strings.Contains(imageName, "timescaledb")) {
		cmd = append(cmd, "-c", "wal_level=logical", "-c", "max_replication_slots=10", "-c", "max_wal_senders=10")
	}

	if len(tmplService.Volumes) > 0 {
		parts := strings.Split(tmplService.Volumes[0], ":")
		if len(parts) >= 2 {
			containerMountPath = parts[1]
		}
	} else {
		containerMountPath = "/data"
	}
	pullResp, err := d.dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
	if err == nil {
		_, _ = io.Copy(io.Discard, pullResp)
		_ = pullResp.Close()
	}
	volumeName := fmt.Sprintf("vessl-db-data-%s", dbConfig.ID)

	if err := utils.EnsureVesslNetwork(ctx, d.dockerClient); err != nil {
		return "", fmt.Errorf("failed to ensure Docker network: %w", err)
	}
	labels := make(map[string]string)
	if dbConfig.ExternalDNS != "" {
		labels["traefik.enable"] = "true"
		labels[fmt.Sprintf("traefik.tcp.routers.%s.rule", containerName)] = fmt.Sprintf("HostSNI(`%s`)", dbConfig.ExternalDNS)
		labels[fmt.Sprintf("traefik.tcp.routers.%s.tls", containerName)] = "true"
		labels[fmt.Sprintf("traefik.tcp.routers.%s.tls.certresolver", containerName)] = "letsencrypt"
		labels[fmt.Sprintf("traefik.tcp.routers.%s.entrypoints", containerName)] = "websecure"
		labels[fmt.Sprintf("traefik.tcp.services.%s.loadbalancer.server.port", containerName)] = fmt.Sprintf("%d", dbConfig.Port)
	}

	containerCfg := &container.Config{
		Image:  imageName,
		Env:    envVars,
		Cmd:    cmd,
		Labels: labels,
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

func (d *DatabaseDeployer) ImportData(ctx context.Context, dbConfig *models.Database, sourceURL string) error {
	if d.dockerClient == nil {
		return fmt.Errorf("docker daemon connection is not available")
	}

	containerName := utils.NormalizeContainerName(fmt.Sprintf("vessl-db-%s", dbConfig.Name))

	switch strings.ToLower(string(dbConfig.Engine)) {
	case "postgres":
		cmd := fmt.Sprintf("pg_dump -Fc \"%s\" | pg_restore -U %s -d %s -1", sourceURL, dbConfig.Username, dbConfig.DatabaseName)
		execConfig := container.ExecOptions{
			Cmd:          []string{"sh", "-c", cmd},
			AttachStderr: true,
			AttachStdout: true,
		}
		execID, err := d.dockerClient.ContainerExecCreate(ctx, containerName, execConfig)
		if err != nil {
			return fmt.Errorf("failed to create exec for pg_dump: %w", err)
		}
		err = d.dockerClient.ContainerExecStart(ctx, execID.ID, container.ExecStartOptions{})
		if err != nil {
			return fmt.Errorf("failed to start exec for pg_dump: %w", err)
		}
		// wait for exec to finish
		for {
			inspect, err := d.dockerClient.ContainerExecInspect(ctx, execID.ID)
			if err != nil {
				return fmt.Errorf("failed to inspect exec: %w", err)
			}
			if !inspect.Running {
				if inspect.ExitCode != 0 {
					return fmt.Errorf("pg_dump/restore failed with exit code %d", inspect.ExitCode)
				}
				break
			}
			time.Sleep(1 * time.Second)
		}
		return nil

	case "redis":
		// Stream to dump.rdb then restart container so it loads it
		cmd := fmt.Sprintf("redis-cli -u \"%s\" --rdb /data/dump.rdb", sourceURL)
		execConfig := container.ExecOptions{
			Cmd:          []string{"sh", "-c", cmd},
			AttachStderr: true,
			AttachStdout: true,
		}
		execID, err := d.dockerClient.ContainerExecCreate(ctx, containerName, execConfig)
		if err != nil {
			return fmt.Errorf("failed to create exec for redis-cli: %w", err)
		}
		err = d.dockerClient.ContainerExecStart(ctx, execID.ID, container.ExecStartOptions{})
		if err != nil {
			return fmt.Errorf("failed to start exec for redis-cli: %w", err)
		}
		for {
			inspect, err := d.dockerClient.ContainerExecInspect(ctx, execID.ID)
			if err != nil {
				return fmt.Errorf("failed to inspect exec: %w", err)
			}
			if !inspect.Running {
				if inspect.ExitCode != 0 {
					return fmt.Errorf("redis-cli failed with exit code %d", inspect.ExitCode)
				}
				break
			}
			time.Sleep(1 * time.Second)
		}
		// restart redis to load rdb
		return d.dockerClient.ContainerRestart(ctx, containerName, container.StopOptions{})

	default:
		return fmt.Errorf("data import not supported for engine %s", dbConfig.Engine)
	}
}
