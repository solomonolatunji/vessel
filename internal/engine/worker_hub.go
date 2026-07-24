package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"github.com/gorilla/websocket"
)

type WorkerHub struct {
	mu          sync.RWMutex
	connections map[string]*WorkerConnection
	serverRepo  repositories.ServerRepository
}

type WorkerConnection struct {
	ServerID string
	Conn     *websocket.Conn
	mu       sync.Mutex
	hub      *WorkerHub
	closeCh  chan struct{}
}

func NewWorkerHub(repo repositories.ServerRepository) *WorkerHub {
	return &WorkerHub{
		connections: make(map[string]*WorkerConnection),
		serverRepo:  repo,
	}
}

func (h *WorkerHub) Register(serverID string, conn *websocket.Conn) *WorkerConnection {
	h.mu.Lock()
	defer h.mu.Unlock()

	if existing, ok := h.connections[serverID]; ok {
		existing.Conn.Close()
	}

	wc := &WorkerConnection{
		ServerID: serverID,
		Conn:     conn,
		hub:      h,
		closeCh:  make(chan struct{}),
	}
	h.connections[serverID] = wc

	if err := h.serverRepo.UpdateStatus(context.Background(), serverID, models.ServerStatusOnline); err != nil {
		slog.Error("failed to update server status to online", "serverID", serverID, "err", err)
	}

	return wc
}

func (h *WorkerHub) Unregister(serverID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.connections[serverID]; ok {
		delete(h.connections, serverID)
		if err := h.serverRepo.UpdateStatus(context.Background(), serverID, models.ServerStatusOffline); err != nil {
			slog.Error("failed to update server status to offline", "serverID", serverID, "err", err)
		}
	}
}

func (h *WorkerHub) SendCommand(serverID string, msg *models.WorkerMessage) error {
	h.mu.RLock()
	wc, ok := h.connections[serverID]
	h.mu.RUnlock()

	if !ok {
		return fmt.Errorf("server %s is not connected", serverID)
	}

	return wc.Send(msg)
}

func (wc *WorkerConnection) Send(msg *models.WorkerMessage) error {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	return wc.Conn.WriteJSON(msg)
}

func (wc *WorkerConnection) Listen() {
	defer func() {
		wc.hub.Unregister(wc.ServerID)
		wc.Conn.Close()
		close(wc.closeCh)
	}()

	for {
		var msg models.WorkerMessage
		err := wc.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("unexpected websocket close", "serverID", wc.ServerID, "err", err)
			}
			break
		}

		go wc.handleMessage(msg)
	}
}

func (wc *WorkerConnection) handleMessage(msg models.WorkerMessage) {
	switch msg.Type {
	case models.WorkerMessageTypeMetrics:
		var metrics models.WorkerMetricsPayload
		if err := json.Unmarshal(msg.Payload, &metrics); err == nil {
			if err := wc.hub.serverRepo.UpdateMetrics(context.Background(), wc.ServerID, msg.Payload); err != nil {
				slog.Error("failed to update metrics", "serverID", wc.ServerID, "err", err)
			}
		}
	case models.WorkerMessageTypeLogStream:
		var logStream models.WorkerLogStreamPayload
		if err := json.Unmarshal(msg.Payload, &logStream); err == nil {
			// In the future, emit to pub/sub or SSE clients watching this container
			slog.Debug("received log from worker", "containerID", logStream.ContainerID, "line", logStream.LogLine)
		}
	case models.WorkerMessageTypeCommandAck:
		var ack models.WorkerCommandAckPayload
		if err := json.Unmarshal(msg.Payload, &ack); err == nil {
			slog.Info("received command ack", "serverID", wc.ServerID, "commandID", ack.CommandID, "success", ack.Success)
		}
	default:
		slog.Warn("unhandled worker message type", "serverID", wc.ServerID, "type", msg.Type)
	}
}
