package proxy

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/solomonolatunji/vessel/internal/types"
	"github.com/solomonolatunji/vessel/internal/utils"
)

// CaddyfileGenerator transforms project definitions and custom domain configurations into valid Caddy v2 configuration rules.
type CaddyfileGenerator struct {
	config *CaddyConfig
}

// NewCaddyfileGenerator creates a CaddyfileGenerator using the provided proxy configuration.
func NewCaddyfileGenerator(config *CaddyConfig) *CaddyfileGenerator {
	return &CaddyfileGenerator{config: config}
}

// Generate formats a Caddyfile string containing global TLS configurations, project/service routes, and custom domain proxies.
func (g *CaddyfileGenerator) Generate(projects []types.ProjectConfig, services []*types.AppServiceConfig, domains []types.DomainConfig) (string, error) {
	var buf bytes.Buffer

	buf.WriteString("{\n")
	if g.config.TLSEmail != "" {
		buf.WriteString(fmt.Sprintf("\temail %s\n", g.config.TLSEmail))
	} else {
		buf.WriteString("\tauto_https disable_redirects\n")
	}
	buf.WriteString("}\n\n")

	projectMap := make(map[string]*types.ProjectConfig)
	for i := range projects {
		projectMap[projects[i].ID] = &projects[i]
		g.writeProjectBlock(&buf, &projects[i])
	}

	serviceMap := make(map[string]*types.AppServiceConfig)
	for _, s := range services {
		if s == nil {
			continue
		}
		serviceMap[s.ID] = s
		g.writeAppServiceBlock(&buf, s)
	}

	for i := range domains {
		domainConfig := &domains[i]
		if s, ok := serviceMap[domainConfig.ProjectID]; ok {
			g.writeCustomDomainServiceBlock(&buf, domainConfig, s)
			continue
		}
		targetProject, ok := projectMap[domainConfig.ProjectID]
		if !ok {
			continue
		}
		g.writeCustomDomainBlock(&buf, domainConfig, targetProject)
	}

	return buf.String(), nil
}

func (g *CaddyfileGenerator) writeProjectBlock(buf *bytes.Buffer, p *types.ProjectConfig) {
	if p.Name == "" {
		return
	}

	containerHost := utils.NormalizeContainerName(p.ID)
	targetPort := 3000

	cleanName := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(p.Name), " ", "-"))
	hostnames := []string{fmt.Sprintf("http://%s.vessel.local", cleanName)}

	buf.WriteString(strings.Join(hostnames, ", ") + " {\n")
	buf.WriteString(fmt.Sprintf("\treverse_proxy %s:%d {\n", containerHost, targetPort))
	buf.WriteString("\t\theader_up Host {upstream_hostport}\n")
	buf.WriteString("\t\theader_up X-Real-IP {remote_host}\n")
	buf.WriteString("\t\theader_up X-Forwarded-Proto {scheme}\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n\n")
}

func (g *CaddyfileGenerator) writeAppServiceBlock(buf *bytes.Buffer, s *types.AppServiceConfig) {
	if s.Domain == "" && s.Name == "" {
		return
	}

	containerHost := utils.NormalizeContainerName(s.ID)
	targetPort := s.InternalPort
	if targetPort <= 0 {
		targetPort = 3000
	}

	var hostnames []string
	if s.Domain != "" {
		hostnames = append(hostnames, strings.TrimSpace(s.Domain))
	}
	if s.Name != "" {
		cleanName := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(s.Name), " ", "-"))
		hostnames = append(hostnames, fmt.Sprintf("http://%s.vessel.local", cleanName))
	}

	if len(hostnames) == 0 {
		return
	}

	buf.WriteString(strings.Join(hostnames, ", ") + " {\n")
	buf.WriteString(fmt.Sprintf("\treverse_proxy %s:%d {\n", containerHost, targetPort))
	buf.WriteString("\t\theader_up Host {upstream_hostport}\n")
	buf.WriteString("\t\theader_up X-Real-IP {remote_host}\n")
	buf.WriteString("\t\theader_up X-Forwarded-Proto {scheme}\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n\n")
}

func (g *CaddyfileGenerator) writeCustomDomainServiceBlock(buf *bytes.Buffer, d *types.DomainConfig, s *types.AppServiceConfig) {
	if d.DomainName == "" {
		return
	}

	containerHost := utils.NormalizeContainerName(s.ID)
	targetPort := s.InternalPort
	if targetPort <= 0 {
		targetPort = 3000
	}

	buf.WriteString(strings.TrimSpace(d.DomainName) + " {\n")
	if d.RedirectTo != "" {
		buf.WriteString(fmt.Sprintf("\tredir %s{uri} permanent\n", strings.TrimSpace(d.RedirectTo)))
		buf.WriteString("}\n\n")
		return
	}

	pathPrefix := d.PathPrefix
	if pathPrefix == "" {
		pathPrefix = "*"
	}

	buf.WriteString(fmt.Sprintf("\treverse_proxy %s %s:%d {\n", pathPrefix, containerHost, targetPort))
	buf.WriteString("\t\theader_up Host {upstream_hostport}\n")
	buf.WriteString("\t\theader_up X-Real-IP {remote_host}\n")
	buf.WriteString("\t\theader_up X-Forwarded-Proto {scheme}\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n\n")
}

func (g *CaddyfileGenerator) writeCustomDomainBlock(buf *bytes.Buffer, d *types.DomainConfig, p *types.ProjectConfig) {
	if d.DomainName == "" {
		return
	}

	containerHost := utils.NormalizeContainerName(p.ID)
	targetPort := 3000
	if targetPort <= 0 {
		targetPort = 3000
	}

	buf.WriteString(strings.TrimSpace(d.DomainName) + " {\n")
	if d.RedirectTo != "" {
		buf.WriteString(fmt.Sprintf("\tredir %s{uri} permanent\n", strings.TrimSpace(d.RedirectTo)))
		buf.WriteString("}\n\n")
		return
	}

	pathPrefix := d.PathPrefix
	if pathPrefix == "" {
		pathPrefix = "*"
	}

	buf.WriteString(fmt.Sprintf("\treverse_proxy %s %s:%d {\n", pathPrefix, containerHost, targetPort))
	buf.WriteString("\t\theader_up Host {upstream_hostport}\n")
	buf.WriteString("\t\theader_up X-Real-IP {remote_host}\n")
	buf.WriteString("\t\theader_up X-Forwarded-Proto {scheme}\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n\n")
}
