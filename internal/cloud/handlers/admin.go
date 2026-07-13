package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"vessl.dev/vessl/internal/cloud/license"
	"vessl.dev/vessl/internal/cloud/repositories"
	"vessl.dev/vessl/internal/utils"
)

type AdminHandler struct {
	cloudRepo repos.CloudRepo
	authRepo  repos.AuthRepo
}

func NewAdminHandler(cloudRepo repos.CloudRepo, authRepo repos.AuthRepo) *AdminHandler {
	return &AdminHandler{
		cloudRepo: cloudRepo,
		authRepo:  authRepo,
	}
}

// GetSystemStats retrieves global platform metrics for staff
// @Summary Get System Stats
// @Description Fetch total users, servers, and active subscriptions (Staff Only)
// @Tags Cloud-Admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/stats [get]
func (h *AdminHandler) GetSystemStats(c echo.Context) error {
	ctx := c.Request().Context()
	totalUsers, _ := h.authRepo.GetTotalUsers(ctx)
	totalServers, _ := h.cloudRepo.GetTotalServers(ctx)
	activeSubs, _ := h.cloudRepo.GetActiveSubscriptions(ctx)

	return utils.Success(c, "Stats fetched successfully", map[string]interface{}{
		"total_users":          totalUsers,
		"total_servers":        totalServers,
		"active_subscriptions": activeSubs,
	})
}

// GetAuditLogs retrieves system audit logs
// @Summary Get Audit Logs
// @Description Fetch paginated audit logs (Staff Only)
// @Tags Cloud-Admin
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /admin/audit-logs [get]
func (h *AdminHandler) GetAuditLogs(c echo.Context) error {
	ctx := c.Request().Context()
	limit := 100
	offset := 0
	if l := c.QueryParam("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	if o := c.QueryParam("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}
	logs, err := h.cloudRepo.ListAuditLogs(ctx, limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to fetch audit logs")
	}
	return utils.Success(c, "Audit logs fetched successfully", logs)
}

// GenerateOfflineLicense generates a signed license for self-hosted instances
// @Summary Generate Offline License
// @Description Generates a signed license for self-hosted enterprise users
// @Tags Cloud-Admin
// @Accept json
// @Produce json
// @Param payload body map[string]interface{} true "Payload"
// @Success 200 {object} map[string]string
// @Router /admin/licenses [post]
func (h *AdminHandler) GenerateOfflineLicense(c echo.Context) error {
	var payload struct {
		TeamID   string `json:"team_id"`
		Plan     string `json:"plan"`
		MaxSeats int    `json:"max_seats"`
		Days     int    `json:"days"`
	}
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	privKey := os.Getenv("LICENSE_PRIVATE_KEY")
	if privKey == "" {
		return utils.Error(c, http.StatusInternalServerError, "license generation not configured")
	}

	expiry := time.Now().Add(time.Duration(payload.Days) * 24 * time.Hour)
	token, err := license.GenerateLicense(privKey, payload.TeamID, payload.Plan, payload.MaxSeats, expiry)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to generate license")
	}

	return utils.Success(c, "License generated successfully", map[string]string{
		"license_key": token,
		"expires_at":  expiry.Format(time.RFC3339),
	})
}
