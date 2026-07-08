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
	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
	"github.com/solomonolatunji/vessel/internal/utils"
)

// CronManager orchestrates background job schedules and executes terminal commands inside active project containers.
type CronManager struct {
	dockerClient *client.Client
	store        *store.Store
	cronEngine   *cron.Cron
	entries      map[string]cron.EntryID
	mu           sync.Mutex
}

// NewCronManager initializes a CronManager with Docker client and store dependencies.
func NewCronManager(dockerClient *client.Client, s *store.Store) *CronManager {
	return &CronManager{
		dockerClient: dockerClient,
		store:        s,
		cronEngine:   cron.New(cron.WithSeconds()),
		entries:      make(map[string]cron.EntryID),
	}
}

// Start launches the background cron loop and loads all active scheduled jobs from the store.
func (cm *CronManager) Start() error {
	jobs, err := cm.store.ListJobs()
	if err != nil {
		return fmt.Errorf("failed to load jobs during cron manager start: %w", err)
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, job := range jobs {
		if job.Status == "active" {
			if err := cm.registerJobLocked(&job); err != nil {
				log.Printf("⚠️ Failed to register cron job %s (%s): %v", job.Name, job.ID, err)
			}
		}
	}

	cm.cronEngine.Start()
	log.Println("⏰ CronManager started and executing background tasks")
	return nil
}

// Stop halts all active cron timers.
func (cm *CronManager) Stop() {
	if cm.cronEngine != nil {
		cm.cronEngine.Stop()
	}
}

// RegisterJob adds or updates a scheduled job in the active cron runner.
func (cm *CronManager) RegisterJob(job *types.JobConfig) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.registerJobLocked(job)
}

func (cm *CronManager) registerJobLocked(job *types.JobConfig) error {
	if entryID, exists := cm.entries[job.ID]; exists {
		cm.cronEngine.Remove(entryID)
		delete(cm.entries, job.ID)
	}

	if job.Status != "active" {
		return nil
	}

	schedule := strings.TrimSpace(job.Schedule)
	if len(strings.Fields(schedule)) == 5 && !strings.HasPrefix(schedule, "@") {
		schedule = "0 " + schedule
	}

	jobID := job.ID
	entryID, err := cm.cronEngine.AddFunc(schedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		_, _ = cm.ExecuteJob(ctx, jobID)
	})
	if err != nil {
		return fmt.Errorf("invalid cron schedule '%s': %w", job.Schedule, err)
	}

	cm.entries[job.ID] = entryID
	return nil
}

// UnregisterJob removes a scheduled job from the cron execution schedule.
func (cm *CronManager) UnregisterJob(jobID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if entryID, exists := cm.entries[jobID]; exists {
		cm.cronEngine.Remove(entryID)
		delete(cm.entries, jobID)
	}
}

// ExecuteJob immediately runs the job's command inside the target project container and persists output logs.
func (cm *CronManager) ExecuteJob(ctx context.Context, jobID string) (string, error) {
	job, err := cm.store.GetJob(jobID)
	if err != nil || job == nil {
		return "", fmt.Errorf("job %s not found: %w", jobID, err)
	}

	project, err := cm.store.GetProject(job.ProjectID)
	if err != nil || project == nil {
		return "", fmt.Errorf("project %s for job %s not found: %w", job.ProjectID, jobID, err)
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
		Cmd:          []string{"sh", "-c", job.Command},
	}

	execIDResp, err := cm.dockerClient.ContainerExecCreate(ctx, inspectResp.ID, execConfig)
	if err != nil {
		now := time.Now()
		_ = cm.store.UpdateJobStatusAndOutput(jobID, "error", &now, err.Error())
		return "", fmt.Errorf("failed to create container exec for job %s: %w", job.Name, err)
	}

	attachResp, err := cm.dockerClient.ContainerExecAttach(ctx, execIDResp.ID, dockertypes.ExecStartCheck{})
	if err != nil {
		now := time.Now()
		_ = cm.store.UpdateJobStatusAndOutput(jobID, "error", &now, err.Error())
		return "", fmt.Errorf("failed to attach to container exec for job %s: %w", job.Name, err)
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
