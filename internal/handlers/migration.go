package handlers

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
)

type MigrationHandler struct {
	service *services.MigrationService
}

func NewMigrationHandler(s *services.MigrationService) *MigrationHandler {
	return &MigrationHandler{service: s}
}

func (h *MigrationHandler) Export(c echo.Context) error {
	var req struct {
		Passphrase string `json:"passphrase"`
	}
	if err := c.Bind(&req); err != nil || req.Passphrase == "" {
		return utils.Error(c, http.StatusBadRequest, "passphrase is required in request body")
	}

	bundleData, err := h.service.Export(c.Request().Context(), req.Passphrase)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, fmt.Sprintf("export failed: %v", err))
	}

	filename := fmt.Sprintf("codedock-bundle-%s.codedock", time.Now().UTC().Format("20060102-150405"))
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().WriteHeader(http.StatusOK)
	_, _ = c.Response().Write(bundleData)
	return nil
}

func (h *MigrationHandler) Import(c echo.Context) error {
	passphrase := c.FormValue("passphrase")
	if passphrase == "" {
		return utils.Error(c, http.StatusBadRequest, "passphrase form value is required")
	}

	file, err := c.FormFile("bundle")
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "bundle file is required (multipart field: bundle)")
	}

	src, err := file.Open()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to open uploaded bundle")
	}
	defer src.Close()

	bundleData, err := io.ReadAll(src)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to read bundle data")
	}

	manifest, err := h.service.Import(c.Request().Context(), bundleData, passphrase)
	if err != nil {
		return utils.Error(c, http.StatusUnprocessableEntity, fmt.Sprintf("import failed: %v", err))
	}

	return utils.Success(c, "Import completed successfully", manifest)
}
