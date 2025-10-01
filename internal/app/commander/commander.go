package commander

import (
	"air-quality-notifyer/internal/app/menu"
	"air-quality-notifyer/internal/config"
	"errors"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type Commander struct {
	bot *tgbotapi.BotAPI
	cfg config.Config
	loc *time.Location
}

func NewCommander(bot *tgbotapi.BotAPI, cfg config.Config) *Commander {
	loc, _ := time.LoadLocation("Asia/Novosibirsk")
	return &Commander{
		bot: bot,
		cfg: cfg,
		loc: loc,
	}
}

type Payload struct {
	Msg         tgbotapi.MessageConfig
	ReplyMarkup interface{}
}

func (c *Commander) Delete(message *tgbotapi.Message) {
	_, err := c.bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
	if err != nil {
		zap.L().Error("Error deleting message", zap.Error(err))
	}
}

func (c *Commander) Send(payload Payload) *tgbotapi.Error {
	payload.Msg.ParseMode = tgbotapi.ModeHTML
	payload.Msg.DisableNotification = c.isNotificationsAllowed()

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

func (c *Commander) isNotificationsAllowed() bool {
	h := time.Now().In(c.loc).Hour()
	return h < 8 && h >= 0
}
