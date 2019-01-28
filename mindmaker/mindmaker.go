package mindmaker

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Config configurations for mindmaker bot
// Token secret token string obtained from Telegram
// SSL indicator whether secured endpoint is being used
// Cert string for location of cert.pem
// Key string for location of key.pem
// WebhookAddr string for webhook address e.g. "https://mybot.herokuapp.com"
// ListeningPath string for path that handles updates e.g. "/"
type Config struct {
	Token         string
	SSL           bool
	Cert          string
	Key           string
	WebhookAddr   string
	ListeningPath string
}

// Initialize initialize bot webhook to handle messages from Telegram
func Initialize(config Config) error {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return err
	}

	var webhook tgbotapi.WebhookConfig
	if config.SSL {
		webhook = tgbotapi.NewWebhookWithCert(config.WebhookAddr+config.ListeningPath, config.Cert)
	} else {
		webhook = tgbotapi.NewWebhook(config.WebhookAddr + config.ListeningPath)
	}
	_, err = bot.SetWebhook(webhook)
	if err != nil {
		return err
	}

	updates := bot.ListenForWebhook(config.ListeningPath)
	go HandleUpdates(updates)

	return nil
}
