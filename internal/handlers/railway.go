package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type RailwayHandler struct {
	service *services.RailwayService
}

func NewRailwayHandler(s *services.RailwayService) *RailwayHandler {
	return &RailwayHandler{service: s}
}

// @Summary List Railway Projects
// @Description Fetches projects from Railway API
// @Tags System
// @Produce json
// @Param token query string true "Railway Personal API Token"
// @Success 200 {object} map[string]any
// @Router /system/migration/railway/projects [get]
func (h *RailwayHandler) GetProjects(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return utils.Error(c, http.StatusBadRequest, "token is required")
	}

	projects, err := h.service.ListProjects(c.Request().Context(), token)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "fetched projects", projects)
}

type railwayImportRequest struct {
	Token              string `json:"token"`
	ProjectID          string `json:"projectId"`
	ExcludeRailwayVars bool   `json:"excludeRailwayVars"`
	RecreateDatabases  bool   `json:"recreateDatabases"`
	ImportData         bool   `json:"importData"`
}

// @Summary Import Railway Project
// @Description Imports a project from Railway
// @Tags System
// @Accept json
// @Produce json
// @Param req body railwayImportRequest true "Import request"
// @Success 200 {object} map[string]any
// @Router /system/migration/railway/import [post]
func (h *RailwayHandler) ImportProject(c echo.Context) error {
	var req railwayImportRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}

	if req.Token == "" || req.ProjectID == "" {
		return utils.Error(c, http.StatusBadRequest, "token and projectId are required")
	}

	err := h.service.ImportProject(c.Request().Context(), req.Token, req.ProjectID, req.ExcludeRailwayVars, req.RecreateDatabases, req.ImportData)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "import started", nil)
}
