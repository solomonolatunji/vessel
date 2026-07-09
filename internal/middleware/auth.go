package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type contextKey string

const userClaimsKey contextKey = "user_claims"

type SettingsProvider interface {
	GetSettings(context.Context) (*models.ServerSettings, error)
}

type AuthGuard struct {
	TokenService *services.TokenService
	Settings     SettingsProvider
}

func NewAuthGuard(ts *services.TokenService, sp SettingsProvider) *AuthGuard {
	return &AuthGuard{TokenService: ts, Settings: sp}
}

func (g *AuthGuard) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if g.Settings != nil {
				settings, _ := g.Settings.GetSettings(c.Request().Context())
				if settings != nil && strings.TrimSpace(settings.IPAllowlist) != "" {
					clientIP := c.RealIP()
					if !IsIPAllowed(clientIP, settings.IPAllowlist) {
						return c.JSON(http.StatusForbidden, map[string]string{"error": fmt.Sprintf("access denied from IP address %s by server allowlist policy", clientIP)})
					}
				}
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

func (g *AuthGuard) RequireRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if g.Settings != nil {
				settings, _ := g.Settings.GetSettings(c.Request().Context())
				if settings != nil && strings.TrimSpace(settings.IPAllowlist) != "" {
					clientIP := c.RealIP()
					if !IsIPAllowed(clientIP, settings.IPAllowlist) {
						return c.JSON(http.StatusForbidden, map[string]string{"error": fmt.Sprintf("access denied from IP address %s by server allowlist policy", clientIP)})
					}
				}
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
