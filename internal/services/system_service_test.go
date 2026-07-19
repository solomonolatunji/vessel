package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadCPUTimes(t *testing.T) {
	tmpDir := t.TempDir()
	statFile := filepath.Join(tmpDir, "stat")

	// Create mock stat file
	content := []byte(`cpu  100 200 300 400 500 600 700 800
cpu0 10 20 30 40 50 60 70 80
intr 12345
`)
	if err := os.WriteFile(statFile, content, 0644); err != nil {
		t.Fatalf("failed to write mock stat file: %v", err)
	}

	// Override variable
	originalProcStatPath := procStatPath
	procStatPath = statFile
	defer func() { procStatPath = originalProcStatPath }()

	idle, total := readCPUTimes()
	if idle != 400 {
		t.Errorf("expected idle 400, got %d", idle)
	}
	// Total = 100+200+300+400+500+600+700+800 = 3600
	if total != 3600 {
		t.Errorf("expected total 3600, got %d", total)
	}

	// Test error path (file does not exist)
	procStatPath = filepath.Join(tmpDir, "does-not-exist")
	idleErr, totalErr := readCPUTimes()
	if idleErr != 0 || totalErr != 0 {
		t.Errorf("expected 0, 0 on error, got %d, %d", idleErr, totalErr)
	}

	// Test error path (invalid content)
	invalidStatFile := filepath.Join(tmpDir, "stat_invalid")
	if err := os.WriteFile(invalidStatFile, []byte("cpu 100"), 0644); err != nil {
		t.Fatalf("failed to write mock stat file: %v", err)
	}
	procStatPath = invalidStatFile
	idleErr, totalErr = readCPUTimes()
	if idleErr != 0 || totalErr != 0 {
		t.Errorf("expected 0, 0 on invalid content, got %d, %d", idleErr, totalErr)
	}
}

func TestGetMemoryStats(t *testing.T) {
	tmpDir := t.TempDir()
	meminfoFile := filepath.Join(tmpDir, "meminfo")

	// Create mock meminfo file
	// Total: 16GB, Free: 2GB, Available: 4GB
	content := []byte(`MemTotal:       16777216 kB
MemFree:         2097152 kB
MemAvailable:    4194304 kB
Buffers:          100000 kB
Cached:           200000 kB
`)
	if err := os.WriteFile(meminfoFile, content, 0644); err != nil {
		t.Fatalf("failed to write mock meminfo file: %v", err)
	}

	// Override variable
	originalProcMeminfoPath := procMeminfoPath
	procMeminfoPath = meminfoFile
	defer func() { procMeminfoPath = originalProcMeminfoPath }()

	stats := getMemoryStats()
	// TotalMB = 16777216 / 1024 = 16384
	if stats.TotalMB != 16384 {
		t.Errorf("expected TotalMB 16384, got %d", stats.TotalMB)
	}
	// FreeMB (which comes from MemAvailable) = 4194304 / 1024 = 4096
	if stats.FreeMB != 4096 {
		t.Errorf("expected FreeMB 4096, got %d", stats.FreeMB)
	}
	// Used = Total - Available = 16777216 - 4194304 = 12582912 kB
	// UsedMB = 12582912 / 1024 = 12288
	if stats.UsedMB != 12288 {
		t.Errorf("expected UsedMB 12288, got %d", stats.UsedMB)
	}
	// Percent = 12288 / 16384 * 100 = 75.0
	if stats.Percent != 75.0 {
		t.Errorf("expected Percent 75.0, got %f", stats.Percent)
	}

	// Test error path (file does not exist)
	procMeminfoPath = filepath.Join(tmpDir, "does-not-exist")
	statsErr := getMemoryStats()
	if statsErr.TotalMB != 0 || statsErr.FreeMB != 0 || statsErr.UsedMB != 0 || statsErr.Percent != 0.0 {
		t.Errorf("expected zero struct on error, got %+v", statsErr)
	}
}
