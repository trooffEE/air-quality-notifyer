package commander

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) BackToMenu(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы вернулись в меню ⬇️")

	if err := c.Send(MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending back message", zap.Error(err))
	}
}
