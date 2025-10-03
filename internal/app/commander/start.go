package commander

import (
	"air-quality-notifyer/internal/service/user"
	"strconv"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) Start(update tgbotapi.Update, service *user.Service) {
	message := update.Message
	chatId, username := message.Chat.ID, message.Chat.UserName
	c.greetNewUser(chatId)

	if !service.IsNewUser(chatId) {
		return
	}

	service.Register(user.User{
		Id:       strconv.Itoa(int(chatId)),
		Username: username,
	})
}

func (c *Commander) greetNewUser(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Данный бот оповещает о плохом качестве воздуха в городе Кемерово.\n\nПросьба настроить уведомления, чтобы бот не беспокоил ночью! 🍵")
	if err := c.Send(MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending faq message", zap.Error(err))
	}
}
