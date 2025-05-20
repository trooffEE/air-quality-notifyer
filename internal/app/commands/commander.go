package commands

import (
	"air-quality-notifyer/internal/config"
	"errors"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type Commander struct {
	bot *tgbotapi.BotAPI //TOOD think about common interface so that Telegram, WhatsApp, VK can be used
	cfg config.ApplicationConfig
}

func NewCommander(bot *tgbotapi.BotAPI, cfg config.ApplicationConfig) *Commander {
	return &Commander{
		bot: bot,
		cfg: cfg,
	}
}

type SendPayload struct {
	Msg                 tgbotapi.MessageConfig
	DisableNotification bool
}

func (c *Commander) Send(payload SendPayload) *tgbotapi.Error {
	payload.Msg.ParseMode = tgbotapi.ModeHTML
	payload.Msg.DisableNotification = payload.DisableNotification

	_, err := c.bot.Send(payload.Msg)
	var tgError *tgbotapi.Error
	if errors.As(err, &tgError) {
		return tgError
	}

	return nil
}
