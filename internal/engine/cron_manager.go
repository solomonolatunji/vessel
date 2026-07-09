package orchestrator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/robfig/cron/v3"
	"vessel.dev/vessel/internal/job"
	"vessel.dev/vessel/internal/utils"
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
	jobs, err := cm.store.ListJobs()
	if err != nil {
		return fmt.Errorf("failed to load jobs during cron manager start: %w", err)
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, j := range jobs {
		if j.Status == "active" {
			jobCopy := j
			if err := cm.registerJobLocked(&jobCopy); err != nil {
				log.Printf("⚠️ Failed to register cron job %s (%s): %v", jobCopy.Name, jobCopy.ID, err)
			}
		}
	}

	cm.cronEngine.Start()
	log.Println("⏰ CronManager started and executing background tasks")
	return nil
}

func (cm *CronManager) Stop() {
	if cm.cronEngine != nil {
		cm.cronEngine.Stop()
	}
}

func (cm *CronManager) RegisterJob(j *job.Job) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.registerJobLocked(j)
}

func (cm *CronManager) registerJobLocked(j *job.Job) error {
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

	jobID := j.ID
	entryID, err := cm.cronEngine.AddFunc(schedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		_, _ = cm.ExecuteJob(ctx, jobID)
	})
	if err != nil {
		return fmt.Errorf("invalid cron schedule '%s': %w", j.Schedule, err)
	}

	cm.entries[j.ID] = entryID
	return nil
}

func (cm *CronManager) UnregisterJob(jobID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if entryID, exists := cm.entries[jobID]; exists {
		cm.cronEngine.Remove(entryID)
		delete(cm.entries, jobID)
	}
}

func (cm *CronManager) ExecuteJob(ctx context.Context, jobID string) (string, error) {
	j, err := cm.store.GetJob(jobID)
	if err != nil || j == nil {
		return "", fmt.Errorf("job %s not found: %w", jobID, err)
	}

	project, err := cm.store.GetProject(j.ProjectID)
	if err != nil || project == nil {
		return "", fmt.Errorf("project %s for job %s not found: %w", j.ProjectID, jobID, err)
	}

	containerName := utils.NormalizeContainerName(project.ID)
	inspectResp, err := cm.dockerClient.ContainerInspect(ctx, containerName)
	if err != nil || !inspectResp.State.Running {
		errMsg := fmt.Sprintf("cannot run job: project container %s (%s) is stopped or not found", project.Name, containerName)
		now := time.Now()
		_ = cm.store.UpdateJobStatusAndOutput(jobID, "error", &now, errMsg)
		return errMsg, errors.New(errMsg)
	}

	execConfig := dockertypes.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"sh", "-c", j.Command},
	}

	execIDResp, err := cm.dockerClient.ContainerExecCreate(ctx, inspectResp.ID, execConfig)
	if err != nil {
		now := time.Now()
		_ = cm.store.UpdateJobStatusAndOutput(jobID, "error", &now, err.Error())
		return "", fmt.Errorf("failed to create container exec for job %s: %w", j.Name, err)
	}

	attachResp, err := cm.dockerClient.ContainerExecAttach(ctx, execIDResp.ID, dockertypes.ExecStartCheck{})
	if err != nil {
		now := time.Now()
		_ = cm.store.UpdateJobStatusAndOutput(jobID, "error", &now, err.Error())
		return "", fmt.Errorf("failed to attach to container exec for job %s: %w", j.Name, err)
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
	_ = cm.store.UpdateJobStatusAndOutput(jobID, "active", &now, output)
	return output, nil
}
