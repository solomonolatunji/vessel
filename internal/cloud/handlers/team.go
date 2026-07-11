package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"vessel.dev/vessel/internal/cloud/repos"
)

type TeamHandler struct {
	repo repos.CloudRepo
}

func NewTeamHandler(repo repos.CloudRepo) *TeamHandler {
	return &TeamHandler{repo: repo}
}

type UpdateBrandingPayload struct {
	CustomDomain string `json:"custom_domain"`
	LogoURL      string `json:"logo_url"`
	PrimaryColor string `json:"primary_color"`
}

// UpdateBranding updates the white-label branding for an enterprise team
// @Summary Update Team Branding
// @Description Updates the white-label branding for an enterprise team
// @Tags Cloud-Team
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} models.CloudTeam
// @Router /cloud/teams/{id}/branding [patch]
func (h *TeamHandler) UpdateBranding(c echo.Context) error {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
	}

	var payload UpdateBrandingPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}

	team, err := h.repo.GetTeamByID(uint(teamID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Team not found"})
	}

	if team.Plan != "enterprise" && team.Plan != "team" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Custom branding requires a Team or Enterprise plan"})
	}

	team.CustomDomain = payload.CustomDomain
	team.LogoURL = payload.LogoURL
	team.PrimaryColor = payload.PrimaryColor

	if err := h.repo.UpdateTeam(team); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update team branding"})
	}

	return c.JSON(http.StatusOK, team)
}
