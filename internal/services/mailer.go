package services

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"

	"vessl.dev/vessl/internal/views"
)

type MailerService struct {
	settingsService *EmailSettingsService
}

func NewMailerService(settingsService *EmailSettingsService) *MailerService {
	return &MailerService{
		settingsService: settingsService,
	}
}

func (s *MailerService) SendTeamEmail(ctx context.Context, teamID, templateName string, toAddress string, subject string, data any) error {
	// Fetch DB settings
	settings, err := s.settingsService.GetTeamEmailSettings(ctx, teamID)
	if err != nil {
		return fmt.Errorf("fetching team email settings: %w", err)
	}

	host := ""
	port := ""
	user := ""
	pass := ""
	from := ""

	if settings != nil {
		if settings.SMTPHost != "" {
			host = settings.SMTPHost
			port = fmt.Sprintf("%d", settings.SMTPPort)
			user = settings.SMTPUser
			pass = settings.SMTPPassword
			from = settings.SMTPFromAddress
		}
	}

	if host == "" || from == "" {
		return fmt.Errorf("SMTP configuration missing for team %s and global env", teamID)
	}

	// Render template
	var buf bytes.Buffer
	if err := views.HTMLTemplates.ExecuteTemplate(&buf, templateName, data); err != nil {
		return fmt.Errorf("executing template %s: %w", templateName, err)
	}

	// Send email
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
		"%s", toAddress, from, subject, buf.String()))

	auth := smtp.PlainAuth("", user, pass, host)
	addr := fmt.Sprintf("%s:%s", host, port)

	err = smtp.SendMail(addr, auth, from, []string{toAddress}, msg)
	if err != nil {
		return fmt.Errorf("smtp.SendMail: %w", err)
	}

	return nil
}
