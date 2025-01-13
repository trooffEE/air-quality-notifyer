package commands

import (
	"fmt"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"log"
)

func (c *Commander) Help(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, `/help - command check`)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := c.bot.Send(msg)

	if err != nil {
		log.Print(fmt.Sprintf("Error appeared upon sending me message %#v", err))
	}
}
