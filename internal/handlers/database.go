package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type DatabaseHandler struct {
	databaseService *services.DatabaseService
}

func NewDatabaseHandler(s *services.DatabaseService) *DatabaseHandler {
	return &DatabaseHandler{databaseService: s}
}

func (h *DatabaseHandler) ListDatabases(c echo.Context) error {
	databases, err := h.databaseService.ListDatabases(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if databases == nil {
		databases = []*models.Database{}
	}
	return c.JSON(http.StatusOK, databases)
}

func (h *DatabaseHandler) CreateDatabase(c echo.Context) error {
	var req models.CreateDatabaseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	db, err := h.databaseService.CreateDatabaseFromRequest(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, db)
}

func (h *DatabaseHandler) GetDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing database id parameter"})
	}
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "database not found"})
	}
	return c.JSON(http.StatusOK, db)
}

func (h *DatabaseHandler) DeleteDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing database id parameter"})
	}
	if err := h.databaseService.DeleteDatabase(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *DatabaseHandler) StartDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing database id parameter"})
	}
	db, err := h.databaseService.StartDatabase(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, db)
}

func (h *DatabaseHandler) StopDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing database id parameter"})
	}
	if err := h.databaseService.StopDatabase(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "stopped"})
}
