package proxy

import (
	"strings"
	"testing"

	"vessel.dev/vessel/internal/models"
)

func TestCaddyfileGenerator(t *testing.T) {
	config := NewCaddyConfig("testdata", "ops@vessel.local")
	gen := NewCaddyfileGenerator(config)

	projects := []models.ProjectConfig{
		{
			ID:   "test-id-123",
			Name: "Frontend App",
		},
		{
			ID:   "api-project-456",
			Name: "Go Backend API",
		},
	}

	services := []models.AppService{
		{
			ID:           "test-id-123",
			ProjectID:    "test-id-123",
			Name:         "Frontend App",
			Domain:       "app.solomon.com",
			InternalPort: 3000,
		},
		{
			ID:           "api-project-456",
			ProjectID:    "api-project-456",
			Name:         "Go Backend API",
			InternalPort: 8080,
		},
	}

	domains := []models.DomainConfig{
		{
			ID:         "dom-1",
			ProjectID:  "test-id-123",
			DomainName: "legacy.solomon.com",
			RedirectTo: "https://app.solomon.com",
		},
		{
			ID:         "dom-2",
			ProjectID:  "api-project-456",
			DomainName: "api.solomon.com",
			PathPrefix: "/v1/*",
		},
	}

	caddyfile, err := gen.Generate(projects, services, domains)
	if err != nil {
		t.Fatalf("expected no error generating caddyfile, got: %v", err)
	}

	if !strings.Contains(caddyfile, "email ops@vessel.local") {
		t.Errorf("expected global email block, got:\n%s", caddyfile)
	}
	if !strings.Contains(caddyfile, "app.solomon.com, http://frontend-app.vessel.local {") {
		t.Errorf("expected project hostnames block, got:\n%s", caddyfile)
	}
	if !strings.Contains(caddyfile, "reverse_proxy vessel-test-id-123:3000") {
		t.Errorf("expected upstream container proxying, got:\n%s", caddyfile)
	}
	if !strings.Contains(caddyfile, "legacy.solomon.com {") || !strings.Contains(caddyfile, "redir https://app.solomon.com{uri} permanent") {
		t.Errorf("expected custom domain redirect block, got:\n%s", caddyfile)
	}
	if !strings.Contains(caddyfile, "api.solomon.com {") || !strings.Contains(caddyfile, "reverse_proxy /v1/* vessel-api-project-456:8080") {
		t.Errorf("expected path prefix reverse proxy block, got:\n%s", caddyfile)
	}
}
