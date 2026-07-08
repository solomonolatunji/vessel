package utils

import (
	"fmt"
	"os"
	"strings"
)

// GenerateSslipDomain synthesizes an automatic sslip.io wildcard domain name using the server host IP and project identifier.
func GenerateSslipDomain(projectNameOrID string, hostIP string) string {
	if hostIP == "" {
		hostIP = os.Getenv("VESSEL_HOST_IP")
	}
	if hostIP == "" {
		hostIP = "127.0.0.1"
	}

	cleanIP := strings.ReplaceAll(strings.TrimSpace(hostIP), ".", "-")
	cleanName := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(projectNameOrID), " ", "-"))
	if len(cleanName) > 32 {
		cleanName = cleanName[:32]
	}
	cleanName = strings.Trim(cleanName, "-")

	return fmt.Sprintf("http://%s.%s.sslip.io", cleanName, cleanIP)
}
