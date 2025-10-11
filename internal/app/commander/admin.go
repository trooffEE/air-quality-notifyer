package commander

import (
	"air-quality-notifyer/internal/service/user"
	"fmt"
	"strings"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (c *Commander) ShowUsers(update tgbotapi.Update, service user.Interface) {
	if !c.isAdmin(update) {
		return
	}

	names := service.GetUsersNames()

	if len(names) == 0 {
		return
	}

	for index, name := range names {
		names[index] = "@" + name
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Bot Users: %d 🙌\n\n%s", len(names), strings.Join(names, ", \n")))
	if err := c.Send(MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending show_users", zap.Error(err))
	}
}

func (c *Commander) Pong(update tgbotapi.Update) {
	if !c.isAdmin(update) {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `pong - 🙌`)
	msg.ReplyParameters.MessageID = update.Message.MessageID

	if err := c.Send(MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending pong message", zap.Error(err))
	}
}
