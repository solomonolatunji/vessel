package handlers

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type AgentHandler struct {
	activeAgents map[string]*yamux.Session
	mu           sync.RWMutex
}

func NewAgentHandler() *AgentHandler {
	return &AgentHandler{
		activeAgents: make(map[string]*yamux.Session),
	}
}

// @Summary Accept Agent Connection
// @Description Accepts an inbound WebSocket tunnel from a remote vessel daemon
// @Tags Cloud-Fleet
// @Success 101
// @Router /cloud/agent/connect [get]
func (h *AgentHandler) AcceptConnection(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
	}

	serverID := token

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Failed to upgrade agent connection: %v", err)
		return err
	}

	netConn := &websocketConn{conn: ws}
	sess, err := yamux.Client(netConn, yamux.DefaultConfig())
	if err != nil {
		ws.Close()
		log.Printf("Failed to establish yamux session: %v", err)
		return err
	}

	h.mu.Lock()
	h.activeAgents[serverID] = sess
	h.mu.Unlock()

	log.Printf("Agent connected: %s", serverID)

	<-sess.CloseChan()

	h.mu.Lock()
	delete(h.activeAgents, serverID)
	h.mu.Unlock()

	log.Printf("Agent disconnected: %s", serverID)
	return nil
}

type FleetDeployRequest struct {
	ImageTag    string   `json:"image_tag"`
	AgentTokens []string `json:"agent_tokens"`
	EnvVars     []string `json:"env_vars"`
}

// @Summary Fleet Deployment
// @Description Dispatch a deployment instruction to a subset of connected Vessel Daemons
// @Tags Cloud-Fleet
// @Accept json
// @Produce json
// @Success 202 {object} map[string]interface{}
// @Router /cloud/fleet/deploy [post]
func (h *AgentHandler) DeployToFleet(c echo.Context) error {
	var req FleetDeployRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	successfulDispatches := 0
	failedAgents := []string{}

	for _, token := range req.AgentTokens {
		if err := h.dispatchToAgent(token, req.ImageTag); err != nil {
			failedAgents = append(failedAgents, token)
			continue
		}
		successfulDispatches++
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"message":   "Fleet deployment dispatched",
		"successes": successfulDispatches,
		"failures":  failedAgents,
	})
}

func (h *AgentHandler) dispatchToAgent(token, imageTag string) error {
	session, exists := h.activeAgents[token]
	if !exists || session.IsClosed() {
		return fmt.Errorf("agent %s not connected", token)
	}

	stream, err := session.Open()
	if err != nil {
		return err
	}
	defer stream.Close()

	_, err = stream.Write([]byte("DEPLOY:" + imageTag + "\n"))
	return err
}

func (h *AgentHandler) GetDockerClient(serverID string) (*client.Client, error) {
	h.mu.RLock()
	session, exists := h.activeAgents[serverID]
	h.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("agent %s is not currently connected", serverID)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return session.Open()
			},
		},
	}

	return client.NewClientWithOpts(
		client.WithHTTPClient(httpClient),
		client.WithHost("http://localhost"),
		client.WithAPIVersionNegotiation(),
	)
}

type websocketConn struct {
	conn *websocket.Conn
}

func (c *websocketConn) Read(p []byte) (int, error) {
	_, r, err := c.conn.NextReader()
	if err != nil {
		return 0, err
	}
	return r.Read(p)
}

func (c *websocketConn) Write(p []byte) (int, error) {
	if err := c.conn.WriteMessage(websocket.BinaryMessage, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *websocketConn) Close() error                       { return c.conn.Close() }
func (c *websocketConn) LocalAddr() net.Addr                { return c.conn.LocalAddr() }
func (c *websocketConn) RemoteAddr() net.Addr               { return c.conn.RemoteAddr() }
func (c *websocketConn) SetDeadline(t time.Time) error      { return c.conn.SetWriteDeadline(t) }
func (c *websocketConn) SetReadDeadline(t time.Time) error  { return c.conn.SetReadDeadline(t) }
func (c *websocketConn) SetWriteDeadline(t time.Time) error { return c.conn.SetWriteDeadline(t) }
