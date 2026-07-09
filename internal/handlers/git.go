package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type GitHandler struct {
	gitService *services.GitService
}

func NewGitHandler(s *services.GitService) *GitHandler {
	return &GitHandler{gitService: s}
}

func (h *GitHandler) Connect(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var req models.GitConnectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	gp, err := h.gitService.ConnectProvider(c.Request().Context(), userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, gp)
}

func (h *GitHandler) Status(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	status, err := h.gitService.GetConnectedProviders(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, status)
}

func (h *GitHandler) Disconnect(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	provider := c.Param("provider")
	if provider == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing provider parameter"})
	}
	if err := h.gitService.DisconnectProvider(c.Request().Context(), userID, provider); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "disconnected"})
}

func (h *GitHandler) ListRepos(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	provider := c.QueryParam("provider")
	if provider == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing provider query parameter"})
	}
	repos, err := h.gitService.ListRepositories(c.Request().Context(), userID, provider)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, repos)
}
