package agent

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
)

// Run connects to the control plane and multiplexes the local Docker socket over a secure WebSocket tunnel.
func Run(ctx context.Context, serverURL, token string) error {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)
	headers.Add("X-Vessel-Agent", "true")

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

	// Wrap the websocket in a net.Conn
	netConn := &websocketConn{conn: conn}

	// Setup yamux server to accept multiplexed connections from the controller
	session, err := yamux.Server(netConn, yamux.DefaultConfig())
	if err != nil {
		return fmt.Errorf("failed to start yamux session: %w", err)
	}
	defer session.Close()

	log.Println(" Secure tunnel established. Listening for controller commands...")

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

	// Dial the local Docker socket
	dockerConn, err := net.Dial("unix", "/var/run/docker.sock")
	if err != nil {
		log.Printf(" Failed to connect to local Docker daemon: %v", err)
		return
	}
	defer dockerConn.Close()

	errc := make(chan error, 2)

	// Copy from stream to docker
	go func() {
		_, err := io.Copy(dockerConn, stream)
		errc <- err
	}()

	// Copy from docker to stream
	go func() {
		_, err := io.Copy(stream, dockerConn)
		errc <- err
	}()

	<-errc
}

// websocketConn adapts a gorilla/websocket.Conn to standard net.Conn
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
		return n, nil // Re-read next frame
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
