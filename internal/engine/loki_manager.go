package engine

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	LokiContainerName = "vessl-loki"
)

type LokiManager struct {
	dockerClient *client.Client
}

func NewLokiManager(cli *client.Client) *LokiManager {
	return &LokiManager{dockerClient: cli}
}

func (m *LokiManager) EnsureLokiRunning(ctx context.Context) error {
	_, err := m.dockerClient.ContainerInspect(ctx, LokiContainerName)
	if err != nil {
		if errdefs.IsNotFound(err) {
			if err := m.createLokiContainer(ctx); err != nil {
				return fmt.Errorf("failed to create loki: %w", err)
			}
		} else {
			return err
		}
	}

	if err := m.dockerClient.ContainerStart(ctx, LokiContainerName, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start loki: %w", err)
	}

	slog.Info("loki is running")
	return nil
}

func (m *LokiManager) createLokiContainer(ctx context.Context) error {
	imageRef := "grafana/loki:latest"
	out, err := m.dockerClient.ImagePull(ctx, imageRef, image.PullOptions{})
	if err == nil {
		defer out.Close()
		io.Copy(io.Discard, out)
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"3100/tcp": []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: "3100"}},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: "vessl-loki-data",
				Target: "/loki",
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
	}

	resp, err := m.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
		ExposedPorts: nat.PortSet{
			"3100/tcp": struct{}{},
		},
	}, hostConfig, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			VesslNetworkName: {},
		},
	}, nil, LokiContainerName)

	if err != nil {
		return err
	}

	slog.Info("created loki container", "containerID", resp.ID)
	return nil
}
