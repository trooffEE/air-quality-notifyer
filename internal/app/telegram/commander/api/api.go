package api

import (
	"air-quality-notifyer/internal/config"
	"context"
	"errors"
	"strconv"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	CommandMenuFaq = "❓ FAQ"
)

type Api struct {
	Bot   *tgbotapi.BotAPI
	cfg   config.Config
	loc   *time.Location
	cache *redis.Client
}

func NewApi(cfg config.Config, bot *tgbotapi.BotAPI, cache *redis.Client) (*Api, error) {
	loc, err := time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		return nil, err
	}
	return &Api{
		Bot:   bot,
		cfg:   cfg,
		loc:   loc,
		cache: cache,
	}, nil
}

type MessageConfig struct {
	Msg              tgbotapi.MessageConfig
	Markup           interface{}
	DisableParseMode bool
}

func (a *Api) Send(ctx context.Context, payload MessageConfig) *tgbotapi.Error {
	if !payload.DisableParseMode && len(payload.Msg.Entities) == 0 {
		payload.Msg.ParseMode = tgbotapi.ModeHTML
	}
	payload.Msg.DisableNotification = a.IsNotificationsAllowed()

	if payload.Markup != nil {
		payload.Msg.ReplyMarkup = payload.Markup
	} else {
		payload.Msg.ReplyMarkup = NewReplyKeyboard()
	}

	response, err := a.Bot.Send(payload.Msg)
	var tgError *tgbotapi.Error
	if errors.As(err, &tgError) {
		return tgError
	}
	if err != nil {
		return nil
	}

	a.trackMessage(ctx, response.Chat.ID, response.MessageID)

	return nil
}

func (a *Api) DeleteRequest(ctx context.Context, message tgbotapi.DeleteMessageConfig) error {
	_, err := a.Bot.Request(message)
	if err != nil {
		zap.L().Error("Error deleting message", zap.Error(err))
		return err
	}
	a.untrackMessage(ctx, message.ChatID, message.MessageID)
	return nil
}

func (a *Api) Delete(ctx context.Context, message *tgbotapi.Message) error {
	_, err := a.Bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
	if err != nil {
		zap.L().Error("Error deleting message", zap.Error(err))
		return err
	}
	a.untrackMessage(ctx, message.Chat.ID, message.MessageID)
	return nil
}

type EditMessageConfig struct {
	Msg    tgbotapi.EditMessageTextConfig
	Markup *tgbotapi.InlineKeyboardMarkup
}

func (a *Api) Edit(ctx context.Context, payload EditMessageConfig) error {
	payload.Msg.ParseMode = tgbotapi.ModeHTML

	if payload.Markup != nil {
		payload.Msg.ReplyMarkup = payload.Markup
	}

	_, err := a.Bot.Send(payload.Msg)
	if err != nil {
		return err
	}
	a.trackMessage(ctx, payload.Msg.ChatID, payload.Msg.MessageID)

	return nil
}

type PollConfig struct {
	Question   string
	Options    []string
	OpenPeriod int
}

func (a *Api) SendPoll(ctx context.Context, chatID int64, config PollConfig) (*tgbotapi.Message, error) {
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
	a.trackMessage(context.Background(), response.Chat.ID, response.MessageID)
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

func NewReplyKeyboard() tgbotapi.ReplyKeyboardMarkup {
	markup := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⚙️ Настройки"),
			tgbotapi.NewKeyboardButton("❓ FAQ"),
		),
	)
	markup.IsPersistent = true

	return markup
}
