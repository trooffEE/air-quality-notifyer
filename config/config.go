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
	Development        bool
	TestTelegramChatID string
	// TODO Убрать
	TestTelegramChatID2 string
	// TODO Убрать
	TestTelegramChatID3 string
}

func InitConfig() ApplicationConfig {
	var config ApplicationConfig = ApplicationConfig{
		TelegramToken:       os.Getenv("TELEGRAM_SECRET"),
		WebhookHost:         os.Getenv("WEBHOOK_HOST"),
		WebhookPort:         os.Getenv("WEBHOOK_PORT"),
		TestTelegramChatID:  os.Getenv("TEST_CHAT_ID"),
		TestTelegramChatID2: os.Getenv("TEST_CHAT_ID2"),
		TestTelegramChatID3: os.Getenv("TEST_CHAT_ID3"),
		Development:         os.Getenv("DEVELOPMENT") == "true",
	}

	return config
}

func (a *ApplicationConfig) GetTestTelegramChatID() int64 {
	testTelegramChatId, err := strconv.Atoi(a.TestTelegramChatID)
	if err != nil {
		log.Panic("Provide me \"Test Telegram Chat Id\" or don't use me if you are not in dev-mode!")
	}
	return int64(testTelegramChatId)
}

// TODO Убрать
func (a *ApplicationConfig) GetTestTelegramChatID2() int64 {
	testTelegramChatId, err := strconv.Atoi(a.TestTelegramChatID2)
	if err != nil {
		log.Panic("Provide me \"Test Telegram Chat Id\" or don't use me if you are not in dev-mode!")
	}
	return int64(testTelegramChatId)
}

// TODO Убрать
func (a *ApplicationConfig) GetTestTelegramChatID3() int64 {
	testTelegramChatId, err := strconv.Atoi(a.TestTelegramChatID3)
	if err != nil {
		log.Panic("Provide me \"Test Telegram Chat Id\" or don't use me if you are not in dev-mode!")
	}
	return int64(testTelegramChatId)
}
