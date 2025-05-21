package menu

import tgbotapi "github.com/OvyFlash/telegram-bot-api"

const (
	FAQ   = "❓ FAQ"
	Setup = "⚙️ Настройки"
)

func IsMenuButton(button string) bool {
	return button == FAQ || button == Setup
}

func NewTelegramMainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Setup),
			tgbotapi.NewKeyboardButton(FAQ),
		),
	)
}
