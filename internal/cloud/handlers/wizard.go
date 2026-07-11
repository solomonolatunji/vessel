package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"vessel.dev/vessel/internal/cloud/repos"
)

type WizardHandler struct {
	repo repos.CloudRepo
}

func NewWizardHandler(repo repos.CloudRepo) *WizardHandler {
	return &WizardHandler{repo: repo}
}

// @Summary Generate an agent connection token
// @Description Generates a unique secure token that the remote server uses to connect to Vessel Cloud
// @Tags Cloud-Wizard
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /cloud/wizard/token [post]
func (h *WizardHandler) GenerateAgentToken(c echo.Context) error {
	token, err := generateSecureToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":           token,
		"install_command": "curl -fsSL https://get.vessel.dev/agent | sh -s -- --token=" + token,
		"expires_at":      time.Now().Add(24 * time.Hour),
	})
}

func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "vsl_live_" + hex.EncodeToString(b), nil
}
