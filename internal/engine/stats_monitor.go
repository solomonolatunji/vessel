package engine

import (
	"context"
	"encoding/json"
	"time"

	"codedock.dev/codedock/internal/utils"
	"github.com/docker/docker/client"
)

type dockerStats struct {
	CPUStats    cpuStats    `json:"cpu_stats"`
	PreCPUStats cpuStats    `json:"precpu_stats"`
	MemoryStats memoryStats `json:"memory_stats"`
}

type cpuStats struct {
	CPUUsage struct {
		TotalUsage  float64   `json:"total_usage"`
		PercpuUsage []float64 `json:"percpu_usage"`
	} `json:"cpu_usage"`
	SystemUsage float64 `json:"system_cpu_usage"`
	OnlineCPUs  float64 `json:"online_cpus"`
}

type memoryStats struct {
	Usage float64            `json:"usage"`
	Limit float64            `json:"limit"`
	Stats map[string]float64 `json:"stats"`
}

type ContainerHealthStatus string

const (
	ContainerHealthStatusRunning     ContainerHealthStatus = "running"
	ContainerHealthStatusStopped     ContainerHealthStatus = "stopped"
	ContainerHealthStatusOffline     ContainerHealthStatus = "offline"
	ContainerHealthStatusNotDeployed ContainerHealthStatus = "not_deployed"
)

type ContainerHealth struct {
	Status             ContainerHealthStatus `json:"status"`
	CPUUsagePercentage float64               `json:"cpuUsagePercentage"`
	MemoryUsageBytes   int64                 `json:"memoryUsageBytes"`
	MemoryLimitBytes   int64                 `json:"memoryLimitBytes"`
	UptimeSeconds      int64                 `json:"uptimeSeconds"`
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
		return &ContainerHealth{Status: ContainerHealthStatusOffline}, utils.NewEngineError("ContainerInspect", err)
	}
	if !inspectResp.State.Running {
		return &ContainerHealth{Status: ContainerHealthStatusStopped}, nil
	}
	statsResp, err := s.dockerClient.ContainerStatsOneShot(ctx, containerIDOrName)
	if err != nil {
		return nil, utils.NewEngineError("ContainerStatsOneShot", err)
	}
	defer statsResp.Body.Close()
	var stats dockerStats
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
		Status:             ContainerHealthStatusRunning,
		CPUUsagePercentage: cpuPercent,
		MemoryUsageBytes:   int64(memoryUsage),
		MemoryLimitBytes:   int64(stats.MemoryStats.Limit),
		UptimeSeconds:      uptimeSeconds,
	}, nil
}

func CalculateCPUPercentage(stats *dockerStats) float64 {
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
