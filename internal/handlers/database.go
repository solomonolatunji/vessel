package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

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
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil || user.Role == "admin" {
		return nil
	}
	p, err := h.projectService.GetProject(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "project not found"})
	}
	if p.TeamID != user.UserID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if databases == nil {
		databases = []*models.Database{}
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" {
		var filtered []*models.Database
		for _, db := range databases {
			p, err := h.projectService.GetProject(c.Request().Context(), db.ProjectID)
			if err == nil && p.TeamID == user.UserID {
				filtered = append(filtered, db)
			}
		}
		return c.JSON(http.StatusOK, filtered)
	}
	return c.JSON(http.StatusOK, databases)
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if err := h.verifyProjectOwnership(c, req.ProjectID); err != nil {
		return err
	}
	db, err := h.databaseService.CreateDatabaseFromRequest(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, db)
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing database id parameter"})
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "database not found"})
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, db)
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing database id parameter"})
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "database not found"})
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	if err := h.databaseService.DeleteDatabase(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing database id parameter"})
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "database not found"})
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	dbStarted, err := h.databaseService.StartDatabase(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dbStarted)
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing database id parameter"})
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "database not found"})
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	if err := h.databaseService.StopDatabase(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "stopped"})
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing database id parameter"})
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "database not found"})
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	var req models.DatabaseQueryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	res, err := h.databaseService.QueryDatabase(c.Request().Context(), id, req.Query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
