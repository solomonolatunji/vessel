package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetProfile retrieves the authenticated user's profile and teams
// @Summary Get User Profile
// @Description Fetch current user details
// @Tags Cloud-Users
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /cloud/users/me [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	// TODO: Extract user ID from JWT, fetch from DB
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    "usr_123",
		"email": "user@example.com",
		"teams": []map[string]string{
			{"id": "team_1", "name": "Personal Team", "role": "owner"},
		},
	})
}
