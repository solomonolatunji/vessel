package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	emails "vessel.dev/vessel/internal/cloud/views/emails"
)

// MailerService sends transactional emails via AWS SES.
type MailerService struct {
	awsClient *ses.Client
	fromEmail string
}

// NewMailerService initialises the SES client from environment variables.
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

// renderTemplate renders an embedded email template with the given data.
func (s *MailerService) renderTemplate(name string, data any) (string, error) {
	tmplContent, err := emails.Templates.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("reading template %s: %w", name, err)
	}

	tmpl, err := template.New(name).Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("parsing template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template %s: %w", name, err)
	}

	return buf.String(), nil
}

// SendWelcomeEmail sends a welcome + verify-email link to a newly registered user.
func (s *MailerService) SendWelcomeEmail(ctx context.Context, toAddress, name, verifyURL string) error {
	log.Printf("[SES] Sending Welcome Email to %s", toAddress)

	htmlBody, err := s.renderTemplate("welcome.tmpl", map[string]string{
		"Name":      name,
		"VerifyURL": verifyURL,
	})
	if err != nil {
		return err
	}

	return s.sendEmail(ctx, toAddress, "Welcome to Vessel Cloud", htmlBody)
}

// SendOTPResetEmail sends a password-reset OTP to the user's email address.
func (s *MailerService) SendOTPResetEmail(ctx context.Context, toAddress, name, otpCode, expiresIn string) error {
	log.Printf("[SES] Sending OTP Reset Email to %s", toAddress)

	htmlBody, err := s.renderTemplate("otp_reset.tmpl", map[string]string{
		"Name":      name,
		"OTPCode":   otpCode,
		"ExpiresIn": expiresIn,
	})
	if err != nil {
		return err
	}

	return s.sendEmail(ctx, toAddress, "Your Vessel Cloud Password Reset Code", htmlBody)
}

// SendBillingAlert sends a payment-failure alert email.
func (s *MailerService) SendBillingAlert(ctx context.Context, toAddress string, amount float64, billingURL string) error {
	log.Printf("[SES] Sending Billing Alert to %s (Amount: %.2f)", toAddress, amount)

	htmlBody, err := s.renderTemplate("billing_alert.tmpl", map[string]interface{}{
		"Amount":     amount,
		"BillingURL": billingURL,
	})
	if err != nil {
		return err
	}

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
