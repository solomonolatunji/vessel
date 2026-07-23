package engine

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type MetricsWorker struct {
	dockerClient *client.Client
	statsMonitor *StatsMonitor
	tsdbURL      string
	httpClient   *http.Client
}

func NewMetricsWorker(cli *client.Client) *MetricsWorker {
	return &MetricsWorker{
		dockerClient: cli,
		statsMonitor: NewStatsMonitor(cli),
		tsdbURL:      "http://127.0.0.1:8428/api/v1/import/prometheus",
		httpClient:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (w *MetricsWorker) Start() {
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			w.collectAndPush()
		}
	}()
}

func (w *MetricsWorker) collectAndPush() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	containers, err := w.dockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		slog.Error("metrics worker failed to list containers", "err", err)
		return
	}

	var sb strings.Builder
	timestamp := time.Now().UnixMilli()

	for _, c := range containers {
		if len(c.Names) > 0 {
			name := strings.TrimPrefix(c.Names[0], "/")
			if name == TraefikContainerName || name == TSDBContainerName {
				continue
			}
		}

		health, err := w.statsMonitor.GetHealth(ctx, c.ID)
		if err != nil {
			continue
		}

		serviceID := c.Labels["codedock.service.id"]
		if serviceID == "" {
			serviceID = "unknown"
		}

		containerName := "unknown"
		if len(c.Names) > 0 {
			containerName = strings.TrimPrefix(c.Names[0], "/")
		}

		tags := fmt.Sprintf(`container_id="%s",container_name="%s",service_id="%s"`, c.ID, containerName, serviceID)

		sb.WriteString(fmt.Sprintf("container_cpu_usage_percent{%s} %f %d\n", tags, health.CPUUsagePercentage, timestamp))
		sb.WriteString(fmt.Sprintf("container_memory_usage_bytes{%s} %d %d\n", tags, health.MemoryUsageBytes, timestamp))
		sb.WriteString(fmt.Sprintf("container_memory_limit_bytes{%s} %d %d\n", tags, health.MemoryLimitBytes, timestamp))
	}

	if sb.Len() == 0 {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.tsdbURL, bytes.NewBufferString(sb.String()))
	if err != nil {
		return
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		slog.Warn("metrics worker failed to push to tsdb", "err", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		slog.Warn("metrics worker received error from tsdb", "status", resp.StatusCode)
	}
}
