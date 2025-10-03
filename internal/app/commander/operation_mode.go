package commander

import (
	"air-quality-notifyer/internal/app/keypads"
	"fmt"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) OperationMode(update tgbotapi.Update) {
	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("Пожалуйста, выберите один из трех режимов работы для его настройки:\n\nЕсли не знайте какой режим выбрать, нажмите на \"%s\", чтобы получить информацию о них", keypads.OperatingModeFAQFromSetupText),
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(keypads.SetOperationModeCityText, keypads.SetOperationModeCityData),
			tgbotapi.NewInlineKeyboardButtonData(keypads.SetOperationModeDistrictText, keypads.SetOperationModeDistrictData),
			tgbotapi.NewInlineKeyboardButtonData(keypads.SetOperationModeHomeText, keypads.SetOperationModeHomeData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(keypads.OperatingModeFAQFromSetupText, keypads.OperatingModeFAQFromSetupData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(keypads.BackToMenuText, keypads.BackToMenuData),
		),
	)

	if err := c.Edit(EditMessageConfig{Msg: msg, Markup: &markup}); err != nil {
		zap.L().Error("Error sending operating_mode message", zap.Error(err))
	}
}
