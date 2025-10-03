package commander

import (
	"air-quality-notifyer/internal/service/user"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) ShowUsers(update tgbotapi.Update, service *user.Service) {
	adminId, err := strconv.Atoi(c.cfg.App.AdminTelegramId)
	chatId := update.Message.Chat.ID
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
	if err = c.Send(MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending show_users", zap.Error(err))
	}
}
