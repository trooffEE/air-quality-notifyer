package commander

import (
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

const (
	announcementHeader = "🤖\n\n"
)

var (
	CommandShowUsers = "users"
	CommandPing      = "ping"
	CommandAnnounce  = "/announce"
)

func NewAdminMessageHandlersRegistry(c *Commander) HandlersRegistry {
	return HandlersRegistry{
		CommandShowUsers: c.ShowUsers,
		CommandPing:      c.Pong,
		CommandAnnounce:  c.Announce,
	}
}

func (c *Commander) Pong(ctx context.Context, update tgbotapi.Update) {
	if !c.API.IsAdmin(update) {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `pong - 🙌`)
	msg.ReplyParameters.MessageID = update.Message.MessageID

	if err := c.API.Send(ctx, api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending pong message", zap.Error(err))
	}
}

func (c *Commander) ShowUsers(ctx context.Context, update tgbotapi.Update) {
	if !c.API.IsAdmin(update) {
		return
	}

	names := c.Services.User.GetUsersNames(ctx)

	if len(names) == 0 {
		return
	}

	for index, name := range names {
		names[index] = "@" + name
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Bot Users: %d 🙌\n\n%s", len(names), strings.Join(names, ", \n")))
	if err := c.API.Send(ctx, api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending show_users", zap.Error(err))
	}
}

func (c *Commander) Announce(ctx context.Context, update tgbotapi.Update) {
	if !c.API.IsAdmin(update) {
		return
	}

	text, entities := announcementPayload(update.Message)
	if text == "" {
		return
	}
	text, entities = announcementMessage(text, entities)

	for _, userID := range c.Services.User.GetUsersIds(ctx) {
		if ctx.Err() != nil {
			return
		}

		msg := tgbotapi.NewMessage(userID, text)
		msg.Entities = entities

		if err := c.API.Send(ctx, api.MessageConfig{Msg: msg}); err != nil {
			if err.Code == 403 {
				c.Services.User.Delete(ctx, userID)
				continue
			}

			zap.L().Error("Error sending announcement", zap.Error(err), zap.Int64("userId", userID))
		}
	}
}

func IsAnnounceCommand(text string) bool {
	return api.IsCommandText(text, CommandAnnounce)
}

func announcementPayload(message *tgbotapi.Message) (string, []tgbotapi.MessageEntity) {
	text, entities, ok := api.CommandPayload(message, CommandAnnounce)
	if !ok || text == "" {
		return "", nil
	}

	return text, entities
}

func announcementMessage(text string, entities []tgbotapi.MessageEntity) (string, []tgbotapi.MessageEntity) {
	return api.PrependText(announcementHeader, text, entities)
}
