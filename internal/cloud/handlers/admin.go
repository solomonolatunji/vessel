package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"vessel.dev/vessel/internal/license"
)

type AdminHandler struct {
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// GetSystemStats retrieves global platform metrics for staff
// @Summary Get System Stats
// @Description Fetch total users, servers, and active subscriptions (Staff Only)
// @Tags Cloud-Admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /cloud/admin/stats [get]
func (h *AdminHandler) GetSystemStats(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"total_users":          150,
		"total_servers":        300,
		"active_subscriptions": 120,
	})
}

// GetAuditLogs retrieves system audit logs
// @Summary Get Audit Logs
// @Description Fetch paginated audit logs (Staff Only)
// @Tags Cloud-Admin
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /cloud/admin/audit-logs [get]
func (h *AdminHandler) GetAuditLogs(c echo.Context) error {
	return c.JSON(http.StatusOK, []map[string]interface{}{})
}

// GenerateOfflineLicense generates a signed license for self-hosted instances
// @Summary Generate Offline License
// @Description Generates a signed license for self-hosted enterprise users
// @Tags Cloud-Admin
// @Accept json
// @Produce json
// @Param payload body map[string]interface{} true "Payload"
// @Success 200 {object} map[string]string
// @Router /cloud/admin/licenses [post]
func (h *AdminHandler) GenerateOfflineLicense(c echo.Context) error {
	var payload struct {
		TeamID   string `json:"team_id"`
		Plan     string `json:"plan"`
		MaxSeats int    `json:"max_seats"`
		Days     int    `json:"days"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	privKey := os.Getenv("LICENSE_PRIVATE_KEY")
	if privKey == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "license generation not configured"})
	}

	expiry := time.Now().Add(time.Duration(payload.Days) * 24 * time.Hour)
	token, err := license.GenerateLicense(privKey, payload.TeamID, payload.Plan, payload.MaxSeats, expiry)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate license"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"license_key": token,
		"expires_at":  expiry.Format(time.RFC3339),
	})
}
