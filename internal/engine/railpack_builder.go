package engine

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
)

type RailpackBuilder struct {
	dockerClient *client.Client
}

func NewRailpackBuilder(dockerClient *client.Client) *RailpackBuilder {
	return &RailpackBuilder{dockerClient: dockerClient}
}

func (r *RailpackBuilder) Build(ctx context.Context, opts BuildOptions, engineName string) (string, error) {
	imageTag := fmt.Sprintf("vessl-app-%s:latest", strings.ToLower(opts.ProjectID))
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "🌟 [Railpack/Nixpacks] Auto-detecting language & framework in %s...\n", opts.SourceDir)
	}
	stack := r.detectLanguageStack(opts.SourceDir)
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "🛠️ [Railpack/Nixpacks] Stack detected: %s\n", stack)
	}
	absSourceDir, err := filepath.Abs(opts.SourceDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute source dir: %w", err)
	}

	if engineName == "buildpacks" {
		if opts.LogWriter != nil {
			fmt.Fprintf(opts.LogWriter, "⚙️ [Buildpacks] Executing pack builder engine via Docker container...\n")
		}
		cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
			"-v", "/var/run/docker.sock:/var/run/docker.sock",
			"-v", fmt.Sprintf("%s:/app", absSourceDir),
			"buildpacksio/pack:latest", "build", imageTag, "--path", "/app", "--builder", "paketobuildpacks/builder:base")
		cmd.Stdout = opts.LogWriter
		cmd.Stderr = opts.LogWriter
		if err := cmd.Run(); err == nil {
			return imageTag, nil
		} else {
			if opts.LogWriter != nil {
				fmt.Fprintf(opts.LogWriter, "⚠️ [Buildpacks] Container execution failed, falling back to synthesized Dockerfile. Error: %v\n", err)
			}
		}
	} else {
		if opts.LogWriter != nil {
			fmt.Fprintf(opts.LogWriter, "⚙️ [Nixpacks] Executing Nixpacks builder engine via Docker container...\n")
		}
		cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
			"-v", "/var/run/docker.sock:/var/run/docker.sock",
			"-v", fmt.Sprintf("%s:/app", absSourceDir),
			"ghcr.io/railwayapp/nixpacks:latest", "build", "/app", "--name", imageTag)
		cmd.Stdout = opts.LogWriter
		cmd.Stderr = opts.LogWriter
		if err := cmd.Run(); err == nil {
			return imageTag, nil
		} else {
			if opts.LogWriter != nil {
				fmt.Fprintf(opts.LogWriter, "⚠️ [Nixpacks] Container execution failed, falling back to synthesized Dockerfile. Error: %v\n", err)
			}
		}
	}
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "⚠️ [Railpack/Nixpacks] Native CLI not found; synthesizing zero-config OCI build plan for %s...\n", stack)
	}
	synthesizedDockerfile, err := r.synthesizeDockerfile(stack, opts)
	if err != nil {
		return "", fmt.Errorf("failed to synthesize fallback dockerfile: %w", err)
	}
	dockerfilePath := filepath.Join(opts.SourceDir, ".vessl.Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(synthesizedDockerfile), 0o644); err != nil {
		return "", fmt.Errorf("failed to write synthesized dockerfile: %w", err)
	}
	defer os.Remove(dockerfilePath)
	fallbackOpts := opts
	fallbackOpts.DockerfilePath = ".vessl.Dockerfile"
	fallbackBuilder := NewDockerfileBuilder(r.dockerClient)
	return fallbackBuilder.Build(ctx, fallbackOpts)
}

func (r *RailpackBuilder) detectLanguageStack(sourceDir string) string {
	if _, err := os.Stat(filepath.Join(sourceDir, "package.json")); err == nil {
		return "Node.js (npm/yarn/pnpm)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "go.mod")); err == nil {
		return "Go (Golang)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "requirements.txt")); err == nil {
		return "Python (pip)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "pyproject.toml")); err == nil {
		return "Python (poetry/pyproject)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "Cargo.toml")); err == nil {
		return "Rust (cargo)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "composer.json")); err == nil {
		return "PHP (composer)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "Gemfile")); err == nil {
		return "Ruby (bundler)"
	}
	return "Static / Universal HTML"
}

//go:embed templates/*.Dockerfile
var templateFS embed.FS

func (r *RailpackBuilder) synthesizeDockerfile(stack string, opts BuildOptions) (string, error) {
	var templateName string
	switch {
	case strings.HasPrefix(stack, "Node.js"):
		templateName = "templates/nodejs.Dockerfile"
	case strings.HasPrefix(stack, "Go"):
		templateName = "templates/go.Dockerfile"
	case strings.HasPrefix(stack, "Python"):
		templateName = "templates/python.Dockerfile"
	case strings.HasPrefix(stack, "Rust"):
		templateName = "templates/rust.Dockerfile"
	default:
		templateName = "templates/static.Dockerfile"
	}
	content, err := templateFS.ReadFile(templateName)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded dockerfile template %s: %w", templateName, err)
	}
	return string(content), nil
}
