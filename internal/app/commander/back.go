package commander

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) Back(callback *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Вы вернулись в меню ⬇️")

	if err := c.Send(Payload{Msg: msg}); err != nil {
		zap.L().Error("Error sending back message", zap.Error(err))
	}
}
