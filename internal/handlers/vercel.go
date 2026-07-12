package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/services"
)

type VercelHandler struct {
	vercelService *services.VercelService
}

func NewVercelHandler(vs *services.VercelService) *VercelHandler {
	return &VercelHandler{vercelService: vs}
}

// @Summary Callback endpoint
// @Description Callback endpoint
// @Tags Auth
// @Accept json
// @Produce json
// @Router /oauth/vercel/callback [get]
func (h *VercelHandler) Callback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return utils.Error(c, http.StatusBadRequest, "authorization code required")
	}

	user := GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	account, err := h.vercelService.HandleCallback(c.Request().Context(), user.UserID, code)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Vercel account linked successfully",
		"account": account,
	})
}

// @Summary ListProjects endpoint
// @Description ListProjects endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Router /vercel/projects [get]
func (h *VercelHandler) ListProjects(c echo.Context) error {
	user := GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	teamID := c.QueryParam("teamId")
	var tID *string
	if teamID != "" {
		tID = &teamID
	}

	projects, err := h.vercelService.ListProjects(c.Request().Context(), user.UserID, tID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"projects": projects,
	})
}

// @Summary GetProjectEnv endpoint
// @Description GetProjectEnv endpoint
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /vercel/projects/{id}/env [get]
func (h *VercelHandler) GetProjectEnv(c echo.Context) error {
	user := GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project ID required")
	}

	teamID := c.QueryParam("teamId")
	var tID *string
	if teamID != "" {
		tID = &teamID
	}

	envs, err := h.vercelService.GetProjectEnvVars(c.Request().Context(), user.UserID, tID, projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"envs": envs,
	})
}
