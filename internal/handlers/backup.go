package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type BackupHandler struct {
	backupService *services.BackupService
}

func NewBackupHandler(s *services.BackupService) *BackupHandler {
	return &BackupHandler{backupService: s}
}

// @Summary List endpoint
// @Description List endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Router /api/workspaces [get]
func (h *BackupHandler) List(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId query parameter"})
	}
	list, err := h.backupService.ListConfigsByProject(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

// @Summary Create endpoint
// @Description Create endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Router /api/workspaces [post]
func (h *BackupHandler) Create(c echo.Context) error {
	var cfg models.BackupConfig
	if err := c.Bind(&cfg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if err := h.backupService.CreateConfig(c.Request().Context(), &cfg); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, cfg)
}

// @Summary Get endpoint
// @Description Get endpoint
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "teamId"
// @Router /api/teams/{teamId}/ai_settings [get]
func (h *BackupHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || cfg == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "backup config not found"})
	}
	return c.JSON(http.StatusOK, cfg)
}

// @Summary Delete endpoint
// @Description Delete endpoint
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/workspaces/{id} [delete]
func (h *BackupHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	projectID := c.QueryParam("projectId")
	if id == "" || projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id or projectId"})
	}
	if err := h.backupService.DeleteConfig(c.Request().Context(), id, projectID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary Trigger endpoint
// @Description Trigger endpoint
// @Tags Backups
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/backups/{id}/trigger [post]
func (h *BackupHandler) Trigger(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	rec, err := h.backupService.TriggerBackup(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rec)
}

// @Summary ListRecords endpoint
// @Description ListRecords endpoint
// @Tags Backups
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/backups/{id}/records [get]
func (h *BackupHandler) ListRecords(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id parameter"})
	}
	recs, err := h.backupService.ListRecordsByConfig(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, recs)
}

// @Summary ListS3Destinations endpoint
// @Description ListS3Destinations endpoint
// @Tags S3-destinations
// @Accept json
// @Produce json
// @Router /api/s3-destinations [get]
func (h *BackupHandler) ListS3Destinations(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing projectId query parameter"})
	}
	list, err := h.backupService.ListS3Destinations(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

// @Summary CreateS3Destination endpoint
// @Description CreateS3Destination endpoint
// @Tags S3-destinations
// @Accept json
// @Produce json
// @Router /api/s3-destinations [post]
func (h *BackupHandler) CreateS3Destination(c echo.Context) error {
	var dest models.S3Destination
	if err := c.Bind(&dest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if err := h.backupService.CreateS3Destination(c.Request().Context(), &dest); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, dest)
}

// @Summary DeleteS3Destination endpoint
// @Description DeleteS3Destination endpoint
// @Tags S3-destinations
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/s3-destinations/{id} [delete]
func (h *BackupHandler) DeleteS3Destination(c echo.Context) error {
	id := c.Param("id")
	projectID := c.QueryParam("projectId")
	if id == "" || projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id or projectId"})
	}
	if err := h.backupService.DeleteS3Destination(c.Request().Context(), id, projectID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
