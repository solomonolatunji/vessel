package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
)

type ServerlessHandler struct {
	serverlessService services.ServerlessService
}

func NewServerlessHandler(serverlessService services.ServerlessService) *ServerlessHandler {
	return &ServerlessHandler{serverlessService: serverlessService}
}

type SaveCodeRequest struct {
	Runtime     string `json:"runtime"`
	CodeContent string `json:"codeContent"`
}

func (h *ServerlessHandler) SaveCode(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "service ID is required")
	}

	var req SaveCodeRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request format")
	}

	code, err := h.serverlessService.SaveCode(c.Request().Context(), serviceID, req.Runtime, req.CodeContent)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "code saved successfully", map[string]interface{}{
		"code": code,
	})
}

func (h *ServerlessHandler) GetCode(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "service ID is required")
	}

	code, err := h.serverlessService.GetCode(c.Request().Context(), serviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.Error(c, http.StatusNotFound, "code not found")
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Code fetched successfully", map[string]interface{}{
		"code": code,
	})
}
