package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/solomonolatunji/vessel/internal/services"
	"github.com/solomonolatunji/vessel/internal/types"
)

type contextKey string

const userClaimsKey contextKey = "user_claims"

type AuthGuard struct {
	TokenService *services.TokenService
}

// NewAuthGuard initializes a new AuthGuard with the provided token service.
func NewAuthGuard(ts *services.TokenService) *AuthGuard {
	return &AuthGuard{TokenService: ts}
}

func (g *AuthGuard) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if g.TokenService == nil {
			userClaims := &types.UserClaims{
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

		userClaims := &types.UserClaims{
			UserID: fmt.Sprintf("%v", claimsMap["sub"]),
			Email:  fmt.Sprintf("%v", claimsMap["email"]),
			Role:   fmt.Sprintf("%v", claimsMap["role"]),
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, userClaims)
		next(w, r.WithContext(ctx))
	}
}

func (g *AuthGuard) RequireRole(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return g.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
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

// ExtractTokenFromRequest extracts a JWT or PAT from the Authorization header, cookie, or query parameters.
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
