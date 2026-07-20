package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

func (h *AppHandler) ListWebhooks(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId")
	}
	existing, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || existing == nil {
		var notFoundErr *utils.NotFoundError
		if err != nil && !errors.As(err, &notFoundErr) {
			return utils.Error(c, http.StatusInternalServerError, "failed to look up app service")
		}
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	list, err := h.appService.ListWebhooks(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to list webhooks")
	}
	return utils.Success(c, "Operation successful", list)
}

func (h *AppHandler) CreateWebhook(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId")
	}
	existing, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || existing == nil {
		var notFoundErr *utils.NotFoundError
		if err != nil && !errors.As(err, &notFoundErr) {
			return utils.Error(c, http.StatusInternalServerError, "failed to look up app service")
		}
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	var req models.CreateWebhookRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	validURL, err := utils.ValidateURL(req.URL)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	req.URL = validURL
	webhook := models.Webhook{
		ServiceID:             serviceID,
		URL:                   req.URL,
		EventTypes:            req.EventTypes,
		IncludePREnvironments: req.IncludePREnvironments,
	}
	created, err := h.appService.CreateWebhook(c.Request().Context(), &webhook)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to create webhook")
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *AppHandler) DeleteWebhook(c echo.Context) error {
	serviceID := c.Param("id")
	webhookID := c.Param("webhookId")
	if serviceID == "" || webhookID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId or webhookId")
	}
	existing, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || existing == nil {
		var notFoundErr *utils.NotFoundError
		if err != nil && !errors.As(err, &notFoundErr) {
			return utils.Error(c, http.StatusInternalServerError, "failed to look up app service")
		}
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	if err := h.appService.DeleteWebhook(c.Request().Context(), webhookID, serviceID); err != nil {
		var notFoundErr *utils.NotFoundError
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, http.StatusNotFound, "webhook not found")
		}
		return utils.Error(c, http.StatusInternalServerError, "failed to delete webhook")
	}
	return c.NoContent(http.StatusNoContent)
}
