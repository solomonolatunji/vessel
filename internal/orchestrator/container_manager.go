package orchestrator

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// ContainerManager wraps the Docker SDK to control container creation, execution lifecycle, and health diagnostics.
type ContainerManager struct {
	dockerClient *client.Client
}

// NewContainerManager initializes a ContainerManager using the supplied Docker client instance.
func NewContainerManager(dockerClient *client.Client) *ContainerManager {
	return &ContainerManager{dockerClient: dockerClient}
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

	var memoryBytes int64 = 512 * 1024 * 1024
	if memoryLimitMB > 0 {
		memoryBytes = int64(memoryLimitMB) * 1024 * 1024
	}
	var nanoCPUs int64 = 500_000_000
	if cpuRequest > 0 {
		nanoCPUs = int64(cpuRequest * 1_000_000_000)
	}

	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "always"},
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: "0"}},
		},
		Resources: container.Resources{
			Memory:   memoryBytes,
			NanoCPUs: nanoCPUs,
		},
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
