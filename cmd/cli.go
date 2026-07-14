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
		fmt.Printf("vessld %s %s/%s\n", vesslVersion, runtime.GOOS, runtime.GOARCH)
	default:
		if strings.Contains(os.Args[1], ":") {
			parts := strings.SplitN(os.Args[1], ":", 2)
			switch parts[0] {
			case "apps":
				runApps(os.Args[2:])
			case "db", "database":
				runDatabases(os.Args[2:])
			default:
				printHelp()
			}
			return
		}
		printHelp()
	}
}

func printHelp() {
	fmt.Printf("Usage: vessld [command]\n\n")
	fmt.Printf("Server commands:\n")
	fmt.Printf("  serve           Start the daemon (default)\n")
	fmt.Printf("  setup           Run interactive setup wizard\n")
	fmt.Printf("  reset-password  Reset admin password\n")
	fmt.Printf("  config          View or update configuration\n\n")
	fmt.Printf("Management commands:\n")
	fmt.Printf("  deploy <url>          Deploy an app from Git URL\n")
	fmt.Printf("  deploy --image <img>  Deploy an app from a Docker image\n")
	fmt.Printf("  deploy --image nginx:latest --port 80\n")
	fmt.Printf("  apps:list       List all apps\n")
	fmt.Printf("  apps:show <id>  Show app details\n")
	fmt.Printf("  apps:create     Create an app\n")
	fmt.Printf("  apps:destroy    Delete an app\n")
	fmt.Printf("  db:list         List all databases\n")
	fmt.Printf("  db:show <id>    Show database details\n")
	fmt.Printf("  db:create       Create a database\n")
	fmt.Printf("  db:destroy      Delete a database\n\n")
	fmt.Printf("Other:\n")
	fmt.Printf("  mcp             Run MCP stdio server\n")
	fmt.Printf("  version         Show version\n")
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
	fmt.Println("🔄 Restarting Vessl daemon...")
	if _, err := os.Stat("/vessl/docker-compose.yml"); err == nil {
		cmd := exec.Command("docker", "compose", "-f", "/vessl/docker-compose.yml", "restart", "vessl-control-plane")
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
