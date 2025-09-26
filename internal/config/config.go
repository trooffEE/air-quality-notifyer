package config

import (
	"os"
)

type ApplicationConfig struct {
	TelegramToken   string
	HttpServerPort  string
	WebhookHost     string
	Development     bool
	AdminTelegramId string
}

// TODO move it to /internal/app directory
func NewApplicationConfig() ApplicationConfig {
	var config = ApplicationConfig{
		TelegramToken:   os.Getenv("TELEGRAM_SECRET"),
		WebhookHost:     os.Getenv("WEBHOOK_HOST"),
		HttpServerPort:  os.Getenv("WEBHOOK_PORT"),
		AdminTelegramId: os.Getenv("ADMIN_TELEGRAM_ID"),
		Development:     os.Getenv("DEVELOPMENT") == "1",
	}

	return config
}
