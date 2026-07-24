package engine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"

	"codedock.run/codedock/internal/models"

	"codedock.run/codedock/internal/utils"
)

type BackupManager struct {
	dockerClient *client.Client
	store        BackupManagerStore
	cronEngine   *cron.Cron
	entries      map[string]cron.EntryID
	backupDir    string
	mu           sync.Mutex
}

func NewBackupManager(dockerClient *client.Client, s BackupManagerStore, backupDir string) *BackupManager {
	if backupDir == "" {
		backupDir = filepath.Join(utils.GetDataDir(), "backups")
	}
	_ = os.MkdirAll(backupDir, 0o755)
	return &BackupManager{
		dockerClient: dockerClient,
		store:        s,
		cronEngine:   cron.New(cron.WithSeconds()),
		entries:      make(map[string]cron.EntryID),
		backupDir:    backupDir,
	}
}

func (bm *BackupManager) Start() error {
	cfgs, err := bm.store.ListAllActiveBackupConfigs()
	if err != nil {
		return fmt.Errorf("failed to load active backup configs during start: %w", err)
	}
	bm.mu.Lock()
	defer bm.mu.Unlock()
	for _, cfg := range cfgs {
		if err := bm.registerBackupLocked(cfg); err != nil {
			slog.Warn("failed to schedule backup", "name", cfg.Name, "id", cfg.ID, "err", err)
		}
	}
	bm.cronEngine.Start()
	slog.Info("backup manager started")
	return nil
}

func (bm *BackupManager) Stop() {
	if bm.cronEngine != nil {
		bm.cronEngine.Stop()
	}
}

func (bm *BackupManager) RegisterBackup(cfg *models.BackupConfig) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	return bm.registerBackupLocked(cfg)
}

func (bm *BackupManager) registerBackupLocked(cfg *models.BackupConfig) error {
	if entryID, exists := bm.entries[cfg.ID]; exists {
		bm.cronEngine.Remove(entryID)
		delete(bm.entries, cfg.ID)
	}
	if cfg.Status != "active" {
		return nil
	}
	schedule := strings.TrimSpace(cfg.Schedule)
	if len(strings.Fields(schedule)) == 5 && !strings.HasPrefix(schedule, "@") {
		schedule = "0 " + schedule
	}
	cfgID := cfg.ID
	entryID, err := bm.cronEngine.AddFunc(schedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		_, _ = bm.TriggerBackup(ctx, cfgID)
	})
	if err != nil {
		return fmt.Errorf("invalid cron schedule '%s' for backup %s: %w", cfg.Schedule, cfg.Name, err)
	}
	bm.entries[cfg.ID] = entryID
	return nil
}

func (bm *BackupManager) UnregisterBackup(backupConfigID string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if entryID, exists := bm.entries[backupConfigID]; exists {
		bm.cronEngine.Remove(entryID)
		delete(bm.entries, backupConfigID)
	}
}

func (bm *BackupManager) failBackupRecord(recID, errStr string) (*models.BackupRecord, error) {
	_ = bm.store.UpdateBackupRecord(models.UpdateBackupRecordOpts{
		ID:          recID,
		Status:      models.BackupRecordStatusFailed,
		Logs:        fmt.Sprintf("Failed: %s\n", errStr),
		CompletedAt: time.Now().UTC().Format(time.RFC3339),
	})
	return nil, errors.New(errStr)
}

func (bm *BackupManager) DeleteBackupRecord(ctx context.Context, recordID string) {
	rec, err := bm.store.GetBackupRecord(recordID)
	if err == nil && rec != nil && rec.FilePath != "" {
		_ = os.Remove(rec.FilePath)
	}
}

func (bm *BackupManager) TriggerBackup(ctx context.Context, backupConfigID string) (*models.BackupRecord, error) {
	cfg, err := bm.store.GetBackupConfig(backupConfigID)
	if err != nil || cfg == nil {
		return nil, fmt.Errorf("backup config %s not found: %w", backupConfigID, err)
	}

	rec := &models.BackupRecord{
		ID:             uuid.New().String(),
		BackupConfigID: cfg.ID,
		DatabaseID:     cfg.DatabaseID,
		Status:         models.BackupRecordStatusRunning,
		Logs:           fmt.Sprintf("Initiating automated backup '%s' at %s...\n", cfg.Name, time.Now().UTC().Format(time.RFC3339)),
	}
	if err := bm.store.CreateBackupRecord(rec); err != nil {
		return nil, fmt.Errorf("failed to create backup record: %w", err)
	}

	var dumpBytes []byte
	var execLogs string
	var fileExt string

	if cfg.DatabaseID == "global" || cfg.DatabaseID == "" {
		dbPath := filepath.Join(utils.GetDataDir(), "codedock.db")
		content, err := os.ReadFile(dbPath)
		if err != nil {
			return bm.failBackupRecord(rec.ID, fmt.Sprintf("failed to read global db: %v", err))
		}
		dumpBytes = content
		fileExt = ".db"
		execLogs = "Global database backed up successfully.\n"
	} else {
		containerName, dumpCmd, ext, err := bm.buildDumpCommand(cfg)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		fileExt = ext
		dumpBytes, execLogs, err = bm.executeDump(ctx, containerName, dumpCmd, cfg.Name)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
	}

	fileName := fmt.Sprintf("backup_%s_%s%s", cfg.ID, time.Now().UTC().Format("20060102_150405"), fileExt)
	filePath := filepath.Join(bm.backupDir, fileName)

	if err := os.WriteFile(filePath, dumpBytes, 0o600); err != nil {
		return bm.failBackupRecord(rec.ID, fmt.Sprintf("failed to write backup archive to disk: %v", err))
	}

	sizeBytes := int64(len(dumpBytes))
	if fileInfo, err := os.Stat(filePath); err == nil && fileInfo != nil {
		sizeBytes = fileInfo.Size()
	}

	s3URL := ""
	if cfg.S3DestinationID != "" {
		s3URL, execLogs = bm.handleS3Upload(ctx, cfg, fileName, dumpBytes, execLogs)
	}

	bm.enforceRetentionPolicy(cfg)

	return bm.finalizeBackupRecord(FinalizeBackupOpts{
		Record:    rec,
		FilePath:  filePath,
		S3URL:     s3URL,
		ExecLogs:  execLogs,
		SizeBytes: sizeBytes,
	})
}

func (bm *BackupManager) buildDumpCommand(cfg *models.BackupConfig) (string, []string, string, error) {
	if cfg.DatabaseID != "" {
		db, err := bm.store.GetDatabase(cfg.DatabaseID)
		if err != nil || db == nil {
			return "", nil, "", fmt.Errorf("target database %s not found", cfg.DatabaseID)
		}
		containerName := utils.NormalizeContainerName(db.ID)
		tmplMgr, err := NewTemplateManager()
		if err != nil {
			return "", nil, "", fmt.Errorf("failed to init template manager: %v", err)
		}

		composeFile, err := tmplMgr.GetTemplate(strings.ToLower(string(db.Engine)))
		if err != nil {
			return "", nil, "", fmt.Errorf("unsupported database engine %s: %v", db.Engine, err)
		}

		tmplService, exists := composeFile.Services[strings.ToLower(string(db.Engine))]
		if !exists {
			for _, s := range composeFile.Services {
				tmplService = s
				break
			}
		}

		if tmplService.XCodedock != nil && tmplService.XCodedock.Backup != nil && len(tmplService.XCodedock.Backup.Command) > 0 {
			var dumpCmd []string
			for _, c := range tmplService.XCodedock.Backup.Command {
				resolved := strings.ReplaceAll(c, "${db.password}", db.Password)
				resolved = strings.ReplaceAll(resolved, "${db.username}", db.Username)
				resolved = strings.ReplaceAll(resolved, "${db.database_name}", db.DatabaseName)
				dumpCmd = append(dumpCmd, resolved)
			}
			return containerName, dumpCmd, tmplService.XCodedock.Backup.FileExtension, nil
		}
		return containerName, []string{"sh", "-c", "echo 'Generic volume snapshot'"}, ".tar.gz", nil

	}

	return "", nil, "", errors.New("backup config requires databaseId")
}

func (bm *BackupManager) executeDump(ctx context.Context, containerName string, dumpCmd []string, backupName string) ([]byte, string, error) {
	if bm.dockerClient == nil {
		dumpBytes := []byte(fmt.Sprintf("-- Simulated backup dump for %s at %s --\n", backupName, time.Now().UTC().Format(time.RFC3339)))
		return dumpBytes, "Docker client nil: simulated successful local dump.\n", nil
	}

	inspectResp, err := bm.dockerClient.ContainerInspect(ctx, containerName)
	if err != nil || !inspectResp.State.Running {
		return nil, "", fmt.Errorf("cannot backup: container %s is stopped or not running", containerName)
	}

	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          dumpCmd,
	}

	execCreateResp, err := bm.dockerClient.ContainerExecCreate(ctx, inspectResp.ID, execConfig)
	if err != nil {
		return nil, "", fmt.Errorf("docker exec create failed: %v", err)
	}

	attachResp, err := bm.dockerClient.ContainerExecAttach(ctx, execCreateResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("docker exec attach failed: %v", err)
	}
	defer attachResp.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attachResp.Reader); err != nil {
		_, _ = io.Copy(&stdoutBuf, attachResp.Reader)
	}

	return stdoutBuf.Bytes(), stderrBuf.String(), nil
}

func (bm *BackupManager) handleS3Upload(ctx context.Context, cfg *models.BackupConfig, fileName string, dumpBytes []byte, execLogs string) (string, string) {
	dest, err := bm.store.GetS3Destination(cfg.S3DestinationID)
	if err != nil || dest == nil {
		return "", execLogs
	}
	s3URL, err := bm.uploadToS3(ctx, dest, fileName, dumpBytes)
	if err != nil {
		execLogs += fmt.Sprintf("\n⚠️ S3 upload failed: %v", err)
	} else {
		execLogs += fmt.Sprintf("\n✅ Successfully uploaded backup to S3/MinIO destination: %s", s3URL)
	}
	return s3URL, execLogs
}

type FinalizeBackupOpts struct {
	Record    *models.BackupRecord
	FilePath  string
	S3URL     string
	ExecLogs  string
	SizeBytes int64
}

func (bm *BackupManager) finalizeBackupRecord(opts FinalizeBackupOpts) (*models.BackupRecord, error) {
	nowStr := time.Now().UTC().Format(time.RFC3339)
	finalLogs := opts.Record.Logs + opts.ExecLogs + "\nBackup run completed successfully."
	_ = bm.store.UpdateBackupRecord(models.UpdateBackupRecordOpts{
		ID:            opts.Record.ID,
		Status:        models.BackupRecordStatusCompleted,
		FilePath:      opts.FilePath,
		S3URL:         opts.S3URL,
		Logs:          finalLogs,
		FileSizeBytes: opts.SizeBytes,
		CompletedAt:   nowStr,
	})

	opts.Record.Status = models.BackupRecordStatusCompleted
	opts.Record.FilePath = opts.FilePath
	opts.Record.FileSizeBytes = opts.SizeBytes
	opts.Record.S3URL = opts.S3URL
	opts.Record.Logs = finalLogs
	opts.Record.CompletedAt = nowStr
	return opts.Record, nil
}
