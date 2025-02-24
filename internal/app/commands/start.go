package commands

import (
	"air-quality-notifyer/internal/service/user"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
	"strconv"
)

func (c *Commander) Start(message *tgbotapi.Message, service *user.Service) {
	chatId, username := message.Chat.ID, message.Chat.UserName

	if service.IsNewUser(chatId) {
		c.greetNewUser(chatId)

		service.Register(user.User{
			Id:       strconv.Itoa(int(chatId)),
			Username: username,
		})
	}
}

func (c *Commander) greetNewUser(chatId int64) {
	text := "Приветствую. Данный бот оповещает о плохом качестве воздуха по районам в городе Кемерово.\n\nПросьба настроить уведомления, чтобы бот не беспокоил ночью! 🍵"
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := c.bot.Send(msg)
	if err != nil {
		zap.L().Error("failed to send message to chatId", zap.Int64("chatId", chatId), zap.Error(err))
	}
}
