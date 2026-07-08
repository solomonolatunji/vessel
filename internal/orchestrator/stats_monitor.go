package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	vesseltypes "github.com/solomonolatunji/vessel/internal/types"
)

// StatsMonitor polls real-time container CPU, RAM, and uptime metrics using the Docker engine statistics API.
type StatsMonitor struct {
	dockerClient *client.Client
}

// NewStatsMonitor initializes a StatsMonitor wired to the provided Docker daemon client.
func NewStatsMonitor(dockerClient *client.Client) *StatsMonitor {
	return &StatsMonitor{dockerClient: dockerClient}
}

// GetHealth fetches one-shot statistical metrics for a running container and returns a formatted ContainerHealth record.
func (s *StatsMonitor) GetHealth(ctx context.Context, containerIDOrName string) (*vesseltypes.ContainerHealth, error) {
	inspectResp, err := s.dockerClient.ContainerInspect(ctx, containerIDOrName)
	if err != nil {
		return &vesseltypes.ContainerHealth{Status: "offline"}, fmt.Errorf("container inspect failed: %w", err)
	}

	if !inspectResp.State.Running {
		return &vesseltypes.ContainerHealth{Status: "stopped"}, nil
	}

	statsResp, err := s.dockerClient.ContainerStatsOneShot(ctx, containerIDOrName)
	if err != nil {
		return nil, fmt.Errorf("container stats failed: %w", err)
	}
	defer statsResp.Body.Close()

	var stats types.StatsJSON
	if err := json.NewDecoder(statsResp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode stats json: %w", err)
	}

	cpuPercent := calculateCPUPercentage(&stats)
	memoryUsage := stats.MemoryStats.Usage - stats.MemoryStats.Stats["cache"]
	if memoryUsage < 0 {
		memoryUsage = stats.MemoryStats.Usage
	}

	startedAt, _ := time.Parse(time.RFC3339Nano, inspectResp.State.StartedAt)
	uptimeSeconds := int64(time.Since(startedAt).Seconds())
	if startedAt.IsZero() {
		uptimeSeconds = 0
	}

	return &vesseltypes.ContainerHealth{
		Status:             "running",
		CPUUsagePercentage: cpuPercent,
		MemoryUsageBytes:   int64(memoryUsage),
		MemoryLimitBytes:   int64(stats.MemoryStats.Limit),
		UptimeSeconds:      uptimeSeconds,
	}, nil
}

func calculateCPUPercentage(stats *types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuCores := float64(stats.CPUStats.OnlineCPUs)
		if cpuCores == 0.0 {
			cpuCores = float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
		}
		if cpuCores == 0.0 {
			cpuCores = 1.0
		}
		return (cpuDelta / systemDelta) * cpuCores * 100.0
	}
	return 0.0
}
