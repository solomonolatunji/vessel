package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.dev/codedock/internal/utils"

	"codedock.dev/codedock/internal/http/middleware"
	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/services"
)

type DatabaseHandler struct {
	databaseService *services.DatabaseService
	projectService  *services.ProjectService
}

func NewDatabaseHandler(s *services.DatabaseService, ps *services.ProjectService) *DatabaseHandler {
	return &DatabaseHandler{databaseService: s, projectService: ps}
}

func (h *DatabaseHandler) verifyProjectOwnership(c echo.Context, projectID string) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	project, err := h.projectService.GetProject(c.Request().Context(), projectID)
	if err != nil || project == nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}

	if !h.projectService.IsMemberOrOwner(c.Request().Context(), projectID, user.UserID, user.Role) {
		return utils.Error(c, http.StatusForbidden, "access denied")
	}
	return nil
}

func (h *DatabaseHandler) ListDatabases(c echo.Context) error {
	databases, err := h.databaseService.ListDatabases(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if databases == nil {
		databases = []*models.Database{}
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" {
		var filtered []*models.Database
		for _, db := range databases {
			_, err := h.projectService.GetProject(c.Request().Context(), db.ProjectID)
			if err == nil {
				filtered = append(filtered, db)
			}
		}
		return utils.Success(c, "Operation successful", filtered)
	}
	return utils.Success(c, "Operation successful", databases)
}

func (h *DatabaseHandler) CreateDatabase(c echo.Context) error {
	var req models.CreateDatabaseRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if err := h.verifyProjectOwnership(c, req.ProjectID); err != nil {
		return err
	}
	db, err := h.databaseService.CreateDatabaseFromRequest(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", db)
}

func (h *DatabaseHandler) GetDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	return utils.Success(c, "Operation successful", db)
}

func (h *DatabaseHandler) UpdateDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	var req models.UpdateDatabaseRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	db.ExternalDNS = req.ExternalDNS
	db.CustomArgs = req.CustomArgs
	db.LogicalReplication = req.LogicalReplication
	db.CPULimit = req.CPULimit
	db.MemoryLimit = req.MemoryLimit
	if err := h.databaseService.UpdateDatabase(c.Request().Context(), db); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", db)
}

func (h *DatabaseHandler) DeleteDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	if err := h.databaseService.DeleteDatabase(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "deleted"})
}

func (h *DatabaseHandler) StartDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	dbStarted, err := h.databaseService.StartDatabase(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", dbStarted)
}

func (h *DatabaseHandler) StopDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	if err := h.databaseService.StopDatabase(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "stopped"})
}

func (h *DatabaseHandler) RestartDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	if err := h.databaseService.StopDatabase(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	dbStarted, err := h.databaseService.StartDatabase(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", dbStarted)
}

func (h *DatabaseHandler) ImportData(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	var req models.ImportDatabaseRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	go func() {
		_ = h.databaseService.ImportData(context.Background(), id, req.SourceURL)
	}()
	return utils.Success(c, "Import started", nil)
}

func (h *DatabaseHandler) QueryDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id parameter")
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	var req models.DatabaseQueryRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	res, err := h.databaseService.QueryDatabase(c.Request().Context(), id, req.Query)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", res)
}
