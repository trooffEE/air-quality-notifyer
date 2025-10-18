package mode

import (
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/constants"
	sDistricts "air-quality-notifyer/internal/service/districts"
	sUser "air-quality-notifyer/internal/service/user"
	"fmt"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type Commander struct {
	api     api.Interface
	service Service
}

type Service struct {
	User     sUser.Interface
	District sDistricts.Interface
}

type Interface interface {
	Setup(update tgbotapi.Update)
	Faq(update tgbotapi.Update)
	SetCity(update tgbotapi.Update)
	SetDistrict(update tgbotapi.Update)
}

func New(api api.Interface, service Service) Interface {
	return &Commander{
		api:     api,
		service: service,
	}
}

func (c *Commander) Setup(update tgbotapi.Update) {
	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("Пожалуйста, выберите один из трех режимов работы для его настройки:\n\nЕсли не знайте какой режим выбрать, нажмите на \"%s\", чтобы получить информацию о них", KeypadFaqFromSetupText),
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(KeypadSetCityText, KeypadSetCityData),
			tgbotapi.NewInlineKeyboardButtonData(KeypadSetDistrictText, KeypadSetDistrictData),
			tgbotapi.NewInlineKeyboardButtonData(KeypadSetHomeText, KeypadSetHomeData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(KeypadFaqFromSetupText, KeypadFaqFromSetupData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(api.KeypadMenuBackText, api.KeypadMenuBackData),
		),
	)

	if err := c.api.Edit(api.EditMessageConfig{Msg: msg, Markup: &markup}); err != nil {
		zap.L().Error("Error sending operating_mode message", zap.Error(err))
	}
}

func (c *Commander) Faq(update tgbotapi.Update) {
	markup := tgbotapi.NewInlineKeyboardMarkup()

	if update.CallbackQuery.Data == KeypadFaqFromSetupData {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(KeypadBackText, KeypadSetupData),
		))
	}

	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(api.KeypadMenuBackText, api.KeypadMenuBackData),
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

	if err := c.api.Edit(api.EditMessageConfig{Msg: msg, Markup: &markup}); err != nil {
		zap.L().Error("Error sending operating_mode_faq message", zap.Error(err))
	}
}

func (c *Commander) SetCity(update tgbotapi.Update) {
	message := update.CallbackQuery.Message
	chatId := message.Chat.ID
	err := c.service.User.SetOperatingMode(chatId, constants.City)
	if err != nil {
		zap.L().Error("Error setting operating mode", zap.Error(err))
		return
	}

	msg := tgbotapi.NewMessage(
		chatId,
		"🏙 Город 🏙\n\nТеперь вы будете получать оповещения с датчиков по всему городу! 🍃",
	)

	if err := c.api.Send(api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending Mode.Set message", zap.Error(err))
	}

	if err = c.api.Delete(message); err != nil {
		zap.L().Error("Error deleting prev message", zap.Error(err))
	}
}

func (c *Commander) SetDistrict(update tgbotapi.Update) {
	chatId := update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(
		chatId,
		"🏘 Район 🏘\n\nДля того чтобы выставить режим работы \"Район 🏘\", выберите перечень интересующих районов:",
	)

	districts := c.service.District.GetAllDistricts()
	//districts := c.service.District.GetOptionForDistrict()
	var buttons []tgbotapi.KeyboardButton
	for _, d := range districts {
		buttons = append(buttons, tgbotapi.KeyboardButton{
			Text: d.Name,
		})
	}

	replyMarkup := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(buttons...),
	)

	if err := c.api.Send(api.MessageConfig{Msg: msg, Markup: replyMarkup}); err != nil {
		zap.L().Error("Error sending Mode.Set message", zap.Error(err))
	}
}
