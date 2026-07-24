package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/services"
)

type ProjectEnvHandler struct {
	envService *services.EnvironmentService
}

func NewProjectEnvHandler(s *services.EnvironmentService) *ProjectEnvHandler {
	return &ProjectEnvHandler{envService: s}
}

func (h *ProjectEnvHandler) GetVars(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	vars, err := h.envService.GetVars(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if vars == nil {
		vars = map[string]string{}
	}
	return utils.Success(c, "Operation successful", vars)
}

func (h *ProjectEnvHandler) SetVars(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	var vars map[string]string
	if err := c.Bind(&vars); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	for k, v := range vars {
		if err := h.envService.SetVar(c.Request().Context(), projectID, k, v); err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
	}
	return utils.Success(c, "Operation successful", vars)
}
