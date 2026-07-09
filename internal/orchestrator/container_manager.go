package orchestrator

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/utils"
)

type ContainerManager struct {
	dockerClient *client.Client
	store        *store.Store
}

func NewContainerManager(dockerClient *client.Client, st *store.Store) *ContainerManager {
	return &ContainerManager{dockerClient: dockerClient, store: st}
}

// CreateAndStart provisions a new container with explicit CPU/RAM boundaries and environment injections.
func (c *ContainerManager) CreateAndStart(ctx context.Context, name, imageTag string, internalPort int, envs []string, memoryLimitMB int, cpuRequest float64) (string, error) {
	containerPort, err := nat.NewPort("tcp", fmt.Sprintf("%d", internalPort))
	if err != nil {
		return "", fmt.Errorf("invalid port definition: %w", err)
	}

	config := &container.Config{
		Image:        imageTag,
		Env:          envs,
		ExposedPorts: nat.PortSet{containerPort: struct{}{}},
	}

	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "always"},
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: "0"}},
		},
		Resources: container.Resources{
			Memory:   utils.MegaBytesToBytes(memoryLimitMB),
			NanoCPUs: utils.CPURequestToNanoCPUs(cpuRequest),
		},
	}

	if c.store != nil {
		settings, _ := c.store.GetServerSettings()
		if settings != nil && strings.TrimSpace(settings.CustomDNSResolvers) != "" {
			parts := strings.Split(settings.CustomDNSResolvers, ",")
			var dnsList []string
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					dnsList = append(dnsList, p)
				}
			}
			if len(dnsList) > 0 {
				hostConfig.DNS = dnsList
			}
		}
	}

	_ = c.StopAndRemove(ctx, name)

	resp, err := c.dockerClient.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return "", fmt.Errorf("docker container create failed: %w", err)
	}

	if err := c.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("docker container start failed: %w", err)
	}

	return resp.ID, nil
}

// StopAndRemove halts and destroys an existing container by name or container ID cleanly.
func (c *ContainerManager) StopAndRemove(ctx context.Context, containerIDOrName string) error {
	stopTimeout := 10
	_ = c.dockerClient.ContainerStop(ctx, containerIDOrName, container.StopOptions{Timeout: &stopTimeout})
	err := c.dockerClient.ContainerRemove(ctx, containerIDOrName, container.RemoveOptions{Force: true})
	if err != nil && !client.IsErrNotFound(err) {
		return err
	}
	return nil
}

// Inspect retrieves the low-level runtime status and mapped host ports of a container.
func (c *ContainerManager) Inspect(ctx context.Context, containerIDOrName string) (types.ContainerJSON, error) {
	return c.dockerClient.ContainerInspect(ctx, containerIDOrName)
}

// StreamLogs pipes live container stdout and stderr to the provided destination io.Writer.
func (c *ContainerManager) StreamLogs(ctx context.Context, containerIDOrName string, out io.Writer) error {
	logsReader, err := c.dockerClient.ContainerLogs(ctx, containerIDOrName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "100",
	})
	if err != nil {
		return fmt.Errorf("failed to open container logs stream: %w", err)
	}
	defer logsReader.Close()

	_, err = io.Copy(out, logsReader)
	return err
}
