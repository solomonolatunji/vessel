package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/services"
)

type ServerHandler struct {
	serverService services.ServerService
}

func NewServerHandler(serverService services.ServerService) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
	}
}

type CreateServerRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ipAddress"`
}

func (h *ServerHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req CreateServerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Server name is required")
	}

	server, err := h.serverService.CreateServer(c.Request().Context(), userID, req.Name, req.IPAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, server)
}

func (h *ServerHandler) List(c echo.Context) error {
	userID := c.Get("user_id").(string)

	servers, err := h.serverService.ListServersByUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, servers)
}
