package orchestrator

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
	"github.com/solomonolatunji/vessel/internal/utils"
)

// DatabaseDeployer manages the automated provisioning, lifecycle, and network isolation of stateful database engines.
type DatabaseDeployer struct {
	dockerClient *client.Client
	store        *store.Store
}

// NewDatabaseDeployer initializes a DatabaseDeployer wired to the Docker SDK and state store.
func NewDatabaseDeployer(dockerClient *client.Client, s *store.Store) *DatabaseDeployer {
	return &DatabaseDeployer{
		dockerClient: dockerClient,
		store:        s,
	}
}

// SpinUp provisions a persistent Docker volume, pulls the engine image, and launches the database container on vessel-net.
func (d *DatabaseDeployer) SpinUp(ctx context.Context, dbConfig *types.DatabaseConfig) (string, error) {
	if d.dockerClient == nil {
		return "", fmt.Errorf("docker daemon connection is not available")
	}

	containerName := utils.NormalizeContainerName(fmt.Sprintf("vessel-db-%s", dbConfig.Name))

	_ = d.dockerClient.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})

	var imageName string
	var envVars []string
	var cmd []string
	var containerMountPath string

	switch strings.ToLower(dbConfig.Engine) {
	case "postgres", "postgresql":
		imageName = "postgres:" + getVersionOrDefault(dbConfig.Version, "16-alpine")
		envVars = []string{
			"POSTGRES_USER=" + dbConfig.Username,
			"POSTGRES_PASSWORD=" + dbConfig.Password,
			"POSTGRES_DB=" + dbConfig.DatabaseName,
		}
		containerMountPath = "/var/lib/postgresql/data"
	case "mysql":
		imageName = "mysql:" + getVersionOrDefault(dbConfig.Version, "8.0")
		envVars = []string{
			"MYSQL_ROOT_PASSWORD=" + dbConfig.Password,
			"MYSQL_USER=" + dbConfig.Username,
			"MYSQL_PASSWORD=" + dbConfig.Password,
			"MYSQL_DATABASE=" + dbConfig.DatabaseName,
		}
		containerMountPath = "/var/lib/mysql"
	case "redis":
		imageName = "redis:" + getVersionOrDefault(dbConfig.Version, "7-alpine")
		if dbConfig.Password != "" {
			cmd = []string{"redis-server", "--requirepass", dbConfig.Password}
		}
		containerMountPath = "/data"
	case "mongodb", "mongo":
		imageName = "mongo:" + getVersionOrDefault(dbConfig.Version, "7.0")
		envVars = []string{
			"MONGO_INITDB_ROOT_USERNAME=" + dbConfig.Username,
			"MONGO_INITDB_ROOT_PASSWORD=" + dbConfig.Password,
			"MONGO_INITDB_DATABASE=" + dbConfig.DatabaseName,
		}
		containerMountPath = "/data/db"
	default:
		return "", fmt.Errorf("unsupported database engine: %s", dbConfig.Engine)
	}

	pullResp, err := d.dockerClient.ImagePull(ctx, imageName, dockertypes.ImagePullOptions{})
	if err == nil {
		_, _ = io.Copy(io.Discard, pullResp)
		_ = pullResp.Close()
	}

	hostVolumeDir, err := filepath.Abs(filepath.Join("data", "databases", dbConfig.ID))
	if err != nil {
		return "", err
	}
	_ = os.MkdirAll(hostVolumeDir, 0755)

	if err := utils.EnsureVesselNetwork(ctx, d.dockerClient); err != nil {
		return "", fmt.Errorf("failed to ensure Docker network: %w", err)
	}

	containerCfg := &container.Config{
		Image: imageName,
		Env:   envVars,
		Cmd:   cmd,
	}

	hostCfg := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "always"},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: hostVolumeDir,
				Target: containerMountPath,
			},
		},
	}

	netCfg := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"vessel-net": {
				Aliases: []string{containerName, dbConfig.Name},
			},
		},
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

// Stop halts and removes the container associated with a managed database instance while preserving volume data.
func (d *DatabaseDeployer) Stop(ctx context.Context, dbID string) error {
	if d.dockerClient == nil {
		return fmt.Errorf("docker daemon connection is not available")
	}

	dbConfig, err := d.store.GetDatabase(dbID)
	if err != nil || dbConfig == nil {
		return fmt.Errorf("database record not found")
	}

	containerName := utils.NormalizeContainerName(fmt.Sprintf("vessel-db-%s", dbConfig.Name))
	_ = d.dockerClient.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})
	return d.store.UpdateDatabaseStatus(dbID, "stopped", "")
}

func getVersionOrDefault(version, defaultVersion string) string {
	if version == "" {
		return defaultVersion
	}
	return version
}
