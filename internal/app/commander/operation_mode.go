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

func (c *Commander) OperatingModeFaq(update tgbotapi.Update) {
	markup := tgbotapi.NewInlineKeyboardMarkup()

	if update.CallbackQuery.Data == keypads.OperatingModeFAQFromSetupData {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(keypads.OperationModeBackText, keypads.OperationModeData),
		))
	}

	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(keypads.BackToMenuText, keypads.BackToMenuData),
	))

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf(
			"⚙️ <strong>Режимы работы</strong> ⚙️\n\n"+
				"🏙 <i>Город</i> 🏙\n\n"+
				"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков по <strong>всему городу</strong>. Данный функционал следует использовать, если вы хотите следить за общим состоянием воздуха в городе 🍃\n\n\n"+
				"🏘 <i>Район</i> 🏘\n\n"+
				"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков по <strong>выбранному району</strong> Кемерово. Данный функционал следует использовать, если вы хотите следить за конекретным районом/районами города 🍃\n\n\n"+
				"🏡 <i>Дом</i> 🏡\n\n"+
				"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков <strong>в пределах километра от выбранного места на карте или выбранных в ручную вами</strong>. Данный функционал следует использовать, если вы хотите следить за конкретными интересующими датчиками 🍃\n\n",
		),
	)

	if err := c.Edit(EditMessageConfig{Msg: msg, Markup: &markup}); err != nil {
		zap.L().Error("Error sending operating_mode_faq message", zap.Error(err))
	}
}
