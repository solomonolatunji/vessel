package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WorkerWSHandler struct {
	hub        *engine.WorkerHub
	serverRepo repositories.ServerRepository
	upgrader   websocket.Upgrader
}

func NewWorkerWSHandler(hub *engine.WorkerHub, serverRepo repositories.ServerRepository) *WorkerWSHandler {
	return &WorkerWSHandler{
		hub:        hub,
		serverRepo: serverRepo,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// Connect upgrades the HTTP request to a WebSocket connection for a codedock-worker.
func (h *WorkerWSHandler) Connect(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	server, err := h.serverRepo.GetByToken(c.Request().Context(), token)
	if err != nil || server == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	conn := h.hub.Register(server.ID, ws)

	// Send an immediate auth acknowledgment
	authAck := models.WorkerAuthResultPayload{Success: true}
	authAckBytes, _ := json.Marshal(authAck)

	_ = conn.Send(&models.WorkerMessage{
		ID:        uuid.New().String(),
		Type:      models.WorkerMessageTypeAuthResult,
		Timestamp: time.Now(),
		Payload:   authAckBytes,
	})

	// Start listening for incoming telemetry and acks from the worker
	conn.Listen()

	return nil
}
