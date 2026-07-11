package services

import (
	"context"
	"log"

	"github.com/growthbook/growthbook-golang"
)

type FeatureFlagsService struct {
	client *growthbook.Client
}

var instance *FeatureFlagsService

func InitGrowthBook() {
	// Initialize the GrowthBook client globally
	client, err := growthbook.NewClient(context.Background())
	if err != nil {
		log.Printf("Failed to init GrowthBook: %v", err)
	}

	instance = &FeatureFlagsService{
		client: client,
	}

	log.Println("GrowthBook feature flag service initialized")
}

func GetFeatures() *FeatureFlagsService {
	if instance == nil {
		InitGrowthBook()
	}
	return instance
}

// GetMaxServers returns the BYOS seat limit based on the team's tier
func (f *FeatureFlagsService) GetMaxServers(teamID string, plan string) int {
	if f.client == nil {
		return 1 // safe fallback
	}

	// Build evaluation context with attributes
	gb, _ := f.client.WithAttributes(growthbook.Attributes{
		"team_id": teamID,
		"plan":    plan,
	})

	// Fallbacks: Hobby=1, Pro=5, Team=unlimited (e.g. 1000)
	defaultLimit := 1
	if plan == "pro" {
		defaultLimit = 5
	} else if plan == "team" {
		defaultLimit = 1000
	}

	res := gb.EvalFeature(context.Background(), "max_byos_servers")
	if res != nil && res.On {
		if val, ok := res.Value.(float64); ok {
			return int(val)
		}
	}

	return defaultLimit
}

// GetDeploymentRateLimit returns the max deployments per hour
func (f *FeatureFlagsService) GetDeploymentRateLimit(teamID string, plan string) int {
	if f.client == nil {
		return 10 // safe fallback
	}

	gb, _ := f.client.WithAttributes(growthbook.Attributes{
		"team_id": teamID,
		"plan":    plan,
	})

	// Fallbacks: Hobby=10/hr, Pro=100/hr, Team=1000/hr
	defaultLimit := 10
	if plan == "pro" {
		defaultLimit = 100
	} else if plan == "team" {
		defaultLimit = 1000
	}

	res := gb.EvalFeature(context.Background(), "max_deployments_per_hour")
	if res != nil && res.On {
		if val, ok := res.Value.(float64); ok {
			return int(val)
		}
	}

	return defaultLimit
}
