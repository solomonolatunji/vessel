package services

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"vessl.dev/vessl/internal/models"
)

type SystemService struct {
	prevIdle  uint64
	prevTotal uint64
	prevTime  time.Time
}

func NewSystemService() *SystemService {
	return &SystemService{}
}

func (s *SystemService) GetStats() (*models.SystemStats, error) {
	cpu := s.getCPUStats()
	mem := getMemoryStats()
	disk := getDiskStats()
	uptime := getUptime()
	load := getLoadAvg()
	procs := getProcessCount()

	return &models.SystemStats{
		CPU:       cpu,
		Memory:    mem,
		Disk:      disk,
		Uptime:    uptime,
		LoadAvg:   load,
		Processes: procs,
	}, nil
}

func (s *SystemService) getCPUStats() models.CPUStats {
	cores := runtime.NumCPU()
	idle, total := readCPUTimes()

	if s.prevTotal > 0 {
		deltaIdle := idle - s.prevIdle
		deltaTotal := total - s.prevTotal
		elapsed := time.Since(s.prevTime).Seconds()

		var percent float64
		if deltaTotal > 0 && elapsed > 0 {
			percent = (1.0 - float64(deltaIdle)/float64(deltaTotal)) * 100.0
			percent = math.Round(percent*10) / 10
		}

		s.prevIdle = idle
		s.prevTotal = total
		s.prevTime = time.Now()

		return models.CPUStats{Percent: percent, Cores: cores}
	}

	s.prevIdle = idle
	s.prevTotal = total
	s.prevTime = time.Now()

	return models.CPUStats{Percent: 0, Cores: cores}
}

func readCPUTimes() (uint64, uint64) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			return 0, 0
		}
		var total uint64
		for _, f := range fields[1:] {
			v, _ := strconv.ParseUint(f, 10, 64)
			total += v
		}
		idle, _ := strconv.ParseUint(fields[4], 10, 64)
		return idle, total
	}
	return 0, 0
}

func getMemoryStats() models.MemoryStats {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return models.MemoryStats{}
	}

	var totalMem, freeMem, availMem uint64
	for _, line := range strings.Split(string(data), "\n") {
		switch {
		case strings.HasPrefix(line, "MemTotal:"):
			fmt.Sscanf(line, "MemTotal: %d kB", &totalMem)
		case strings.HasPrefix(line, "MemFree:"):
			fmt.Sscanf(line, "MemFree: %d kB", &freeMem)
		case strings.HasPrefix(line, "MemAvailable:"):
			fmt.Sscanf(line, "MemAvailable: %d kB", &availMem)
		}
	}

	used := totalMem - availMem
	if used > totalMem {
		used = totalMem - freeMem
	}

	totalMB := int64(totalMem / 1024)
	usedMB := int64(used / 1024)
	freeMB := int64(availMem / 1024)

	var percent float64
	if totalMB > 0 {
		percent = float64(usedMB) / float64(totalMB) * 100
		percent = math.Round(percent*10) / 10
	}

	return models.MemoryStats{
		TotalMB: totalMB,
		UsedMB:  usedMB,
		FreeMB:  freeMB,
		Percent: percent,
	}
}

func getDiskStats() models.DiskStats {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return models.DiskStats{}
	}

	total := int64(stat.Blocks * uint64(stat.Bsize) / (1024 * 1024 * 1024))
	free := int64(stat.Bavail * uint64(stat.Bsize) / (1024 * 1024 * 1024))
	used := total - free

	var percent float64
	if total > 0 {
		percent = float64(used) / float64(total) * 100
		percent = math.Round(percent*10) / 10
	}

	return models.DiskStats{
		TotalGB: total,
		UsedGB:  used,
		FreeGB:  free,
		Percent: percent,
	}
}

func getUptime() int64 {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	parts := strings.Fields(string(data))
	if len(parts) == 0 {
		return 0
	}
	secs, _ := strconv.ParseFloat(parts[0], 64)
	return int64(secs)
}

func getLoadAvg() [3]float64 {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return [3]float64{}
	}
	parts := strings.Fields(string(data))
	if len(parts) < 3 {
		return [3]float64{}
	}
	var load [3]float64
	for i := 0; i < 3; i++ {
		load[i], _ = strconv.ParseFloat(parts[i], 64)
	}
	return load
}

func getProcessCount() int {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if e.IsDir() {
			name := e.Name()
			if len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
				count++
			}
		}
	}
	return count
}
