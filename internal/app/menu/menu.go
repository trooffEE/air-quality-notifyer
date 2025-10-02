package menu

import (
	"slices"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

const (
	FAQ   = "❓ FAQ"
	Setup = "⚙️ Настройки"
	Users = "users"
	Ping  = "ping"
)

var options = []string{FAQ, Setup, Users, Ping}

func IsMenuButton(button string) bool {
	return slices.Contains(options, button)
}

func NewTelegramMainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Setup),
			tgbotapi.NewKeyboardButton(FAQ),
		),
	)
}
