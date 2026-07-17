package engine

import (
	"context"
	"fmt"
	"os"
	"time"

	"vessl.dev/vessl/internal/models"
)

func (bm *BackupManager) uploadToS3(ctx context.Context, dest *models.S3Destination, fileName string, data []byte) (string, error) {
	resp, err := signedS3Request(ctx, dest, "PUT", fileName, data, "application/octet-stream")
	if err != nil {
		return "", err
	}
	resp.Body.Close()
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
		if rec.Status == models.BackupRecordStatusCompleted && rec.FilePath != "" {
			started, err := time.Parse(time.RFC3339, rec.StartedAt)
			if err == nil && started.Before(cutoff) {
				_ = os.Remove(rec.FilePath)
				_ = bm.store.UpdateBackupRecord(models.UpdateBackupRecordOpts{
					ID:          rec.ID,
					Status:      models.BackupRecordStatusExpired,
					S3URL:       rec.S3URL,
					Logs:        rec.Logs + "\nFile pruned by retention policy.",
					CompletedAt: time.Now().UTC().Format(time.RFC3339),
				})
			}
		}
	}
}
