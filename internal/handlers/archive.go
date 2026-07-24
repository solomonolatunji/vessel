package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
)

type ArchiveHandler struct {
	service *services.ArchiveService
}

func NewArchiveHandler(s *services.ArchiveService) *ArchiveHandler {
	return &ArchiveHandler{service: s}
}

func (h *ArchiveHandler) DeployArchive(c echo.Context) error {
	projectID := c.FormValue("projectId")
	appName := c.FormValue("name")

	file, err := c.FormFile("file")
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "archive file is required")
	}

	if appName == "" {
		base := strings.TrimSuffix(file.Filename, ".tar.gz")
		base = strings.TrimSuffix(base, ".tar")
		appName = base
	}

	src, err := file.Open()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to read uploaded file")
	}
	defer src.Close()

	tmpPath := filepath.Join(os.TempDir(), "codedock-upload", uuid.New().String()+".tar.gz")
	if err := writeFile(tmpPath, src); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to save archive")
	}
	defer os.Remove(tmpPath)

	result, err := h.service.Deploy(c.Request().Context(), projectID, appName, tmpPath)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "Archive deployed", result)
}

func writeFile(path string, r io.Reader) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}
