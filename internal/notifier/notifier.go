package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/repositories"
)

type NotifierService struct {
	repo repositories.NotificationRepository
}

func NewNotifierService(repo repositories.NotificationRepository) *NotifierService {
	return &NotifierService{repo: repo}
}

func (n *NotifierService) Dispatch(event *models.NotificationEvent) {
	go func() {
		if err := n.Send(event); err != nil {
			log.Printf("[Notifier] Failed to dispatch event '%s': %v", event.Title, err)
		}
	}()
}

func (n *NotifierService) Send(event *models.NotificationEvent) error {
	integ, err := n.repo.GetIntegration(context.Background())
	if err != nil || integ == nil {
		return fmt.Errorf("could not load notification integrations: %w", err)
	}

	var pref *models.ProjectNotificationPref
	if event.ProjectID != "" {
		p, err := n.repo.GetProjectPref(context.Background(), event.ProjectID)
		if err == nil {
			pref = p
		}
	}

	if pref == nil || pref.EmailEnabled {
		if integ.SMTPEnabled && integ.SMTPHost != "" {
			_ = n.sendSMTP(integ, event)
		} else if integ.ResendEnabled && integ.ResendAPIKey != "" {
			_ = n.sendResend(integ, event)
		}
	}

	if (pref == nil || pref.SlackEnabled) && integ.SlackEnabled && integ.SlackWebhookURL != "" {
		_ = n.sendSlack(integ.SlackWebhookURL, event)
	}

	if (pref == nil || pref.DiscordEnabled) && integ.DiscordEnabled && integ.DiscordWebhookURL != "" {
		_ = n.sendDiscord(integ.DiscordWebhookURL, integ.DiscordPingEnabled, event)
	}

	if (pref == nil || pref.TelegramEnabled) && integ.TelegramEnabled && integ.TelegramBotToken != "" && integ.TelegramChatID != "" {
		_ = n.sendTelegram(integ.TelegramBotToken, integ.TelegramChatID, event)
	}

	if (pref == nil || pref.PushoverEnabled) && integ.PushoverEnabled && integ.PushoverUserKey != "" && integ.PushoverAPIToken != "" {
		_ = n.sendPushover(integ.PushoverUserKey, integ.PushoverAPIToken, event)
	}

	if (pref == nil || pref.WebhookEnabled) && integ.WebhookEnabled && integ.WebhookURL != "" {
		_ = n.sendWebhook(integ.WebhookURL, event)
	}

	return nil
}

func (n *NotifierService) sendSMTP(integ *models.NotificationIntegration, event *models.NotificationEvent) error {
	var auth smtp.Auth
	if integ.SMTPUser != "" && integ.SMTPPassword != "" {
		auth = smtp.PlainAuth("", integ.SMTPUser, integ.SMTPPassword, integ.SMTPHost)
	}
	addr := fmt.Sprintf("%s:%d", integ.SMTPHost, integ.SMTPPort)
	fromAddr := integ.SMTPFromAddress
	if fromAddr == "" {
		fromAddr = integ.SMTPUser
	}
	if fromAddr == "" {
		return fmt.Errorf("SMTP from address is required")
	}
	toAddr := integ.SMTPUser
	if toAddr == "" {
		toAddr = fromAddr
	}
	to := []string{toAddr}

	fromHeader := fromAddr
	if integ.SMTPFromName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", integ.SMTPFromName, fromAddr)
	}

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: [Vessel %s] %s\r\n\r\n%s\r\n\r\nURL: %s\r\n",
		fromHeader, toAddr, strings.ToUpper(event.Level), event.Title, event.Message, event.URL))

	return smtp.SendMail(addr, auth, fromAddr, to, msg)
}

func (n *NotifierService) sendResend(integ *models.NotificationIntegration, event *models.NotificationEvent) error {
	fromStr := "Vessel Notifications <alerts@vessel.dev>"
	if integ.SMTPFromAddress != "" {
		if integ.SMTPFromName != "" {
			fromStr = fmt.Sprintf("%s <%s>", integ.SMTPFromName, integ.SMTPFromAddress)
		} else {
			fromStr = integ.SMTPFromAddress
		}
	}
	payload := map[string]interface{}{
		"from":    fromStr,
		"to":      []string{"admin@localhost"},
		"subject": fmt.Sprintf("[Vessel] %s", event.Title),
		"html":    fmt.Sprintf("<p><strong>%s</strong></p><p>%s</p><p><a href=\"%s\">View in Dashboard</a></p>", event.Title, event.Message, event.URL),
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+integ.ResendAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (n *NotifierService) sendSlack(webhookURL string, event *models.NotificationEvent) error {
	payload := map[string]string{
		"text": fmt.Sprintf("🚀 *%s*\n%s\n<%s|View Details>", event.Title, event.Message, event.URL),
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (n *NotifierService) sendDiscord(webhookURL string, ping bool, event *models.NotificationEvent) error {
	content := fmt.Sprintf("**%s**\n%s\n[View Details](%s)", event.Title, event.Message, event.URL)
	if ping && event.Level == "error" {
		content = "@everyone " + content
	}
	payload := map[string]string{"content": content}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (n *NotifierService) sendTelegram(botToken, chatID string, event *models.NotificationEvent) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	text := fmt.Sprintf("🛰 *%s*\n%s\n%s", event.Title, event.Message, event.URL)
	payload := map[string]string{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (n *NotifierService) sendPushover(userKey, apiToken string, event *models.NotificationEvent) error {
	values := url.Values{
		"token":   {apiToken},
		"user":    {userKey},
		"title":   {event.Title},
		"message": {fmt.Sprintf("%s\n\n%s", event.Message, event.URL)},
	}
	resp, err := http.PostForm("https://api.pushover.net/1/messages.json", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (n *NotifierService) sendWebhook(webhookURL string, event *models.NotificationEvent) error {
	body, _ := json.Marshal(event)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
