package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ApplicationConfig struct {
	TelegramToken   string
	HttpServerPort  string
	WebhookHost     string
	Development     bool
	AdminTelegramId string
}

func NewApplicationConfig() ApplicationConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var config = ApplicationConfig{
		TelegramToken:   os.Getenv("TELEGRAM_SECRET"),
		WebhookHost:     os.Getenv("WEBHOOK_HOST"),
		HttpServerPort:  os.Getenv("WEBHOOK_PORT"), // TODO rename env
		AdminTelegramId: os.Getenv("ADMIN_TELEGRAM_ID"),
		Development:     os.Getenv("DEVELOPMENT") == "true",
	}

	return config
}
