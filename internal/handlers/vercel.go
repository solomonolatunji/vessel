package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessel.dev/vessel/internal/services"
)

type VercelHandler struct {
	vercelService *services.VercelService
}

func NewVercelHandler(vs *services.VercelService) *VercelHandler {
	return &VercelHandler{vercelService: vs}
}

// Callback endpoint (GET /api/v1/oauth/vercel/callback)
func (h *VercelHandler) Callback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "authorization code required"})
	}

	user := GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	account, err := h.vercelService.HandleCallback(c.Request().Context(), user.UserID, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Vercel account linked successfully",
		"account": account,
	})
}

// List Projects (GET /api/v1/vercel/projects)
func (h *VercelHandler) ListProjects(c echo.Context) error {
	user := GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	teamID := c.QueryParam("teamId")
	var tID *string
	if teamID != "" {
		tID = &teamID
	}

	projects, err := h.vercelService.ListProjects(c.Request().Context(), user.UserID, tID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"projects": projects,
	})
}

// Get Env Vars (GET /api/v1/vercel/projects/:id/env)
func (h *VercelHandler) GetProjectEnv(c echo.Context) error {
	user := GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	projectID := c.Param("id")
	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "project ID required"})
	}

	teamID := c.QueryParam("teamId")
	var tID *string
	if teamID != "" {
		tID = &teamID
	}

	envs, err := h.vercelService.GetProjectEnvVars(c.Request().Context(), user.UserID, tID, projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"envs": envs,
	})
}
