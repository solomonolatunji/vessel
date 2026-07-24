package engine

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/docker/docker/client"
)

type RailpackBuilder struct {
	dockerClient *client.Client
}

func NewRailpackBuilder(dockerClient *client.Client) *RailpackBuilder {
	return &RailpackBuilder{dockerClient: dockerClient}
}

func (r *RailpackBuilder) Build(ctx context.Context, opts BuildOptions, engineName string) (string, error) {
	imageTag := fmt.Sprintf("codedock-app-%s:latest", strings.ToLower(opts.ProjectID))
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

	dockerSock := os.Getenv("DOCKER_SOCKET_PATH")
	if dockerSock == "" {
		dockerSock = "/var/run/docker.sock"
	}

	if engineName == "buildpacks" {
		if opts.LogWriter != nil {
			fmt.Fprintf(opts.LogWriter, "⚙️ [Buildpacks] Executing pack builder engine via Docker container...\n")
		}
		packImage := os.Getenv("CODEDOCK_PACK_IMAGE")
		if packImage == "" {
			packImage = "buildpacksio/pack:latest"
		}
		builderImage := os.Getenv("CODEDOCK_BUILDER_IMAGE")
		if builderImage == "" {
			builderImage = "paketobuildpacks/builder:base"
		}
		args := []string{"run", "--rm",
			"-v", fmt.Sprintf("%s:/var/run/docker.sock", dockerSock),
			"-v", fmt.Sprintf("%s:/app", absSourceDir),
			packImage, "build", imageTag, "--path", "/app", "--builder", builderImage}
		for k, v := range opts.EnvVars {
			args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
		}
		cmd := exec.CommandContext(ctx, "docker", args...)
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
		nixpacksImage := os.Getenv("CODEDOCK_NIXPACKS_IMAGE")
		if nixpacksImage == "" {
			nixpacksImage = "ghcr.io/railwayapp/nixpacks:latest"
		}
		args := []string{"run", "--rm",
			"-v", fmt.Sprintf("%s:/var/run/docker.sock", dockerSock),
			"-v", fmt.Sprintf("%s:/app", absSourceDir),
			nixpacksImage, "build", "/app", "--name", imageTag}
		for k, v := range opts.EnvVars {
			args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
		}
		if opts.AppConfig != nil {
			if opts.AppConfig.InstallCommand != "" {
				args = append(args, "--install-cmd", opts.AppConfig.InstallCommand)
			}
			if opts.AppConfig.BuildCommand != "" {
				args = append(args, "--build-cmd", opts.AppConfig.BuildCommand)
			}
			if opts.AppConfig.StartCommand != "" {
				args = append(args, "--start-cmd", opts.AppConfig.StartCommand)
			}
		}
		cmd := exec.CommandContext(ctx, "docker", args...)
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
	dockerfilePath := filepath.Join(opts.SourceDir, ".codedock.Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(synthesizedDockerfile), 0o644); err != nil {
		return "", fmt.Errorf("failed to write synthesized dockerfile: %w", err)
	}
	defer os.Remove(dockerfilePath)
	fallbackOpts := opts
	fallbackOpts.DockerfilePath = ".codedock.Dockerfile"
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

type dockerfileOverrides struct {
	InstallCommand string
	BuildCommand   string
	StartCommand   string
}

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

	if opts.AppConfig == nil {
		return string(content), nil
	}

	overrides := dockerfileOverrides{
		InstallCommand: opts.AppConfig.InstallCommand,
		BuildCommand:   opts.AppConfig.BuildCommand,
		StartCommand:   opts.AppConfig.StartCommand,
	}

	if overrides.InstallCommand == "" && overrides.BuildCommand == "" && overrides.StartCommand == "" {
		return string(content), nil
	}

	result := applyDockerfileOverrides(string(content), overrides)
	return result, nil
}

func applyDockerfileOverrides(dockerfile string, o dockerfileOverrides) string {
	lines := strings.Split(dockerfile, "\n")
	tmplFuncs := template.FuncMap{}
	_ = tmplFuncs

	var out []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if o.InstallCommand != "" && isInstallLine(trimmed) {
			out = append(out, "RUN "+o.InstallCommand)
			continue
		}
		if o.BuildCommand != "" && isBuildLine(trimmed) {
			out = append(out, "RUN "+o.BuildCommand)
			continue
		}
		if o.StartCommand != "" && strings.HasPrefix(trimmed, "CMD ") {
			out = append(out, buildCMDLine(o.StartCommand))
			continue
		}
		out = append(out, line)
	}

	var buf bytes.Buffer
	for _, l := range out {
		buf.WriteString(l)
		buf.WriteByte('\n')
	}
	return buf.String()
}

func isInstallLine(line string) bool {
	installPatterns := []string{
		"RUN npm install", "RUN npm ci", "RUN yarn install", "RUN yarn",
		"RUN pnpm install", "RUN pip install", "RUN pip3 install",
		"RUN bundle install", "RUN composer install", "RUN cargo fetch",
	}
	for _, p := range installPatterns {
		if strings.HasPrefix(line, p) {
			return true
		}
	}
	return false
}

func isBuildLine(line string) bool {
	buildPatterns := []string{
		"RUN npm run build", "RUN yarn build", "RUN pnpm build",
		"RUN go build", "RUN cargo build",
	}
	for _, p := range buildPatterns {
		if strings.HasPrefix(line, p) {
			return true
		}
	}
	return false
}

func buildCMDLine(startCmd string) string {
	parts := strings.Fields(startCmd)
	if len(parts) == 0 {
		return ""
	}
	quoted := make([]string, len(parts))
	for i, p := range parts {
		quoted[i] = `"` + p + `"`
	}
	return "CMD [" + strings.Join(quoted, ", ") + "]"
}
