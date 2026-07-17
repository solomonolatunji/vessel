package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/utils"

	"vessl.dev/vessl/internal/http/middleware"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type DatabaseHandler struct {
	databaseService *services.DatabaseService
	projectService  *services.ProjectService
}

func NewDatabaseHandler(s *services.DatabaseService, ps *services.ProjectService) *DatabaseHandler {
	return &DatabaseHandler{databaseService: s, projectService: ps}
}

func (h *DatabaseHandler) verifyProjectOwnership(c echo.Context, projectID string) error {
	_, err := h.projectService.GetProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}
	return nil
}

// @Summary ListDatabases endpoint
// @Description ListDatabases endpoint
// @Tags Databases
// @Accept json
// @Produce json
// @Router /databases [get]
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

// @Summary CreateDatabase endpoint
// @Description CreateDatabase endpoint
// @Tags Databases
// @Accept json
// @Produce json
// @Param request body models.CreateDatabaseRequest true "Payload"
// @Router /databases [post]
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

// @Summary GetDatabase endpoint
// @Description GetDatabase endpoint
// @Tags Databases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /databases/{id} [get]
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

// @Summary UpdateDatabase endpoint
// @Description UpdateDatabase endpoint
// @Tags Databases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param request body models.UpdateDatabaseRequest true "Payload"
// @Router /databases/{id} [put]
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
	if err := h.databaseService.UpdateDatabase(c.Request().Context(), db); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", db)
}

// @Summary DeleteDatabase endpoint
// @Description DeleteDatabase endpoint
// @Tags Databases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /databases/{id} [delete]
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

// @Summary StartDatabase endpoint
// @Description StartDatabase endpoint
// @Tags Databases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /databases/{id}/start [post]
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

// @Summary StopDatabase endpoint
// @Description StopDatabase endpoint
// @Tags Databases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /databases/{id}/stop [post]
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

// @Summary RestartDatabase endpoint
// @Description RestartDatabase endpoint
// @Tags Databases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /databases/{id}/restart [post]
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

// @Summary ImportData endpoint
// @Description ImportData endpoint
// @Tags Databases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param request body models.ImportDatabaseRequest true "Payload"
// @Router /databases/{id}/import [post]
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

// @Summary QueryDatabase endpoint
// @Description QueryDatabase endpoint
// @Tags Databases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param request body models.DatabaseQueryRequest true "Payload"
// @Router /databases/{id}/query [post]
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
