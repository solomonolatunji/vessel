package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type MailerService struct {
	awsClient *ses.Client
	fromEmail string
}

func NewMailerService(ctx context.Context) (*MailerService, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	fromEmail := os.Getenv("SES_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "noreply@vessel.dev"
	}

	var cfg aws.Config
	var err error

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if accessKey != "" && secretKey != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := ses.NewFromConfig(cfg)

	return &MailerService{
		awsClient: client,
		fromEmail: fromEmail,
	}, nil
}

func (s *MailerService) SendWelcomeEmail(ctx context.Context, toAddress string, name string) error {
	log.Printf("[SES] Sending Welcome Email to %s", toAddress)

	htmlBody := fmt.Sprintf(`
		<h1>Welcome to Vessel Cloud, %s!</h1>
		<p>We're thrilled to have you. You can now deploy apps to your own connected servers or use our managed regions.</p>
		<p>Happy shipping!</p>
	`, name)

	return s.sendEmail(ctx, toAddress, "Welcome to Vessel Cloud", htmlBody)
}

func (s *MailerService) SendBillingAlert(ctx context.Context, toAddress string, amount float64) error {
	log.Printf("[SES] Sending Billing Alert to %s (Amount: %.2f)", toAddress, amount)

	htmlBody := fmt.Sprintf(`
		<h1>Payment Failed</h1>
		<p>Your recent payment of $%.2f failed. Please update your payment method to avoid service interruption.</p>
	`, amount)

	return s.sendEmail(ctx, toAddress, "Action Required: Payment Failed", htmlBody)
}

func (s *MailerService) sendEmail(ctx context.Context, toAddress, subject, htmlBody string) error {
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{toAddress},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(s.fromEmail),
	}

	_, err := s.awsClient.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send SES email: %w", err)
	}

	return nil
}
