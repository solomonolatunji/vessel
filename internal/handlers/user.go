package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

type UpdateProfileRequest struct {
	Email string          `json:"email"`
	Role  models.UserRole `json:"role"`
}

type CreatePATRequest struct {
	Name string `json:"name"`
}

// @Summary ListUsers endpoint
// @Description ListUsers endpoint
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Router /users [get]
func (h *UserHandler) ListUsers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	users, total, err := h.userService.ListUsers(c.Request().Context(), limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	var out []models.User
	for _, u := range users {
		u.PasswordHash = ""
		out = append(out, u)
	}
	return utils.Paginated(c, "Users retrieved", out, total, page, limit)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

// @Summary ChangePassword endpoint
// @Description ChangePassword endpoint
// @Tags Users
// @Accept json
// @Produce json
// @Param request body handlers.ChangePasswordRequest true "Payload"
// @Router /profile/password [put]
func (h *UserHandler) ChangePassword(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	var payload ChangePasswordRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if payload.OldPassword == "" || payload.NewPassword == "" {
		return utils.Error(c, http.StatusBadRequest, "old and new passwords are required")
	}
	if err := h.userService.ChangePassword(c.Request().Context(), userID, payload.OldPassword, payload.NewPassword); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Password changed successfully", nil)
}

// @Summary GetProfile endpoint
// @Description GetProfile endpoint
// @Tags Users
// @Accept json
// @Produce json
// @Router /profile [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	u, err := h.userService.GetUserByID(c.Request().Context(), userID)
	if err != nil || u == nil {
		return utils.Error(c, http.StatusNotFound, "user profile not found")
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	return utils.Success(c, "Profile retrieved", &uCopy)
}

// @Summary UpdateProfile endpoint
// @Description UpdateProfile endpoint
// @Tags Users
// @Accept json
// @Produce json
// @Param request body handlers.UpdateProfileRequest true "Payload"
// @Router /profile [put]
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	u, err := h.userService.GetUserByID(c.Request().Context(), userID)
	if err != nil || u == nil {
		return utils.Error(c, http.StatusNotFound, "user profile not found")
	}
	var payload UpdateProfileRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if payload.Email != "" {
		u.Email = payload.Email
	}
	if payload.Role != "" {
		u.Role = payload.Role
	}
	if err := h.userService.UpdateUser(c.Request().Context(), u); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	return utils.Success(c, "Profile updated", &uCopy)
}

// @Summary CreatePAT endpoint
// @Description CreatePAT endpoint
// @Tags Users
// @Accept json
// @Produce json
// @Param request body handlers.CreatePATRequest true "Payload"
// @Router /profile/tokens [post]
func (h *UserHandler) CreatePAT(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	var payload CreatePATRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "token name is required")
	}
	pat, rawToken, err := h.userService.CreatePAT(c.Request().Context(), userID, payload.Name, nil)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Token created successfully", map[string]any{
		"token": rawToken,
		"pat":   pat,
	})
}

// @Summary ListPATs endpoint
// @Description ListPATs endpoint
// @Tags Users
// @Accept json
// @Produce json
// @Router /profile/tokens [get]
func (h *UserHandler) ListPATs(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	pats, err := h.userService.ListPATs(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Tokens retrieved successfully", pats)
}

// @Summary DeletePAT endpoint
// @Description DeletePAT endpoint
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /profile/tokens/{id} [delete]
func (h *UserHandler) DeletePAT(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	tokenID := c.Param("id")
	if tokenID == "" {
		tokenID = strings.TrimPrefix(c.Request().URL.Path, "/api/auth/pat/")
	}
	if tokenID == "" || tokenID == c.Request().URL.Path {
		return utils.Error(c, http.StatusBadRequest, "invalid personal access token id")
	}
	if err := h.userService.DeletePAT(c.Request().Context(), tokenID, userID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Token deleted successfully", nil)
}
