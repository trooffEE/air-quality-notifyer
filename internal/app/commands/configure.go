package commands

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) Configure(message *tgbotapi.Message) {
	loc := tgbotapi.NewKeyboardButtonLocation("Предоставить доступ к геолокации")
	keyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{loc})

	msg := tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста предоставьте доступ до геолокации, чтобы мы могли сформировать список датчиков, за которыми мы будем следить")
	msg.ReplyMarkup = keyboard

	err := c.Send(SendPayload{Msg: msg})

	if err != nil {
		zap.L().Error("Error sending configure message", zap.Error(err))
	}
}
