package engine

import (
	"strings"

	"github.com/docker/docker/api/types/container"
)

func ApplyCustomDNS(hostCfg *container.HostConfig, customDNS string) {
	if strings.TrimSpace(customDNS) == "" {
		return
	}
	parts := strings.Split(customDNS, ",")
	var dnsList []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			dnsList = append(dnsList, p)
		}
	}
	if len(dnsList) > 0 {
		hostCfg.DNS = dnsList
	}
}
