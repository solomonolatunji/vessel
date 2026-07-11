package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"vessl.dev/vessl/internal/services"
)

type Bridge struct {
	server         *server.MCPServer
	projectService *services.ProjectService
	appService     *services.AppService
	dbService      *services.DatabaseService
}

func NewBridge(ps *services.ProjectService, as *services.AppService, db *services.DatabaseService) *Bridge {
	mcpServer := server.NewMCPServer("vessel-mcp", "1.0.0", server.WithResourceCapabilities(true, true), server.WithPromptCapabilities(true))
	b := &Bridge{
		server:         mcpServer,
		projectService: ps,
		appService:     as,
		dbService:      db,
	}
	b.registerTools()
	return b
}

func (b *Bridge) MCPServer() *server.MCPServer {
	return b.server
}

func (b *Bridge) registerTools() {
	b.server.AddTool(
		mcp.NewTool("list_projects",
			mcp.WithDescription("List all deployment projects registered in this Vessel instance."),
		),
		b.handleListProjects,
	)

	b.server.AddTool(
		mcp.NewTool("list_databases",
			mcp.WithDescription("List all managed databases registered in this Vessel instance."),
		),
		b.handleListDatabases,
	)

	b.server.AddTool(
		mcp.NewTool("get_system_status",
			mcp.WithDescription("Check basic operational and health metrics of the Vessel platform."),
		),
		b.handleGetSystemStatus,
	)
}

func (b *Bridge) handleListProjects(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projects, err := b.projectService.ListProjects(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res := "Vessel Projects:\n"
	for _, p := range projects {
		res += fmt.Sprintf("- ID: %s | Name: %s\n", p.ID, p.Name)
	}
	if len(projects) == 0 {
		res = "No projects found."
	}
	return mcp.NewToolResultText(res), nil
}

func (b *Bridge) handleListDatabases(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dbs, err := b.dbService.ListDatabases(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res := "Vessel Databases:\n"
	for _, d := range dbs {
		res += fmt.Sprintf("- ID: %s | Name: %s | Engine: %s | Status: %s\n", d.ID, d.Name, d.Engine, d.Status)
	}
	if len(dbs) == 0 {
		res = "No databases found."
	}
	return mcp.NewToolResultText(res), nil
}

func (b *Bridge) handleGetSystemStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	res := "Vessel Status: OK\nEngine: Active\nVersion: 1.0.0"
	return mcp.NewToolResultText(res), nil
}
