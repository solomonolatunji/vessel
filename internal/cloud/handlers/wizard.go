package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type WizardHandler struct {
	// In reality we would inject a CloudDB repo here to store tokens
	// db *repos.CloudDB
}

func NewWizardHandler() *WizardHandler {
	return &WizardHandler{}
}

// GenerateAgentToken creates a new token and stores it in the db for a server
// @Summary Generate an agent connection token
// @Description Generates a unique secure token that the remote server uses to connect to Vessel Cloud
// @Tags Cloud-Wizard
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /cloud/wizard/token [post]
func (h *WizardHandler) GenerateAgentToken(c echo.Context) error {
	// 1. Validate user's team ID from session/JWT
	// teamID := "team_default" // Mocked

	// 2. Generate secure token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}
	token := "vsl_live_" + hex.EncodeToString(bytes)

	// 3. Save to database (mocked for now)
	// err := h.db.CreateServerToken(teamID, "My New Server", token)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":           token,
		"install_command": "curl -fsSL https://get.vessel.dev/agent | sh -s -- --token=" + token,
		"expires_at":      time.Now().Add(24 * time.Hour), // The token is valid for 24h to connect for the first time
	})
}
