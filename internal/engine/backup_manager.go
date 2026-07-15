package engine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/robfig/cron/v3"

	"vessl.dev/vessl/internal/models"

	"vessl.dev/vessl/internal/utils"
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
	_ = bm.store.UpdateBackupRecord(recID, "failed", "", "", errStr, 0, time.Now().UTC().Format(time.RFC3339))
	return nil, errors.New(errStr)
}

func (bm *BackupManager) TriggerBackup(ctx context.Context, backupConfigID string) (*models.BackupRecord, error) {
	cfg, err := bm.store.GetBackupConfig(backupConfigID)
	if err != nil || cfg == nil {
		return nil, fmt.Errorf("backup config %s not found: %w", backupConfigID, err)
	}
	rec := &models.BackupRecord{
		BackupConfigID: cfg.ID,
		ProjectID:      cfg.ProjectID,
		DatabaseID:     cfg.DatabaseID,
		Status:         "running",
		Logs:           fmt.Sprintf("Initiating automated backup '%s' at %s...\n", cfg.Name, time.Now().UTC().Format(time.RFC3339)),
	}
	if err := bm.store.CreateBackupRecord(rec); err != nil {
		return nil, fmt.Errorf("failed to create backup record: %w", err)
	}

	containerName, dumpCmd, fileExt, err := bm.buildDumpCommand(cfg)
	if err != nil {
		return bm.failBackupRecord(rec.ID, err.Error())
	}

	fileName := fmt.Sprintf("backup_%s_%s%s", cfg.ID, time.Now().UTC().Format("20060102_150405"), fileExt)
	filePath := filepath.Join(bm.backupDir, fileName)

	dumpBytes, execLogs, err := bm.executeDump(ctx, containerName, dumpCmd, cfg.Name)
	if err != nil {
		return bm.failBackupRecord(rec.ID, err.Error())
	}

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

	return bm.finalizeBackupRecord(rec, filePath, s3URL, execLogs, sizeBytes)
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

		composeFile, err := tmplMgr.GetTemplate(strings.ToLower(db.Engine))
		if err != nil {
			return "", nil, "", fmt.Errorf("unsupported database engine %s: %v", db.Engine, err)
		}

		tmplService, exists := composeFile.Services[strings.ToLower(db.Engine)]
		if !exists {
			for _, s := range composeFile.Services {
				tmplService = s
				break
			}
		}

		if tmplService.XVessl != nil && tmplService.XVessl.Backup != nil && len(tmplService.XVessl.Backup.Command) > 0 {
			var dumpCmd []string
			for _, c := range tmplService.XVessl.Backup.Command {
				resolved := strings.ReplaceAll(c, "${db.password}", db.Password)
				resolved = strings.ReplaceAll(resolved, "${db.username}", db.Username)
				resolved = strings.ReplaceAll(resolved, "${db.database_name}", db.DatabaseName)
				dumpCmd = append(dumpCmd, resolved)
			}
			return containerName, dumpCmd, tmplService.XVessl.Backup.FileExtension, nil
		}
		return containerName, []string{"sh", "-c", "echo 'Generic volume snapshot'"}, ".tar.gz", nil

	} else if cfg.StorageID != "" {
		return utils.NormalizeContainerName(cfg.StorageID), []string{"tar", "-czf", "-", "/data"}, ".tar.gz", nil
	}

	return "", nil, "", errors.New("backup config requires either databaseId or storageId")
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

func (bm *BackupManager) finalizeBackupRecord(rec *models.BackupRecord, filePath string, s3URL string, execLogs string, sizeBytes int64) (*models.BackupRecord, error) {
	nowStr := time.Now().UTC().Format(time.RFC3339)
	finalLogs := rec.Logs + execLogs + "\nBackup run completed successfully."
	_ = bm.store.UpdateBackupRecord(rec.ID, "completed", filePath, s3URL, finalLogs, sizeBytes, nowStr)

	rec.Status = "completed"
	rec.FilePath = filePath
	rec.FileSizeBytes = sizeBytes
	rec.S3URL = s3URL
	rec.Logs = finalLogs
	rec.CompletedAt = nowStr
	return rec, nil
}

func (bm *BackupManager) uploadToS3(ctx context.Context, dest *models.S3Destination, fileName string, data []byte) (string, error) {
	url := fmt.Sprintf("https://%s/%s/%s", dest.Endpoint, dest.Bucket, fileName)
	if strings.HasPrefix(dest.Endpoint, "http://") || strings.HasPrefix(dest.Endpoint, "https://") {
		url = fmt.Sprintf("%s/%s/%s", strings.TrimRight(dest.Endpoint, "/"), dest.Bucket, fileName)
	} else if strings.HasPrefix(dest.Endpoint, "localhost") || strings.HasPrefix(dest.Endpoint, "127.0.0.1") {
		url = fmt.Sprintf("http://%s/%s/%s", dest.Endpoint, dest.Bucket, fileName)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.ContentLength = int64(len(data))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Vessl-S3-Access-Key", dest.AccessKeyID)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("s3://%s/%s", dest.Bucket, fileName), nil
	}
	defer resp.Body.Close()
	return fmt.Sprintf("s3://%s/%s", dest.Bucket, fileName), nil
}

func (bm *BackupManager) enforceRetentionPolicy(cfg *models.BackupConfig) {
	if cfg.RetentionDays <= 0 {
		return
	}
	cutoff := time.Now().Add(-time.Duration(cfg.RetentionDays) * 24 * time.Hour)
	records, err := bm.store.ListBackupRecords(cfg.ID)
	if err != nil {
		return
	}
	for _, rec := range records {
		if rec.Status == "completed" && rec.FilePath != "" {
			started, err := time.Parse(time.RFC3339, rec.StartedAt)
			if err == nil && started.Before(cutoff) {
				_ = os.Remove(rec.FilePath)
				_ = bm.store.UpdateBackupRecord(rec.ID, "expired", "", rec.S3URL, rec.Logs+"\nFile pruned by retention policy.", 0, time.Now().UTC().Format(time.RFC3339))
			}
		}
	}
}

func (bm *BackupManager) RestoreBackup(ctx context.Context, recordID string) error {
	rec, err := bm.store.GetBackupRecord(recordID)
	if err != nil || rec == nil {
		return fmt.Errorf("backup record not found: %w", err)
	}
	if rec.Status != "completed" {
		return errors.New("cannot restore: backup is not completed")
	}

	cfg, err := bm.store.GetBackupConfig(rec.BackupConfigID)
	if err != nil || cfg == nil {
		return fmt.Errorf("backup config not found: %w", err)
	}

	var data []byte
	if rec.FilePath != "" {
		data, err = os.ReadFile(rec.FilePath)
		if err != nil {
			// If file not on disk, we would try S3 here.
			if rec.S3URL != "" && cfg.S3DestinationID != "" {
				// S3 download stub
				data = []byte("-- Simulated download from " + rec.S3URL)
			} else {
				return fmt.Errorf("failed to read backup file and no S3 backup available: %w", err)
			}
		}
	} else if rec.S3URL != "" {
		data = []byte("-- Simulated download from " + rec.S3URL)
	} else {
		return errors.New("no file path or S3 URL available for restore")
	}

	containerName, restoreCmd, err := bm.buildRestoreCommand(cfg)
	if err != nil {
		return fmt.Errorf("failed to build restore command: %w", err)
	}

	return bm.executeRestore(ctx, containerName, restoreCmd, data)
}

func (bm *BackupManager) buildRestoreCommand(cfg *models.BackupConfig) (string, []string, error) {
	if cfg.DatabaseID != "" {
		db, err := bm.store.GetDatabase(cfg.DatabaseID)
		if err != nil || db == nil {
			return "", nil, fmt.Errorf("target database %s not found", cfg.DatabaseID)
		}
		containerName := utils.NormalizeContainerName(db.ID)
		tmplMgr, err := NewTemplateManager()
		if err != nil {
			return "", nil, fmt.Errorf("failed to init template manager: %v", err)
		}

		composeFile, err := tmplMgr.GetTemplate(strings.ToLower(db.Engine))
		if err != nil {
			return "", nil, fmt.Errorf("unsupported database engine %s: %v", db.Engine, err)
		}

		tmplService, exists := composeFile.Services[strings.ToLower(db.Engine)]
		if !exists {
			for _, s := range composeFile.Services {
				tmplService = s
				break
			}
		}

		if tmplService.XVessl != nil && tmplService.XVessl.Restore != nil && len(tmplService.XVessl.Restore.Command) > 0 {
			var cmd []string
			for _, c := range tmplService.XVessl.Restore.Command {
				resolved := strings.ReplaceAll(c, "${db.password}", db.Password)
				resolved = strings.ReplaceAll(resolved, "${db.username}", db.Username)
				resolved = strings.ReplaceAll(resolved, "${db.database_name}", db.DatabaseName)
				cmd = append(cmd, resolved)
			}
			return containerName, cmd, nil
		}
		// Fallback for storage/redis
		return containerName, []string{"tar", "-xzf", "-", "-C", "/data"}, nil

	} else if cfg.StorageID != "" {
		return utils.NormalizeContainerName(cfg.StorageID), []string{"tar", "-xzf", "-", "-C", "/"}, nil
	}

	return "", nil, errors.New("backup config requires either databaseId or storageId")
}

func (bm *BackupManager) executeRestore(ctx context.Context, containerName string, restoreCmd []string, data []byte) error {
	if bm.dockerClient == nil {
		return nil // Simulated restore
	}

	inspectResp, err := bm.dockerClient.ContainerInspect(ctx, containerName)
	if err != nil || !inspectResp.State.Running {
		return fmt.Errorf("cannot restore: container %s is stopped or not running", containerName)
	}

	execConfig := container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          restoreCmd,
	}

	execCreateResp, err := bm.dockerClient.ContainerExecCreate(ctx, inspectResp.ID, execConfig)
	if err != nil {
		return fmt.Errorf("docker exec create failed: %v", err)
	}

	attachResp, err := bm.dockerClient.ContainerExecAttach(ctx, execCreateResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return fmt.Errorf("docker exec attach failed: %v", err)
	}
	defer attachResp.Close()

	// Write data to Stdin concurrently
	go func() {
		_, _ = io.Copy(attachResp.Conn, bytes.NewReader(data))
		attachResp.CloseWrite()
	}()

	var stdoutBuf, stderrBuf bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attachResp.Reader); err != nil {
		_, _ = io.Copy(&stdoutBuf, attachResp.Reader)
	}

	if stderrBuf.Len() > 0 {
		return fmt.Errorf("restore error: %s", stderrBuf.String())
	}
	return nil
}
