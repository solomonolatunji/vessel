package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
)

type Mailer interface {
	SendSystemEmail(ctx context.Context, templateName string, toAddress string, subject string, data any) error
}

type UserHandler struct {
	userService *services.UserService
	mailer      Mailer
}

func NewUserHandler(s *services.UserService, mailer Mailer) *UserHandler {
	return &UserHandler{userService: s, mailer: mailer}
}

type UpdateProfileRequest struct {
	Name string          `json:"name"`
	Role models.UserRole `json:"role"`
}

type RequestEmailChangeRequest struct {
	NewEmail string `json:"newEmail" validate:"required"`
}

type VerifyEmailChangeRequest struct {
	OTP string `json:"otp" validate:"required"`
}

type CreatePATRequest struct {
	Name            string     `json:"name"`
	AccessLevel     string     `json:"accessLevel"`
	ProjectScope    string     `json:"projectScope"`
	AllowedProjects []string   `json:"allowedProjects"`
	ExpiresAt       *time.Time `json:"expiresAt"`
}

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

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "user id required")
	}
	if err := h.userService.DeleteUser(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "User deleted", nil)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

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
	if payload.Name != "" {
		u.Name = payload.Name
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

func (h *UserHandler) RequestEmailChange(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	var payload RequestEmailChangeRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if payload.NewEmail == "" {
		return utils.Error(c, http.StatusBadRequest, "new email is required")
	}

	otp, err := services.GenerateEmailOTP(userID, payload.NewEmail)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to generate OTP")
	}

	err = h.mailer.SendSystemEmail(c.Request().Context(), "email_change", payload.NewEmail, "Verify Your Email Change", map[string]string{
		"OTP": otp,
	})
	if err != nil {

		return utils.Error(c, http.StatusInternalServerError, "failed to send verification email")
	}

	return utils.Success(c, "OTP sent to new email address", nil)
}

func (h *UserHandler) VerifyEmailChange(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	var payload VerifyEmailChangeRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	newEmail, ok := services.VerifyEmailOTP(userID, payload.OTP)
	if !ok {
		return utils.Error(c, http.StatusBadRequest, "invalid or expired OTP")
	}

	u, err := h.userService.GetUserByID(c.Request().Context(), userID)
	if err != nil || u == nil {
		return utils.Error(c, http.StatusNotFound, "user not found")
	}

	u.Email = newEmail
	if err := h.userService.UpdateUser(c.Request().Context(), u); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to update email")
	}

	return utils.Success(c, "Email updated successfully", nil)
}

func (h *UserHandler) CreatePAT(c echo.Context) error {
	userID := ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized access")
	}
	var payload CreatePATRequest
	if err := c.Bind(&payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	pat, rawToken, err := h.userService.CreatePAT(c.Request().Context(), userID, payload.Name, payload.AccessLevel, payload.ProjectScope, payload.AllowedProjects, payload.ExpiresAt)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Token created successfully", map[string]any{
		"token": rawToken,
		"pat":   pat,
	})
}

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
