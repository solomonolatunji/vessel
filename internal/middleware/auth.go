package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

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

func (g *AuthGuard) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if g.Settings != nil {
			settings, _ := g.Settings.GetSettings(r.Context())
			if settings != nil && strings.TrimSpace(settings.IPAllowlist) != "" {
				clientIP := ExtractClientIP(r)
				if !IsIPAllowed(clientIP, settings.IPAllowlist) {
					writeError(w, http.StatusForbidden, fmt.Sprintf("access denied from IP address %s by server allowlist policy", clientIP))
					return
				}
			}
		}

		if g.TokenService == nil {
			userClaims := &models.UserClaims{
				UserID: "default",
				Email:  "default@vessel.dev",
				Role:   "admin",
			}
			ctx := context.WithValue(r.Context(), userClaimsKey, userClaims)
			next(w, r.WithContext(ctx))
			return
		}

		tokenStr := ExtractTokenFromRequest(r)
		if tokenStr == "" {
			writeError(w, http.StatusUnauthorized, "missing authentication token")
			return
		}

		claimsMap, err := g.TokenService.ValidateToken(tokenStr)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid authentication token: "+err.Error())
			return
		}

		totpEnabled, _ := claimsMap["totpEnabled"].(bool)
		userClaims := &models.UserClaims{
			UserID:      fmt.Sprintf("%v", claimsMap["sub"]),
			Email:       fmt.Sprintf("%v", claimsMap["email"]),
			Role:        fmt.Sprintf("%v", claimsMap["role"]),
			TOTPEnabled: totpEnabled,
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, userClaims)
		next(w, r.WithContext(ctx))
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

func ExtractClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (g *AuthGuard) RequireRole(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if g.Settings != nil {
			settings, _ := g.Settings.GetSettings(r.Context())
			if settings != nil && strings.TrimSpace(settings.IPAllowlist) != "" {
				clientIP := ExtractClientIP(r)
				if !IsIPAllowed(clientIP, settings.IPAllowlist) {
					writeError(w, http.StatusForbidden, fmt.Sprintf("access denied from IP address %s by server allowlist policy", clientIP))
					return
				}
			}
		}

		if g.TokenService == nil {
			userClaims := &models.UserClaims{
				UserID: "default",
				Email:  "default@vessel.dev",
				Role:   "admin",
			}
			ctx := context.WithValue(r.Context(), userClaimsKey, userClaims)
			next(w, r.WithContext(ctx))
			return
		}

		tokenStr := ExtractTokenFromRequest(r)
		if tokenStr == "" {
			writeError(w, http.StatusUnauthorized, "missing authentication token")
			return
		}

		claimsMap, err := g.TokenService.ValidateToken(tokenStr)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid authentication token: "+err.Error())
			return
		}

		role := fmt.Sprintf("%v", claimsMap["role"])
		if role != requiredRole && role != "admin" {
			writeError(w, http.StatusForbidden, "insufficient permissions")
			return
		}

		totpEnabled, _ := claimsMap["totpEnabled"].(bool)
		userClaims := &models.UserClaims{
			UserID:      fmt.Sprintf("%v", claimsMap["sub"]),
			Email:       fmt.Sprintf("%v", claimsMap["email"]),
			Role:        role,
			TOTPEnabled: totpEnabled,
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, userClaims)
		next(w, r.WithContext(ctx))
	}
}

func GetUserClaimsFromContext(ctx context.Context) *models.UserClaims {
	if c, ok := ctx.Value(userClaimsKey).(*models.UserClaims); ok {
		return c
	}
	return nil
}

func ExtractTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}

	cookie, err := r.Cookie("vessel_token")
	if err == nil && cookie.Value != "" {
		return strings.TrimSpace(cookie.Value)
	}

	queryToken := r.URL.Query().Get("token")
	if queryToken != "" {
		return strings.TrimSpace(queryToken)
	}

	return ""
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, `{"error":"%s"}`, message)
}
