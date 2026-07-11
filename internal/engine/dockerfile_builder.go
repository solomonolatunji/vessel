package engine

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"vessl.dev/vessl/internal/utils"
)

type DockerfileBuilder struct {
	dockerClient *client.Client
}

func NewDockerfileBuilder(dockerClient *client.Client) *DockerfileBuilder {
	return &DockerfileBuilder{dockerClient: dockerClient}
}

func (d *DockerfileBuilder) Build(ctx context.Context, opts BuildOptions) (string, error) {
	imageTag := fmt.Sprintf("vessel-app-%s:latest", strings.ToLower(opts.ProjectID))
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "📦 [Dockerfile] Packaging build context from %s...\n", opts.SourceDir)
	}
	tarContext, err := utils.CreateTarContext(opts.SourceDir)
	if err != nil {
		return "", fmt.Errorf("failed to create tar context: %w", err)
	}
	dockerfilePath := opts.DockerfilePath
	if dockerfilePath == "" {
		dockerfilePath = "Dockerfile"
	}
	buildOptions := types.ImageBuildOptions{
		Tags:       []string{imageTag},
		Dockerfile: dockerfilePath,
		Remove:     true,
		CacheFrom:  []string{imageTag},
	}
	resp, err := d.dockerClient.ImageBuild(ctx, tarContext, buildOptions)
	if err != nil {
		return "", fmt.Errorf("docker image build request failed: %w", err)
	}
	defer resp.Body.Close()
	if opts.LogWriter != nil {
		_, _ = io.Copy(opts.LogWriter, resp.Body)
	} else {
		_, _ = io.Copy(io.Discard, resp.Body)
	}
	return imageTag, nil
}
