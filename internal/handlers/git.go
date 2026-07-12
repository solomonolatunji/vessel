package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type GitHandler struct {
	gitService *services.GitService
}

func NewGitHandler(s *services.GitService) *GitHandler {
	return &GitHandler{gitService: s}
}

// @Summary Connect endpoint
// @Description Connect endpoint
// @Tags Git
// @Accept json
// @Produce json
// @Param request body models.GitConnectRequest true "Payload"
// @Router /git/connect [post]
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

// @Summary Status endpoint
// @Description Status endpoint
// @Tags Git
// @Accept json
// @Produce json
// @Router /git/status [get]
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

// @Summary Disconnect endpoint
// @Description Disconnect endpoint
// @Tags Git
// @Accept json
// @Produce json
// @Param provider path string true "provider"
// @Router /git/connect/{provider} [delete]
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

// @Summary ListRepos endpoint
// @Description ListRepos endpoint
// @Tags Git
// @Accept json
// @Produce json
// @Router /git/repos [get]
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
