package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"vessl.dev/vessl/internal/cloud/repositories"
	"vessl.dev/vessl/internal/utils"
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
// @Router /teams/{id}/branding [patch]
func (h *TeamHandler) UpdateBranding(c echo.Context) error {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid team ID")
	}

	var payload UpdateBrandingPayload
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid payload")
	}

	team, err := h.repo.GetTeamByID(uint(teamID))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "Team not found")
	}

	if team.Plan != "enterprise" && team.Plan != "team" {
		return utils.Error(c, http.StatusForbidden, "Custom branding requires a Team or Enterprise plan")
	}

	team.CustomDomain = payload.CustomDomain
	team.LogoURL = payload.LogoURL
	team.PrimaryColor = payload.PrimaryColor

	if err := h.repo.UpdateTeam(team); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "Failed to update team branding")
	}

	return utils.Success(c, "Branding updated", team)
}

// @Summary List Team Servers
// @Description Fetch all active servers associated with a team
// @Tags Cloud-Team
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} []models.CloudServer
// @Router /teams/{id}/servers [get]
func (h *TeamHandler) ListServers(c echo.Context) error {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid team ID")
	}

	servers, err := h.repo.GetTeamServers(uint(teamID))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "Failed to list servers")
	}

	return utils.Success(c, "Servers retrieved", servers)
}
