package types

import "time"

// DomainConfig manages custom domain routing, SSL certificate issuance state, and Caddy integration.
type DomainConfig struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	DomainName    string    `json:"domainName"`
	RedirectTo    string    `json:"redirectTo,omitempty"`
	SSLCertStatus string    `json:"sslCertStatus"`
	PathPrefix    string    `json:"pathPrefix"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
