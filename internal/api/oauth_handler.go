package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"vessel.dev/vessel/internal/middleware"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/services/oauth"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
)

type OAuthHandler struct {
	store        *store.Store
	oauthService *oauth.OAuthService
	tokenService *services.TokenService
}

func NewOAuthHandler(s *store.Store, os *oauth.OAuthService, ts *services.TokenService) *OAuthHandler {
	return &OAuthHandler{store: s, oauthService: os, tokenService: ts}
}

func (h *OAuthHandler) ListProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	providers, err := h.store.ListOAuthProviders()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, providers)
}

func (h *OAuthHandler) SaveProvider(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var p types.OAuthProvider
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if p.ID == "" && p.ProviderName != "" {
		p.ID = strings.ToLower(p.ProviderName)
	}

	if err := h.store.SaveOAuthProvider(&p); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (h *OAuthHandler) OAuthRedirect(w http.ResponseWriter, r *http.Request) {
	providerName := strings.TrimPrefix(r.URL.Path, "/api/auth/oauth/")
	if idx := strings.Index(providerName, "/"); idx != -1 {
		providerName = providerName[:idx]
	}

	p, err := h.store.GetOAuthProvider(providerName)
	if err != nil || p == nil {
		writeError(w, http.StatusNotFound, "oauth provider not found or not enabled: "+providerName)
		return
	}

	stateBytes := make([]byte, 16)
	_, _ = rand.Read(stateBytes)
	state := hex.EncodeToString(stateBytes)

	authURL, err := h.oauthService.GetAuthorizationURL(p, state)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	providerName := strings.TrimPrefix(r.URL.Path, "/api/auth/oauth/")
	providerName = strings.TrimSuffix(providerName, "/callback")

	p, err := h.store.GetOAuthProvider(providerName)
	if err != nil || p == nil {
		writeError(w, http.StatusNotFound, "oauth provider not found: "+providerName)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "missing authorization code parameter")
		return
	}

	email, err := h.oauthService.ExchangeCode(p, code)
	if err != nil || email == "" {
		writeError(w, http.StatusUnauthorized, "failed oauth code exchange: "+err.Error())
		return
	}

	user, err := h.store.GetUserByEmail(email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if user == nil {
		settings, _ := h.store.GetServerSettings()
		if settings != nil && !settings.RegistrationEnabled {
			writeError(w, http.StatusForbidden, "new account registration is disabled by the administrator")
			return
		}
		// Auto register via OAuth
		user = &types.User{
			Email:         email,
			PasswordHash:  "oauth-login-no-password",
			Role:          "member",
			OAuthProvider: p.ProviderName,
		}
		if err := h.store.CreateUser(user); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create user account from oauth: "+err.Error())
			return
		}
	}

	token, err := h.tokenService.GenerateToken(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed generating token")
		return
	}

	setAuthCookie(w, token)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) Setup2FA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims := middleware.GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	secret, err := oauth.GenerateTOTPSecret()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed generating totp secret")
		return
	}

	recoveryCodes, err := oauth.GenerateRecoveryCodes(8)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed generating recovery codes")
		return
	}

	if err := h.store.UpdateUserTOTP(claims.UserID, false, secret, recoveryCodes); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	qrURI := oauth.GenerateTOTPQRUri(claims.Email, secret)
	writeJSON(w, http.StatusOK, types.TwoFASetupResponse{
		Secret:        secret,
		QRCodeURI:     qrURI,
		RecoveryCodes: recoveryCodes,
	})
}

func (h *OAuthHandler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims := middleware.GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	var payload struct {
		Passcode string `json:"passcode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Passcode == "" {
		writeError(w, http.StatusBadRequest, "missing 6-digit passcode")
		return
	}

	secret, recoveryCodes, err := h.store.GetUserTOTPSecret(claims.UserID)
	if err != nil || secret == "" {
		writeError(w, http.StatusBadRequest, "totp setup has not been initiated for this user")
		return
	}

	if !oauth.ValidateTOTP(secret, payload.Passcode) {
		writeError(w, http.StatusUnauthorized, "invalid 6-digit totp verification code")
		return
	}

	if err := h.store.UpdateUserTOTP(claims.UserID, true, secret, recoveryCodes); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "totp_enabled"})
}

func (h *OAuthHandler) Disable2FA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims := middleware.GetUserClaimsFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	if err := h.store.UpdateUserTOTP(claims.UserID, false, "", nil); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "totp_disabled"})
}
