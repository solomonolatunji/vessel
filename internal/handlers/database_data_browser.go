package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

func (h *DatabaseHandler) GetSchemas(c echo.Context) error {
	id := c.Param("id")
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	schemas, err := h.databaseService.GetSchemas(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", schemas)
}

func (h *DatabaseHandler) GetTableData(c echo.Context) error {
	id := c.Param("id")
	table := c.Param("table")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 100
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	data, err := h.databaseService.GetTableData(c.Request().Context(), id, table, limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", data)
}

func (h *DatabaseHandler) InsertTableRow(c echo.Context) error {
	id := c.Param("id")
	table := c.Param("table")
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	var data map[string]any
	if err := c.Bind(&data); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	res, err := h.databaseService.InsertTableRow(c.Request().Context(), id, table, data)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", res)
}

func (h *DatabaseHandler) UpdateTableRow(c echo.Context) error {
	id := c.Param("id")
	table := c.Param("table")
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	var req map[string]map[string]any
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	opts := services.UpdateTableRowOpts{
		ID:    id,
		Table: table,
		Keys:  req["keys"],
		Data:  req["data"],
	}
	res, err := h.databaseService.UpdateTableRow(c.Request().Context(), opts)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", res)
}

func (h *DatabaseHandler) DeleteTableRow(c echo.Context) error {
	id := c.Param("id")
	table := c.Param("table")
	db, err := h.databaseService.GetDatabase(c.Request().Context(), id)
	if err != nil || db == nil {
		return utils.Error(c, http.StatusNotFound, "database not found")
	}
	if err := h.verifyProjectOwnership(c, db.ProjectID); err != nil {
		return err
	}
	var keys map[string]any
	if err := c.Bind(&keys); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	res, err := h.databaseService.DeleteTableRow(c.Request().Context(), id, table, keys)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", res)
}
