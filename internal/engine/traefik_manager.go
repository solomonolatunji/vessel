package engine

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	TraefikContainerName = "codedock-traefik"
	CodedockNetworkName  = "codedock-network"
)

type TraefikManager struct {
	dockerClient *client.Client
	tlsEmail     string
}

func NewTraefikManager(cli *client.Client, tlsEmail string) *TraefikManager {
	return &TraefikManager{dockerClient: cli, tlsEmail: tlsEmail}
}

func (m *TraefikManager) EnsureTraefikRunning(ctx context.Context) error {
	if err := m.ensureNetwork(ctx); err != nil {
		return fmt.Errorf("failed to ensure network: %w", err)
	}

	_, err := m.dockerClient.ContainerInspect(ctx, TraefikContainerName)
	if err != nil {
		if errdefs.IsNotFound(err) {
			if err := m.createTraefikContainer(ctx); err != nil {
				return fmt.Errorf("failed to create traefik: %w", err)
			}
		} else {
			return err
		}
	}

	if err := m.dockerClient.ContainerStart(ctx, TraefikContainerName, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start traefik: %w", err)
	}

	slog.Info("traefik reverse proxy is running")
	return nil
}

func (m *TraefikManager) ensureNetwork(ctx context.Context) error {
	_, err := m.dockerClient.NetworkInspect(ctx, CodedockNetworkName, network.InspectOptions{})
	if err != nil {
		if errdefs.IsNotFound(err) {
			_, err = m.dockerClient.NetworkCreate(ctx, CodedockNetworkName, network.CreateOptions{
				Driver: "bridge",
			})
			return err
		}
		return err
	}
	return nil
}

func traefikImage() string {
	if img := os.Getenv("CODEDOCK_TRAEFIK_IMAGE"); img != "" {
		return img
	}
	return "traefik:v3.6"
}

func dockerSocketPath() string {
	if p := os.Getenv("DOCKER_SOCKET_PATH"); p != "" {
		return p
	}
	return "/var/run/docker.sock"
}

func (m *TraefikManager) createTraefikContainer(ctx context.Context) error {
	imageRef := traefikImage()
	out, err := m.dockerClient.ImagePull(ctx, imageRef, image.PullOptions{})
	if err == nil {
		defer out.Close()
		io.Copy(io.Discard, out)
	}

	cmdArgs := m.buildTraefikCmdArgs()
	hostConfig := &container.HostConfig{
		PortBindings: m.buildPortBindings(),
		Mounts:       m.buildTraefikMounts(),
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
	}

	resp, err := m.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
		Cmd:   cmdArgs,
		ExposedPorts: nat.PortSet{
			"80/tcp":   struct{}{},
			"443/tcp":  struct{}{},
			"443/udp":  struct{}{},
			"8080/tcp": struct{}{},
		},
		Labels: map[string]string{
			"traefik.enable": "true",
			"traefik.http.routers.traefik.entrypoints":               "http",
			"traefik.http.routers.traefik.service":                   "api@internal",
			"traefik.http.services.traefik.loadbalancer.server.port": "8080",
		},
		Healthcheck: &container.HealthConfig{
			Test:     []string{"CMD", "wget", "-qO-", "http://localhost:80/ping"},
			Interval: 4 * time.Second,
			Timeout:  2 * time.Second,
			Retries:  5,
		},
	}, hostConfig, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			CodedockNetworkName: {},
		},
	}, nil, TraefikContainerName)

	if err != nil {
		return err
	}

	slog.Info("created traefik container", "containerID", resp.ID)
	return nil
}

func (m *TraefikManager) buildTraefikCmdArgs() []string {
	cmdArgs := []string{
		"--ping=true",
		"--ping.entrypoint=http",
		"--api.insecure=true",
		"--providers.docker=true",
		"--providers.docker.exposedbydefault=false",
		"--providers.docker.network=" + CodedockNetworkName,
		"--entrypoints.web.address=:80",
		"--entrypoints.websecure.address=:443",
		"--entrypoints.https.http3=true",
		"--entrypoints.web.http.redirections.entryPoint.to=websecure",
		"--entrypoints.web.http.redirections.entryPoint.scheme=https",
	}

	if m.tlsEmail != "" {
		cmdArgs = append(cmdArgs,
			"--certificatesresolvers.letsencrypt.acme.tlschallenge=true",
			"--certificatesresolvers.letsencrypt.acme.email="+m.tlsEmail,
			"--certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json",
		)
	}
	return cmdArgs
}

func (m *TraefikManager) buildTraefikMounts() []mount.Mount {
	sockPath := dockerSocketPath()
	mounts := []mount.Mount{
		{
			Type:     mount.TypeBind,
			Source:   sockPath,
			Target:   "/var/run/docker.sock",
			ReadOnly: true,
		},
	}
	if m.tlsEmail != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: "codedock-traefik-acme",
			Target: "/letsencrypt",
		})
	}
	return mounts
}

func (m *TraefikManager) buildPortBindings() nat.PortMap {
	httpPort := os.Getenv("CODEDOCK_TRAEFIK_HTTP_PORT")
	if httpPort == "" {
		httpPort = "80"
	}
	httpsPort := os.Getenv("CODEDOCK_TRAEFIK_HTTPS_PORT")
	if httpsPort == "" {
		httpsPort = "443"
	}
	apiPort := os.Getenv("CODEDOCK_TRAEFIK_API_PORT")
	if apiPort == "" {
		apiPort = "8080"
	}
	return nat.PortMap{
		"80/tcp":   []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: httpPort}},
		"443/tcp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: httpsPort}},
		"443/udp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: httpsPort}},
		"8080/tcp": []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: apiPort}},
	}
}
