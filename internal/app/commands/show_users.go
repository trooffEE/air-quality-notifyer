package commands

import (
	"air-quality-notifyer/internal/service/user"
	"fmt"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func (c *Commander) ShowUsers(message *tgbotapi.Message, service *user.Service) {
	adminId, err := strconv.Atoi(c.cfg.AdminTelegramId)
	chatId := message.Chat.ID
	if err != nil {
		zap.L().Error("failed to convert admin telegram id to int", zap.Error(err))
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

	msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("Bot Users: %d 🙌\n\n%s", len(names), strings.Join(names, ", \n")))
	err = c.Send(SendPayload{Msg: msg})

	if err != nil {
		zap.L().Error("Error sending show_users", zap.Error(err))
	}
}
