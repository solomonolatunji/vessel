package handlers

import (
	"github.com/labstack/echo/v4"

	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type OAuthHandler struct {
	oauthService *services.OAuthService
}

func NewOAuthHandler(s *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{oauthService: s}
}

func (h *OAuthHandler) ListProviders(c echo.Context) error {
	providers, err := h.oauthService.ListProviders(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, providers)
}

func (h *OAuthHandler) SaveProvider(c echo.Context) error {
	var p models.OAuthProviderConfig
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if err := h.oauthService.SaveProvider(c.Request().Context(), &p); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (h *OAuthHandler) OAuthRedirect(c echo.Context) error {
	providerName := strings.TrimPrefix(c.Request().URL.Path, "/api/auth/oauth/")
	if idx := strings.Index(providerName, "/"); idx != -1 {
		providerName = providerName[:idx]
	}
	p, err := h.oauthService.GetProvider(c.Request().Context(), providerName)
	if err != nil || p == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "oauth provider not found or not enabled: " + providerName})
	}
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate secure state token"})
	}
	state := hex.EncodeToString(stateBytes)
	authURL, err := services.GetAuthorizationURL(p, state)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *OAuthHandler) OAuthCallback(c echo.Context) error {
	providerName := strings.TrimPrefix(c.Request().URL.Path, "/api/auth/oauth/")
	providerName = strings.TrimSuffix(providerName, "/callback")
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing authorization code parameter"})
	}
	token, _, err := h.oauthService.HandleCallback(c.Request().Context(), providerName, code)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	SetAuthCookie(c, token)
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *OAuthHandler) Setup2FA(c echo.Context) error {
	claims := ExtractClaims(c)
	if claims == nil || claims.UserID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized access"})
	}
	res, err := h.oauthService.Setup2FA(c.Request().Context(), claims.UserID, claims.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *OAuthHandler) Verify2FA(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized access"})
	}
	var payload struct {
		Passcode string `json:"passcode"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing 6-digit passcode"})
	}
	if err := h.oauthService.Verify2FA(c.Request().Context(), userID, payload.Passcode); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "totp_enabled"})
}

func (h *OAuthHandler) Disable2FA(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized access"})
	}
	if err := h.oauthService.Disable2FA(c.Request().Context(), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "totp_disabled"})
}
