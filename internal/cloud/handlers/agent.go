package handlers

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
	"github.com/labstack/echo/v4"
	repos "vessl.dev/vessl/internal/cloud/repositories"
	"vessl.dev/vessl/internal/cloud/services"
	"vessl.dev/vessl/internal/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type AgentHandler struct {
	repo            repos.CloudRepo
	meteringService services.MeteringService
	activeAgents    map[string]*yamux.Session
	mu              sync.RWMutex
}

func NewAgentHandler(repo repos.CloudRepo, meteringService services.MeteringService) *AgentHandler {
	return &AgentHandler{
		repo:            repo,
		meteringService: meteringService,
		activeAgents:    make(map[string]*yamux.Session),
	}
}

// @Summary Accept Agent Connection
// @Description Accepts an inbound WebSocket tunnel from a remote vessl daemon
// @Tags Cloud-Fleet
// @Success 101
// @Router /agent/connect [get]
func (h *AgentHandler) AcceptConnection(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return utils.Error(c, http.StatusUnauthorized, "Missing Authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {

		token = authHeader
	}

	server, err := h.repo.GetServerByToken(token)
	if err != nil || server == nil {
		return utils.Error(c, http.StatusUnauthorized, "Invalid or unknown Agent token")
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
	DryRun      bool     `json:"dry_run"`
}

// @Summary Fleet Deployment
// @Description Dispatch a deployment instruction to a subset of connected Vessl Daemons
// @Tags Cloud-Fleet
// @Accept json
// @Produce json
// @Success 202 {object} map[string]interface{}
// @Router /fleet/deploy [post]
func (h *AgentHandler) DeployToFleet(c echo.Context) error {
	var req FleetDeployRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request")
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	successfulDispatches := 0
	failedAgents := []string{}

	for _, token := range req.AgentTokens {
		if err := h.dispatchToAgent(token, req.ImageTag, req.DryRun); err != nil {
			failedAgents = append(failedAgents, token)
			continue
		}

		server, err := h.repo.GetServerByToken(token)
		if err == nil && server != nil {

			_ = h.meteringService.RecordUsage(server.WorkspaceID, 1, 0, 0)
		}

		successfulDispatches++
	}

	return utils.Accepted(c, "Fleet deployment dispatched", map[string]interface{}{
		"successes": successfulDispatches,
		"failures":  failedAgents,
	})
}

func (h *AgentHandler) dispatchToAgent(token, imageTag string, dryRun bool) error {
	session, exists := h.activeAgents[token]
	if !exists || session.IsClosed() {
		return fmt.Errorf("agent %s not connected", token)
	}

	if os.Getenv("DEPLOY_DRY_RUN") == "true" || dryRun {
		log.Printf("Dry-run mode is enabled. Skipping dispatch to agent %s for image %s", token, imageTag)
		return nil
	}

	dockerClient, err := h.GetDockerClient(token)
	if err != nil {
		return err
	}

	reader, err := dockerClient.ImagePull(context.Background(), imageTag, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	buf := make([]byte, 1024)
	for {
		_, err := reader.Read(buf)
		if err != nil {
			break
		}
	}

	log.Printf("Successfully pulled %s on agent %s", imageTag, token)
	return nil
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

// ProxyToServer intercepts an HTTP request meant for a user's server and forwards it down the Yamux tunnel.
func (h *AgentHandler) ProxyToServer(c echo.Context) error {
	serverID := c.Param("serverId")

	h.mu.RLock()
	session, exists := h.activeAgents[serverID]
	h.mu.RUnlock()

	if !exists || session.IsClosed() {
		return utils.Error(c, http.StatusBadGateway, "Server is not currently connected to the cloud")
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = "localhost" // The local Vessl instance inside the tunnel
			// Strip the proxy prefix so the local server just sees the raw API route
			req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api/cloud/servers/"+serverID+"/proxy")
		},
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return session.Open() // Send the HTTP request down the Yamux tunnel
			},
		},
	}

	proxy.ServeHTTP(c.Response().Writer, c.Request())
	return nil
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
