package engine

import (
	"context"
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"vessl.dev/vessl/internal/utils"
)

type ContainerHealth struct {
	Status             string  `json:"status"`
	CPUUsagePercentage float64 `json:"cpuUsagePercentage"`
	MemoryUsageBytes   int64   `json:"memoryUsageBytes"`
	MemoryLimitBytes   int64   `json:"memoryLimitBytes"`
	UptimeSeconds      int64   `json:"uptimeSeconds"`
}

type StatsMonitor struct {
	dockerClient *client.Client
}

func NewStatsMonitor(dockerClient *client.Client) *StatsMonitor {
	return &StatsMonitor{dockerClient: dockerClient}
}

func (s *StatsMonitor) GetHealth(ctx context.Context, containerIDOrName string) (*ContainerHealth, error) {
	inspectResp, err := s.dockerClient.ContainerInspect(ctx, containerIDOrName)
	if err != nil {
		return &ContainerHealth{Status: "offline"}, utils.NewEngineError("ContainerInspect", err)
	}
	if !inspectResp.State.Running {
		return &ContainerHealth{Status: "stopped"}, nil
	}
	statsResp, err := s.dockerClient.ContainerStatsOneShot(ctx, containerIDOrName)
	if err != nil {
		return nil, utils.NewEngineError("ContainerStatsOneShot", err)
	}
	defer statsResp.Body.Close()
	var stats types.StatsJSON
	if err := json.NewDecoder(statsResp.Body).Decode(&stats); err != nil {
		return nil, utils.NewEngineError("DecodeStats", err)
	}
	cpuPercent := CalculateCPUPercentage(&stats)
	memoryUsage := stats.MemoryStats.Usage
	if cache, exists := stats.MemoryStats.Stats["cache"]; exists && cache <= memoryUsage {
		memoryUsage -= cache
	}
	startedAt, _ := time.Parse(time.RFC3339Nano, inspectResp.State.StartedAt)
	uptimeSeconds := int64(time.Since(startedAt).Seconds())
	if startedAt.IsZero() {
		uptimeSeconds = 0
	}
	return &ContainerHealth{
		Status:             "running",
		CPUUsagePercentage: cpuPercent,
		MemoryUsageBytes:   int64(memoryUsage),
		MemoryLimitBytes:   int64(stats.MemoryStats.Limit),
		UptimeSeconds:      uptimeSeconds,
	}, nil
}

func CalculateCPUPercentage(stats *types.StatsJSON) float64 {
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
