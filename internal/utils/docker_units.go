package utils

import (
	"fmt"
	"os"
	"strings"
)

func DefaultDBMemoryMB() int {
	if m := os.Getenv("CODEDOCK_DEFAULT_DB_MEMORY_MB"); m != "" {
		if v, err := parseInt(m); err == nil && v > 0 {
			return v
		}
	}
	return 1024
}

func DefaultDBCPURequest() float64 {
	if c := os.Getenv("CODEDOCK_DEFAULT_DB_CPU"); c != "" {
		if v, err := parseFloat(c); err == nil && v > 0 {
			return v
		}
	}
	return 1.0
}

func MegaBytesToBytes(mb int) int64 {
	if mb <= 0 {
		return 512 * 1024 * 1024
	}
	return int64(mb) * 1024 * 1024
}

func CPURequestToNanoCPUs(cores float64) int64 {
	if cores <= 0 {
		return 500_000_000
	}
	return int64(cores * 1_000_000_000)
}

func NormalizeContainerName(projectID string) string {
	return fmt.Sprintf("codedock-%s", strings.ToLower(strings.TrimSpace(projectID)))
}

func parseInt(s string) (int, error) {
	var v int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number: %s", s)
		}
		v = v*10 + int(c-'0')
	}
	return v, nil
}

func parseFloat(s string) (float64, error) {
	var v float64
	var dec bool
	var div float64 = 1
	for _, c := range s {
		if c == '.' && !dec {
			dec = true
			continue
		}
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number: %s", s)
		}
		v = v*10 + float64(c-'0')
		if dec {
			div *= 10
		}
	}
	return v / div, nil
}
