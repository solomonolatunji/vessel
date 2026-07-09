package handlers

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"vessel.dev/vessel/internal/middleware"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/utils"
)

var terminalUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type TerminalHandler struct {
	dockerClient  *client.Client
	tokenService  *services.TokenService
	appService    *services.AppService
	normalizeName func(id string) string
}

func NewTerminalHandler(
	dockerClient *client.Client,
	tokenService *services.TokenService,
	appService *services.AppService,
) *TerminalHandler {
	return &TerminalHandler{
		dockerClient:  dockerClient,
		tokenService:  tokenService,
		appService:    appService,
		normalizeName: utils.NormalizeContainerName,
	}
}

func (h *TerminalHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if h.tokenService != nil {
		tokenStr := middleware.ExtractTokenFromRequest(r)
		if tokenStr == "" {
			WriteError(w, http.StatusUnauthorized, "missing authentication token for terminal access")
			return
		}
		if _, err := h.tokenService.ValidateToken(tokenStr); err != nil {
			WriteError(w, http.StatusUnauthorized, "invalid authentication token for terminal access")
			return
		}
	}

	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	containerName := h.normalizeName(id)
	if h.appService != nil {
		if svc, err := h.appService.GetAppService(r.Context(), id); err == nil && svc != nil {
			if svc.ContainerID != "" && svc.ContainerID != "-" {
				containerName = svc.ContainerID
			} else {
				containerName = h.normalizeName(svc.ID)
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

	if h.dockerClient == nil {
		WriteError(w, http.StatusInternalServerError, "docker client unavailable")
		return
	}

	resp, err := h.dockerClient.ContainerExecCreate(context.Background(), containerName, execConfig)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to create exec instance: "+err.Error())
		return
	}

	hijackedResp, err := h.dockerClient.ContainerExecAttach(context.Background(), resp.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to attach to exec instance: "+err.Error())
		return
	}
	defer hijackedResp.Close()

	ws, err := terminalUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	errChan := make(chan error, 2)
	go func() {
		wsReader := h.wsToReader(ws)
		_, err := io.Copy(hijackedResp.Conn, wsReader)
		errChan <- err
	}()

	go func() {
		wsWriter := h.wsToWriter(ws)
		_, err := io.Copy(wsWriter, hijackedResp.Reader)
		errChan <- err
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				return
			}
		}
	}()

	<-errChan
}

func (h *TerminalHandler) wsToReader(ws *websocket.Conn) io.Reader {
	r, w := io.Pipe()
	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				w.CloseWithError(err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				return
			}
		}
	}()
	return r
}

func (h *TerminalHandler) wsToWriter(ws *websocket.Conn) io.Writer {
	r, w := io.Pipe()
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		}
	}()
	return w
}
