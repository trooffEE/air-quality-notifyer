package commands

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) Help(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, `/help - command check`)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := c.bot.Send(msg)

	if err != nil {
		zap.L().Error("Error sending help message", zap.Error(err))
	}
}
