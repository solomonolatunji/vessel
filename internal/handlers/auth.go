package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
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

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Register endpoint
// @Description Register endpoint
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration credentials"
// @Router /auth/signup [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var payload RegisterRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	u, token, err := h.authService.Register(c.Request().Context(), payload.Name, payload.Email, payload.Password)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	SetAuthCookie(c, token)
	return utils.Success(c, "Registration successful", map[string]any{
		"token": token,
		"user":  u,
	})
}

// @Summary Login endpoint
// @Description Login endpoint
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Login credentials"
// @Router /auth/signin [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var payload AuthRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	u, token, err := h.authService.Login(c.Request().Context(), payload.Email, payload.Password)
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, err.Error())
	}
	SetAuthCookie(c, token)
	return utils.Success(c, "Login successful", map[string]any{
		"token": token,
		"user":  u,
	})
}

// @Summary Logout endpoint
// @Description Logout endpoint
// @Tags Auth
// @Accept json
// @Produce json
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	ClearAuthCookie(c)
	return utils.Success(c, "Logged out successfully", nil)
}
