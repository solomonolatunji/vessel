package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

func runDeployments(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vessld deployment:<command> [args]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  list --service <id>     List deployments for a service")
		fmt.Println("  show <id>               Show deployment details")
		fmt.Println("  logs <id>               View build logs for a deployment")
		return
	}

	_, db, _ := initDataDir()
	defer db.Close()

	deployRepo := repositories.NewDeploymentRepo(db)

	serviceID := ""
	for i := 1; i < len(args); i++ {
		if args[i] == "--service" && i+1 < len(args) {
			serviceID = args[i+1]
		}
	}

	cmd := args[0]

	switch cmd {
	case "list":
		if serviceID == "" {
			exitError("--service <id> is required")
		}
		deps, total, err := deployRepo.ListByService(context.Background(), serviceID, 20, 0)
		if err != nil {
			exitError("Failed to list deployments: %v", err)
		}
		fmt.Printf("  %d deployments for service %s:\n", total, serviceID[:8])
		for _, d := range deps {
			sha := d.CommitHash
			if len(sha) > 8 {
				sha = sha[:8]
			}
			fmt.Printf("  %s  %-10s  %s  %s  %s\n", d.ID[:8], d.Status, d.CreatedAt.Format("2006-01-02 15:04"), d.Branch, sha)
		}

	case "show":
		if len(args) < 2 {
			exitError("Usage: vessld deployment:show <id>")
		}
		d, err := deployRepo.GetByID(context.Background(), args[1])
		if err != nil {
			exitError("Deployment not found: %v", err)
		}
		fmt.Printf("  ID:         %s\n", d.ID)
		fmt.Printf("  Service:    %s\n", d.ServiceID)
		fmt.Printf("  Status:     %s\n", d.Status)
		fmt.Printf("  Branch:     %s\n", d.Branch)
		fmt.Printf("  Commit:     %s\n", d.CommitHash)
		fmt.Printf("  Message:    %s\n", d.CommitMessage)
		fmt.Printf("  Created:    %s\n", d.CreatedAt.Format("2006-01-02 15:04:05"))
		if d.BuildLogs != "" {
			fmt.Printf("  Build Logs: (use 'logs %s' to view)\n", d.ID[:8])
		}

	case "logs":
		if len(args) < 2 {
			exitError("Usage: vessld deployment:logs <deployment-id>")
		}
		d, err := deployRepo.GetByID(context.Background(), args[1])
		if err != nil {
			exitError("Deployment not found: %v", err)
		}
		if d.BuildLogs == "" {
			fmt.Println("  No build logs available for this deployment.")
			return
		}
		lines := strings.Split(d.BuildLogs, "\n")
		for _, line := range lines {
			fmt.Println(line)
		}

	default:
		fmt.Printf("Unknown deployment command: %s\n", cmd)
		fmt.Println("Try: list, show <id>, logs <id>")
	}
}

func runDomains(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vessld domain:<command> [args]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  list --project <id>           List all domains for a project")
		fmt.Println("  add <domain> --project <id>   Add a custom domain")
		fmt.Println("  remove <id>                   Remove a domain")
		return
	}

	_, db, _ := initDataDir()
	defer db.Close()

	domainRepo := repositories.NewDomainRepo(db)

	cmd := args[0]

	switch cmd {
	case "list":
		serviceID := ""
		for i := 1; i < len(args); i++ {
			if args[i] == "--service" && i+1 < len(args) {
				serviceID = args[i+1]
			}
		}
		if serviceID == "" {
			exitError("--service <id> is required")
		}
		domains, err := domainRepo.ListByService(context.Background(), serviceID)
		if err != nil {
			exitError("Failed to list domains: %v", err)
		}
		if len(domains) == 0 {
			fmt.Println("  No custom domains configured.")
			return
		}
		for _, d := range domains {
			fmt.Printf("  %s  %s\n", d.ID[:8], d.DomainName)
		}

	case "add":
		if len(args) < 2 {
			exitError("Usage: vessld domain:add <domain> --service <id>")
		}
		domain := args[1]
		serviceID := ""
		for i := 2; i < len(args); i++ {
			if args[i] == "--service" && i+1 < len(args) {
				serviceID = args[i+1]
			}
		}
		if serviceID == "" {
			exitError("--service <id> is required")
		}
		cfg := &models.DomainConfig{
			ID:         uuid.New().String(),
			ServiceID:  serviceID,
			DomainName: domain,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := domainRepo.Create(context.Background(), cfg); err != nil {
			exitError("Failed to add domain: %v", err)
		}
		fmt.Printf("✅ Domain added: %s (%s)\n", cfg.DomainName, cfg.ID[:8])

	case "remove":
		if len(args) < 2 {
			exitError("Usage: vessld domain:remove <id>")
		}
		if err := domainRepo.Delete(context.Background(), args[1]); err != nil {
			exitError("Failed to remove domain: %v", err)
		}
		fmt.Printf("✅ Domain %s removed.\n", args[1])

	default:
		fmt.Printf("Unknown domain command: %s\n", cmd)
		fmt.Println("Try: list, add <domain>, remove <id>")
	}
}
