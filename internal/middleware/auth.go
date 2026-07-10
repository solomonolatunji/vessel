package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type contextKey string

const userClaimsKey contextKey = "user_claims"

type SettingsProvider interface {
	GetSettings(context.Context) (*models.ServerSettings, error)
}

type ProjectTokenProvider interface {
	GetTokenByHash(ctx context.Context, tokenHash string) (*models.ProjectToken, error)
	UpdateTokenLastUsed(ctx context.Context, id string) error
}

type AuthGuard struct {
	TokenService  *services.TokenService
	Settings      SettingsProvider
	ProjectTokens ProjectTokenProvider
}

func NewAuthGuard(ts *services.TokenService, sp SettingsProvider, pt ProjectTokenProvider) *AuthGuard {
	return &AuthGuard{TokenService: ts, Settings: sp, ProjectTokens: pt}
}

func (g *AuthGuard) checkIPAllowlist(c echo.Context) error {
	if g.Settings == nil {
		return nil
	}
	settings, _ := g.Settings.GetSettings(c.Request().Context())
	if settings == nil || strings.TrimSpace(settings.IPAllowlist) == "" {
		return nil
	}
	clientIP := c.RealIP()
	if !IsIPAllowed(clientIP, settings.IPAllowlist) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": fmt.Sprintf("access denied from IP address %s by server allowlist policy", clientIP)})
	}
	return nil
}

func (g *AuthGuard) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := g.checkIPAllowlist(c); err != nil {
				return err
			}
			tokenStr := ExtractTokenFromRequest(c)
			if tokenStr == "" {
				if g.TokenService == nil {
					userClaims := &models.UserClaims{
						UserID: "default",
						Email:  "default@vessel.dev",
						Role:   "admin",
					}
					c.Set("user", userClaims)
					return next(c)
				}
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authentication token"})
			}
			if strings.HasPrefix(tokenStr, "vsl_tok_") {
				if g.ProjectTokens == nil {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "API tokens not supported"})
				}
				pt, err := g.ProjectTokens.GetTokenByHash(c.Request().Context(), tokenStr)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or revoked API token"})
				}
				importTime := time.Now()
				if pt.ExpiresAt != nil && pt.ExpiresAt.Before(importTime) {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "API token has expired"})
				}
				if len(pt.IPAllowlist) > 0 {
					clientIP := c.RealIP()
					if !IsIPAllowed(clientIP, strings.Join(pt.IPAllowlist, ",")) {
						return c.JSON(http.StatusForbidden, map[string]string{"error": "IP address not allowed for this API token"})
					}
				}
				_ = g.ProjectTokens.UpdateTokenLastUsed(c.Request().Context(), pt.ID)
				userClaims := &models.UserClaims{
					UserID: "api-token-" + pt.ID,
					Email:  "api@" + pt.ProjectID + ".vessel.local",
					Role:   "api",
				}
				c.Set("user", userClaims)
				c.Set("api_scopes", pt.Scopes)
				c.Set("project_id", pt.ProjectID)
				c.Set("environment_id", pt.EnvironmentID)
				return next(c)
			}
			claimsMap, err := g.TokenService.ValidateToken(tokenStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authentication token: " + err.Error()})
			}
			totpEnabled, _ := claimsMap["totpEnabled"].(bool)
			userClaims := &models.UserClaims{
				UserID:      fmt.Sprintf("%v", claimsMap["sub"]),
				Email:       fmt.Sprintf("%v", claimsMap["email"]),
				Role:        fmt.Sprintf("%v", claimsMap["role"]),
				TOTPEnabled: totpEnabled,
			}
			c.Set("user", userClaims)
			return next(c)
		}
	}
}

func IsIPAllowed(clientIPStr string, allowlistStr string) bool {
	clientIP := net.ParseIP(clientIPStr)
	if clientIP == nil {
		return false
	}
	entries := strings.Split(allowlistStr, ",")
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if strings.Contains(entry, "/") {
			_, cidrNet, err := net.ParseCIDR(entry)
			if err == nil && cidrNet.Contains(clientIP) {
				return true
			}
		} else {
			if clientIPStr == entry {
				return true
			}
		}
	}
	return false
}

func ExtractClientIP(c echo.Context) string {
	return c.RealIP()
}

func (g *AuthGuard) RequireScope(requiredScope string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims, ok := c.Get("user").(*models.UserClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			if userClaims.Role == "api" {
				scopes, ok := c.Get("api_scopes").([]string)
				if !ok {
					return c.JSON(http.StatusForbidden, map[string]string{"error": "insufficient scopes"})
				}
				hasScope := false
				for _, s := range scopes {
					if s == requiredScope || s == "admin" || s == "*" {
						hasScope = true
						break
					}
				}
				if !hasScope {
					return c.JSON(http.StatusForbidden, map[string]string{"error": "missing required scope: " + requiredScope})
				}
			}
			return next(c)
		}
	}
}

func (g *AuthGuard) RequireRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := g.checkIPAllowlist(c); err != nil {
				return err
			}
			if g.TokenService == nil {
				userClaims := &models.UserClaims{
					UserID: "default",
					Email:  "default@vessel.dev",
					Role:   "admin",
				}
				c.Set("user", userClaims)
				return next(c)
			}
			tokenStr := ExtractTokenFromRequest(c)
			if tokenStr == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authentication token"})
			}
			if strings.HasPrefix(tokenStr, "vsl_tok_") {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "API tokens cannot access role-restricted endpoints"})
			}
			claimsMap, err := g.TokenService.ValidateToken(tokenStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authentication token: " + err.Error()})
			}
			role := fmt.Sprintf("%v", claimsMap["role"])
			if role != requiredRole && role != "admin" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "insufficient permissions"})
			}
			totpEnabled, _ := claimsMap["totpEnabled"].(bool)
			userClaims := &models.UserClaims{
				UserID:      fmt.Sprintf("%v", claimsMap["sub"]),
				Email:       fmt.Sprintf("%v", claimsMap["email"]),
				Role:        role,
				TOTPEnabled: totpEnabled,
			}
			c.Set("user", userClaims)
			return next(c)
		}
	}
}

func GetUserClaimsFromContext(ctx context.Context) *models.UserClaims {
	if c, ok := ctx.Value(userClaimsKey).(*models.UserClaims); ok {
		return c
	}
	return nil
}

func ExtractTokenFromRequest(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}
	cookie, err := c.Cookie("vessel_token")
	if err == nil && cookie.Value != "" {
		return strings.TrimSpace(cookie.Value)
	}
	queryToken := c.QueryParam("token")
	if queryToken != "" {
		return strings.TrimSpace(queryToken)
	}
	return ""
}
