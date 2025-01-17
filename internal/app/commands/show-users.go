package commands

import (
	"air-quality-notifyer/internal/lib"
	"air-quality-notifyer/internal/service/user"
	"fmt"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"strconv"
	"strings"
)

func (c *Commander) ShowUsers(message *tgbotapi.Message, service *user.Service) {
	adminId, err := strconv.Atoi(c.cfg.AdminTelegramId)
	chatId := message.Chat.ID
	if err != nil {
		lib.LogError("ShowUsers", "failed to convert admin telegram id to int", err)
		return
	}
	if chatId != int64(adminId) {
		return
	}

	names := service.GetUsersNames()

	if len(names) == 0 {
		return
	}

	for index, name := range names {
		names[index] = "@" + name
	}

	msgString := fmt.Sprintf("Bot Users: %d 🙌\n\n%s", len(names), strings.Join(names, ", \n"))
	msg := tgbotapi.NewMessage(chatId, msgString)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = c.bot.Send(msg)

	if err != nil {
		lib.LogError("ShowUsers", "failed to send message", err)
	}
}
