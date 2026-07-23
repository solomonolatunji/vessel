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
	TSDBContainerName = "codedock-tsdb"
)

type TSDBManager struct {
	dockerClient *client.Client
}

func NewTSDBManager(cli *client.Client) *TSDBManager {
	return &TSDBManager{dockerClient: cli}
}

func (m *TSDBManager) EnsureTSDBRunning(ctx context.Context) error {
	_, err := m.dockerClient.ContainerInspect(ctx, TSDBContainerName)
	if err != nil {
		if errdefs.IsNotFound(err) {
			if err := m.createTSDBContainer(ctx); err != nil {
				return fmt.Errorf("failed to create tsdb: %w", err)
			}
		} else {
			return err
		}
	}

	if err := m.dockerClient.ContainerStart(ctx, TSDBContainerName, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start tsdb: %w", err)
	}

	slog.Info("tsdb (victoriametrics) is running")
	return nil
}

func (m *TSDBManager) createTSDBContainer(ctx context.Context) error {
	imageRef := "victoriametrics/victoria-metrics:latest"
	out, err := m.dockerClient.ImagePull(ctx, imageRef, image.PullOptions{})
	if err == nil {
		defer out.Close()
		io.Copy(io.Discard, out)
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8428/tcp": []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: "8428"}},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: "codedock-tsdb-data",
				Target: "/victoria-metrics-data",
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
	}

	resp, err := m.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
		ExposedPorts: nat.PortSet{
			"8428/tcp": struct{}{},
		},
		Cmd: []string{
			"-retentionPeriod=1y",
		},
	}, hostConfig, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			CodedockNetworkName: {},
		},
	}, nil, TSDBContainerName)

	if err != nil {
		return err
	}

	slog.Info("created tsdb container", "containerID", resp.ID)
	return nil
}
