package api

import (
	"fmt"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (a *Api) MenuBack(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы вернулись в меню ⬇️")

	if err := a.Send(MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending back message", zap.Error(err))
	}
}

func (a *Api) MenuFaq(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
		"⚙️<strong>Ответы на вопросы</strong>\n\n"+
			"<i>- Связан ли данный бот с https://airkemerovo.ru ?</i>\n"+
			"Данный проект берет информацию именно с этого сервиса, используя публичный API, однако разработчик бота <strong>никак не связан</strong> с https://airkemerovo.ru\n\n"+
			"<i>- Бесплатно ли использование?</i>\n"+
			"Да, проект доступен для использования всем желающим, кто хочет следить за состоянием воздуху и получать уведомления об этом 🌿\n\n"+
			"<i>- Что умеет данный бот?</i>\n"+
			"Оповещать об опасном воздухе в разных режимах\n"+
			"1. по всему городу (режим - \"город\")\n"+
			"2. * по районам города (режим - \"район\")\n"+
			"3. * по датчикам (режим - \"дом\")\n\n"+
			"* - функционал в работе\n"+
			"ℹ Подробнее о режимах можно прочитать в меню - \n\"%s\"\n\n"+
			"<i>- Почему на https://airkemerovo.ru больше датчиков, чем в предложенном перечне?</i>\n"+
			"Потому что на данном этапе проекта реализована работа только с <strong>датчиками города Кемерово</strong>. Это сознательное решение автора проекта, однако это вполне может поменяться в будущем\n\n",
		KeypadFaqText,
	),
	)
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(KeypadFaqText, KeypadFaqData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(KeypadMenuBackText, KeypadMenuBackData),
		),
	)

	if err := a.Send(MessageConfig{Msg: msg, Markup: markup}); err != nil {
		zap.L().Error("Error sending faq message", zap.Error(err))
	}
}
