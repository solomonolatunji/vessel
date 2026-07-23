package engine

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"codedock.dev/codedock/internal/utils"
)

func GetAvailablePort() (int, error) {
	startStr := os.Getenv("DEPLOY_HOST_PORT_START")
	endStr := os.Getenv("DEPLOY_HOST_PORT_END")

	start := 4100
	end := 4999

	if startStr != "" {
		if s, err := strconv.Atoi(startStr); err == nil {
			start = s
		}
	}
	if endStr != "" {
		if e, err := strconv.Atoi(endStr); err == nil {
			end = e
		}
	}

	if start > end {
		start, end = end, start
	}

	for port := start; port <= end; port++ {
		addr := fmt.Sprintf("0.0.0.0:%d", port)
		l, err := net.Listen("tcp", addr)
		if err == nil {
			l.Close()
			return port, nil
		}
	}

	return 0, &utils.NonReportableError{Message: fmt.Sprintf("no available ports found between %d and %d", start, end)}
}
