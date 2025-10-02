package commander

import (
	"air-quality-notifyer/internal/app/keypads"
	"fmt"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) OperationMode(update tgbotapi.Update) {
	fmt.Println(update.Message)
	msg := tgbotapi.NewMessage(
		update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Пожалуйста, выберите один из трех режимов работы для его настройки:\n\nЕсли не знайте какой режим выбрать, нажмите на \"%s\", чтобы получить информацию о них", keypads.CommonKnowMoreText),
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(keypads.SetOperationModeCityText, keypads.SetOperationModeCityData),
			tgbotapi.NewInlineKeyboardButtonData(keypads.SetOperationModeDistrictText, keypads.SetOperationModeDistrictData),
			tgbotapi.NewInlineKeyboardButtonData(keypads.SetOperationModeHomeText, keypads.SetOperationModeHomeData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(keypads.CommonKnowMoreText, keypads.OperationModeFAQData),
		),
	)

	if err := c.Send(Payload{Msg: msg, ReplyMarkup: markup}); err != nil {
		zap.L().Error("Error sending operating_mode message", zap.Error(err))
	}
}
