package config

import (
	"log"
	"os"
	"strconv"
)

type ApplicationConfig struct {
	TelegramToken      string
	WebhookPort        string
	WebhookHost        string
	TestTelegramChatID string
}

func InitConfig() ApplicationConfig {
	return ApplicationConfig{
		TelegramToken:      os.Getenv("TELEGRAM_SECRET"),
		WebhookHost:        os.Getenv("WEBHOOK_HOST"),
		WebhookPort:        os.Getenv("WEBHOOK_PORT"),
		TestTelegramChatID: os.Getenv("TEST_CHAT_ID"),
	}
}

func (a *ApplicationConfig) GetTestTelegramChatID() int64 {
	testTelegramChatId, err := strconv.Atoi(a.TestTelegramChatID)
	if err != nil {
		log.Panic("Provide me \"Test Telegram Chat Id\" or don't use me if you are not in dev-mode!")
	}
	return int64(testTelegramChatId)
}
