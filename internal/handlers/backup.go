package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type BackupHandler struct {
	backupService *services.BackupService
}

func NewBackupHandler(s *services.BackupService) *BackupHandler {
	return &BackupHandler{backupService: s}
}

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
