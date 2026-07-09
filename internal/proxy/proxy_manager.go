package proxy

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"vessel.dev/vessel/internal/models"
)

type ProjectLister interface {
	ListProjects() ([]models.ProjectConfig, error)
}

type AppServiceLister interface {
	ListServices() ([]models.AppService, error)
}

type DomainLister interface {
	ListAllDomains() ([]models.DomainConfig, error)
}

type ProxyManager struct {
	config    *CaddyConfig
	generator *CaddyfileGenerator
	projects  ProjectLister
	services  AppServiceLister
	domains   DomainLister
	docker    *client.Client
}

func NewProxyManager(config *CaddyConfig, projects ProjectLister, services AppServiceLister, domains DomainLister, docker *client.Client) *ProxyManager {
	return &ProxyManager{
		config:    config,
		generator: NewCaddyfileGenerator(config),
		projects:  projects,
		services:  services,
		domains:   domains,
		docker:    docker,
	}
}

func (m *ProxyManager) Reload(ctx context.Context) error {
	projects, err := m.projects.ListProjects()
	if err != nil {
		return fmt.Errorf("failed to load active projects for caddy reload: %w", err)
	}

	services, _ := m.services.ListServices()

	domains, err := m.domains.ListAllDomains()
	if err != nil {
		return fmt.Errorf("failed to load custom domains for caddy reload: %w", err)
	}

	caddyfileContent, err := m.generator.Generate(projects, services, domains)
	if err != nil {
		return fmt.Errorf("failed to generate caddyfile syntax: %w", err)
	}

	if err := os.WriteFile(m.config.CaddyfilePath, []byte(caddyfileContent), 0644); err != nil {
		return fmt.Errorf("failed to persist caddyfile to %s: %w", m.config.CaddyfilePath, err)
	}

	if err := m.reloadAdminAPI(ctx, []byte(caddyfileContent)); err == nil {
		return nil
	}

	if m.docker != nil {
		if err := m.reloadDockerContainer(ctx); err == nil {
			return nil
		}
	}

	return nil
}

func (m *ProxyManager) reloadAdminAPI(ctx context.Context, caddyfileBytes []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.config.AdminAPIEndpoint, bytes.NewReader(caddyfileBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/caddyfile")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("caddy admin api returned unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (m *ProxyManager) reloadDockerContainer(ctx context.Context) error {
	execConfig := types.ExecConfig{
		Cmd:          []string{"caddy", "reload", "--config", m.config.CaddyfilePath, "--adapter", "caddyfile"},
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := m.docker.ContainerExecCreate(ctx, m.config.DockerContainer, execConfig)
	if err != nil {
		return err
	}

	return m.docker.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{})
}
