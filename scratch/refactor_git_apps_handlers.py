import re

with open("internal/handlers/git_apps.go", "r") as f:
    content = f.read()

# Add the factory functions after NewGitAppsHandler
factory_funcs = """
// --- Handler Factories ---

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
"""

content = content.replace("func NewGitAppsHandler(gs *services.GitAppsService) *GitAppsHandler {\n\treturn &GitAppsHandler{gitAppsService: gs}\n}",
                          "func NewGitAppsHandler(gs *services.GitAppsService) *GitAppsHandler {\n\treturn &GitAppsHandler{gitAppsService: gs}\n}\n" + factory_funcs)

# Now, we replace the body of the handlers using regex.

replacements = {
    r"func \(h \*GitAppsHandler\) ListGithubApps\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) ListGithubApps(c echo.Context) error {\n\treturn listAppsHandler(h.gitAppsService.ListGithubApps)(c)\n}",
    r"func \(h \*GitAppsHandler\) GetGithubApp\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) GetGithubApp(c echo.Context) error {\n\treturn getAppHandler(h.gitAppsService.GetGithubApp)(c)\n}",
    r"func \(h \*GitAppsHandler\) SaveGithubApp\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) SaveGithubApp(c echo.Context) error {\n\treturn saveAppHandler(h.gitAppsService.SaveGithubApp, func(a *models.GithubApp, t string) {\n\t\tif a.TeamID == \"\" {\n\t\t\ta.TeamID = t\n\t\t}\n\t})(c)\n}",
    r"func \(h \*GitAppsHandler\) DeleteGithubApp\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) DeleteGithubApp(c echo.Context) error {\n\treturn deleteAppHandler(h.gitAppsService.DeleteGithubApp)(c)\n}",

    r"func \(h \*GitAppsHandler\) ListGitlabApps\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) ListGitlabApps(c echo.Context) error {\n\treturn listAppsHandler(h.gitAppsService.ListGitlabApps)(c)\n}",
    r"func \(h \*GitAppsHandler\) GetGitlabApp\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) GetGitlabApp(c echo.Context) error {\n\treturn getAppHandler(h.gitAppsService.GetGitlabApp)(c)\n}",
    r"func \(h \*GitAppsHandler\) SaveGitlabApp\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) SaveGitlabApp(c echo.Context) error {\n\treturn saveAppHandler(h.gitAppsService.SaveGitlabApp, func(a *models.GitlabApp, t string) {\n\t\tif a.TeamID == \"\" {\n\t\t\ta.TeamID = t\n\t\t}\n\t})(c)\n}",
    r"func \(h \*GitAppsHandler\) DeleteGitlabApp\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) DeleteGitlabApp(c echo.Context) error {\n\treturn deleteAppHandler(h.gitAppsService.DeleteGitlabApp)(c)\n}",

    r"func \(h \*GitAppsHandler\) ListBitbucketApps\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) ListBitbucketApps(c echo.Context) error {\n\treturn listAppsHandler(h.gitAppsService.ListBitbucketApps)(c)\n}",
    r"func \(h \*GitAppsHandler\) GetBitbucketApp\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) GetBitbucketApp(c echo.Context) error {\n\treturn getAppHandler(h.gitAppsService.GetBitbucketApp)(c)\n}",
    r"func \(h \*GitAppsHandler\) SaveBitbucketApp\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) SaveBitbucketApp(c echo.Context) error {\n\treturn saveAppHandler(h.gitAppsService.SaveBitbucketApp, func(a *models.BitbucketApp, t string) {\n\t\tif a.TeamID == \"\" {\n\t\t\ta.TeamID = t\n\t\t}\n\t})(c)\n}",
    r"func \(h \*GitAppsHandler\) DeleteBitbucketApp\(c echo.Context\) error \{[\s\S]*?\n\}": "func (h *GitAppsHandler) DeleteBitbucketApp(c echo.Context) error {\n\treturn deleteAppHandler(h.gitAppsService.DeleteBitbucketApp)(c)\n}",
}

import contextlib

for pattern, replacement in replacements.items():
    content = re.sub(pattern, replacement, content)
    
if "context" not in content[:100]:
    content = content.replace('"net/http"', '"context"\n\t"net/http"')

with open("internal/handlers/git_apps.go", "w") as f:
    f.write(content)
