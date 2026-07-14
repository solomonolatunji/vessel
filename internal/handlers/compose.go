package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/http/middleware"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type ComposeHandler struct {
	composeDeployer *engine.ComposeDeployer
	projectService  *services.ProjectService
	appService      *services.AppService
	envRepo         repositories.EnvironmentRepository
	appRepo         repositories.AppServiceRepository
}

func NewComposeHandler(
	cd *engine.ComposeDeployer,
	ps *services.ProjectService,
	as *services.AppService,
	er repositories.EnvironmentRepository,
	ar repositories.AppServiceRepository,
) *ComposeHandler {
	return &ComposeHandler{
		composeDeployer: cd,
		projectService:  ps,
		appService:      as,
		envRepo:         er,
		appRepo:         ar,
	}
}

type ComposeDeployRequest struct {
	ProjectID string `json:"projectId"`
}

// @Summary Deploy a docker-compose file
// @Description Parses and deploys all services defined in a docker-compose.yml
// @Tags Compose
// @Accept multipart/form-data
// @Produce json
// @Param projectId formData string false "Project ID (optional, uses default if empty)"
// @Param file formData file true "docker-compose.yml file"
// @Success 200 {object} map[string]any
// @Router /compose/deploy [post]
func (h *ComposeHandler) Deploy(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.FormValue("projectId")
	if projectID == "" {
		projectID = c.FormValue("project_id")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "compose file is required")
	}

	src, err := file.Open()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to read uploaded file")
	}
	defer src.Close()

	tmpDir := filepath.Join(os.TempDir(), "vessl-compose", uuid.New().String())
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to create temp directory")
	}
	defer os.RemoveAll(tmpDir)

	tmpPath := filepath.Join(tmpDir, "docker-compose.yml")
	dst, err := os.Create(tmpPath)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to save compose file")
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return utils.Error(c, http.StatusInternalServerError, "failed to write compose file")
	}
	dst.Close()

	services, err := h.composeDeployer.Deploy(c.Request().Context(), tmpPath, projectID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "compose deploy failed: "+err.Error())
	}

	return utils.Success(c, "Compose file deployed", map[string]any{
		"services": services,
		"count":    len(services),
	})
}
