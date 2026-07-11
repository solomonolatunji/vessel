package services

import (
	"fmt"
	"log"
	"time"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/billing/meterevent"
)

type MeteringService interface {
	RecordUsage(teamID string, deployments int, containerHours int, bandwidthGB int) error
	ReportToStripe(customerID string, eventName string, value int) error
}

type DefaultMeteringService struct{}

func NewMeteringService() *DefaultMeteringService {
	return &DefaultMeteringService{}
}

// RecordUsage stores usage internally in our database and conditionally pushes it to external billing providers
func (s *DefaultMeteringService) RecordUsage(teamID string, deployments int, containerHours int, bandwidthGB int) error {
	// TODO: Save the raw metrics to Postgres table `cloud_usage_logs`

	// Example logic:
	// Find customer subscription details
	// customerStripeID := repo.GetStripeCustomer(teamID)
	customerStripeID := "cus_12345" // mock

	// Report metric events to Stripe if they are on a metered plan
	if deployments > 0 {
		err := s.ReportToStripe(customerStripeID, "deployments_meter", deployments)
		if err != nil {
			log.Printf("Failed to report deployment usage to Stripe for %s: %v", teamID, err)
		}
	}
	
	if containerHours > 0 {
		err := s.ReportToStripe(customerStripeID, "container_hours_meter", containerHours)
		if err != nil {
			log.Printf("Failed to report container usage to Stripe for %s: %v", teamID, err)
		}
	}

	if bandwidthGB > 0 {
		err := s.ReportToStripe(customerStripeID, "bandwidth_gb_meter", bandwidthGB)
		if err != nil {
			log.Printf("Failed to report bandwidth usage to Stripe for %s: %v", teamID, err)
		}
	}

	return nil
}

// ReportToStripe pushes a single usage record to Stripe's v2 metered billing API
func (s *DefaultMeteringService) ReportToStripe(customerID string, eventName string, value int) error {
	// Create a new metering event in Stripe
	params := &stripe.BillingMeterEventParams{
		EventName: stripe.String(eventName),
		Payload: map[string]string{
			"stripe_customer_id": customerID,
			"value":              fmt.Sprintf("%d", value),
		},
		Timestamp: stripe.Int64(time.Now().Unix()),
	}

	_, err := meterevent.New(params)
	if err != nil {
		return err
	}

	return nil
}
