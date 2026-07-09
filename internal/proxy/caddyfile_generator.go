package proxy

import (
	"bytes"
	"fmt"
	"strings"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/utils"
)

type CaddyfileGenerator struct {
	config *CaddyConfig
}

func NewCaddyfileGenerator(config *CaddyConfig) *CaddyfileGenerator {
	return &CaddyfileGenerator{config: config}
}

func (g *CaddyfileGenerator) Generate(projects []models.ProjectConfig, services []models.AppService, domains []models.DomainConfig) (string, error) {
	var buf bytes.Buffer

	buf.WriteString("{\n")
	if g.config.TLSEmail != "" {
		buf.WriteString(fmt.Sprintf("\temail %s\n", g.config.TLSEmail))
	} else {
		buf.WriteString("\tauto_https disable_redirects\n")
	}
	buf.WriteString("}\n\n")

	projectMap := make(map[string]*models.ProjectConfig)
	for i := range projects {
		projectMap[projects[i].ID] = &projects[i]
		g.writeProjectBlock(&buf, &projects[i])
	}

	serviceMap := make(map[string]*models.AppService)
	for i := range services {
		s := &services[i]
		serviceMap[s.ID] = s
		g.writeAppServiceBlock(&buf, s)
	}

	for i := range domains {
		domainConfig := &domains[i]
		var s *models.AppService
		if val, ok := serviceMap[domainConfig.ProjectID]; ok {
			s = val
		}

		targetProject, ok := projectMap[domainConfig.ProjectID]
		if !ok && s == nil {
			continue
		}
		g.writeDomainBlock(&buf, domainConfig, targetProject, s)
	}

	return buf.String(), nil
}

func (g *CaddyfileGenerator) writeProjectBlock(buf *bytes.Buffer, p *models.ProjectConfig) {
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

func (g *CaddyfileGenerator) writeAppServiceBlock(buf *bytes.Buffer, s *models.AppService) {
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

func (g *CaddyfileGenerator) writeDomainBlock(buf *bytes.Buffer, d *models.DomainConfig, p *models.ProjectConfig, s *models.AppService) {
	if d.DomainName == "" {
		return
	}

	var containerHost string
	var targetPort int

	if s != nil {
		containerHost = utils.NormalizeContainerName(s.ID)
		targetPort = s.InternalPort
	} else if p != nil {
		containerHost = utils.NormalizeContainerName(p.ID)
		targetPort = 3000
	} else {
		return
	}

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
	buf.WriteString("}\n")
}
