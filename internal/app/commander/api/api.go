package api

import (
	"air-quality-notifyer/internal/config"
	"errors"
	"strconv"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type Api struct {
	bot *tgbotapi.BotAPI
	cfg config.Config
	loc *time.Location
}

type Interface interface {
	Send(payload MessageConfig) *tgbotapi.Error
	Delete(update *tgbotapi.Message) error
	Edit(payload EditMessageConfig) *tgbotapi.Error
	IsAdmin(update tgbotapi.Update) bool
	IsNotificationsAllowed() bool
	MenuBack(update tgbotapi.Update)
	MenuFaq(update tgbotapi.Update)
}

func NewApi(bot *tgbotapi.BotAPI, cfg config.Config) (Interface, error) {
	loc, err := time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		return nil, err
	}
	return &Api{
		bot: bot,
		cfg: cfg,
		loc: loc,
	}, nil
}

type MessageConfig struct {
	Msg    tgbotapi.MessageConfig
	Markup interface{}
}

func (a *Api) Send(payload MessageConfig) *tgbotapi.Error {
	payload.Msg.ParseMode = tgbotapi.ModeHTML
	payload.Msg.DisableNotification = a.IsNotificationsAllowed()

	if payload.Markup != nil {
		payload.Msg.ReplyMarkup = payload.Markup
	} else {
		payload.Msg.ReplyMarkup = NewReplyKeyboard()
	}

	_, err := a.bot.Send(payload.Msg)
	var tgError *tgbotapi.Error
	if errors.As(err, &tgError) {
		return tgError
	}

	return nil
}

func (a *Api) Delete(message *tgbotapi.Message) error {
	_, err := a.bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
	if err != nil {
		zap.L().Error("Error deleting message", zap.Error(err))
		return err
	}
	return nil
}

type EditMessageConfig struct {
	Msg    tgbotapi.EditMessageTextConfig
	Markup *tgbotapi.InlineKeyboardMarkup
}

func (a *Api) Edit(payload EditMessageConfig) *tgbotapi.Error {
	payload.Msg.ParseMode = tgbotapi.ModeHTML

	if payload.Markup != nil {
		payload.Msg.ReplyMarkup = payload.Markup
	}

	_, err := a.bot.Send(payload.Msg)
	var tgError *tgbotapi.Error
	if errors.As(err, &tgError) {
		return tgError
	}

	return nil
}

func (a *Api) IsAdmin(update tgbotapi.Update) bool {
	adminId, err := strconv.Atoi(a.cfg.App.AdminTelegramId)
	if err != nil {
		zap.L().Error("conversion error", zap.Error(err))
		return false
	}
	return int64(adminId) == update.Message.Chat.ID
}

func (a *Api) IsNotificationsAllowed() bool {
	h := time.Now().In(a.loc).Hour()
	return h < 8 && h >= 0
}
