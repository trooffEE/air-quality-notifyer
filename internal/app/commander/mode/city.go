package mode

import (
	"air-quality-notifyer/internal/app/commander/api"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) SetCity(update tgbotapi.Update) {
	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"🏙 Город 🏙\n\nТеперь вы будете получать оповещения с датчиков по всему городу! 🍃",
	)

	if err := c.api.Edit(api.EditMessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending SetCity message", zap.Error(err))
	}
}
