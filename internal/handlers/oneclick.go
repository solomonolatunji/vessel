package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.dev/codedock/internal/services"
	"codedock.dev/codedock/internal/utils"
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

func (h *OneClickHandler) List(c echo.Context) error {
	return utils.Success(c, "Available one-click apps", h.service.ListApps())
}

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
