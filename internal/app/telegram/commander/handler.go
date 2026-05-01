package commander

import (
	"context"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message != nil {
		c.handleMessageUpdate(ctx, update)
	} else if update.CallbackQuery != nil {
		c.handleCallbackQueryUpdate(ctx, update)
	} else if update.Poll != nil {
		c.handlePollUpdate(ctx, update)
	}
}

func (c *Commander) handleMessageUpdate(ctx context.Context, update tgbotapi.Update) {
	message := update.Message
	if message == nil {
		return
	}

	logClientMessage(message)

	if commandName, exists := c.PendingCommand(ctx, message.Chat.ID); exists {
		c.DeletePendingCommand(ctx, update.Message.Chat.ID)
		command, ok := c.messageHandlersRegistry[commandName]
		if !ok {
			zap.L().Error("pending command not found", zap.String("commandName", commandName))
			return
		}
		command(ctx, update)
		return
	}

	if command, ok := c.messageHandlersRegistry[message.Text]; ok {
		command(ctx, update)
	}

	if c.API.IsMenuButton(message.Text) {
		err := c.API.Delete(ctx, message)
		if err != nil {
			zap.L().Error("failed to delete commander menu item", zap.Error(err))
		}
	}
}

func (c *Commander) handleCallbackQueryUpdate(ctx context.Context, update tgbotapi.Update) {
	callbackQuery := update.CallbackQuery
	if callbackQuery == nil {
		return
	}

	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := c.API.Bot.Request(callback); err != nil {
		zap.L().Error("failed to answer callback query", zap.Error(err))
	}

	if command, ok := c.messageHandlersRegistry[callbackQuery.Data]; ok {
		command(ctx, update)
		return
	}
}

func logClientMessage(message *tgbotapi.Message) {
	if message.From == nil || message.From.IsBot {
		return
	}

	zap.L().Info(
		"client message",
		zap.String("msg", message.Text),
		zap.String("username", message.From.UserName),
	)
}

func (c *Commander) handlePollUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Poll == nil {
		return
	}

	c.HandleDistrictsOptionsResult(ctx, update.Poll)
}
