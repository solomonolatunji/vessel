package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/services"
)

type GitAppsHandler struct {
	gitAppsService *services.GitAppsService
}

func NewGitAppsHandler(gs *services.GitAppsService) *GitAppsHandler {
	return &GitAppsHandler{gitAppsService: gs}
}

type getFunc[T any] func(ctx context.Context, id string) (*T, error)
type listFunc[T any] func(ctx context.Context) ([]T, error)
type saveFunc[T any] func(ctx context.Context, app *T) error
type deleteFunc func(ctx context.Context, id string) error

type GitAppsManifestRequest struct {
	Code string `json:"code"`
}

func listAppsHandler[T any](list listFunc[T]) echo.HandlerFunc {
	return func(c echo.Context) error {
		apps, err := list(c.Request().Context())
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		if apps == nil {
			apps = []T{}
		}
		return utils.Success(c, "Operation successful", apps)
	}
}

func getAppHandler[T any](get getFunc[T]) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		app, err := get(c.Request().Context(), id)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		if app == nil {
			return utils.Error(c, http.StatusNotFound, "App not found")
		}
		return utils.Success(c, "Operation successful", app)
	}
}

func saveAppHandler[T any](save saveFunc[T]) echo.HandlerFunc {
	return func(c echo.Context) error {
		var app T
		if err := c.Bind(&app); err != nil {
			return utils.Error(c, http.StatusBadRequest, "invalid payload")
		}
		if err := save(c.Request().Context(), &app); err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		return utils.Success(c, "Operation successful", app)
	}
}

func deleteAppHandler(del deleteFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if err := del(c.Request().Context(), id); err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		return utils.Success(c, "Operation successful", map[string]string{"status": "deleted"})
	}
}

func (h *GitAppsHandler) ExchangeGithubManifestCode(c echo.Context) error {
	var payload GitAppsManifestRequest

	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if payload.Code == "" {
		return utils.Error(c, http.StatusBadRequest, "code is required")
	}

	app, err := h.gitAppsService.ExchangeGithubManifestCode(c.Request().Context(), payload.Code)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Operation successful", app)
}

func (h *GitAppsHandler) ListGithubApps(c echo.Context) error {
	return listAppsHandler(h.gitAppsService.ListGithubApps)(c)
}

func (h *GitAppsHandler) GetGithubApp(c echo.Context) error {
	return getAppHandler(h.gitAppsService.GetGithubApp)(c)
}

func (h *GitAppsHandler) SaveGithubApp(c echo.Context) error {
	return saveAppHandler(h.gitAppsService.SaveGithubApp)(c)
}

func (h *GitAppsHandler) DeleteGithubApp(c echo.Context) error {
	return deleteAppHandler(h.gitAppsService.DeleteGithubApp)(c)
}
