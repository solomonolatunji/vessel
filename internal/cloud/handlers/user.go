package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	repos "vessl.dev/vessl/internal/cloud/repositories"
	"vessl.dev/vessl/internal/utils"
)

type UserHandler struct {
	cloudRepo repos.CloudRepo
	authRepo  repos.AuthRepo
}

func NewUserHandler(cloudRepo repos.CloudRepo, authRepo repos.AuthRepo) *UserHandler {
	return &UserHandler{cloudRepo: cloudRepo, authRepo: authRepo}
}

// @Summary Get User Profile
// @Description Fetch current user details
// @Tags Cloud-Users
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "Unauthorized")
	}

	user, err := h.authRepo.GetUserByID(c.Request().Context(), userID)
	if err != nil || user == nil {
		return utils.Error(c, http.StatusNotFound, "User not found")
	}

	teams, err := h.cloudRepo.GetUserTeams(userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "Failed to load teams")
	}

	teamsData := make([]map[string]interface{}, 0)
	for _, t := range teams {
		servers, _ := h.cloudRepo.GetTeamServers(t.ID)
		serversData := make([]map[string]interface{}, 0)
		for _, s := range servers {
			serversData = append(serversData, map[string]interface{}{
				"id":   s.ID,
				"name": s.Name,
			})
		}

		teamsData = append(teamsData, map[string]interface{}{
			"id":            t.ID,
			"name":          t.Name,
			"plan":          t.Plan,
			"custom_domain": t.CustomDomain,
			"logo_url":      t.LogoURL,
			"servers":       serversData,
		})
	}

	return utils.Success(c, "Success", map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.FullName,
		"teams": teamsData,
	})
}
