package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/solomonolatunji/vessel/internal/types"
)

type contextKey string

const userClaimsKey contextKey = "user_claims"

func (s *Server) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractTokenFromRequest(r)
		if tokenStr == "" {
			writeError(w, http.StatusUnauthorized, "missing authentication token")
			return
		}

		claimsMap, err := s.tokenService.ValidateToken(tokenStr)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid authentication token: "+err.Error())
			return
		}

		userClaims := &types.UserClaims{
			UserID: fmt.Sprintf("%v", claimsMap["sub"]),
			Email:  fmt.Sprintf("%v", claimsMap["email"]),
			Role:   fmt.Sprintf("%v", claimsMap["role"]),
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, userClaims)
		next(w, r.WithContext(ctx))
	}
}

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

func GetUserClaimsFromContext(ctx context.Context) *types.UserClaims {
	claims, ok := ctx.Value(userClaimsKey).(*types.UserClaims)
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
