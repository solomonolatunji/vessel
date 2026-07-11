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
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, restrict this
	},
}

type AgentHandler struct {
	// A map to hold active connections by token or tenant ID
	activeAgents map[string]*yamux.Session
	mu           sync.RWMutex
}

func NewAgentHandler() *AgentHandler {
	return &AgentHandler{
		activeAgents: make(map[string]*yamux.Session),
	}
}

// AcceptConnection handles incoming websocket connections from remote agents
func (h *AgentHandler) AcceptConnection(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
	}

	// TODO: Validate token against PostgreSQL (CloudDB) to get tenant/server ID
	// For now, we just use the raw token as the identifier
	serverID := token

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Failed to upgrade agent connection: %v", err)
		return err
	}

	netConn := &websocketConn{conn: ws}
	session, err := yamux.Client(netConn, yamux.DefaultConfig())
	if err != nil {
		ws.Close()
		log.Printf("Failed to establish yamux client session: %v", err)
		return err
	}

	h.mu.Lock()
	h.activeAgents[serverID] = session
	h.mu.Unlock()

	log.Printf("Agent connected securely for server/tenant: %s", serverID)

	// Keep the connection alive until the agent disconnects
	<-session.CloseChan()

	h.mu.Lock()
	delete(h.activeAgents, serverID)
	h.mu.Unlock()

	log.Printf("Agent disconnected: %s", serverID)
	return nil
}

// GetDockerClient returns a Docker client that routes traffic over the yamux tunnel
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

	cli, err := client.NewClientWithOpts(
		client.WithHTTPClient(httpClient),
		client.WithHost("http://localhost"), // Host is ignored because of DialContext
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client over tunnel: %w", err)
	}

	return cli, nil
}

// websocketConn wraps a gorilla/websocket to implement net.Conn
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
	err := c.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *websocketConn) Close() error {
	return c.conn.Close()
}

func (c *websocketConn) LocalAddr() net.Addr                { return c.conn.LocalAddr() }
func (c *websocketConn) RemoteAddr() net.Addr               { return c.conn.RemoteAddr() }
func (c *websocketConn) SetDeadline(t time.Time) error      { return c.conn.SetWriteDeadline(t) }
func (c *websocketConn) SetReadDeadline(t time.Time) error  { return c.conn.SetReadDeadline(t) }
func (c *websocketConn) SetWriteDeadline(t time.Time) error { return c.conn.SetWriteDeadline(t) }
