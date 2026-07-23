package engine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/utils"
)

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
	var fetchErr error
	if rec.FilePath != "" {
		data, fetchErr = os.ReadFile(rec.FilePath)
	}

	if (rec.FilePath == "" || fetchErr != nil) && rec.S3URL != "" && cfg.S3DestinationID != "" {
		dest, err := bm.store.GetS3Destination(cfg.S3DestinationID)
		if err != nil || dest == nil {
			return fmt.Errorf("failed to retrieve S3 destination for restore: %w", err)
		}

		prefix := fmt.Sprintf("s3://%s/", dest.Bucket)
		key := strings.TrimPrefix(rec.S3URL, prefix)

		resp, err := signedS3Request(ctx, dest, "GET", key, nil, "")
		if err != nil {
			return fmt.Errorf("failed to download backup from S3: %w", err)
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read downloaded S3 backup: %w", err)
		}
	} else if len(data) == 0 {
		if fetchErr != nil {
			return fmt.Errorf("failed to read local backup file: %w", fetchErr)
		}
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

		composeFile, err := tmplMgr.GetTemplate(strings.ToLower(string(db.Engine)))
		if err != nil {
			return "", nil, fmt.Errorf("unsupported database engine %s: %v", db.Engine, err)
		}

		tmplService, exists := composeFile.Services[strings.ToLower(string(db.Engine))]
		if !exists {
			for _, s := range composeFile.Services {
				tmplService = s
				break
			}
		}

		if tmplService.XCodedock != nil && tmplService.XCodedock.Restore != nil && len(tmplService.XCodedock.Restore.Command) > 0 {
			var cmd []string
			for _, c := range tmplService.XCodedock.Restore.Command {
				resolved := strings.ReplaceAll(c, "${db.password}", db.Password)
				resolved = strings.ReplaceAll(resolved, "${db.username}", db.Username)
				resolved = strings.ReplaceAll(resolved, "${db.database_name}", db.DatabaseName)
				cmd = append(cmd, resolved)
			}
			return containerName, cmd, nil
		}
		return containerName, []string{"tar", "-xzf", "-", "-C", "/data"}, nil

	}

	return "", nil, errors.New("backup config requires databaseId")
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
