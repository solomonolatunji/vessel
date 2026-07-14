package models

type SystemStats struct {
	CPU       CPUStats    `json:"cpu"`
	Memory    MemoryStats `json:"memory"`
	Disk      DiskStats   `json:"disk"`
	Uptime    int64       `json:"uptimeSeconds"`
	LoadAvg   [3]float64  `json:"loadAvg"`
	Processes int         `json:"processes"`
}

type CPUStats struct {
	Percent float64 `json:"percent"`
	Cores   int     `json:"cores"`
}

type MemoryStats struct {
	TotalMB int64   `json:"totalMb"`
	UsedMB  int64   `json:"usedMb"`
	FreeMB  int64   `json:"freeMb"`
	Percent float64 `json:"percent"`
}

type DiskStats struct {
	TotalGB int64   `json:"totalGb"`
	UsedGB  int64   `json:"usedGb"`
	FreeGB  int64   `json:"freeGb"`
	Percent float64 `json:"percent"`
}
