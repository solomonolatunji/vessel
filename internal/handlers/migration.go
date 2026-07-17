package handlers

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type MigrationHandler struct {
	service *services.MigrationService
}

func NewMigrationHandler(s *services.MigrationService) *MigrationHandler {
	return &MigrationHandler{service: s}
}

// @Summary Export server bundle
// @Description Exports the full server state (SQLite database + all DB container dumps) into an AES-256-GCM encrypted .vessl bundle. Admin only.
// @Tags System
// @Accept json
// @Produce application/octet-stream
// @Param request body map[string]string true "JSON with passphrase"
// @Success 200 {file} binary "Encrypted .vessl bundle file"
// @Failure 400 {object} map[string]any "Missing passphrase"
// @Failure 500 {object} map[string]any "Export failed"
// @Router /system/export [post]
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

	filename := fmt.Sprintf("vessl-bundle-%s.vessl", time.Now().UTC().Format("20060102-150405"))
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().WriteHeader(http.StatusOK)
	_, _ = c.Response().Write(bundleData)
	return nil
}

// @Summary Import server bundle
// @Description Imports an encrypted .vessl bundle, restoring the SQLite database and all DB container data. Admin only. Triggers a server restart after restore.
// @Tags System
// @Accept multipart/form-data
// @Produce json
// @Param passphrase formData string true "Passphrase used to decrypt the bundle"
// @Param bundle formData file true "The .vessl bundle file to import"
// @Success 200 {object} map[string]any "Import completed with manifest details"
// @Failure 400 {object} map[string]any "Missing passphrase or bundle file"
// @Failure 422 {object} map[string]any "Decryption or restore failed"
// @Router /system/import [post]
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
