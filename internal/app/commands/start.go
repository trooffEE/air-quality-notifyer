package commands

import (
	"air-quality-notifyer/internal/service/user"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func (c *Commander) Start(message *tgbotapi.Message, service *user.Service) {
	chatId, username := message.Chat.ID, message.Chat.UserName

	fmt.Println(chatId, username)
	if service.IsNewUser(chatId) {
		service.Register(user.User{
			Id:       strconv.Itoa(int(chatId)),
			Username: username,
		})
	}
}

func (c *Commander) greeting(chatId int64) {
	text := "Привествую. Данный бот оповещает о плохом качестве воздуха по районам в городе Кемерово.\n\nПросьба настроить уведомления, чтобы бот не беспокоил ночью! 🍵"
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := c.bot.Send(msg)
	if err != nil {
		log.Print(fmt.Sprintf("Error appeared upon sending message to user %d with message %s", chatId, text))
	}
}