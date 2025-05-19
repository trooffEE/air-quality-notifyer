package commands

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) Pong(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, `pong - 🙌`)
	msg.ParseMode = tgbotapi.ModeHTML

	msg.ReplyParameters.MessageID = message.MessageID
	_, err := c.bot.Send(msg)

	if err != nil {
		zap.L().Error("Error sending pong message", zap.Error(err))
	}
}
