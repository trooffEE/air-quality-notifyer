package commands

import (
	"air-quality-notifyer/internal/service/user"
	"fmt"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"log"
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
		log.Print(fmt.Sprintf("Error appeared upon sending message to user %d with message %s, %#v", chatId, text, err))
	}
}
