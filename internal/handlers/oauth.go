package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type OAuthHandler struct {
	oauthService *services.OAuthService
}

func NewOAuthHandler(s *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{oauthService: s}
}

type Verify2FARequest struct {
	Passcode string `json:"passcode"`
}

// @Summary ListProviders endpoint
// @Description ListProviders endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/oauth/providers [get]
func (h *OAuthHandler) ListProviders(c echo.Context) error {
	providers, err := h.oauthService.ListProviders(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, providers)
}

// @Summary SaveProvider endpoint
// @Description SaveProvider endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body models.OAuthProviderConfig true "Payload"
// @Router /api/settings/oauth/providers [put]
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

// @Summary OAuthRedirect endpoint
// @Description OAuthRedirect endpoint
// @Tags Auth
// @Accept json
// @Produce json
// @Param provider path string true "provider"
// @Router /api/auth/oauth/{provider} [get]
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

// @Summary OAuthCallback endpoint
// @Description OAuthCallback endpoint
// @Tags Auth
// @Accept json
// @Produce json
// @Param provider path string true "provider"
// @Router /api/auth/oauth/{provider}/callback [get]
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

// @Summary Setup2FA endpoint
// @Description Setup2FA endpoint
// @Tags Auth
// @Accept json
// @Produce json
// @Router /api/auth/2fa/setup [post]
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

// @Summary Verify2FA endpoint
// @Description Verify2FA endpoint
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body handlers.Verify2FARequest true "Payload"
// @Router /api/auth/2fa/verify [post]
func (h *OAuthHandler) Verify2FA(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized access"})
	}
	var payload Verify2FARequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing 6-digit passcode"})
	}
	if err := h.oauthService.Verify2FA(c.Request().Context(), userID, payload.Passcode); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "totp_enabled"})
}

// @Summary Disable2FA endpoint
// @Description Disable2FA endpoint
// @Tags Auth
// @Accept json
// @Produce json
// @Router /api/auth/2fa/disable [post]
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
