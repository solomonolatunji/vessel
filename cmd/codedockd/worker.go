package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"codedock.run/codedock/internal/worker"
)

func runWorker() {
	token := os.Getenv("CODEDOCK_WORKER_TOKEN")
	if token == "" {
		fmt.Println("Error: CODEDOCK_WORKER_TOKEN environment variable is required")
		os.Exit(1)
	}
	serverURL := os.Getenv("CODEDOCK_SERVER_URL")
	if serverURL == "" {
		fmt.Println("Error: CODEDOCK_SERVER_URL environment variable is required (e.g. wss://control-plane.example.com)")
		os.Exit(1)
	}

	daemon := worker.NewWorkerDaemon(serverURL, token)
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := daemon.Start(ctx); err != nil {
			fmt.Printf("Worker daemon failed: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("Worker daemon started successfully. Listening for tasks...")
	
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	
	fmt.Println("Shutting down worker daemon...")
	cancel()
}
