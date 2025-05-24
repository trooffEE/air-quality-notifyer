package commander

import (
	"air-quality-notifyer/internal/app/menu"
	"air-quality-notifyer/internal/config"
	"errors"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
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

type SendPayload struct {
	Msg                 tgbotapi.MessageConfig
	ReplyMarkup         interface{}
	DisableNotification bool
}

func (c *Commander) Delete(message *tgbotapi.Message) {
	_, err := c.bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
	if err != nil {
		zap.L().Error("Error deleting message", zap.Error(err))
	}
}

func (c *Commander) Send(payload SendPayload) *tgbotapi.Error {
	payload.Msg.ParseMode = tgbotapi.ModeHTML
	payload.Msg.DisableNotification = payload.DisableNotification

	if payload.ReplyMarkup != nil {
		payload.Msg.ReplyMarkup = payload.ReplyMarkup
	} else {
		payload.Msg.ReplyMarkup = menu.NewTelegramMainMenu()
	}

	_, err := c.bot.Send(payload.Msg)
	var tgError *tgbotapi.Error
	if errors.As(err, &tgError) {
		return tgError
	}

	return nil
}
