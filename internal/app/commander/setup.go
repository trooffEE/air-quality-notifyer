package commander

import (
	"air-quality-notifyer/internal/app/keypads"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) Setup(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"⚙️ <strong>Настройки</strong>\n"+
			"Здесь вы можете настроить нужный функционал бота",
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(keypads.OperationModeText, keypads.OperationModeData),
			//TODO will be back soon
			//tgbotapi.NewInlineKeyboardButtonData(keypads.SensorsText, keypads.SensorsData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(keypads.BackToMenuText, keypads.BackToMenuData),
		),
	)

	if err := c.Send(Payload{Msg: msg, ReplyMarkup: markup}); err != nil {
		zap.L().Error("Error sending configure message", zap.Error(err))
	}
}
