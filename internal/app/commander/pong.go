package commander

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) Pong(update tgbotapi.Update) {
	message := update.Message
	msg := tgbotapi.NewMessage(message.Chat.ID, `pong - 🙌`)
	msg.ParseMode = tgbotapi.ModeHTML

	msg.ReplyParameters.MessageID = message.MessageID

	if err := c.Send(MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending pong message", zap.Error(err))
	}
}
