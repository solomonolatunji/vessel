package engine

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/utils"
)

type StorageDeployer struct {
	dockerClient *client.Client
	store        StorageDeployerStore
}

func NewStorageDeployer(dockerClient *client.Client, s StorageDeployerStore) *StorageDeployer {
	return &StorageDeployer{
		dockerClient: dockerClient,
		store:        s,
	}
}

func (d *StorageDeployer) SpinUp(ctx context.Context, sc *models.Storage) (string, error) {
	if d.dockerClient == nil {
		return "", fmt.Errorf("docker daemon connection is not available")
	}
	containerName := utils.NormalizeContainerName(fmt.Sprintf("vessel-storage-%s", sc.Name))
	_ = d.dockerClient.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})
	imageName := "minio/minio:latest"
	envVars := []string{
		"MINIO_ROOT_USER=" + sc.AccessKey,
		"MINIO_ROOT_PASSWORD=" + sc.SecretKey,
	}
	cmd := []string{"server", "/data", "--console-address", fmt.Sprintf(":%d", sc.ConsolePort)}
	pullResp, err := d.dockerClient.ImagePull(ctx, imageName, dockertypes.ImagePullOptions{})
	if err == nil {
		_, _ = io.Copy(io.Discard, pullResp)
		_ = pullResp.Close()
	}
	hostVolumeDir, err := filepath.Abs(filepath.Join("data", "storage", sc.ID))
	if err != nil {
		return "", err
	}
	_ = os.MkdirAll(hostVolumeDir, 0o755)
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
				Target: "/data",
			},
		},
	}
	netCfg := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"vessel-net": {
				Aliases: []string{containerName, sc.Name},
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
		return "", fmt.Errorf("failed to create object storage container: %w", err)
	}
	if err := d.dockerClient.ContainerStart(ctx, created.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start object storage container: %w", err)
	}
	internalDNS := fmt.Sprintf("%s:%d", containerName, sc.APIPort)
	_ = d.store.UpdateStorageStatus(sc.ID, "running", created.ID)
	sc.ContainerID = created.ID
	sc.Status = "running"
	sc.InternalDNS = internalDNS
	return created.ID, nil
}

func (d *StorageDeployer) Stop(ctx context.Context, storageID string) error {
	if d.dockerClient == nil {
		return fmt.Errorf("docker daemon connection is not available")
	}
	sc, err := d.store.GetStorage(storageID)
	if err != nil || sc == nil {
		return fmt.Errorf("storage record not found")
	}
	containerName := utils.NormalizeContainerName(fmt.Sprintf("vessel-storage-%s", sc.Name))
	_ = d.dockerClient.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})
	return d.store.UpdateStorageStatus(storageID, "stopped", "")
}
