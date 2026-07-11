package proxy

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	TraefikContainerName = "vessel-traefik"
	VesselNetworkName    = "vessel-network"
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
		if client.IsErrNotFound(err) {
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

	log.Println("Traefik reverse proxy is running")
	return nil
}

func (m *TraefikManager) ensureNetwork(ctx context.Context) error {
	_, err := m.dockerClient.NetworkInspect(ctx, VesselNetworkName, types.NetworkInspectOptions{})
	if err != nil {
		if client.IsErrNotFound(err) {
			_, err = m.dockerClient.NetworkCreate(ctx, VesselNetworkName, types.NetworkCreate{
				Driver: "bridge",
			})
			return err
		}
		return err
	}
	return nil
}

func (m *TraefikManager) createTraefikContainer(ctx context.Context) error {
	imageRef := "traefik:v3.0"
	out, err := m.dockerClient.ImagePull(ctx, imageRef, types.ImagePullOptions{})
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
	}

	resp, err := m.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
		Cmd:   cmdArgs,
		ExposedPorts: nat.PortSet{
			"80/tcp":   struct{}{},
			"443/tcp":  struct{}{},
			"8080/tcp": struct{}{},
		},
	}, hostConfig, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			VesselNetworkName: {},
		},
	}, nil, TraefikContainerName)

	if err != nil {
		return err
	}

	log.Printf("Created Traefik container ID: %s", resp.ID)
	return nil
}

func (m *TraefikManager) buildTraefikCmdArgs() []string {
	cmdArgs := []string{
		"--api.insecure=true", // Enable dashboard (do not expose in prod)
		"--providers.docker=true",
		"--providers.docker.exposedbydefault=false",
		"--providers.docker.network=" + VesselNetworkName,
		"--entrypoints.web.address=:80",
		"--entrypoints.websecure.address=:443",
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
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: "/var/run/docker.sock",
			Target: "/var/run/docker.sock",
		},
	}
	if m.tlsEmail != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: "vessel-traefik-acme",
			Target: "/letsencrypt",
		})
	}
	return mounts
}

func (m *TraefikManager) buildPortBindings() nat.PortMap {
	return nat.PortMap{
		"80/tcp":   []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "80"}},
		"443/tcp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "443"}},
		"8080/tcp": []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: "8080"}}, // Dashboard (local only)
	}
}
