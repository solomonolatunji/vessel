package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type BackupHandler struct {
	backupService *services.BackupService
}

func NewBackupHandler(s *services.BackupService) *BackupHandler {
	return &BackupHandler{backupService: s}
}

// @Summary List Backups
// @Description List Backups
// @Tags Backups
// @Accept json
// @Produce json
// @Router /backups [get]
func (h *BackupHandler) List(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId query parameter")
	}
	list, err := h.backupService.ListConfigsByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", list)
}

// @Summary Create Backup
// @Description Create Backup
// @Tags Backups
// @Accept json
// @Produce json
// @Param request body models.BackupConfig true "Payload"
// @Router /backups [post]
func (h *BackupHandler) Create(c echo.Context) error {
	var cfg models.BackupConfig
	if err := c.Bind(&cfg); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if err := h.backupService.CreateConfig(c.Request().Context(), &cfg); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", cfg)
}

// @Summary Get Backup
// @Description Get Backup
// @Tags Backups
// @Accept json
// @Produce json
// @Param id path string true "Backup ID"
// @Router /backups/{id} [get]
func (h *BackupHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || cfg == nil {
		return utils.Error(c, http.StatusNotFound, "backup config not found")
	}
	return utils.Success(c, "Operation successful", cfg)
}

// @Summary Delete Backup
// @Description Delete Backup
// @Tags Backups
// @Accept json
// @Produce json
// @Param id path string true "Backup ID"
// @Router /backups/{id} [delete]
func (h *BackupHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	projectID := c.QueryParam("projectId")
	if id == "" || projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id or projectId")
	}
	if err := h.backupService.DeleteConfig(c.Request().Context(), id, projectID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary Trigger endpoint
// @Description Trigger endpoint
// @Tags Backups
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /backups/{id}/trigger [post]
func (h *BackupHandler) Trigger(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	rec, err := h.backupService.TriggerBackup(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", rec)
}

// @Summary ListRecords endpoint
// @Description ListRecords endpoint
// @Tags Backups
// @Accept json
// @Produce json
// @Param id path string true "id"
func (h *BackupHandler) ListRecords(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	recs, err := h.backupService.ListRecordsByConfig(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", recs)
}

// @Summary Restore endpoint
// @Description Restore endpoint
// @Tags Backups
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /backups/{id}/restore [post]
func (h *BackupHandler) Restore(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id parameter")
	}
	err := h.backupService.RestoreBackup(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Backup successfully restored", nil)
}

// @Summary ListS3Destinations endpoint
// @Description ListS3Destinations endpoint
// @Tags Backups
// @Accept json
// @Produce json
// @Router /s3-destinations [get]
func (h *BackupHandler) ListS3Destinations(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId query parameter")
	}
	list, err := h.backupService.ListS3Destinations(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", list)
}

// @Summary CreateS3Destination endpoint
// @Description CreateS3Destination endpoint
// @Tags Backups
// @Accept json
// @Produce json
// @Param request body models.S3Destination true "Payload"
func (h *BackupHandler) CreateS3Destination(c echo.Context) error {
	var dest models.S3Destination
	if err := c.Bind(&dest); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if err := h.backupService.CreateS3Destination(c.Request().Context(), &dest); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", dest)
}

// @Summary DeleteS3Destination endpoint
// @Description DeleteS3Destination endpoint
// @Tags Backups
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /s3-destinations/{id} [delete]
func (h *BackupHandler) DeleteS3Destination(c echo.Context) error {
	id := c.Param("id")
	projectID := c.QueryParam("projectId")
	if id == "" || projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id or projectId")
	}
	if err := h.backupService.DeleteS3Destination(c.Request().Context(), id, projectID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
