package admin

import (
	"air-quality-notifyer/internal/app/telegram/commander/api"
	tgmessage "air-quality-notifyer/internal/app/telegram/commander/message"
	"context"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

const (
	CommandAnnounce    = "/announce"
	announcementHeader = "🤖\n\n"
)

func (c *Commander) Announce(ctx context.Context, update tgbotapi.Update) {
	if !c.api.IsAdmin(update) {
		return
	}

	text, entities := announcementPayload(update.Message)
	if text == "" {
		return
	}
	text, entities = announcementMessage(text, entities)

	for _, userID := range c.service.User.GetUsersIds(ctx) {
		if ctx.Err() != nil {
			return
		}

		msg := tgbotapi.NewMessage(userID, text)
		msg.Entities = entities

		if err := c.api.Send(api.MessageConfig{Msg: msg}); err != nil {
			if err.Code == 403 {
				c.service.User.Delete(ctx, userID)
				continue
			}

			zap.L().Error("Error sending announcement", zap.Error(err), zap.Int64("userId", userID))
		}
	}
}

func IsAnnounceCommand(text string) bool {
	return tgmessage.IsCommandText(text, CommandAnnounce)
}

func announcementPayload(message *tgbotapi.Message) (string, []tgbotapi.MessageEntity) {
	text, entities, ok := tgmessage.CommandPayload(message, CommandAnnounce)
	if !ok || text == "" {
		return "", nil
	}

	return text, entities
}

func announcementMessage(text string, entities []tgbotapi.MessageEntity) (string, []tgbotapi.MessageEntity) {
	return tgmessage.Prepend(announcementHeader, text, entities)
}
