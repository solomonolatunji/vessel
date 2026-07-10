package utils

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func GetRuntimeNetwork() string {
	net := os.Getenv("VESSEL_RUNTIME_NETWORK")
	if net == "" {
		return "vessel-network"
	}
	return net
}

func EnsureVesselNetwork(ctx context.Context, cli *client.Client) error {
	if cli == nil {
		return nil
	}
	
	netName := GetRuntimeNetwork()
	
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return err
	}
	for _, net := range networks {
		if net.Name == netName {
			return nil
		}
	}
	_, err = cli.NetworkCreate(ctx, netName, types.NetworkCreate{
		Driver: "bridge",
	})
	return err
}
