package commander

import (
	"context"
	"strconv"

	"air-quality-notifyer/internal/app/telegram/commander/admin"
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/app/telegram/commander/mode"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/service/districts"
	"air-quality-notifyer/internal/service/sensor"
	"air-quality-notifyer/internal/service/user"
	"air-quality-notifyer/internal/service/user/model"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Commander struct {
	API      *api.Api
	Admin    admin.Interface
	Mode     mode.Interface
	Services *Services
}

type Services struct {
	User     user.Interface
	District districts.Interface
	Sensor   sensor.Interface
	Cache    *redis.Client
}

func New(cfg config.Config, bot *tgbotapi.BotAPI, s *Services) *Commander {
	apiCmder, err := api.NewApi(cfg, bot, s.Cache)
	if err != nil {
		zap.S().Fatalw("Failed to create api interface", "error", err)
		return nil
	}

	return &Commander{
		API:      apiCmder,
		Admin:    admin.New(apiCmder, admin.Service{User: s.User}),
		Mode:     mode.New(apiCmder, mode.Service{User: s.User, District: s.District}),
		Services: s,
	}
}

func (c *Commander) HandleUpdate(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			c.handleUpdate(ctx, update)
		}
	}
}

func (c *Commander) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message != nil {
		if !update.Message.From.IsBot {
			zap.L().Info(
				"client message",
				zap.String("msg", update.Message.Text),
				zap.String("username", update.Message.From.UserName),
			)
		}

		if c.HandlePendingFeedback(ctx, update) {
			return
		}

		switch update.Message.Text {
		case "/start":
			c.Start(ctx, update)
		case api.KeypadUsersText:
			c.Admin.ShowUsers(ctx, update)
		case api.KeypadFaqText:
			c.API.MenuFaq(update)
		case api.KeypadSettingsText:
			c.Settings(ctx, update)
		case api.KeypadPingText:
			c.Admin.Pong(update)
		default:
			switch {
			case admin.IsAnnounceCommand(update.Message.Text):
				c.Admin.Announce(ctx, update)
			case isFeedbackCommand(update.Message):
				c.Feedback(ctx, update)
			}
		}

		if api.IsMenuButton(update.Message.Text) {
			err := c.API.Delete(update.Message)
			if err != nil {
				zap.L().Error("failed to delete commander menu item", zap.Error(err))
			}
		}
	}

	if update.Poll != nil {
		c.Mode.HandleDistrictsOptionsResult(ctx, update.Poll)
	}

	if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		if _, err := c.API.Bot.Request(callback); err != nil {
			zap.L().Error("Error receiving response from callback with id", zap.Error(err), zap.String("id", update.CallbackQuery.ID))
			return
		}

		switch update.CallbackQuery.Data {
		case api.KeypadMenuBackData:
			c.API.MenuBack(update)
		case api.KeypadModeFaqData, mode.KeypadFaqFromSetupData:
			c.Mode.Faq(update)
		case mode.KeypadSetupData:
			c.Mode.Setup(update)
		case mode.KeypadSetCityData:
			c.Mode.SetCity(ctx, update)
		case mode.KeypadAskForDistrictOptionsData:
			c.Mode.AskForDistrictOptions(ctx, update)
		}
	}
}

func (c *Commander) Start(ctx context.Context, update tgbotapi.Update) {
	message := update.Message
	chatId, username := message.Chat.ID, message.Chat.UserName

	msg := tgbotapi.NewMessage(chatId, "Данный бот оповещает о плохом качестве воздуха в городе Кемерово.\n\nПросьба настроить уведомления, чтобы бот не беспокоил ночью! 🍵")
	if err := c.API.Send(api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending faq message", zap.Error(err))
	}

	if !c.Services.User.IsNew(ctx, chatId) {
		return
	}

	c.Services.User.Register(ctx, model.User{
		Id:       strconv.Itoa(int(chatId)),
		Username: username,
	})
}

func (c *Commander) Settings(ctx context.Context, update tgbotapi.Update) {
	if ctx.Err() != nil {
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"⚙️ <strong>Настройки</strong>\n"+
			"Здесь вы можете настроить нужный функционал бота",
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(mode.KeypadSetupText, mode.KeypadSetupData),
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
