package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"vessl.dev/vessl/internal/services"
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

// @Summary SaveCode endpoint
// @Description SaveCode endpoint
// @Tags Services
// @Accept json
// @Produce json
// @Param serviceId path string true "serviceId"
// @Router /api/services/{serviceId}/serverless/code [post]
func (h *ServerlessHandler) SaveCode(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "service ID is required"})
	}

	var req SaveCodeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	code, err := h.serverlessService.SaveCode(c.Request().Context(), serviceID, req.Runtime, req.CodeContent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "code saved successfully",
		"code":    code,
	})
}

// @Summary GetCode endpoint
// @Description GetCode endpoint
// @Tags Services
// @Accept json
// @Produce json
// @Param serviceId path string true "serviceId"
// @Router /api/services/{serviceId}/serverless/code [get]
func (h *ServerlessHandler) GetCode(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "service ID is required"})
	}

	code, err := h.serverlessService.GetCode(c.Request().Context(), serviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "code not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
	})
}
