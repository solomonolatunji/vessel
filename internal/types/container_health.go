package types

// ContainerHealth tracks real-time resource utilization and status for running containers.
type ContainerHealth struct {
	Status             string  `json:"status"`
	CPUUsagePercentage float64 `json:"cpuUsagePercentage"`
	MemoryUsageBytes   int64   `json:"memoryUsageBytes"`
	MemoryLimitBytes   int64   `json:"memoryLimitBytes"`
	UptimeSeconds      int64   `json:"uptimeSeconds"`
}
