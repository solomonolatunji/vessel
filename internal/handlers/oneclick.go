package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type OneClickHandler struct {
	service *services.OneClickService
}

func NewOneClickHandler(s *services.OneClickService) *OneClickHandler {
	return &OneClickHandler{service: s}
}

type oneClickDeployRequest struct {
	AppID     string `json:"appId" form:"appId"`
	ProjectID string `json:"projectId" form:"projectId"`
	Name      string `json:"name" form:"name"`
}

// @Summary List one-click apps
// @Description Returns available one-click deployable applications
// @Tags OneClick
// @Produce json
// @Success 200 {object} map[string]any
// @Router /one-click [get]
func (h *OneClickHandler) List(c echo.Context) error {
	return utils.Success(c, "Available one-click apps", h.service.ListApps())
}

// @Summary Deploy a one-click app
// @Description Deploys a pre-configured application from template
// @Tags OneClick
// @Accept json
// @Produce json
// @Param request body oneClickDeployRequest true "Deployment details"
// @Success 200 {object} map[string]any
// @Router /one-click/deploy [post]
func (h *OneClickHandler) Deploy(c echo.Context) error {
	var req oneClickDeployRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	db, err := h.service.DeployApp(c.Request().Context(), req.AppID, req.ProjectID, req.Name)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "App deployed", db)
}
