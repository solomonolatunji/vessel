package handlers

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net/http"

	"vessel.dev/vessel/internal/cloud/services/sso"

	"github.com/crewjam/saml/samlsp"
	"github.com/labstack/echo/v4"
)

type SSOHandler struct {
	samlMiddleware *samlsp.Middleware
}

func NewSSOHandler(baseURL string, idpMetadataURL string, key *rsa.PrivateKey, cert *x509.Certificate) (*SSOHandler, error) {
	svc := sso.NewSAMLService(baseURL, key, cert)

	// Create middleware
	sp, err := svc.ConfigureMiddleware(idpMetadataURL)
	if err != nil {
		return nil, fmt.Errorf("failed to configure SAML middleware: %v", err)
	}

	return &SSOHandler{
		samlMiddleware: sp,
	}, nil
}

// RegisterRoutes sets up the SAML endpoints:
// - GET /sso/saml/metadata
// - POST /sso/saml/acs
func (h *SSOHandler) RegisterRoutes(g *echo.Group) {
	// Crewjam's samlsp.Middleware provides its own ServeHTTP for /saml/acs and /saml/metadata
	// We wrap it for Echo
	g.Any("/saml/*", echo.WrapHandler(h.samlMiddleware))
}

// RequireSAML is an Echo middleware protecting routes that require SSO
func (h *SSOHandler) RequireSAML() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			handler := h.samlMiddleware.RequireAccount(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				// The session information is stored in samlsp.Session
				session := samlsp.SessionFromContext(r.Context())
				if session != nil {
					c.Set("saml_session", session)
				}
				next(c)
			}))
			handler.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}
	}
}
