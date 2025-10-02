package commander

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) OperatingModeFaq(update tgbotapi.Update) {
	//msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(
	//	"⚙️ <strong>Режимы работы</strong> ⚙️\n\n"+
	//		"🏙 <i>Город</i> 🏙\n\n"+
	//		"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков по <strong>всему городу</strong>. Данный функционал следует использовать, если вы хотите следить за общим состоянием воздуха в городе 🍃\n\n\n"+
	//		"🏘 <i>Район</i> 🏘\n\n"+
	//		"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков по <strong>выбранному району</strong> Кемерово. Данный функционал следует использовать, если вы хотите следить за конекретным районом/районами города 🍃\n\n\n"+
	//		"🏡 <i>Дом</i> 🏡\n\n"+
	//		"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков <strong>в пределах километра от выбранного места на карте или выбранных в ручную вами</strong>. Данный функционал следует использовать, если вы хотите следить за конкретными интересующими датчиками 🍃\n\n",
	//),
	//)

	msgEdit := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "test")

	if err := c.Edit(PayloadEdit{Msg: msgEdit}); err != nil {
		zap.L().Error("Error sending operating_mode_faq message", zap.Error(err))
	}
}
