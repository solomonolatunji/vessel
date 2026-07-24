package utils

import (
	"context"
	"os"

	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func GetDataDir() string {
	dir := os.Getenv("CODEDOCK_DATA_DIR")
	if dir == "" {
		return "data"
	}
	return dir
}

func GetRuntimeNetwork() string {
	net := os.Getenv("CODEDOCK_RUNTIME_NETWORK")
	if net == "" {
		return "codedock-network"
	}
	return net
}

func EnsureCodedockNetwork(ctx context.Context, cli *client.Client) error {
	if cli == nil {
		return nil
	}

	netName := GetRuntimeNetwork()

	networks, err := cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return err
	}
	for _, net := range networks {
		if net.Name == netName {
			return nil
		}
	}
	_, err = cli.NetworkCreate(ctx, netName, network.CreateOptions{
		Driver: "bridge",
	})
	return err
}
