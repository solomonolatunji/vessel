package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type DNSProviderService struct {
	settingsRepo repositories.SettingsRepository
}

func NewDNSProviderService(repo repositories.SettingsRepository) *DNSProviderService {
	return &DNSProviderService{settingsRepo: repo}
}

func (s *DNSProviderService) ProvisionARecord(ctx context.Context, domain string) error {
	cfg, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil || cfg == nil {
		return err
	}
	targetIP := cfg.PublicIPv4
	if targetIP == "" {
		return fmt.Errorf("PublicIPv4 is not set in server settings")
	}

	if cfg.CloudflareAPIToken != "" {
		if err := s.provisionCloudflare(ctx, cfg.CloudflareAPIToken, domain, targetIP); err == nil {
			return nil
		}
	}

	if cfg.NamecheapAPIKey != "" && cfg.NamecheapAPIUser != "" {
		if err := s.provisionNamecheap(ctx, cfg, domain, targetIP); err == nil {
			return nil
		}
	}

	if cfg.SpaceshipAPIKey != "" {
		if err := s.provisionSpaceship(ctx, cfg.SpaceshipAPIKey, domain, targetIP); err == nil {
			return nil
		}
	}

	return nil
}

func getRootDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return domain
	}
	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}

func (s *DNSProviderService) provisionCloudflare(ctx context.Context, token, domain, targetIP string) error {
	rootDomain := getRootDomain(domain)
	client := &http.Client{}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.cloudflare.com/client/v4/zones?name="+rootDomain, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var zoneRes struct {
		Result []struct {
			ID string `json:"id"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&zoneRes); err != nil || len(zoneRes.Result) == 0 {
		return fmt.Errorf("cloudflare zone not found for %s", rootDomain)
	}
	zoneID := zoneRes.Result[0].ID

	payload := map[string]any{
		"type":    "A",
		"name":    domain,
		"content": targetIP,
		"proxied": false,
	}
	b, _ := json.Marshal(payload)
	req2, _ := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID), bytes.NewBuffer(b))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := client.Do(req2)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()
	if resp2.StatusCode >= 400 {
		return fmt.Errorf("cloudflare returned status %d", resp2.StatusCode)
	}
	return nil
}

func (s *DNSProviderService) provisionNamecheap(ctx context.Context, cfg *models.ServerSettings, domain, targetIP string) error {
	rootDomain := getRootDomain(domain)
	subDomain := strings.TrimSuffix(domain, "."+rootDomain)
	if subDomain == domain {
		subDomain = "@"
	}
	parts := strings.Split(rootDomain, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid root domain for namecheap")
	}

	client := &http.Client{}
	u := "https://api.namecheap.com/xml.response"
	q := url.Values{}
	q.Set("ApiUser", cfg.NamecheapAPIUser)
	q.Set("ApiKey", cfg.NamecheapAPIKey)
	q.Set("UserName", cfg.NamecheapAPIUser)
	q.Set("Command", "namecheap.domains.dns.addHost")
	q.Set("ClientIp", cfg.NamecheapClientIP)
	q.Set("SLD", parts[0])
	q.Set("TLD", parts[1])
	q.Set("HostName1", subDomain)
	q.Set("RecordType1", "A")
	q.Set("Address1", targetIP)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u+"?"+q.Encode(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *DNSProviderService) provisionSpaceship(ctx context.Context, key, domain, targetIP string) error {
	payload := map[string]any{
		"type":    "A",
		"name":    domain,
		"address": targetIP,
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://spaceship.dev/api/v1/dns", bytes.NewBuffer(b))
	req.Header.Set("X-Api-Key", key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
