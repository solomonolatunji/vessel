package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// CloudClaims holds the JWT claims for authenticated Vessel Cloud users.
type CloudClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// UserID returns the user's ID from the Subject claim (populated from "sub").
func (c *CloudClaims) UserID() string {
	return c.Subject
}

// RequireCloudAuth validates a Bearer JWT and injects CloudClaims into the echo context as "cloud_user".
func RequireCloudAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := parseCloudJWT(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			c.Set("cloud_user", claims)
			return next(c)
		}
	}
}

// RequireAdmin validates the token and enforces role == "admin".
func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := parseCloudJWT(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			if claims.Role != "admin" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden: admin access required"})
			}
			c.Set("cloud_user", claims)
			return next(c)
		}
	}
}

// RequireStaff validates the token and enforces role == "admin" or "staff".
func RequireStaff() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := parseCloudJWT(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			if claims.Role != "admin" && claims.Role != "staff" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden: staff access required"})
			}
			c.Set("cloud_user", claims)
			return next(c)
		}
	}
}

func parseCloudJWT(c echo.Context) (*CloudClaims, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, echo.ErrUnauthorized
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	secret := os.Getenv("VESSEL_CLOUD_JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}

	claims := &CloudClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.ErrUnauthorized
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, echo.ErrUnauthorized
	}

	return claims, nil
}
