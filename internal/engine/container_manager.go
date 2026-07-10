package engine

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"vessel.dev/vessel/internal/utils"
)

type ContainerManager struct {
	dockerClient *client.Client
	store        ContainerManagerStore
}

func NewContainerManager(dockerClient *client.Client, st ContainerManagerStore) *ContainerManager {
	return &ContainerManager{dockerClient: dockerClient, store: st}
}

func (c *ContainerManager) CreateAndStart(ctx context.Context, name, imageTag, serviceID, domain string, internalPort int, envs []string, memoryLimitMB int, cpuRequest float64) (string, error) {
	containerPort, err := nat.NewPort("tcp", fmt.Sprintf("%d", internalPort))
	if err != nil {
		return "", fmt.Errorf("invalid port definition: %w", err)
	}
	labels := map[string]string{}
	if serviceID != "" && domain != "" {
		labels = map[string]string{
			"traefik.enable": "true",
			fmt.Sprintf("traefik.http.routers.%s.rule", serviceID):                      fmt.Sprintf("Host(`%s`)", domain),
			fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", serviceID): fmt.Sprintf("%d", internalPort),
		}
	}

	config := &container.Config{
		Image:        imageTag,
		Env:          envs,
		ExposedPorts: nat.PortSet{containerPort: struct{}{}},
		Labels:       labels,
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
		NetworkMode: container.NetworkMode(utils.GetRuntimeNetwork()), // Connect to the shared Traefik network
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

func (c *ContainerManager) StopAndRemove(ctx context.Context, containerIDOrName string) error {
	stopTimeout := 10
	_ = c.dockerClient.ContainerStop(ctx, containerIDOrName, container.StopOptions{Timeout: &stopTimeout})
	err := c.dockerClient.ContainerRemove(ctx, containerIDOrName, container.RemoveOptions{Force: true})
	if err != nil && !client.IsErrNotFound(err) {
		return err
	}
	return nil
}

func (c *ContainerManager) Inspect(ctx context.Context, containerIDOrName string) (types.ContainerJSON, error) {
	return c.dockerClient.ContainerInspect(ctx, containerIDOrName)
}

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
