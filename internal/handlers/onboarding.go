package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
)

type OnboardingHandler struct {
	userService       *services.UserService
	onboardingService *services.OnboardingService
}

func NewOnboardingHandler(
	userService *services.UserService,
	onboardingService *services.OnboardingService,
) *OnboardingHandler {
	return &OnboardingHandler{
		userService:       userService,
		onboardingService: onboardingService,
	}
}

func (h *OnboardingHandler) SetupStatus(c echo.Context) error {
	count, err := h.userService.CountUsers(c.Request().Context())
	if err != nil {
		return utils.Error(c, 500, "failed to check user count")
	}
	cwd, _ := os.Getwd()
	return utils.Success(c, "Setup status", map[string]any{
		"setupRequired": count == 0,
		"cwd":           cwd,
	})
}

func (h *OnboardingHandler) Setup(c echo.Context) error {
	ctx := c.Request().Context()

	var req services.SetupRequest
	if err := c.Bind(&req); err != nil {
		fmt.Printf("Setup Error: Failed to bind request: %v\n", err)
		return utils.Error(c, 400, "invalid request")
	}

	u, token, refreshToken, err := h.onboardingService.CompleteSetup(ctx, req)
	if err != nil {
		if err.Error() == "setup has already been completed" {
			return utils.Error(c, 403, err.Error())
		}
		return utils.Error(c, 400, err.Error())
	}

	SetAuthCookie(c, token)
	c.SetCookie(&http.Cookie{
		Name:     "codedock_refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	res := map[string]any{
		"user":         u,
		"token":        token,
		"refreshToken": refreshToken,
	}

	return utils.Success(c, "Setup completed successfully", res)
}
