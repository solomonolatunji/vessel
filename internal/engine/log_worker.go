package engine

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type LogWorker struct {
	dockerClient *client.Client
	lokiURL      string
	httpClient   *http.Client
	mu           sync.Mutex
	activeTails  map[string]context.CancelFunc
}

func NewLogWorker(cli *client.Client) *LogWorker {
	return &LogWorker{
		dockerClient: cli,
		lokiURL:      "http://127.0.0.1:3100/loki/api/v1/push",
		httpClient:   &http.Client{Timeout: 5 * time.Second},
		activeTails:  make(map[string]context.CancelFunc),
	}
}

func (w *LogWorker) Start(ctx context.Context) {
	containers, err := w.dockerClient.ContainerList(ctx, container.ListOptions{})
	if err == nil {
		for _, c := range containers {
			w.startTailing(c.ID, c.Labels["vessl.service.id"])
		}
	} else {
		slog.Error("log worker failed to list initial containers", "err", err)
	}

	go w.listenForEvents(ctx)
}

func (w *LogWorker) listenForEvents(ctx context.Context) {
	msgs, errs := w.dockerClient.Events(ctx, events.ListOptions{
		Filters: filters.NewArgs(filters.Arg("type", "container"), filters.Arg("event", "start")),
	})

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errs:
			if err != nil {
				slog.Error("docker events stream error", "err", err)
				time.Sleep(2 * time.Second)
				msgs, errs = w.dockerClient.Events(ctx, events.ListOptions{
					Filters: filters.NewArgs(filters.Arg("type", "container"), filters.Arg("event", "start")),
				})
			}
		case msg := <-msgs:
			if msg.Action == "start" {
				serviceID := msg.Actor.Attributes["vessl.service.id"]
				w.startTailing(msg.Actor.ID, serviceID)
			}
		}
	}
}

func (w *LogWorker) startTailing(containerID, serviceID string) {
	if serviceID == "" {
		serviceID = "unknown"
	}

	w.mu.Lock()
	if _, exists := w.activeTails[containerID]; exists {
		w.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	w.activeTails[containerID] = cancel
	w.mu.Unlock()

	go w.tailContainerLogs(ctx, containerID, serviceID)
}

func (w *LogWorker) tailContainerLogs(ctx context.Context, containerID, serviceID string) {
	defer func() {
		w.mu.Lock()
		delete(w.activeTails, containerID)
		w.mu.Unlock()
	}()

	reader, err := w.dockerClient.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "0", // Only new logs
	})
	if err != nil {
		slog.Warn("failed to start container log stream", "container", containerID, "err", err)
		return
	}
	defer reader.Close()

	pipeReader, pipeWriter := io.Pipe()
	go func() {
		defer pipeWriter.Close()
		_, _ = stdcopy.StdCopy(pipeWriter, pipeWriter, reader)
	}()

	scanner := bufio.NewScanner(pipeReader)
	for scanner.Scan() {
		line := scanner.Text()
		w.pushToLoki(containerID, serviceID, line)
	}
}

func (w *LogWorker) pushToLoki(containerID, serviceID, line string) {
	payload := map[string]any{
		"streams": []map[string]any{
			{
				"stream": map[string]string{
					"container_id": containerID,
					"service_id":   serviceID,
				},
				"values": [][]string{
					{
						fmt.Sprintf("%d", time.Now().UnixNano()),
						strings.TrimSpace(line),
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, w.lokiURL, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
