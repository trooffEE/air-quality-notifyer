package commander

import (
	"fmt"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) FAQ(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(
		"⚙️<strong>Ответы на вопросы</strong>\n\n"+
			"<i>- Связан ли данный бот с https://airkemerovo.ru ?</i>\n"+
			"Данный проект берет информацию именно с этого сервиса, используя публичный API, однако разработчик бота <strong>никак не связан</strong> с https://airkemerovo.ru\n\n"+
			"<i>- Бесплатно ли использование?</i>\n"+
			"Да, проект доступен для использования всем желающим, кто хочет следить за состоянием воздуху и получать уведомления об этом\n\n"+
			"<i>- Что умеет данный бот?</i>\n"+
			"Оповещать об опасном воздухе в разных режимах\n"+
			"1. по всему городу (режим - \"город\")\n"+
			"2. * по районам города (режим - \"район\")\n"+
			"3. * по датчикам (режим - \"дом\")\n\n"+
			"* - функционал в работе\n"+
			"ℹ Подробнее о режимах можно прочитать в меню - \"Режимы работы\"\n\n"+
			"Если не нашли ответ на свой вопрос, то welcome - @adorable_internet_friend",
	),
	)
	err := c.Send(SendPayload{Msg: msg})

	if err != nil {
		zap.L().Error("Error sending faq message", zap.Error(err))
	}
}
