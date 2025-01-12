package commands

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (c *Commander) DefaultSend(chatId int64, text string) *tgbotapi.Error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeHTML

	_, err := c.bot.Send(msg)
	if tgError, ok := err.(*tgbotapi.Error); ok {
		return tgError
	}

	return nil
}
