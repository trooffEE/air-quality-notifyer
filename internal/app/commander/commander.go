package commander

import (
	"air-quality-notifyer/internal/app/commander/admin"
	"air-quality-notifyer/internal/app/commander/api"
	"air-quality-notifyer/internal/app/commander/mode"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/service/user"
	"air-quality-notifyer/internal/service/user/model"
	"strconv"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type Commander struct {
	API   api.Interface
	Admin admin.Interface
	Mode  mode.Interface
}

func New(bot *tgbotapi.BotAPI, cfg config.Config) *Commander {
	apiCmder, err := api.NewApi(bot, cfg)
	if err != nil {
		zap.S().Fatalw("Failed to create api interface", "error", err)
		return nil
	}

	return &Commander{
		API:   apiCmder,
		Admin: admin.New(apiCmder),
		Mode:  mode.New(apiCmder),
	}
}

func (c *Commander) Start(update tgbotapi.Update, service user.Interface) {
	message := update.Message
	chatId, username := message.Chat.ID, message.Chat.UserName

	msg := tgbotapi.NewMessage(chatId, "Данный бот оповещает о плохом качестве воздуха в городе Кемерово.\n\nПросьба настроить уведомления, чтобы бот не беспокоил ночью! 🍵")
	if err := c.API.Send(api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending faq message", zap.Error(err))
	}

	if !service.IsNewUser(chatId) {
		return
	}

	service.Register(model.User{
		Id:       strconv.Itoa(int(chatId)),
		Username: username,
	})
}

func (c *Commander) Settings(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"⚙️ <strong>Настройки</strong>\n"+
			"Здесь вы можете настроить нужный функционал бота",
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(mode.KeypadText, mode.KeypadData),
			//TODO will be back soon
			//tgbotapi.NewInlineKeyboardButtonData(keypads.SensorsText, keypads.SensorsData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(api.KeypadMenuBackText, api.KeypadMenuBackData),
		),
	)

	if err := c.API.Send(api.MessageConfig{Msg: msg, Markup: markup}); err != nil {
		zap.L().Error("Error sending configure message", zap.Error(err))
	}
}
