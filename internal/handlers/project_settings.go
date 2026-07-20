package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type ProjectSettingsHandler struct {
	settingsService *services.ProjectSettingsService
}

func NewProjectSettingsHandler(s *services.ProjectSettingsService) *ProjectSettingsHandler {
	return &ProjectSettingsHandler{settingsService: s}
}

func (h *ProjectSettingsHandler) ListTokens(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId")
	}
	list, err := h.settingsService.ListTokens(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", list)
}

func (h *ProjectSettingsHandler) CreateToken(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId")
	}
	var req models.CreateTokenRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	t := &models.ProjectToken{
		ProjectID:     projectID,
		Name:          req.Name,
		EnvironmentID: req.EnvironmentID,
		Scopes:        req.Scopes,
		IPAllowlist:   req.IPAllowlist,
		ExpiresAt:     req.ExpiresAt,
	}
	token, raw, err := h.settingsService.CreateToken(c.Request().Context(), t)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Token created successfully", map[string]any{
		"id":          token.ID,
		"name":        token.Name,
		"token":       raw,
		"scopes":      token.Scopes,
		"ipAllowlist": token.IPAllowlist,
		"expiresAt":   token.ExpiresAt,
		"createdAt":   token.CreatedAt,
	})
}

func (h *ProjectSettingsHandler) DeleteToken(c echo.Context) error {
	projectID := c.Param("projectId")
	id := c.Param("id")
	if projectID == "" || id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId or id")
	}
	if err := h.settingsService.DeleteToken(c.Request().Context(), id, projectID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ProjectSettingsHandler) ListMembers(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId")
	}
	list, err := h.settingsService.ListMembers(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", list)
}

func (h *ProjectSettingsHandler) AddMember(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId")
	}
	var req models.AddMemberRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Permission == "" {
		req.Permission = models.MemberPermissionMember
	}

	scheme := "http"
	if c.Request().TLS != nil || c.Request().Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	originUrl := scheme + "://" + c.Request().Host

	opts := services.AddMemberOpts{
		ProjectID:  projectID,
		Email:      req.Email,
		Permission: req.Permission,
		OriginURL:  originUrl,
	}
	added, err := h.settingsService.AddMemberByEmail(c.Request().Context(), opts)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", added)
}

func (h *ProjectSettingsHandler) RemoveMember(c echo.Context) error {
	projectID := c.Param("projectId")
	id := c.Param("id")
	if projectID == "" || id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId or id")
	}
	if err := h.settingsService.RemoveMember(c.Request().Context(), id, projectID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
