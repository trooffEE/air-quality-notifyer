package commands

import (
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/service/user"
	"fmt"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

func (c *Commander) ShowUsers(message *tgbotapi.Message, service *user.Service) {
	adminId, err := strconv.Atoi(config.Cfg.AdminTelegramId)
	chatId := message.Chat.ID
	if err != nil {
		log.Println(err)
		return
	}
	if chatId != int64(adminId) {
		return
	}

	names := *service.GetUsersNames()

	if len(names) == 0 {
		return
	}

	for index, name := range names {
		names[index] = "@" + name
	}

	test := fmt.Sprintf("Bot Users: %d 🙌\n\n%s", len(names), strings.Join(names, ", \n"))
	msg := tgbotapi.NewMessage(chatId, test)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = c.bot.Send(msg)

	if err != nil {
		log.Print(fmt.Sprintf("Error appeared upon sending me message %#v", err))
	}
}
