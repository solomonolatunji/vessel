package orchestrator

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DockerfileBuilder compiles OCI images from repositories that contain an explicit Dockerfile.
type DockerfileBuilder struct {
	dockerClient *client.Client
}

// NewDockerfileBuilder creates a new DockerfileBuilder using the provided Docker daemon client.
func NewDockerfileBuilder(dockerClient *client.Client) *DockerfileBuilder {
	return &DockerfileBuilder{dockerClient: dockerClient}
}

// Build archives the source context and invokes the Docker daemon ImageBuild API.
func (d *DockerfileBuilder) Build(ctx context.Context, opts BuildOptions) (string, error) {
	imageTag := fmt.Sprintf("vessel-app-%s:latest", strings.ToLower(opts.ProjectID))
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "📦 [Dockerfile] Packaging build context from %s...\n", opts.SourceDir)
	}

	tarContext, err := d.createTarContext(opts.SourceDir)
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

func (d *DockerfileBuilder) createTarContext(sourceDir string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, file); err != nil {
				file.Close()
				return err
			}
			file.Close()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}
