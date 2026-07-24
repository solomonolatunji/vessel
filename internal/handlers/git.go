package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
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
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req models.GitConnectRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	gp, err := h.gitService.ConnectProvider(c.Request().Context(), userID, &req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Created(c, "Created successfully", gp)
}

func (h *GitHandler) Status(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	status, err := h.gitService.GetConnectedProviders(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", status)
}

func (h *GitHandler) Disconnect(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	provider := c.Param("provider")
	if provider == "" {
		return utils.Error(c, http.StatusBadRequest, "missing provider parameter")
	}
	if err := h.gitService.DisconnectProvider(c.Request().Context(), userID, provider); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "disconnected"})
}

func (h *GitHandler) ListRepos(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	provider := c.QueryParam("provider")
	repos, err := h.gitService.ListRepositories(c.Request().Context(), userID, provider)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Operation successful", repos)
}
