package handlers

import (
	"net/http"
	"strconv"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/services"
	"codedock.dev/codedock/internal/utils"
	"github.com/labstack/echo/v4"
)

type AuditLogHandler struct {
	auditService *services.AuditService
}

func NewAuditLogHandler(as *services.AuditService) *AuditLogHandler {
	return &AuditLogHandler{auditService: as}
}

func (h *AuditLogHandler) List(c echo.Context) error {
	limitParam := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 100
	}

	offsetParam := c.QueryParam("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	logs, err := h.auditService.ListLogs(c.Request().Context(), limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if logs == nil {
		logs = []models.AuditLog{}
	}

	return utils.Success(c, "Audit logs fetched", logs)
}
