package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"

	"codedock.run/codedock/internal/models"
)

type WorkerDaemon struct {
	serverURL string
	token     string
	conn      *websocket.Conn
}

func NewWorkerDaemon(serverURL, token string) *WorkerDaemon {
	return &WorkerDaemon{
		serverURL: serverURL,
		token:     token,
	}
}

func (d *WorkerDaemon) Start(ctx context.Context) error {
	u, err := url.Parse(d.serverURL)
	if err != nil {
		return err
	}
	u.Path = "/api/ws/worker"
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	}

	header := http.Header{}
	header.Add("Authorization", "Bearer "+d.token)

	slog.Info("dialing control plane", "url", u.String())
	conn, resp, err := websocket.DefaultDialer.DialContext(ctx, u.String(), header)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("dial failed with status %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("dial failed: %w", err)
	}
	d.conn = conn

	// Send auth message
	authPayload := []byte(fmt.Sprintf(`{"token":"%s"}`, d.token))
	err = d.conn.WriteJSON(models.WorkerMessage{
		Type:      models.WorkerMessageTypeAuth,
		Timestamp: time.Now().UTC(),
		Payload:   authPayload,
	})
	if err != nil {
		return fmt.Errorf("failed to send auth message: %w", err)
	}

	go d.listen(ctx)
	return nil
}

func (d *WorkerDaemon) listen(ctx context.Context) {
	defer d.conn.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var msg models.WorkerMessage
		if err := d.conn.ReadJSON(&msg); err != nil {
			slog.Error("failed to read websocket message", "err", err)
			time.Sleep(2 * time.Second) // basic backoff
			continue
		}

		go d.handleMessage(ctx, msg)
	}
}

func (d *WorkerDaemon) handleMessage(ctx context.Context, msg models.WorkerMessage) {
	switch msg.Type {
	case models.WorkerMessageTypeDeployApp:
		var payload models.WorkerDeployAppPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			slog.Error("failed to unmarshal deploy payload", "err", err)
			return
		}
		slog.Info("received deploy app command", "app_id", payload.AppID, "deployment_id", payload.DeploymentID)

		go d.executeDeployment(ctx, msg.ID, payload)
	default:
		slog.Warn("unhandled message type", "type", msg.Type)
	}
}

func (d *WorkerDaemon) executeDeployment(ctx context.Context, commandID string, payload models.WorkerDeployAppPayload) {
	slog.Info("executing deployment", "app_id", payload.AppID)

	err := d.processDeployment(ctx, commandID, payload)

	ack := models.WorkerCommandAckPayload{
		CommandID: commandID,
		Success:   err == nil,
	}
	if err != nil {
		slog.Error("deployment failed", "app_id", payload.AppID, "err", err)
		ack.Error = err.Error()
	} else {
		slog.Info("deployment succeeded", "app_id", payload.AppID)
	}

	ackBytes, _ := json.Marshal(ack)
	_ = d.conn.WriteJSON(models.WorkerMessage{
		Type:      models.WorkerMessageTypeCommandAck,
		Timestamp: time.Now().UTC(),
		Payload:   ackBytes,
	})
}
