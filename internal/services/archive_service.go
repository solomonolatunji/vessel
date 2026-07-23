package services

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock.dev/codedock/internal/models"
)

type ArchiveService struct {
	appService        *AppService
	deploymentService *DeploymentService
}

func NewArchiveService(as *AppService, ds *DeploymentService) *ArchiveService {
	return &ArchiveService{
		appService:        as,
		deploymentService: ds,
	}
}

type ArchiveDeployResult struct {
	ContainerID string `json:"containerId"`
	AppID       string `json:"appId"`
	AppName     string `json:"appName"`
}

func (s *ArchiveService) Deploy(ctx context.Context, projectID, appName, archivePath string) (*ArchiveDeployResult, error) {
	app, err := s.resolveOrCreateApp(ctx, projectID, appName)
	if err != nil {
		return nil, err
	}

	srcDir, err := extractToTemp(archivePath)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(filepath.Dir(srcDir))

	containerID, err := s.deploymentService.DeployAppService(ctx, app.ID, srcDir, nil)
	if err != nil {
		return nil, fmt.Errorf("deployment failed: %w", err)
	}

	return &ArchiveDeployResult{
		ContainerID: containerID,
		AppID:       app.ID,
		AppName:     app.Name,
	}, nil
}

func (s *ArchiveService) resolveOrCreateApp(ctx context.Context, projectID, appName string) (*models.AppService, error) {
	apps, err := s.appService.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list apps: %w", err)
	}

	for _, a := range apps {
		if a.Name == appName {
			return a, nil
		}
	}

	app := &models.AppService{
		ID:           uuid.New().String(),
		ProjectID:    projectID,
		Name:         appName,
		InternalPort: 3000,
		Status:       "created",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.appService.CreateAppService(ctx, app)
}

func extractToTemp(archivePath string) (string, error) {
	tmpDir := filepath.Join(os.TempDir(), "codedock-archive", uuid.New().String())
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	src, err := os.Open(archivePath)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("failed to open archive: %w", err)
	}
	defer src.Close()

	if err := untarStream(tmpDir, src); err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("failed to extract: %w", err)
	}

	return findSourceDir(tmpDir), nil
}

func untarStream(dest string, src io.Reader) error {
	gzr, err := gzip.NewReader(src)
	if err != nil {
		return untarReader(dest, tar.NewReader(src))
	}
	defer gzr.Close()
	return untarReader(dest, tar.NewReader(gzr))
}

func untarReader(dest string, tr *tar.Reader) error {
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dest, filepath.Clean(hdr.Name))
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			continue
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(path, 0o755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(path), 0o755)
			f, _ := os.Create(path)
			if f != nil {
				_, _ = io.Copy(f, tr)
				f.Close()
			}
		}
	}
	return nil
}

func findSourceDir(dir string) string {
	entries, _ := os.ReadDir(dir)
	if len(entries) == 1 && entries[0].IsDir() {
		return filepath.Join(dir, entries[0].Name())
	}
	return dir
}
