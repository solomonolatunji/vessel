package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"vessel.dev/vessel/internal/cloud/repos"
	"vessel.dev/vessel/internal/cloud/services"
)

func DeploymentRateLimiter(repo repos.CloudRepo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			teamID := uint(1)
			teamStrID := "team_1"
			plan := "hobby"

			if team, err := repo.GetTeamByID(teamID); err == nil && team != nil {
				plan = team.Plan
			}

			limit := services.GetFeatures().GetDeploymentRateLimit(teamStrID, plan)
			currentUsage, _ := repo.GetDeploymentsInLastHour(teamID)

			if int(currentUsage) >= limit {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Deployment rate limit exceeded for your tier",
				})
			}

			return next(c)
		}
	}
}

func SeatLimitGuard(repo repos.CloudRepo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			teamID := uint(1)
			teamStrID := "team_1"
			plan := "hobby"

			if team, err := repo.GetTeamByID(teamID); err == nil && team != nil {
				plan = team.Plan
			}

			limit := services.GetFeatures().GetMaxServers(teamStrID, plan)
			currentServers, _ := repo.GetActiveServerCount(teamID)

			if int(currentServers) >= limit {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Bring Your Own Server (BYOS) seat limit reached. Please upgrade your plan.",
				})
			}

			return next(c)
		}
	}
}
