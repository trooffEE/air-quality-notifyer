package menu

import tgbotapi "github.com/OvyFlash/telegram-bot-api"

const (
	FAQ               = "❓ FAQ"
	OperationModeInfo = "❓ Режимы работы"
	Setup             = "⚙️ Настройки"
)

func IsMenuButton(button string) bool {
	return button == FAQ || button == OperationModeInfo || button == Setup
}

func NewTelegramMainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Setup),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(FAQ),
			tgbotapi.NewKeyboardButton(OperationModeInfo),
		),
	)
}
