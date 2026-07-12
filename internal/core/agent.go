package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
)

func Run(ctx context.Context, serverURL, token string) error {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)
	headers.Add("X-Vessl-Agent", "true")
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}
	conn, resp, err := dialer.DialContext(ctx, serverURL, headers)
	if err != nil {
		status := "unknown"
		if resp != nil {
			status = resp.Status
		}
		return fmt.Errorf("failed to connect to server %s (status: %s): %w", serverURL, status, err)
	}
	defer conn.Close()
	log.Printf(" Successfully connected to controller at %s", serverURL)
	netConn := &websocketConn{conn: conn}
	session, err := yamux.Server(netConn, yamux.DefaultConfig())
	if err != nil {
		return fmt.Errorf("failed to start yamux session: %w", err)
	}
	defer session.Close()
	log.Println(" Secure tunnel established. Listening for controller commands...")

	// Start telemetry background task
	go startTelemetryLoop(ctx, serverURL, token)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			stream, err := session.AcceptStream()
			if err != nil {
				if err == io.EOF {
					return fmt.Errorf("controller closed the tunnel")
				}
				log.Printf(" Yamux accept error: %v", err)
				continue
			}
			go handleDockerStream(stream)
		}
	}
}

func handleDockerStream(stream *yamux.Stream) {
	defer stream.Close()
	dockerConn, err := net.Dial("unix", "/var/run/docker.sock")
	if err != nil {
		log.Printf(" Failed to connect to local Docker daemon: %v", err)
		return
	}
	defer dockerConn.Close()
	errc := make(chan error, 2)
	go func() {
		_, err := io.Copy(dockerConn, stream)
		errc <- err
	}()
	go func() {
		_, err := io.Copy(stream, dockerConn)
		errc <- err
	}()
	<-errc
}

type websocketConn struct {
	conn *websocket.Conn
	r    io.Reader
}

func (c *websocketConn) Read(p []byte) (int, error) {
	if c.r == nil {
		var err error
		var msgType int
		msgType, c.r, err = c.conn.NextReader()
		if err != nil {
			return 0, err
		}
		if msgType != websocket.BinaryMessage {
			return 0, fmt.Errorf("unexpected message type: %d", msgType)
		}
	}
	n, err := c.r.Read(p)
	if err == io.EOF {
		c.r = nil
		return n, nil
	}
	return n, err
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

func (c *websocketConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *websocketConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *websocketConn) SetDeadline(t time.Time) error {
	if err := c.conn.SetReadDeadline(t); err != nil {
		return err
	}
	return c.conn.SetWriteDeadline(t)
}

func (c *websocketConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *websocketConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func startTelemetryLoop(ctx context.Context, serverURL, token string) {
	apiURL := strings.Replace(serverURL, "wss://", "https://", 1)
	apiURL = strings.Replace(apiURL, "ws://", "http://", 1)
	apiURL = strings.Replace(apiURL, "/agent/connect", "/billing/usage/report", 1)

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
			if err != nil {
				log.Printf(" Telemetry: failed to connect to docker: %v", err)
				continue
			}

			containers, err := cli.ContainerList(ctx, container.ListOptions{})
			if err != nil {
				log.Printf(" Telemetry: failed to list containers: %v", err)
				cli.Close()
				continue
			}

			vesslContainers := 0
			for _, c := range containers {
				if _, ok := c.Labels["vessl"]; ok {
					vesslContainers++
				}
			}
			cli.Close()

			payload := map[string]interface{}{
				"deployments":     0,
				"container_hours": vesslContainers, // 1 hour elapsed * N containers
				"bandwidth_gb":    0,
			}
			body, _ := json.Marshal(payload)

			req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(body))
			if err == nil {
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json")
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Printf(" Telemetry: failed to report usage: %v", err)
				} else {
					resp.Body.Close()
				}
			}
		}
	}
}
