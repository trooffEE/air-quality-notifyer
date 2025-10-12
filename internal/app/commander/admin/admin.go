package admin

import (
	"air-quality-notifyer/internal/app/commander/api"
	"air-quality-notifyer/internal/service/user"
	"fmt"
	"strings"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type Commander struct {
	api api.Interface
}

type Interface interface {
	Pong(update tgbotapi.Update)
	ShowUsers(update tgbotapi.Update, service user.Interface)
}

func New(api api.Interface) Interface {
	return &Commander{
		api: api,
	}
}

func (c *Commander) Pong(update tgbotapi.Update) {
	if !c.api.IsAdmin(update) {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `pong - 🙌`)
	msg.ReplyParameters.MessageID = update.Message.MessageID

	if err := c.api.Send(api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending pong message", zap.Error(err))
	}
}

func (c *Commander) ShowUsers(update tgbotapi.Update, service user.Interface) {
	if !c.api.IsAdmin(update) {
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
	if err := c.api.Send(api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending show_users", zap.Error(err))
	}
}
