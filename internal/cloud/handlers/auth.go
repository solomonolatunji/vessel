package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// @Summary Register
// @Description Register a new user for Vessel Cloud
// @Tags Cloud-Auth
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string
// @Router /cloud/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

// @Summary Login
// @Description Authenticate a cloud user and return a JWT
// @Tags Cloud-Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /cloud/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"token": "mock_jwt_token"})
}
