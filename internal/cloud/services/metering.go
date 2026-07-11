package services

import (
	"fmt"
	"log"
	"time"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/billing/meterevent"
	"vessel.dev/vessel/internal/cloud/models"
	"vessel.dev/vessel/internal/cloud/repos"
)

type MeteringService interface {
	RecordUsage(teamID uint, deployments int, containerHours int, bandwidthGB int) error
}

type DefaultMeteringService struct {
	repo repos.CloudRepo
}

func NewMeteringService(repo repos.CloudRepo) *DefaultMeteringService {
	return &DefaultMeteringService{repo: repo}
}

func (s *DefaultMeteringService) RecordUsage(teamID uint, deployments int, containerHours int, bandwidthGB int) error {
	err := s.repo.LogUsage(&models.CloudUsageLog{
		TeamID:         teamID,
		Deployments:    deployments,
		ContainerHours: containerHours,
		BandwidthGB:    bandwidthGB,
		ReportedAt:     time.Now(),
	})
	if err != nil {
		log.Printf("Failed to record usage in DB: %v", err)
	}

	team, err := s.repo.GetTeamByID(teamID)
	if err != nil || team == nil {
		return err
	}

	if team.StripeCustomerID == "" {
		return nil
	}

	if deployments > 0 {
		if err := s.reportToStripe(team.StripeCustomerID, "deployments_meter", deployments); err != nil {
			log.Printf("Failed to report deployment usage to Stripe for %d: %v", teamID, err)
		}
	}

	if containerHours > 0 {
		if err := s.reportToStripe(team.StripeCustomerID, "container_hours_meter", containerHours); err != nil {
			log.Printf("Failed to report container usage to Stripe for %d: %v", teamID, err)
		}
	}

	if bandwidthGB > 0 {
		if err := s.reportToStripe(team.StripeCustomerID, "bandwidth_gb_meter", bandwidthGB); err != nil {
			log.Printf("Failed to report bandwidth usage to Stripe for %d: %v", teamID, err)
		}
	}

	return nil
}

func (s *DefaultMeteringService) reportToStripe(customerID string, eventName string, value int) error {
	params := &stripe.BillingMeterEventParams{
		EventName: stripe.String(eventName),
		Payload: map[string]string{
			"stripe_customer_id": customerID,
			"value":              fmt.Sprintf("%d", value),
		},
		Timestamp: stripe.Int64(time.Now().Unix()),
	}

	_, err := meterevent.New(params)
	return err
}
