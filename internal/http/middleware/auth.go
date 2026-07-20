package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
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

type ProjectMemberProvider interface {
	GetMember(ctx context.Context, projectID, userID string) (*models.ProjectMember, error)
}

type AuthGuard struct {
	TokenService   *services.TokenService
	Settings       SettingsProvider
	ProjectTokens  ProjectTokenProvider
	ProjectMembers ProjectMemberProvider
}

func NewAuthGuard(ts *services.TokenService, sp SettingsProvider, pt ProjectTokenProvider, pm ProjectMemberProvider) *AuthGuard {
	return &AuthGuard{TokenService: ts, Settings: sp, ProjectTokens: pt, ProjectMembers: pm}
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
		return utils.Error(c, http.StatusForbidden, fmt.Sprintf("access denied from IP address %s by server allowlist policy", clientIP))
	}
	return nil
}

func (g *AuthGuard) validateAPIToken(c echo.Context, tokenStr string, denyAPITokens bool) (*models.UserClaims, error) {
	if denyAPITokens {
		return nil, utils.Error(c, http.StatusForbidden, "API tokens cannot access role-restricted endpoints")
	}
	if g.ProjectTokens == nil {
		return nil, utils.Error(c, http.StatusUnauthorized, "API tokens not supported")
	}
	pt, err := g.ProjectTokens.GetTokenByHash(c.Request().Context(), tokenStr)
	if err != nil {
		return nil, utils.Error(c, http.StatusUnauthorized, "invalid or revoked API token")
	}
	if pt.ExpiresAt != nil && pt.ExpiresAt.Before(time.Now()) {
		return nil, utils.Error(c, http.StatusUnauthorized, "API token has expired")
	}
	if len(pt.IPAllowlist) > 0 {
		if !IsIPAllowed(c.RealIP(), strings.Join(pt.IPAllowlist, ",")) {
			return nil, utils.Error(c, http.StatusForbidden, "IP address not allowed for this API token")
		}
	}
	_ = g.ProjectTokens.UpdateTokenLastUsed(c.Request().Context(), pt.ID)

	c.Set("api_scopes", pt.Scopes)
	c.Set("project_id", pt.ProjectID)
	c.Set("environment_id", pt.EnvironmentID)

	return &models.UserClaims{
		UserID: "api-token-" + pt.ID,
		Email:  "api@" + pt.ProjectID + ".vessl.local",
		Role:   "api",
	}, nil
}

func (g *AuthGuard) validateJWT(c echo.Context, tokenStr string) (*models.UserClaims, error) {
	claimsMap, err := g.TokenService.ValidateToken(tokenStr)
	if err != nil {
		return nil, utils.Error(c, http.StatusUnauthorized, "invalid authentication token: "+err.Error())
	}

	totpEnabled, _ := claimsMap["totpEnabled"].(bool)
	return &models.UserClaims{
		UserID:      fmt.Sprintf("%v", claimsMap["sub"]),
		Email:       fmt.Sprintf("%v", claimsMap["email"]),
		Role:        models.UserRole(fmt.Sprintf("%v", claimsMap["role"])),
		TOTPEnabled: totpEnabled,
	}, nil
}

func (g *AuthGuard) baseAuth(c echo.Context, denyAPITokens bool) (*models.UserClaims, error) {
	if err := g.checkIPAllowlist(c); err != nil {
		return nil, err
	}
	tokenStr := ExtractTokenFromRequest(c)
	if tokenStr == "" {
		if g.TokenService == nil {
			return &models.UserClaims{
				UserID: "default",
				Email:  "default@vessl.dev",
				Role:   "admin",
			}, nil
		}
		return nil, utils.Error(c, http.StatusUnauthorized, "missing authentication token")
	}

	if strings.HasPrefix(tokenStr, "vsl_tok_") {
		return g.validateAPIToken(c, tokenStr, denyAPITokens)
	}

	return g.validateJWT(c, tokenStr)
}

func (g *AuthGuard) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims, err := g.baseAuth(c, false)
			if err != nil {
				return err
			}
			c.Set("user", userClaims)
			ctx := context.WithValue(c.Request().Context(), userClaimsKey, userClaims)
			c.SetRequest(c.Request().WithContext(ctx))
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
			if !ok || userClaims == nil {
				return utils.Error(c, http.StatusUnauthorized, "unauthorized")
			}
			if userClaims.Role == "api" {
				scopes, ok := c.Get("api_scopes").([]string)
				if !ok {
					return utils.Error(c, http.StatusForbidden, "insufficient scopes")
				}
				hasScope := false
				for _, s := range scopes {
					if s == requiredScope || s == "admin" || s == "*" {
						hasScope = true
						break
					}
				}
				if !hasScope {
					return utils.Error(c, http.StatusForbidden, "missing required scope: "+requiredScope)
				}
			}
			return next(c)
		}
	}
}

func (g *AuthGuard) RequireProjectRole(minPermission models.MemberPermission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims, ok := c.Get("user").(*models.UserClaims)
			if !ok || userClaims == nil {
				return utils.Error(c, http.StatusUnauthorized, "unauthorized")
			}

			// Admin API keys or instance admins bypass project-level checks
			if userClaims.Role == "admin" {
				return next(c)
			}

			projectID := c.Param("projectId")
			if projectID == "" {
				projectID = c.Param("id") // fallback if route uses :id instead of :projectId
			}
			if projectID == "" {
				return utils.Error(c, http.StatusBadRequest, "missing project id")
			}

			if g.ProjectMembers == nil {
				return utils.Error(c, http.StatusInternalServerError, "project members provider not configured")
			}

			member, err := g.ProjectMembers.GetMember(c.Request().Context(), projectID, userClaims.UserID)
			if err != nil {
				return utils.Error(c, http.StatusInternalServerError, "failed to verify project membership")
			}
			if member == nil {
				return utils.Error(c, http.StatusForbidden, "you do not have access to this project")
			}

			// Validate permission level if necessary
			// For now, if they are a member, we allow them. If minPermission is specific, we could enforce it.
			if minPermission != "" && member.Permission != minPermission && member.Permission != models.MemberPermissionAdmin {
				return utils.Error(c, http.StatusForbidden, "insufficient project permissions")
			}

			return next(c)
		}
	}
}

func (g *AuthGuard) RequireRole(requiredRole models.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims, err := g.baseAuth(c, true)
			if err != nil {
				return err
			}
			if userClaims.Role != requiredRole && userClaims.Role != models.UserRoleAdmin {
				return utils.Error(c, http.StatusForbidden, "insufficient permissions")
			}
			c.Set("user", userClaims)
			ctx := context.WithValue(c.Request().Context(), userClaimsKey, userClaims)
			c.SetRequest(c.Request().WithContext(ctx))
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
	cookie, err := c.Cookie("vessl_token")
	if err == nil && cookie.Value != "" {
		return strings.TrimSpace(cookie.Value)
	}
	return ""
}
