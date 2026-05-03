// Core commands and handlers

package commander

import (
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/service/user/model"
	"context"
	"strconv"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

const (
	CommandStart                 = "/start"
	CommandFeedback              = "/feedback"
	InternalCommandApplyFeedback = "/internal/apply_feedback"
)

func NewCoreMessageHandlersRegistry(c *Commander) HandlersRegistry {
	return HandlersRegistry{
		CommandStart:                 c.Start,
		CommandFeedback:              c.Feedback,
		InternalCommandApplyFeedback: c.ApplyFeedback,
	}
}

func (c *Commander) ApplyFeedback(ctx context.Context, update tgbotapi.Update) {
	message := update.Message
	if message == nil {
		return
	}
	chatID := message.Chat.ID

	c.SendFeedbackToAdmin(ctx, message, message.Text, message.Entities)
	c.ConfirmFeedback(ctx, chatID)
}

func (c *Commander) Start(ctx context.Context, update tgbotapi.Update) {
	message := update.Message
	chatId, username := message.Chat.ID, message.Chat.UserName

	msg := tgbotapi.NewMessage(chatId, "Данный бот оповещает о плохом качестве воздуха в городе Кемерово.\n\nПросьба настроить уведомления, чтобы бот не беспокоил ночью! 🍵")
	if err := c.API.Send(ctx, api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending faq message", zap.Error(err))
	}

	if !c.Services.User.IsNew(ctx, chatId) {
		return
	}

	c.Services.User.Register(ctx, model.User{
		Id:       strconv.Itoa(int(chatId)),
		Username: username,
	})
}

func isFeedbackCommand(message *tgbotapi.Message) bool {
	return message != nil && api.IsCommandText(message.Text, CommandFeedback)
}

func (c *Commander) Feedback(ctx context.Context, update tgbotapi.Update) {
	message := update.Message
	if message == nil {
		return
	}
	chatID := message.Chat.ID

	c.API.SetPendingCommand(ctx, chatID, InternalCommandApplyFeedback)
	c.AskForFeedback(ctx, chatID)
}

func (c *Commander) AskForFeedback(ctx context.Context, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Напишите обратную связь для админа ниже (баги, пожелания, предложения):")
	if err := c.API.Send(ctx, api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending feedback prompt", zap.Error(err))
	}
}

func (c *Commander) ConfirmFeedback(ctx context.Context, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Спасибо! Ваша обратная связь отправлена разработчику")
	if err := c.API.Send(ctx, api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending feedback confirmation", zap.Error(err))
	}
}

func (c *Commander) SendFeedbackToAdmin(ctx context.Context, message *tgbotapi.Message, text string, entities []tgbotapi.MessageEntity) {
	adminID, ok := c.API.AdminChatID()
	if !ok || message == nil {
		return
	}

	text, entities = api.PrependText(feedbackHeader(message), text, entities)
	msg := tgbotapi.NewMessage(adminID, text)
	msg.Entities = entities

	if err := c.API.Send(ctx, api.MessageConfig{Msg: msg, DisableParseMode: len(msg.Entities) == 0}); err != nil {
		zap.L().Error("Error sending feedback message", zap.Error(err))
	}
}

func feedbackHeader(message *tgbotapi.Message) string {
	username := "@unknown"
	userID := message.Chat.ID

	if message.From != nil {
		userID = message.From.ID
		if message.From.UserName != "" {
			username = "@" + message.From.UserName
		}
	}

	return "/feedback from " + username + " (chat_id: " + strconv.FormatInt(userID, 10) + ")\n\n"
}
