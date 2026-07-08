package types

type SystemInfo struct {
	Version         string `json:"version"`
	GoVersion       string `json:"goVersion"`
	DockerVersion   string `json:"dockerVersion"`
	CaddyVersion    string `json:"caddyVersion"`
	OS              string `json:"os"`
	Arch            string `json:"arch"`
	TotalMemoryMB   int64  `json:"totalMemoryMB"`
	FreeMemoryMB    int64  `json:"freeMemoryMB"`
	CPUCores        int    `json:"cpuCores"`
	UpdateAvailable bool   `json:"updateAvailable"`
	LatestVersion   string `json:"latestVersion,omitempty"`
}
