package commander

import (
	"air-quality-notifyer/internal/app/keypads"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) Setup(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		"⚙️ <strong>Настройки</strong>\n"+
			"Здесь вы можете настроить нужный функционал бота",
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(keypads.OperationModeText, keypads.OperationModeData),
			tgbotapi.NewInlineKeyboardButtonData(keypads.SensorsText, keypads.SensorsData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(keypads.BackText, keypads.BackData),
		),
	)
	err := c.Send(SendPayload{Msg: msg, ReplyMarkup: markup})

	if err != nil {
		zap.L().Error("Error sending configure message", zap.Error(err))
	}
}
