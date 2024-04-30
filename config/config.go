package config

import (
	"os"
)

type ApplicationConfig struct {
	TelegramToken string
	WebhookPort   string
	WebhookHost   string
}

func InitConfig() ApplicationConfig {
	return ApplicationConfig{
		TelegramToken: os.Getenv("TELEGRAM_SECRET"),
		WebhookHost:   os.Getenv("WEBHOOK_HOST"),
		WebhookPort:   os.Getenv("WEBHOOK_PORT"),
	}
}
