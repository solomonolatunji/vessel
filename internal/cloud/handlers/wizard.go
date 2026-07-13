package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	repos "vessl.dev/vessl/internal/cloud/repositories"
	"vessl.dev/vessl/internal/utils"
)

type WizardHandler struct {
	repo repos.CloudRepo
}

func NewWizardHandler(repo repos.CloudRepo) *WizardHandler {
	return &WizardHandler{repo: repo}
}

// @Summary Generate an agent connection token
// @Description Generates a unique secure token that the remote server uses to connect to Vessl Cloud
// @Tags Cloud-Wizard
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /wizard/token [post]
func (h *WizardHandler) GenerateAgentToken(c echo.Context) error {
	token, err := generateSecureToken()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "Failed to generate token")
	}

	teamIDRaw := c.Get("team_id")
	if teamIDRaw != nil {
		teamID := teamIDRaw.(uint)
		err = h.repo.RegisterServer(context.Background(), teamID, token, "New Server", c.RealIP())
		if err != nil {
			log.Printf("Failed to register server: %v", err)
		}
	}

	return utils.Success(c, "Success", map[string]interface{}{
		"token":           token,
		"install_command": "curl -fsSL https://get.vessl.dev/agent | sh -s -- --token=" + token,
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
