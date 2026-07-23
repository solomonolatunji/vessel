package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"codedock.dev/codedock/internal/models"
)

type SettingsRepository interface {
	GetServerSettings(ctx context.Context) (*models.ServerSettings, error)
	UpdateServerSettings(ctx context.Context, cfg *models.ServerSettings) error
	ListProjects(ctx context.Context) ([]map[string]any, error)
}

type SettingsRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewSettingsRepo(db *sql.DB) *SettingsRepo {
	return &SettingsRepo{db: sqlx.NewDb(db, "sqlite")}
}

const serverSettingsColumns = `id, traefik_wildcard_ip, registration_enabled, registration_domain_allowlist, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, default_wildcard_domain, panel_domain, site_name, public_ipv4, public_ipv6, show_sponsorship_popup, disable_two_step_confirmation, cloudflare_api_token, namecheap_api_user, namecheap_api_key, namecheap_client_ip, spaceship_api_key, update_check_cron, auto_update_enabled, concurrent_builds, deployment_timeout, server_timezone, docker_cleanup_cron, disk_usage_threshold, disk_usage_cron, current_version, latest_version, last_update_check, updated_at`

func serverSettingsPlaceholders() string {
	columns := strings.Split(serverSettingsColumns, ",")
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}

func scanServerSettings(scanner interface{ Scan(dest ...any) error }, cfg *models.ServerSettings) error {
	return scanner.Scan(
		&cfg.ID, &cfg.TraefikWildcardIP,
		&cfg.RegistrationEnabled, &cfg.RegistrationDomainAllowlist, &cfg.CustomDNSResolvers, &cfg.DNSValidationEnabled, &cfg.IPAllowlist, &cfg.MCPServerEnabled, &cfg.DefaultWildcardDomain, &cfg.PanelDomain,
		&cfg.SiteName, &cfg.PublicIPv4, &cfg.PublicIPv6, &cfg.ShowSponsorshipPopup, &cfg.DisableTwoStepConfirmation,
		&cfg.CloudflareAPIToken, &cfg.NamecheapAPIUser, &cfg.NamecheapAPIKey, &cfg.NamecheapClientIP, &cfg.SpaceshipAPIKey,
		&cfg.UpdateCheckCron, &cfg.AutoUpdateEnabled,
		&cfg.ConcurrentBuilds, &cfg.DeploymentTimeout, &cfg.ServerTimezone, &cfg.DockerCleanupCron, &cfg.DiskUsageThreshold, &cfg.DiskUsageCron,
		&cfg.CurrentVersion, &cfg.LatestVersion, &cfg.LastUpdateCheck, &cfg.UpdatedAt,
	)
}

func serverSettingsArgs(cfg *models.ServerSettings) []any {
	return []any{
		cfg.ID, cfg.TraefikWildcardIP,
		cfg.RegistrationEnabled, cfg.RegistrationDomainAllowlist, cfg.CustomDNSResolvers, cfg.DNSValidationEnabled, cfg.IPAllowlist, cfg.MCPServerEnabled, cfg.DefaultWildcardDomain, cfg.PanelDomain,
		cfg.SiteName, cfg.PublicIPv4, cfg.PublicIPv6, cfg.ShowSponsorshipPopup, cfg.DisableTwoStepConfirmation,
		cfg.CloudflareAPIToken, cfg.NamecheapAPIUser, cfg.NamecheapAPIKey, cfg.NamecheapClientIP, cfg.SpaceshipAPIKey,
		cfg.UpdateCheckCron, cfg.AutoUpdateEnabled, cfg.ConcurrentBuilds, cfg.DeploymentTimeout, cfg.ServerTimezone, cfg.DockerCleanupCron, cfg.DiskUsageThreshold, cfg.DiskUsageCron, cfg.CurrentVersion, cfg.LatestVersion, cfg.LastUpdateCheck, cfg.UpdatedAt,
	}
}

func (r *SettingsRepo) GetServerSettings(ctx context.Context) (*models.ServerSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg models.ServerSettings
	err := scanServerSettings(r.db.QueryRowContext(ctx, `SELECT `+serverSettingsColumns+` FROM server_settings WHERE id = 'global' LIMIT 1`), &cfg)
	if errors.Is(err, sql.ErrNoRows) {
		defaultSettings := &models.ServerSettings{
			ID:                   "global",
			RegistrationEnabled:  true,
			DNSValidationEnabled: true,
			CustomDNSResolvers:   "1.1.1.1",
			MCPServerEnabled:     true,
			UpdateCheckCron:      "0 * * * *",
			AutoUpdateEnabled:    false,
			CurrentVersion:       "0.1.0",
			LatestVersion:        "0.1.0",
			UpdatedAt:            time.Now().UTC().Format(time.RFC3339),
			ShowSponsorshipPopup: true,
		}
		query := fmt.Sprintf(`INSERT INTO server_settings (%s) VALUES (%s)`, serverSettingsColumns, serverSettingsPlaceholders())
		_, _ = r.db.ExecContext(ctx, query, serverSettingsArgs(defaultSettings)...)
		return defaultSettings, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get server settings: %w", err)
	}
	return &cfg, nil
}

func (r *SettingsRepo) UpdateServerSettings(ctx context.Context, cfg *models.ServerSettings) error {
	if cfg.ID == "" {
		cfg.ID = "global"
	}
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.mu.Lock()
	defer r.mu.Unlock()
	query := fmt.Sprintf(`INSERT INTO server_settings (%s)
	          VALUES (%s)
	          ON CONFLICT(id) DO UPDATE SET
	          traefik_wildcard_ip = excluded.traefik_wildcard_ip,
	          registration_enabled = excluded.registration_enabled,
	          registration_domain_allowlist = excluded.registration_domain_allowlist,
	          custom_dns_resolvers = excluded.custom_dns_resolvers,
	          dns_validation_enabled = excluded.dns_validation_enabled,
	          ip_allowlist = excluded.ip_allowlist,
	          mcp_server_enabled = excluded.mcp_server_enabled,
	          default_wildcard_domain = excluded.default_wildcard_domain,
	          panel_domain = excluded.panel_domain,
	          site_name = excluded.site_name,
	          public_ipv4 = excluded.public_ipv4,
	          public_ipv6 = excluded.public_ipv6,
	          show_sponsorship_popup = excluded.show_sponsorship_popup,
	          disable_two_step_confirmation = excluded.disable_two_step_confirmation,
	          cloudflare_api_token = excluded.cloudflare_api_token,
	          namecheap_api_user = excluded.namecheap_api_user,
	          namecheap_api_key = excluded.namecheap_api_key,
	          namecheap_client_ip = excluded.namecheap_client_ip,
	          spaceship_api_key = excluded.spaceship_api_key,
	          update_check_cron = excluded.update_check_cron,
	          auto_update_enabled = excluded.auto_update_enabled,
	          concurrent_builds = excluded.concurrent_builds,
	          deployment_timeout = excluded.deployment_timeout,
	          server_timezone = excluded.server_timezone,
	          docker_cleanup_cron = excluded.docker_cleanup_cron,
	          disk_usage_threshold = excluded.disk_usage_threshold,
	          disk_usage_cron = excluded.disk_usage_cron,
	          current_version = excluded.current_version,
	          latest_version = excluded.latest_version,
	          last_update_check = excluded.last_update_check,
	          updated_at = excluded.updated_at`, serverSettingsColumns, serverSettingsPlaceholders())
	_, err := r.db.ExecContext(ctx, query, serverSettingsArgs(cfg)...)
	if err != nil {
		return fmt.Errorf("failed to update server settings: %w", err)
	}
	return nil
}

func (r *SettingsRepo) ListProjects(ctx context.Context) ([]map[string]any, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, repo_url, branch, status, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()
	var projects []map[string]any
	for rows.Next() {
		var id, name, repoURL, branch, status, updatedAt string
		if err := rows.Scan(&id, &name, &repoURL, &branch, &status, &updatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, map[string]any{
			"id":        id,
			"name":      name,
			"repoUrl":   repoURL,
			"branch":    branch,
			"status":    status,
			"updatedAt": updatedAt,
		})
	}
	return projects, nil
}
