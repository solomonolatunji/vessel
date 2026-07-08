package utils

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// EnsureVesselNetwork verifies that the isolated Docker bridge network vessel-net exists, creating it if missing.
func EnsureVesselNetwork(ctx context.Context, cli *client.Client) error {
	if cli == nil {
		return nil
	}

	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return err
	}

	for _, net := range networks {
		if net.Name == "vessel-net" {
			return nil
		}
	}

	_, err = cli.NetworkCreate(ctx, "vessel-net", types.NetworkCreate{
		Driver: "bridge",
	})
	return err
}
