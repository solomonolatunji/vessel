package handlers

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
	"github.com/labstack/echo/v4"
)

type RegistryHandler struct {
	service services.RegistryService
}

func NewRegistryHandler(service services.RegistryService) *RegistryHandler {
	return &RegistryHandler{service: service}
}

func (h *RegistryHandler) Create(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, 400, "projectId is required")
	}

	var req struct {
		Name          string `json:"name"`
		RegistryURL   string `json:"registryUrl"`
		Username      string `json:"username"`
		PasswordToken string `json:"passwordToken"`
	}

	if err := c.Bind(&req); err != nil {
		return utils.Error(c, 400, "invalid request")
	}
	if req.Name == "" || req.RegistryURL == "" {
		return utils.Error(c, 400, "name and registryUrl are required")
	}

	registry := &models.Registry{
		ProjectID:     projectID,
		Name:          req.Name,
		RegistryURL:   req.RegistryURL,
		Username:      req.Username,
		PasswordToken: req.PasswordToken,
	}

	if err := h.service.CreateRegistry(c.Request().Context(), registry); err != nil {
		return utils.Error(c, 500, "failed to create registry")
	}

	return utils.Success(c, "registry created", registry)
}

func (h *RegistryHandler) List(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, 400, "projectId is required")
	}

	registries, err := h.service.ListRegistriesByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, 500, "failed to list registries")
	}

	return utils.Success(c, "registries listed", registries)
}

func (h *RegistryHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, 400, "id is required")
	}

	if err := h.service.DeleteRegistry(c.Request().Context(), id); err != nil {
		return utils.Error(c, 500, "failed to delete registry")
	}

	return utils.Success(c, "registry deleted", map[string]bool{"success": true})
}
