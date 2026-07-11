package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"vessel.dev/vessel/internal/cloud/services"
)

// DeploymentRateLimiter intercepts deployment requests to check if the team has exceeded their hourly limit
func DeploymentRateLimiter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Mock extracting user plan and team ID from JWT
			// In production, this is set by earlier AuthMiddleware
			teamID := "team_1"
			plan := "hobby"

			limit := services.GetFeatures().GetDeploymentRateLimit(teamID, plan)

			// TODO: Check Redis or Database to see how many deployments have occurred in the last hour
			// Mocking current usage
			currentUsage := 5

			if currentUsage >= limit {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Deployment rate limit exceeded for your tier",
				})
			}

			return next(c)
		}
	}
}

// SeatLimitGuard intercepts server connection (BYOS) requests to check if they can add another server
func SeatLimitGuard() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Mock extracting user plan and team ID from JWT
			teamID := "team_1"
			plan := "hobby"

			limit := services.GetFeatures().GetMaxServers(teamID, plan)

			// TODO: Query cloud_servers database to count currently connected active servers
			// Mocking current server count
			currentServers := 1

			if currentServers >= limit {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Bring Your Own Server (BYOS) seat limit reached. Please upgrade your plan.",
				})
			}

			return next(c)
		}
	}
}
