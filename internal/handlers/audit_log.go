package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type AuditLogHandler struct {
	auditService *services.AuditService
}

func NewAuditLogHandler(as *services.AuditService) *AuditLogHandler {
	return &AuditLogHandler{auditService: as}
}

// @Summary List Audit Logs
// @Description Fetches audit logs for the dashboard
// @Tags Audit
// @Accept json
// @Produce json
// @Param limit query int false "Max lines to fetch"
// @Param offset query int false "Offset for pagination"
// @Router /audit-logs [get]
func (h *AuditLogHandler) List(c echo.Context) error {
	limitParam := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 100 // default limit
	}

	offsetParam := c.QueryParam("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0 // default offset
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
