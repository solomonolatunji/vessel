package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func mainCLI() {
	if len(os.Args) < 2 || os.Args[1] == "serve" {
		startServer()
		return
	}

	switch os.Args[1] {
	case "serve":
		startServer()
	case "setup":
		runSetup()
	case "reset-password":
		runResetPassword()
	case "config":
		runConfig()
	case "deploy":
		runDeploy(os.Args[2:])
	case "restart":
		runRestart()
	case "mcp":
		runMCP()
	case "version", "--version", "-v":
		fmt.Printf("codedockd %s %s/%s\n", codedockVersion, runtime.GOOS, runtime.GOARCH)
	default:
		if strings.Contains(os.Args[1], ":") {
			parts := strings.SplitN(os.Args[1], ":", 2)
			switch parts[0] {
			case "apps", "app":
				runApps(append([]string{parts[1]}, os.Args[2:]...))
			case "db", "database":
				runDatabases(append([]string{parts[1]}, os.Args[2:]...))
			case "project":
				runProjects(append([]string{parts[1]}, os.Args[2:]...))
			case "env", "vars":
				runEnvVars(append([]string{parts[1]}, os.Args[2:]...))
			case "deployment":
				runDeployments(append([]string{parts[1]}, os.Args[2:]...))
			case "domain":
				runDomains(append([]string{parts[1]}, os.Args[2:]...))
			default:
				printHelp()
			}
			return
		}
		printHelp()
	}
}

func printHelp() {
	fmt.Printf("Usage: codedockd [command]\n\n")
	fmt.Printf("Server commands:\n")
	fmt.Printf("  serve              Start the daemon (default)\n")
	fmt.Printf("  setup              Run interactive setup wizard\n")
	fmt.Printf("  reset-password     Reset admin password\n")
	fmt.Printf("  deploy <git-url>           Deploy a new service from Git\n")
	fmt.Printf("  deploy --template <t>      Deploy from a template (e.g. redis, postgres)\n")
	fmt.Printf("  deploy --image <img>        Deploy from a Docker image\n")
	fmt.Printf("  deploy --archive <file>     Deploy from a tar archive\n\n")
	fmt.Printf("Project management (no auth required):\n")
	fmt.Printf("  project:list                       List all projects\n")
	fmt.Printf("  project:show <id>                  Show project details\n")
	fmt.Printf("  project:create <name>              Create a project\n")
	fmt.Printf("  project:destroy <id>               Delete a project\n\n")
	fmt.Printf("App management:\n")
	fmt.Printf("  apps:list                          List all apps\n")
	fmt.Printf("  apps:show <id>                     Show app details\n")
	fmt.Printf("  apps:create <name> --project <id>  Create an app\n")
	fmt.Printf("  apps:destroy <id>                  Delete an app\n\n")
	fmt.Printf("Database management:\n")
	fmt.Printf("  db:list                            List all databases\n")
	fmt.Printf("  db:show <id>                       Show database details\n")
	fmt.Printf("  db:create <name> <engine>          Create a database\n")
	fmt.Printf("  db:destroy <id>                    Delete a database\n\n")
	fmt.Printf("Environment variables:\n")
	fmt.Printf("  env:list --project <id>            List env vars\n")
	fmt.Printf("  env:set KEY=VALUE --project <id>   Set env var(s)\n")
	fmt.Printf("  env:unset KEY --project <id>       Remove an env var\n\n")
	fmt.Printf("Deployments & logs:\n")
	fmt.Printf("  deployment:list --service <id>     List deployments\n")
	fmt.Printf("  deployment:show <id>               Show deployment\n")
	fmt.Printf("  deployment:logs <id>               View build logs\n\n")
	fmt.Printf("Custom domains:\n")
	fmt.Printf("  domain:list --project <id>         List domains\n")
	fmt.Printf("  domain:add <host> --project <id>   Add a domain\n")
	fmt.Printf("  domain:remove <id>                 Remove a domain\n\n")
	fmt.Printf("Other:\n")
	fmt.Printf("  mcp                                Run MCP stdio server\n")
	fmt.Printf("  version                            Show version\n")
	os.Exit(1)
}

func prompt(msg string) string {
	fmt.Print(msg)
	var input string
	fmt.Scanln(&input)
	return input
}

func promptOptional(msg string) string {
	fmt.Print(msg)
	var input string
	fmt.Scanln(&input)
	return input
}

func runRestart() {
	fmt.Println("🔄 Restarting Codedock daemon...")
	if _, err := os.Stat("/codedock/docker-compose.yml"); err == nil {
		cmd := exec.Command("docker", "compose", "-f", "/codedock/docker-compose.yml", "restart", "codedock-control-plane")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			exitError("Restart failed: %v", err)
		}
	} else {
		fmt.Println("   Standalone mode — exiting. Use systemd/supervisor to restart.")
		os.Exit(0)
	}
}

func parseUint(s string) (int, error) {
	var v int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number: %s", s)
		}
		v = v*10 + int(c-'0')
	}
	return v, nil
}

func promptPassword(msg string) string {
	fmt.Print(msg)
	var input string
	fmt.Scanln(&input)
	fmt.Println()
	return input
}
