package api

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/gorilla/websocket"
	"github.com/solomonolatunji/vessel/internal/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) handleTerminalWebSocket(w http.ResponseWriter, r *http.Request) {
	if s.tokenService != nil {
		tokenStr := extractTokenFromRequest(r)
		if tokenStr == "" {
			writeError(w, http.StatusUnauthorized, "missing authentication token for terminal access")
			return
		}
		if _, err := s.tokenService.ValidateToken(tokenStr); err != nil {
			writeError(w, http.StatusUnauthorized, "invalid authentication token for terminal access")
			return
		}
	}

	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	containerName := utils.NormalizeContainerName(id)
	if s.store != nil {
		if appService, err := s.store.GetAppService(id); err == nil && appService != nil {
			if appService.ContainerID != "" && appService.ContainerID != "-" {
				containerName = appService.ContainerID
			} else {
				containerName = utils.NormalizeContainerName(appService.ID)
			}
		}
	}
	execConfig := types.ExecConfig{
		Cmd:          []string{"/bin/sh"},
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}

	if s.dockerClient == nil {
		writeError(w, http.StatusServiceUnavailable, "docker client unavailable")
		return
	}

	ctx := r.Context()
	execID, err := s.dockerClient.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		execConfig.Cmd = []string{"/bin/sh"}
		execID, err = s.dockerClient.ContainerExecCreate(ctx, containerName, execConfig)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create interactive container session: "+err.Error())
			return
		}
	}

	execResp, err := s.dockerClient.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to attach to container shell: "+err.Error())
		return
	}
	defer execResp.Close()

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer wsConn.Close()

	_, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer cancel()
		buf := make([]byte, 4096)
		for {
			n, err := execResp.Reader.Read(buf)
			if err != nil {
				if err != io.EOF {
					_ = wsConn.WriteMessage(websocket.TextMessage, []byte("\r\n[session terminated]\r\n"))
				}
				break
			}
			if n > 0 {
				_ = wsConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := wsConn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					break
				}
			}
		}
	}()

	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
		if len(message) > 0 {
			if _, err := execResp.Conn.Write(message); err != nil {
				break
			}
		}
	}
}
