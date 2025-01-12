package commands

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type Commander struct {
	bot *tgbotapi.BotAPI
}

func NewCommander(bot *tgbotapi.BotAPI) *Commander {
	return &Commander{
		bot,
	}
}
