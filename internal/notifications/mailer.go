package notifications

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"

	"vessl.dev/vessl/internal/services"
)

type MailerService struct {
	globalSettingsService *services.SettingsService
}

func NewMailerService(globalSettings *services.SettingsService) (*MailerService, error) {
	if err := LoadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load email templates: %w", err)
	}
	return &MailerService{
		globalSettingsService: globalSettings,
	}, nil
}

func (s *MailerService) SendSystemEmail(ctx context.Context, templateName string, toAddress string, subject string, data any) error {
	settings, err := s.globalSettingsService.GetSettings(ctx)
	if err != nil {
		return fmt.Errorf("fetching server settings: %w", err)
	}

	if settings == nil || !settings.SMTPEnabled {
		return fmt.Errorf("global SMTP is not configured or enabled")
	}

	host := settings.SMTPHost
	port := fmt.Sprintf("%d", settings.SMTPPort)
	user := settings.SMTPUser
	pass := settings.SMTPPassword
	from := settings.SMTPFromAddress

	if host == "" || from == "" {
		return fmt.Errorf("global SMTP configuration missing")
	}

	var buf bytes.Buffer
	if err := HTMLTemplates.ExecuteTemplate(&buf, templateName, data); err != nil {
		return fmt.Errorf("executing template %s: %w", templateName, err)
	}

	msg := fmt.Appendf(nil, "To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
		"%s", toAddress, from, subject, buf.String())

	auth := smtp.PlainAuth("", user, pass, host)
	addr := fmt.Sprintf("%s:%s", host, port)

	err = smtp.SendMail(addr, auth, from, []string{toAddress}, msg)
	if err != nil {
		return fmt.Errorf("smtp.SendMail: %w", err)
	}

	return nil
}
