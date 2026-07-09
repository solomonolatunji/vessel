package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"
	"strings"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.userService.ListUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	var out []models.User
	for _, u := range users {
		u.PasswordHash = ""
		out = append(out, u)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized access"})
	}
	u, err := h.userService.GetUserByID(c.Request().Context(), userID)
	if err != nil || u == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user profile not found"})
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	return c.JSON(http.StatusOK, &uCopy)
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized access"})
	}
	u, err := h.userService.GetUserByID(c.Request().Context(), userID)
	if err != nil || u == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user profile not found"})
	}
	var payload struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if payload.Email != "" {
		u.Email = payload.Email
	}
	if payload.Role != "" {
		u.Role = payload.Role
	}
	if err := h.userService.UpdateUser(c.Request().Context(), u); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	return c.JSON(http.StatusOK, &uCopy)
}

func (h *UserHandler) CreatePAT(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized access"})
	}
	var payload struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "token name is required"})
	}
	pat, rawToken, err := h.userService.CreatePAT(c.Request().Context(), userID, payload.Name, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"token": rawToken,
		"pat":   pat,
	})
}

func (h *UserHandler) ListPATs(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized access"})
	}
	pats, err := h.userService.ListPATs(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, pats)
}

func (h *UserHandler) DeletePAT(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized access"})
	}
	tokenID := c.Param("id")
	if tokenID == "" {
		tokenID = strings.TrimPrefix(c.Request().URL.Path, "/api/auth/pat/")
	}
	if tokenID == "" || tokenID == c.Request().URL.Path {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid personal access token id"})
	}
	if err := h.userService.DeletePAT(c.Request().Context(), tokenID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
