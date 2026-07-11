package sso

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"
	"github.com/labstack/echo/v4"
)

type SAMLService struct {
	BaseURL string
	Key     *rsa.PrivateKey
	Cert    *x509.Certificate
}

func NewSAMLService(baseURL string, key *rsa.PrivateKey, cert *x509.Certificate) *SAMLService {
	return &SAMLService{
		BaseURL: baseURL,
		Key:     key,
		Cert:    cert,
	}
}

func (s *SAMLService) ConfigureMiddleware(idpMetadataURL string) (*samlsp.Middleware, error) {
	idpMetadata, err := url.Parse(idpMetadataURL)
	if err != nil {
		return nil, err
	}

	idpMetadataResp, err := http.Get(idpMetadata.String())
	if err != nil {
		return nil, err
	}
	defer idpMetadataResp.Body.Close()

	idpEntityDescriptor, err := samlsp.ParseMetadata(idpMetadataResp.Body)
	if err != nil {
		return nil, err
	}

	rootURL, err := url.Parse(s.BaseURL)
	if err != nil {
		return nil, err
	}

	return samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         s.Key,
		Certificate: s.Cert,
		IDPMetadata: idpEntityDescriptor,
	})
}

// Wrap is an Echo middleware wrapper around samlsp.Middleware
func (s *SAMLService) RequireSAML(samlSP *samlsp.Middleware) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// A quick wrapper passing the standard http.Handler
			handler := samlSP.RequireAccount(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				next(c)
			}))
			handler.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}
	}
}
