package handlers

import (
	"errors"
	"net/http"
	"path/filepath"

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
	list, err := h.backupService.ListConfigs(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	for i := range list {
		list[i].DbPassword = "********"
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

// @Summary Update Backup
// @Description Update Backup
// @Tags Backups
// @Accept json
// @Produce json
// @Param id path string true "Backup ID"
// @Param request body models.BackupConfig true "Payload"
// @Router /backups/{id} [put]
func (h *BackupHandler) Update(c echo.Context) error {
	id := c.Param("id")

	existing, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || existing == nil {
		var notFoundErr *utils.NotFoundError
		if err != nil && !errors.As(err, &notFoundErr) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup config")
		}
		return utils.Error(c, http.StatusNotFound, "backup config not found")
	}



	var req struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		DbUser          string `json:"dbUser"`
		DbPassword      string `json:"dbPassword"`
		Schedule        string `json:"schedule"`
		Timezone        string `json:"timezone"`
		Timeout         int    `json:"timeout"`
		RetentionDays   int    `json:"retentionDays"`
		MaxBackups      int    `json:"maxBackups"`
		MaxStorageGB    int    `json:"maxStorageGB"`
		S3DestinationID string `json:"s3DestinationId"`
		DatabaseID      string `json:"databaseId"`
		BackupEnabled   *bool  `json:"backupEnabled"`
		S3Enabled       *bool  `json:"s3Enabled"`
		DisableLocal    *bool  `json:"disableLocal"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.DbUser != "" {
		existing.DbUser = req.DbUser
	}
	if req.DbPassword != "" {
		existing.DbPassword = req.DbPassword
	}
	if req.Schedule != "" {
		existing.Schedule = req.Schedule
	}
	if req.Timezone != "" {
		existing.Timezone = req.Timezone
	}
	if req.Timeout != 0 {
		existing.Timeout = req.Timeout
	}
	if req.RetentionDays != 0 {
		existing.RetentionDays = req.RetentionDays
	}
	if req.MaxBackups != 0 {
		existing.MaxBackups = req.MaxBackups
	}
	if req.MaxStorageGB != 0 {
		existing.MaxStorageGB = req.MaxStorageGB
	}
	if req.S3DestinationID != "" {
		existing.S3DestinationID = req.S3DestinationID
	}
	if req.DatabaseID != "" {
		existing.DatabaseID = req.DatabaseID
	}

	if req.BackupEnabled != nil {
		existing.BackupEnabled = *req.BackupEnabled
	}
	if req.S3Enabled != nil {
		existing.S3Enabled = *req.S3Enabled
	}
	if req.DisableLocal != nil {
		existing.DisableLocal = *req.DisableLocal
	}

	if err := h.backupService.UpdateConfig(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	existing.DbPassword = "********"
	return utils.Success(c, "Updated successfully", existing)
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
		var notFoundErr *utils.NotFoundError
		if err != nil && !errors.As(err, &notFoundErr) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup config")
		}
		return utils.Error(c, http.StatusNotFound, "backup config not found")
	}
	cfg.DbPassword = "********"
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
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id")
	}
	if err := h.backupService.DeleteConfig(c.Request().Context(), id); err != nil {
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
// @Router /backups/{id}/records [get]
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

// @Summary Download Backup Record
// @Description Download Backup Record
// @Tags Backups
// @Accept json
// @Produce application/octet-stream
// @Param id path string true "Backup ID"
// @Param recordId path string true "Record ID"
// @Router /backups/{id}/records/{recordId}/download [get]
func (h *BackupHandler) DownloadRecord(c echo.Context) error {
	id := c.Param("id")
	recordID := c.Param("recordId")
	if id == "" || recordID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id or recordId parameter")
	}
	rec, err := h.backupService.GetRecord(c.Request().Context(), recordID)
	if err != nil {
		var notFound *utils.NotFoundError
		if !errors.As(err, &notFound) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup record")
		}
		return utils.Error(c, http.StatusNotFound, "record not found")
	}
	if rec == nil {
		return utils.Error(c, http.StatusNotFound, "record not found")
	}
	if rec.BackupConfigID != id {
		return utils.Error(c, http.StatusNotFound, "record not found")
	}

	if rec.FilePath == "" {
		return utils.Error(c, http.StatusNotFound, "local backup file not available")
	}
	return c.Attachment(rec.FilePath, filepath.Base(rec.FilePath))
}

// @Summary Delete Backup Record
// @Description Delete Backup Record
// @Tags Backups
// @Accept json
// @Produce json
// @Param id path string true "Backup ID"
// @Param recordId path string true "Record ID"
// @Router /backups/{id}/records/{recordId} [delete]
func (h *BackupHandler) DeleteRecord(c echo.Context) error {
	id := c.Param("id")
	recordID := c.Param("recordId")
	if id == "" || recordID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id or recordId parameter")
	}
	rec, err := h.backupService.GetRecord(c.Request().Context(), recordID)
	if err != nil {
		var notFound *utils.NotFoundError
		if !errors.As(err, &notFound) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup record")
		}
		return utils.Error(c, http.StatusNotFound, "record not found")
	}
	if rec == nil {
		return utils.Error(c, http.StatusNotFound, "record not found")
	}
	if rec.BackupConfigID != id {
		return utils.Error(c, http.StatusNotFound, "record not found")
	}

	if err := h.backupService.DeleteRecord(c.Request().Context(), recordID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
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
	list, err := h.backupService.ListS3Destinations(c.Request().Context())
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
// @Router /s3-destinations [post]
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
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id")
	}
	if err := h.backupService.DeleteS3Destination(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
