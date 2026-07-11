package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	// db *repos.CloudDB
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Register handles new cloud user registration
// @Summary Register
// @Description Register a new user for Vessel Cloud
// @Tags Cloud-Auth
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string
// @Router /cloud/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	// TODO: Parse request, hash password, create user in cloud_users
	// TODO: Create default team in cloud_teams
	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

// Login handles cloud user authentication
// @Summary Login
// @Description Authenticate a cloud user and return a JWT
// @Tags Cloud-Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /cloud/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	// TODO: Verify credentials, generate JWT
	return c.JSON(http.StatusOK, map[string]string{"token": "mock_jwt_token"})
}
