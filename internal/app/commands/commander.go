package commands

import (
	"air-quality-notifyer/internal/config"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type Commander struct {
	bot *tgbotapi.BotAPI
	cfg config.ApplicationConfig
}

func NewCommander(bot *tgbotapi.BotAPI, cfg config.ApplicationConfig) *Commander {
	return &Commander{
		bot: bot,
		cfg: cfg,
	}
}
