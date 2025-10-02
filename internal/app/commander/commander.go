package commander

import (
	"air-quality-notifyer/internal/app/menu"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/service/user"
	"errors"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type Commander struct {
	bot *tgbotapi.BotAPI
	cfg config.Config
	loc *time.Location
}

type Interface interface {
	Send(payload Payload) *tgbotapi.Error
	Delete(update tgbotapi.Update)
	Start(update tgbotapi.Update, service *user.Service)
	BackToMenu(update tgbotapi.Update)
	FAQ(update tgbotapi.Update)
	OperationMode(update tgbotapi.Update)
	OperatingModeFaq(update tgbotapi.Update)
	Pong(update tgbotapi.Update)
	Setup(update tgbotapi.Update)
	ShowUsers(update tgbotapi.Update, service *user.Service)
}

func NewCommander(bot *tgbotapi.BotAPI, cfg config.Config) Interface {
	loc, _ := time.LoadLocation("Asia/Novosibirsk")
	return &Commander{
		bot: bot,
		cfg: cfg,
		loc: loc,
	}
}

type Payload struct {
	Msg         tgbotapi.MessageConfig
	ReplyMarkup interface{}
}

type PayloadEdit struct {
	Msg         tgbotapi.EditMessageTextConfig
	ReplyMarkup *tgbotapi.InlineKeyboardMarkup
}

func (c *Commander) Delete(update tgbotapi.Update) {
	message := update.Message
	_, err := c.bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
	if err != nil {
		zap.L().Error("Error deleting message", zap.Error(err))
	}
}

func (c *Commander) Edit(payload PayloadEdit) *tgbotapi.Error {
	payload.Msg.ParseMode = tgbotapi.ModeHTML

	//if payload.ReplyMarkup != nil {
	//	payload.Msg.ReplyMarkup = payload.ReplyMarkup
	//} else {
	//	payload.Msg.ReplyMarkup = menu.NewTelegramMainMenu()
	//}
	_, err := c.bot.Send(payload.Msg)
	var tgError *tgbotapi.Error
	if errors.As(err, &tgError) {
		return tgError
	}

	return nil
}

func (c *Commander) Send(payload Payload) *tgbotapi.Error {
	payload.Msg.ParseMode = tgbotapi.ModeHTML
	payload.Msg.DisableNotification = c.isNotificationsAllowed()

	if payload.ReplyMarkup != nil {
		payload.Msg.ReplyMarkup = payload.ReplyMarkup
	} else {
		payload.Msg.ReplyMarkup = menu.NewTelegramMainMenu()
	}

	_, err := c.bot.Send(payload.Msg)
	var tgError *tgbotapi.Error
	if errors.As(err, &tgError) {
		return tgError
	}

	return nil
}

func (c *Commander) isNotificationsAllowed() bool {
	h := time.Now().In(c.loc).Hour()
	return h < 8 && h >= 0
}
