package commands

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) Back(callback *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Вы вернулись в меню ⬇️")
	err := c.Send(SendPayload{Msg: msg})

	if err != nil {
		zap.L().Error("Error sending back message", zap.Error(err))
	}
}
