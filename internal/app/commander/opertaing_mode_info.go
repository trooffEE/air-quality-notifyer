package commander

import (
	"fmt"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) OperatingModeInfo(callback *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf(
		"⚙️ <strong>Режимы работы</strong> ⚙️\n\n"+
			"🏙 <i>Город</i> 🏙\n\n"+
			"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков по <strong>всему городу</strong>. Данный функционал следует использовать, если вы хотите следить за общим состоянием воздуха в городе 🍃\n\n\n"+
			"🏘 <i>Район</i> 🏘\n\n"+
			"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков по <strong>выбранному району</strong> Кемерово. Данный функционал следует использовать, если вы хотите следить за конекретным районом/районами города 🍃\n\n\n"+
			"🏡 <i>Дом</i> 🏡\n\n"+
			"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков <strong>в пределах километра от выбранного места на карте или выбранных в ручную вами</strong>. Данный функционал следует использовать, если вы хотите следить за конкретными интересующими датчиками 🍃\n\n",
	),
	)

	if err := c.Send(Payload{Msg: msg}); err != nil {
		zap.L().Error("Error sending operating_mode message", zap.Error(err))
	}
}
