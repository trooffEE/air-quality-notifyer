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
	Bot *tgbotapi.BotAPI
	cfg config.Config
	loc *time.Location
}

type Interface interface {
	Send(payload MessageConfig) *tgbotapi.Error
	SendPoll(chatID int64, config PollConfig) (*tgbotapi.Message, error)
	AdminChatID() (int64, bool)
	Delete(update *tgbotapi.Message) error
	DeleteRequest(message tgbotapi.DeleteMessageConfig) error
	Edit(payload EditMessageConfig) error
	IsAdmin(update tgbotapi.Update) bool
	IsNotificationsAllowed() bool
	MenuBack(update tgbotapi.Update)
	MenuFaq(update tgbotapi.Update)
}

func NewApi(cfg config.Config, bot *tgbotapi.BotAPI) (*Api, error) {
	loc, err := time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		return nil, err
	}
	return &Api{
		Bot: bot,
		cfg: cfg,
		loc: loc,
	}, nil
}

type MessageConfig struct {
	Msg              tgbotapi.MessageConfig
	Markup           interface{}
	DisableParseMode bool
}

func (a *Api) Send(payload MessageConfig) *tgbotapi.Error {
	if !payload.DisableParseMode && len(payload.Msg.Entities) == 0 {
		payload.Msg.ParseMode = tgbotapi.ModeHTML
	}
	payload.Msg.DisableNotification = a.IsNotificationsAllowed()

	if payload.Markup != nil {
		payload.Msg.ReplyMarkup = payload.Markup
	} else {
		payload.Msg.ReplyMarkup = NewReplyKeyboard()
	}

	_, err := a.Bot.Send(payload.Msg)
	var tgError *tgbotapi.Error
	if errors.As(err, &tgError) {
		return tgError
	}

	return nil
}

func (a *Api) DeleteRequest(message tgbotapi.DeleteMessageConfig) error {
	_, err := a.Bot.Request(message)
	if err != nil {
		zap.L().Error("Error deleting message", zap.Error(err))
		return err
	}
	return nil
}

func (a *Api) Delete(message *tgbotapi.Message) error {
	_, err := a.Bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
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

func (a *Api) Edit(payload EditMessageConfig) error {
	payload.Msg.ParseMode = tgbotapi.ModeHTML

	if payload.Markup != nil {
		payload.Msg.ReplyMarkup = payload.Markup
	}

	_, err := a.Bot.Send(payload.Msg)
	if err != nil {
		return err
	}

	return nil
}

type PollConfig struct {
	Question   string
	Options    []string
	OpenPeriod int
}

func (a *Api) SendPoll(chatID int64, config PollConfig) (*tgbotapi.Message, error) {
	var options []tgbotapi.InputPollOption
	for _, option := range config.Options {
		options = append(options, tgbotapi.NewPollOption(option))
	}

	poll := tgbotapi.NewPoll(chatID, config.Question, options...)
	poll.AllowsMultipleAnswers = true
	poll.ProtectContent = true
	poll.Type = "regular"
	if config.OpenPeriod != 0 {
		poll.OpenPeriod = config.OpenPeriod
	}

	response, err := a.Bot.Send(poll)
	if err != nil {
		zap.L().Error("poll: error sending error", zap.Error(err))
		return nil, err
	}
	return &response, nil
}

func (a *Api) IsAdmin(update tgbotapi.Update) bool {
	adminId, ok := a.AdminChatID()
	return ok && update.Message != nil && adminId == update.Message.Chat.ID
}

func (a *Api) AdminChatID() (int64, bool) {
	adminId, err := strconv.ParseInt(a.cfg.App.AdminTelegramId, 10, 64)
	if err != nil {
		zap.L().Error("conversion error", zap.Error(err))
		return 0, false
	}
	return adminId, true
}

func (a *Api) IsNotificationsAllowed() bool {
	h := time.Now().In(a.loc).Hour()
	return h < 8 && h >= 0
}
