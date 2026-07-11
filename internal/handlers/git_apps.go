package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type GitAppsHandler struct {
	gitAppsService *services.GitAppsService
}

func NewGitAppsHandler(gs *services.GitAppsService) *GitAppsHandler {
	return &GitAppsHandler{gitAppsService: gs}
}

type getFunc[T any] func(ctx context.Context, id string) (*T, error)
type listFunc[T any] func(ctx context.Context, teamID string) ([]T, error)
type saveFunc[T any] func(ctx context.Context, app *T) error
type deleteFunc func(ctx context.Context, id string) error

func listAppsHandler[T any](list listFunc[T]) echo.HandlerFunc {
	return func(c echo.Context) error {
		teamID := c.QueryParam("teamId")
		if teamID == "" {
			teamID = "default"
		}
		apps, err := list(c.Request().Context(), teamID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if apps == nil {
			apps = []T{}
		}
		return c.JSON(http.StatusOK, apps)
	}
}

func getAppHandler[T any](get getFunc[T]) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		app, err := get(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if app == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "App not found"})
		}
		return c.JSON(http.StatusOK, app)
	}
}

func saveAppHandler[T any](save saveFunc[T], setTeamID func(*T, string)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var app T
		if err := c.Bind(&app); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		}
		setTeamID(&app, "default")
		if err := save(c.Request().Context(), &app); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, app)
	}
}

func deleteAppHandler(del deleteFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if err := del(c.Request().Context(), id); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// @Summary ExchangeGithubManifestCode endpoint
// @Description ExchangeGithubManifestCode endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/github/manifest-callback [post]
func (h *GitAppsHandler) ExchangeGithubManifestCode(c echo.Context) error {
	var payload struct {
		Code   string `json:"code"`
		TeamID string `json:"teamId"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if payload.Code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "code is required"})
	}

	app, err := h.gitAppsService.ExchangeGithubManifestCode(c.Request().Context(), payload.Code, payload.TeamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, app)
}

// @Summary ListGithubApps endpoint
// @Description ListGithubApps endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/github [get]
func (h *GitAppsHandler) ListGithubApps(c echo.Context) error {
	return listAppsHandler(h.gitAppsService.ListGithubApps)(c)
}

// @Summary GetGithubApp endpoint
// @Description GetGithubApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/github/{id} [get]
func (h *GitAppsHandler) GetGithubApp(c echo.Context) error {
	return getAppHandler(h.gitAppsService.GetGithubApp)(c)
}

// @Summary SaveGithubApp endpoint
// @Description SaveGithubApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/github [put]
func (h *GitAppsHandler) SaveGithubApp(c echo.Context) error {
	return saveAppHandler(h.gitAppsService.SaveGithubApp, func(a *models.GithubApp, t string) {
		if a.TeamID == "" {
			a.TeamID = t
		}
	})(c)
}

// @Summary DeleteGithubApp endpoint
// @Description DeleteGithubApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/github/{id} [delete]
func (h *GitAppsHandler) DeleteGithubApp(c echo.Context) error {
	return deleteAppHandler(h.gitAppsService.DeleteGithubApp)(c)
}

// @Summary ListGitlabApps endpoint
// @Description ListGitlabApps endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/gitlab [get]
func (h *GitAppsHandler) ListGitlabApps(c echo.Context) error {
	return listAppsHandler(h.gitAppsService.ListGitlabApps)(c)
}

// @Summary GetGitlabApp endpoint
// @Description GetGitlabApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/gitlab/{id} [get]
func (h *GitAppsHandler) GetGitlabApp(c echo.Context) error {
	return getAppHandler(h.gitAppsService.GetGitlabApp)(c)
}

// @Summary SaveGitlabApp endpoint
// @Description SaveGitlabApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/gitlab [put]
func (h *GitAppsHandler) SaveGitlabApp(c echo.Context) error {
	return saveAppHandler(h.gitAppsService.SaveGitlabApp, func(a *models.GitlabApp, t string) {
		if a.TeamID == "" {
			a.TeamID = t
		}
	})(c)
}

// @Summary DeleteGitlabApp endpoint
// @Description DeleteGitlabApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/gitlab/{id} [delete]
func (h *GitAppsHandler) DeleteGitlabApp(c echo.Context) error {
	return deleteAppHandler(h.gitAppsService.DeleteGitlabApp)(c)
}

// @Summary ListBitbucketApps endpoint
// @Description ListBitbucketApps endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/bitbucket [get]
func (h *GitAppsHandler) ListBitbucketApps(c echo.Context) error {
	return listAppsHandler(h.gitAppsService.ListBitbucketApps)(c)
}

// @Summary GetBitbucketApp endpoint
// @Description GetBitbucketApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/bitbucket/{id} [get]
func (h *GitAppsHandler) GetBitbucketApp(c echo.Context) error {
	return getAppHandler(h.gitAppsService.GetBitbucketApp)(c)
}

// @Summary SaveBitbucketApp endpoint
// @Description SaveBitbucketApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/bitbucket [put]
func (h *GitAppsHandler) SaveBitbucketApp(c echo.Context) error {
	return saveAppHandler(h.gitAppsService.SaveBitbucketApp, func(a *models.BitbucketApp, t string) {
		if a.TeamID == "" {
			a.TeamID = t
		}
	})(c)
}

// @Summary DeleteBitbucketApp endpoint
// @Description DeleteBitbucketApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/bitbucket/{id} [delete]
func (h *GitAppsHandler) DeleteBitbucketApp(c echo.Context) error {
	return deleteAppHandler(h.gitAppsService.DeleteBitbucketApp)(c)
}
