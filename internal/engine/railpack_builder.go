package engine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"vessl.dev/vessl/internal/utils"
)

type RailpackBuilder struct {
	dockerClient *client.Client
}

func NewRailpackBuilder(dockerClient *client.Client) *RailpackBuilder {
	return &RailpackBuilder{dockerClient: dockerClient}
}

func (r *RailpackBuilder) Build(ctx context.Context, opts BuildOptions, engineName string) (string, error) {
	imageTag := fmt.Sprintf("vessel-app-%s:latest", strings.ToLower(opts.ProjectID))
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "🌟 [Railpack/Nixpacks] Auto-detecting language & framework in %s...\n", opts.SourceDir)
	}
	stack := r.detectLanguageStack(opts.SourceDir)
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "🛠️ [Railpack/Nixpacks] Stack detected: %s\n", stack)
	}
	if engineName == "buildpacks" {
		packPath, err := exec.LookPath("pack")
		if err == nil {
			if opts.LogWriter != nil {
				fmt.Fprintf(opts.LogWriter, "⚙️ [Buildpacks] Executing pack builder engine (%s)...\n", packPath)
			}
			cmd := exec.CommandContext(ctx, packPath, "build", imageTag, "--path", opts.SourceDir, "--builder", "paketobuildpacks/builder:base")
			cmd.Stdout = opts.LogWriter
			cmd.Stderr = opts.LogWriter
			if err := cmd.Run(); err != nil {
				return "", utils.NewDeploymentError("buildpacks execution failed", err)
			}
			return imageTag, nil
		}
	} else {
		nixpacksPath, err := exec.LookPath("nixpacks")
		if err == nil {
			if opts.LogWriter != nil {
				fmt.Fprintf(opts.LogWriter, "⚙️ [Nixpacks] Executing Nixpacks builder engine (%s)...\n", nixpacksPath)
			}
			cmd := exec.CommandContext(ctx, nixpacksPath, "build", opts.SourceDir, "--name", imageTag)
			cmd.Stdout = opts.LogWriter
			cmd.Stderr = opts.LogWriter
			if err := cmd.Run(); err != nil {
				return "", utils.NewDeploymentError("nixpacks execution failed", err)
			}
			return imageTag, nil
		}
	}
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "⚠️ [Railpack/Nixpacks] Native CLI not found; synthesizing zero-config OCI build plan for %s...\n", stack)
	}
	synthesizedDockerfile, err := r.synthesizeDockerfile(stack, opts)
	if err != nil {
		return "", fmt.Errorf("failed to synthesize fallback dockerfile: %w", err)
	}
	dockerfilePath := filepath.Join(opts.SourceDir, ".vessel.Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(synthesizedDockerfile), 0o644); err != nil {
		return "", fmt.Errorf("failed to write synthesized dockerfile: %w", err)
	}
	defer os.Remove(dockerfilePath)
	fallbackOpts := opts
	fallbackOpts.DockerfilePath = ".vessel.Dockerfile"
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

func (r *RailpackBuilder) synthesizeDockerfile(stack string, opts BuildOptions) (string, error) {
	switch {
	case strings.HasPrefix(stack, "Node.js"):
		return `FROM node:22-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install --production
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
`, nil
	case strings.HasPrefix(stack, "Go"):
		return `FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server /app/server
EXPOSE 3000
CMD ["/app/server"]
`, nil
	case strings.HasPrefix(stack, "Python"):
		return `FROM python:3.12-slim
WORKDIR /app
COPY requirements.txt* pyproject.toml* ./
RUN if [ -f requirements.txt ]; then pip install --no-cache-dir -r requirements.txt; fi
COPY . .
EXPOSE 3000
CMD ["python3", "-m", "http.server", "3000"]
`, nil
	case strings.HasPrefix(stack, "Rust"):
		return `FROM rust:1.83-alpine AS builder
WORKDIR /app
COPY . .
RUN cargo build --release
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/target/release/* /app/server
EXPOSE 3000
CMD ["/app/server"]
`, nil
	default:
		return `FROM nginx:alpine
COPY . /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
`, nil
	}
}
