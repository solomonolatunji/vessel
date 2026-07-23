package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/services"
	"codedock.dev/codedock/internal/telemetry"
	"codedock.dev/codedock/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var payload ForgotPasswordRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	scheme := "http"
	if c.Request().TLS != nil || c.Request().Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	originUrl := scheme + "://" + c.Request().Host

	err := h.authService.ForgotPassword(c.Request().Context(), payload.Email, originUrl)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "If an account with that email exists, a password reset link has been sent.", nil)
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var payload ResetPasswordRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	err := h.authService.ResetPassword(c.Request().Context(), payload.Token, payload.NewPassword)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Password reset successful", nil)
}

func (h *AuthHandler) Register(c echo.Context) error {
	var payload RegisterRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	u, token, refreshToken, err := h.authService.Register(c.Request().Context(), payload.Name, payload.Email, payload.Password)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	SetAuthCookie(c, token)
	telemetry.Track(u.Email, "user_signed_up", map[string]interface{}{
		"email": u.Email,
		"name":  u.Name,
	})
	return utils.Success(c, "Registration successful", map[string]any{
		"token":        token,
		"refreshToken": refreshToken,
		"user":         u,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var payload AuthRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	u, token, refreshToken, err := h.authService.Login(c.Request().Context(), payload.Email, payload.Password)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err.Error())
	}
	SetAuthCookie(c, token)
	telemetry.Track(u.Email, "user_logged_in", map[string]interface{}{
		"email": u.Email,
	})
	return utils.Success(c, "Login successful", map[string]any{
		"token":        token,
		"refreshToken": refreshToken,
		"user":         u,
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var payload RefreshRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	u, token, newRefreshToken, err := h.authService.RefreshToken(c.Request().Context(), payload.RefreshToken)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err.Error())
	}
	SetAuthCookie(c, token)
	return utils.Success(c, "Token refreshed successfully", map[string]any{
		"token":        token,
		"refreshToken": newRefreshToken,
		"user":         u,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	ClearAuthCookie(c)
	return utils.Success(c, "Logged out successfully", nil)
}

type AdminInviteUserRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *AuthHandler) AdminInviteUser(c echo.Context) error {
	var req AdminInviteUserRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}
	if req.Email == "" {
		return utils.Error(c, http.StatusBadRequest, "email is required")
	}
	if req.Role == "" {
		req.Role = string(models.UserRoleMember)
	}

	origin := c.Request().Header.Get("Origin")
	if origin == "" {
		origin = c.Request().Header.Get("Referer")
	}
	if origin == "" {
		origin = "http://localhost:3000"
	}

	u, err := h.authService.InviteUser(c.Request().Context(), req.Email, models.UserRole(req.Role), origin)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	u.PasswordHash = ""
	return utils.Created(c, "User invited", u)
}
