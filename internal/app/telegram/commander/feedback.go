package commander

import (
	"context"
	"strconv"
	"time"

	"air-quality-notifyer/internal/app/telegram/commander/api"
	tgmessage "air-quality-notifyer/internal/app/telegram/commander/message"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

const (
	CommandFeedback    = "feedback"
	feedbackPendingTTL = 24 * time.Hour
)

func isFeedbackCommand(message *tgbotapi.Message) bool {
	return tgmessage.IsCommand(message, CommandFeedback)
}

func (c *Commander) Feedback(update tgbotapi.Update) {
	message := update.Message
	if message == nil {
		return
	}

	text, entities, ok := tgmessage.CommandPayload(message, CommandFeedback)
	if !ok {
		return
	}

	if text == "" {
		c.SetFeedbackPending(message.Chat.ID)
		c.AskForFeedback(message.Chat.ID)
		return
	}

	c.DeleteFeedbackPending(message.Chat.ID)
	c.SendFeedbackToAdmin(message, text, entities)
	c.ConfirmFeedback(message.Chat.ID)
}

func (c *Commander) HandlePendingFeedback(update tgbotapi.Update) bool {
	message := update.Message
	if message == nil {
		return false
	}

	if message.IsCommand() || api.IsMenuButton(message.Text) || !c.ConsumeFeedbackPending(message.Chat.ID) {
		return false
	}

	c.SendFeedbackToAdmin(message, message.Text, message.Entities)
	c.ConfirmFeedback(message.Chat.ID)

	return true
}

func (c *Commander) SetFeedbackPending(chatID int64) {
	if c.Services.Cache == nil {
		zap.L().Error("feedback cache is not configured")
		return
	}

	err := c.Services.Cache.Set(context.Background(), feedbackPendingKey(chatID), "1", feedbackPendingTTL).Err()
	if err != nil {
		zap.L().Error("failed to set pending feedback state", zap.Error(err), zap.Int64("chatId", chatID))
	}
}

func (c *Commander) ConsumeFeedbackPending(chatID int64) bool {
	if c.Services.Cache == nil {
		zap.L().Error("feedback cache is not configured")
		return false
	}

	deleted, err := c.Services.Cache.Del(context.Background(), feedbackPendingKey(chatID)).Result()
	if err != nil {
		zap.L().Error("failed to consume pending feedback state", zap.Error(err), zap.Int64("chatId", chatID))
		return false
	}

	return deleted > 0
}

func (c *Commander) DeleteFeedbackPending(chatID int64) {
	if c.Services.Cache == nil {
		zap.L().Error("feedback cache is not configured")
		return
	}

	err := c.Services.Cache.Del(context.Background(), feedbackPendingKey(chatID)).Err()
	if err != nil {
		zap.L().Error("failed to delete pending feedback state", zap.Error(err), zap.Int64("chatId", chatID))
	}
}

func (c *Commander) AskForFeedback(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Напишите обратную связь для админа ниже (баги, пожелания, предложения):")
	if err := c.API.Send(api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending feedback prompt", zap.Error(err))
	}
}

func (c *Commander) ConfirmFeedback(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Спасибо! Ваша обратная связь отправлена разработчику")
	if err := c.API.Send(api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending feedback confirmation", zap.Error(err))
	}
}

func (c *Commander) SendFeedbackToAdmin(message *tgbotapi.Message, text string, entities []tgbotapi.MessageEntity) {
	adminID, ok := c.API.AdminChatID()
	if !ok || message == nil {
		return
	}

	text, entities = tgmessage.Prepend(feedbackHeader(message), text, entities)
	msg := tgbotapi.NewMessage(adminID, text)
	msg.Entities = entities

	if err := c.API.Send(api.MessageConfig{Msg: msg, DisableParseMode: len(msg.Entities) == 0}); err != nil {
		zap.L().Error("Error sending feedback message", zap.Error(err))
	}
}

func feedbackPendingKey(chatID int64) string {
	return "telegram:feedback:pending:" + strconv.FormatInt(chatID, 10)
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
