package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.dev/codedock/internal/models"
)

func SetAuthCookie(c echo.Context, token string) {
	c.SetCookie(&http.Cookie{
		Name:     "codedock_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   72 * 3600,
	})
}

func ClearAuthCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "codedock_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
}

func ExtractClaims(c echo.Context) *models.UserClaims {
	if claims, ok := c.Get("user").(*models.UserClaims); ok {
		return claims
	}
	return nil
}

func ExtractUserID(c echo.Context) string {
	if claims := ExtractClaims(c); claims != nil {
		return claims.UserID
	}
	return ""
}
