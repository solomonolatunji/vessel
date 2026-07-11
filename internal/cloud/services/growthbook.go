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

func (f *FeatureFlagsService) GetMaxServers(teamID string, plan string) int {
	if f.client == nil {
		return 1
	}

	gb, _ := f.client.WithAttributes(growthbook.Attributes{
		"team_id": teamID,
		"plan":    plan,
	})

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

func (f *FeatureFlagsService) GetDeploymentRateLimit(teamID string, plan string) int {
	if f.client == nil {
		return 10
	}

	gb, _ := f.client.WithAttributes(growthbook.Attributes{
		"team_id": teamID,
		"plan":    plan,
	})

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
