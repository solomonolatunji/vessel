package api

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userClaimsKey contextKey = "user_claims"

// RequireAuth intercepts HTTP requests to verify the presence and validity of a JWT Bearer token or cookie.
func (s *Server) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractTokenFromRequest(r)
		if tokenStr == "" {
			writeError(w, http.StatusUnauthorized, "missing authentication token")
			return
		}

		claims, err := s.tokenService.ValidateToken(tokenStr)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid authentication token: "+err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

// RequireRole intercepts HTTP requests to ensure the authenticated user holds a specific authorization role or is an admin.
func (s *Server) RequireRole(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return s.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserClaimsFromContext(r.Context())
		if claims == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized access")
			return
		}

		if claims.Role != requiredRole && claims.Role != "admin" {
			writeError(w, http.StatusForbidden, "insufficient role privileges for this operation")
			return
		}

		next(w, r)
	})
}

// GetUserClaimsFromContext extracts verified JWT UserClaims stored in the request context during authentication.
func GetUserClaimsFromContext(ctx context.Context) *UserClaims {
	claims, ok := ctx.Value(userClaimsKey).(*UserClaims)
	if !ok {
		return nil
	}
	return claims
}

func extractTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}

	cookie, err := r.Cookie("vessel_token")
	if err == nil && cookie.Value != "" {
		return strings.TrimSpace(cookie.Value)
	}

	return ""
}
