package engine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/robfig/cron/v3"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/utils"
)

type CronManager struct {
	dockerClient *client.Client
	store        CronManagerStore
	cronEngine   *cron.Cron
	entries      map[string]cron.EntryID
	mu           sync.Mutex
}

func NewCronManager(dockerClient *client.Client, s CronManagerStore) *CronManager {
	return &CronManager{
		dockerClient: dockerClient,
		store:        s,
		cronEngine:   cron.New(cron.WithSeconds()),
		entries:      make(map[string]cron.EntryID),
	}
}

func (cm *CronManager) Start() error {
	scheduledTasks, err := cm.store.ListScheduledTasks()
	if err != nil {
		return fmt.Errorf("failed to load scheduledTasks during cron manager start: %w", err)
	}
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for _, j := range scheduledTasks {
		if j.Status == "active" {
			scheduledTaskCopy := j
			if err := cm.registerScheduledTaskLocked(&scheduledTaskCopy); err != nil {
				slog.Warn("failed to register cron scheduledTask", "name", scheduledTaskCopy.Name, "id", scheduledTaskCopy.ID, "err", err)
			}
		}
	}
	cm.cronEngine.Start()
	slog.Info("cron manager started")
	return nil
}

func (cm *CronManager) Stop() {
	if cm.cronEngine != nil {
		cm.cronEngine.Stop()
	}
}

func (cm *CronManager) RegisterScheduledTask(j *models.ScheduledTask) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.registerScheduledTaskLocked(j)
}

func (cm *CronManager) registerScheduledTaskLocked(j *models.ScheduledTask) error {
	if entryID, exists := cm.entries[j.ID]; exists {
		cm.cronEngine.Remove(entryID)
		delete(cm.entries, j.ID)
	}
	if j.Status != "active" {
		return nil
	}
	schedule := strings.TrimSpace(j.Schedule)
	if len(strings.Fields(schedule)) == 5 && !strings.HasPrefix(schedule, "@") {
		schedule = "0 " + schedule
	}
	scheduledTaskID := j.ID
	entryID, err := cm.cronEngine.AddFunc(schedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		_, _ = cm.ExecuteScheduledTask(ctx, scheduledTaskID)
	})
	if err != nil {
		return fmt.Errorf("invalid cron schedule '%s': %w", j.Schedule, err)
	}
	cm.entries[j.ID] = entryID
	return nil
}

func (cm *CronManager) UnregisterScheduledTask(scheduledTaskID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if entryID, exists := cm.entries[scheduledTaskID]; exists {
		cm.cronEngine.Remove(entryID)
		delete(cm.entries, scheduledTaskID)
	}
}

func (cm *CronManager) ExecuteScheduledTask(ctx context.Context, scheduledTaskID string) (string, error) {
	j, err := cm.store.GetScheduledTask(scheduledTaskID)
	if err != nil || j == nil {
		return "", fmt.Errorf("scheduledTask %s not found: %w", scheduledTaskID, err)
	}
	app, err := cm.store.GetAppService(j.ServiceID)
	if err != nil || app == nil {
		return "", fmt.Errorf("service %s for scheduledTask %s not found: %w", j.ServiceID, scheduledTaskID, err)
	}
	containerName := utils.NormalizeContainerName(app.ID)
	inspectResp, err := cm.dockerClient.ContainerInspect(ctx, containerName)
	if err != nil || !inspectResp.State.Running {
		errMsg := fmt.Sprintf("cannot run scheduledTask: service container %s (%s) is stopped or not found", app.Name, containerName)
		now := time.Now()
		_ = cm.store.UpdateScheduledTaskStatusAndOutput(scheduledTaskID, models.ScheduledTaskStatusError, &now, errMsg)
		return errMsg, errors.New(errMsg)
	}
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"sh", "-c", j.Command},
	}
	execCreateResp, err := cm.dockerClient.ContainerExecCreate(ctx, inspectResp.ID, execConfig)
	if err != nil {
		now := time.Now()
		_ = cm.store.UpdateScheduledTaskStatusAndOutput(scheduledTaskID, models.ScheduledTaskStatusError, &now, err.Error())
		return "", fmt.Errorf("failed to create container exec for scheduledTask %s: %w", j.Name, err)
	}
	attachResp, err := cm.dockerClient.ContainerExecAttach(ctx, execCreateResp.ID, container.ExecAttachOptions{})
	if err != nil {
		now := time.Now()
		_ = cm.store.UpdateScheduledTaskStatusAndOutput(scheduledTaskID, models.ScheduledTaskStatusError, &now, err.Error())
		return "", fmt.Errorf("failed to attach to container exec for scheduledTask %s: %w", j.Name, err)
	}
	defer attachResp.Close()
	var stdoutBuf, stderrBuf bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attachResp.Reader); err != nil {
		_, _ = io.Copy(&stdoutBuf, attachResp.Reader)
	}
	output := stdoutBuf.String()
	if stderrBuf.Len() > 0 {
		if output != "" {
			output += "\nSTDERR:\n"
		}
		output += stderrBuf.String()
	}
	now := time.Now()
	_ = cm.store.UpdateScheduledTaskStatusAndOutput(scheduledTaskID, models.ScheduledTaskStatusActive, &now, output)
	return output, nil
}

func (cm *CronManager) ScheduleDockerCleanup(schedule string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cleanSchedule := strings.TrimSpace(schedule)
	if len(strings.Fields(cleanSchedule)) == 5 && !strings.HasPrefix(cleanSchedule, "@") {
		cleanSchedule = "0 " + cleanSchedule
	}

	if entryID, exists := cm.entries["docker-cleanup"]; exists {
		cm.cronEngine.Remove(entryID)
		delete(cm.entries, "docker-cleanup")
	}

	entryID, err := cm.cronEngine.AddFunc(cleanSchedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		cm.DockerCleanup(ctx)
	})
	if err != nil {
		return fmt.Errorf("invalid docker cleanup schedule '%s': %w", schedule, err)
	}
	cm.entries["docker-cleanup"] = entryID
	slog.Info("docker cleanup scheduled", "schedule", schedule)
	return nil
}

func (cm *CronManager) DockerCleanup(ctx context.Context) {
	if cm.dockerClient == nil {
		slog.Warn("docker cleanup skipped", "reason", "no Docker client")
		return
	}
	slog.Info("running docker cleanup")

	reclaimed := uint64(0)

	report, err := cm.dockerClient.ContainersPrune(ctx, filters.NewArgs(filters.Arg("until", "24h")))
	if err == nil {
		reclaimed += report.SpaceReclaimed
	}

	imgReport, err := cm.dockerClient.ImagesPrune(ctx, filters.NewArgs(filters.Arg("dangling", "true")))
	if err == nil {
		reclaimed += imgReport.SpaceReclaimed
	}

	volReport, err := cm.dockerClient.VolumesPrune(ctx, filters.NewArgs(filters.Arg("until", "24h")))
	if err == nil {
		reclaimed += volReport.SpaceReclaimed
	}

	if reclaimed > 0 {
		slog.Info("docker cleanup reclaimed space", "bytes", reclaimed)
	} else {
		slog.Info("docker cleanup completed", "result", "nothing to clean")
	}
}

func (cm *CronManager) ScheduleDiskUsageCheck(schedule string, threshold int) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cleanSchedule := strings.TrimSpace(schedule)
	if len(strings.Fields(cleanSchedule)) == 5 && !strings.HasPrefix(cleanSchedule, "@") {
		cleanSchedule = "0 " + cleanSchedule
	}

	if entryID, exists := cm.entries["disk-usage"]; exists {
		cm.cronEngine.Remove(entryID)
		delete(cm.entries, "disk-usage")
	}

	entryID, err := cm.cronEngine.AddFunc(cleanSchedule, func() {
		usageCmd := exec.Command("df", "--output=pcent", "/")
		out, err := usageCmd.Output()
		if err != nil {
			return
		}
		pctStr := strings.TrimSpace(string(out))
		pctStr = strings.TrimSuffix(strings.TrimSpace(strings.Split(pctStr, "\n")[1]), "%")
		usage, err := strconv.Atoi(pctStr)
		if err != nil {
			return
		}
		if usage > threshold {
			slog.Warn("disk usage above threshold", "usage_pct", usage, "threshold_pct", threshold)
		}
	})
	if err != nil {
		return fmt.Errorf("invalid disk usage check schedule '%s': %w", schedule, err)
	}
	cm.entries["disk-usage"] = entryID
	return nil
}
