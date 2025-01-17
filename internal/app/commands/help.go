package commands

import (
	"air-quality-notifyer/internal/lib"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (c *Commander) Help(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, `/help - command check`)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := c.bot.Send(msg)

	if err != nil {
		lib.LogError("Help", "failed to send message", err)
	}
}
