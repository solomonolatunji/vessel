package utils

import (
	"fmt"
	"os"
	"strings"
)

func GenerateAppDomain(projectNameOrID string, hostIP string, wildcardDomain string) string {
	if wildcardDomain == "" {
		wildcardDomain = os.Getenv("VESSL_WILDCARD_DOMAIN")
	}

	if wildcardDomain != "" {
		cleanName := SanitizeDomainName(projectNameOrID)
		base := strings.TrimPrefix(wildcardDomain, "*.")
		if !strings.HasPrefix(base, "http") {
			return fmt.Sprintf("https://%s.%s", cleanName, base)
		}
		return fmt.Sprintf("%s://%s.%s", parseScheme(base), cleanName, parseHost(base))
	}

	if hostIP == "" {
		hostIP = os.Getenv("VESSL_HOST_IP")
	}

	if hostIP == "" {
		hostIP = "127.0.0.1"
	}

	cleanIP := strings.ReplaceAll(strings.TrimSpace(hostIP), ".", "-")
	cleanName := SanitizeDomainName(projectNameOrID)
	magicDomain := os.Getenv("VESSL_MAGIC_DOMAIN")
	if magicDomain == "" {
		magicDomain = "sslip.io"
	}

	return fmt.Sprintf("http://%s.%s.%s", cleanName, cleanIP, magicDomain)
}

func parseScheme(url string) string {
	if strings.HasPrefix(url, "https://") {
		return "https"
	}
	return "http"
}

func parseHost(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	return url
}

func SanitizeDomainName(name string) string {
	clean := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-"))
	if len(clean) > 32 {
		clean = clean[:32]
	}
	return strings.Trim(clean, "-")
}
