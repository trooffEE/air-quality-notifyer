package commander

import (
	"context"

	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/app/telegram/commander/mode"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/helper"
	"air-quality-notifyer/internal/service/districts"
	"air-quality-notifyer/internal/service/sensor"
	"air-quality-notifyer/internal/service/user"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Handler func(ctx context.Context, update tgbotapi.Update)
type HandlersRegistry map[string]Handler

type Commander struct {
	API                      *api.Api
	Services                 *Services
	messageHandlersRegistry  HandlersRegistry
	callbackHandlersRegistry HandlersRegistry
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

	commander := &Commander{
		API:      apiCmder,
		Services: s,
	}

	commander.RegisterMessageHandlers()
	commander.RegisterCallbackHandlers()

	return commander
}

func (c *Commander) RegisterMessageHandlers() {
	if c.messageHandlersRegistry != nil {
		zap.L().Warn("message handlers registry already registered, can not override it")
		return
	}

	adminMessageHandlersRegistry := NewAdminMessageHandlersRegistry(c)

	registries := []HandlersRegistry{adminMessageHandlersRegistry}

	if helper.HasOverlappingKeys(registries...) {
		zap.L().Error("message handlers registry has overlapping keys")
		return
	}

	c.messageHandlersRegistry = helper.MergeMaps(registries...)
}

func (c *Commander) RegisterCallbackHandlers() {
	if c.callbackHandlersRegistry != nil {
		zap.L().Warn("callback handlers registry already registered, can not override it")
		return
	}

	modeCallbackHandlersRegistry := NewModeCallbackHandlersRegistry(c)

	registries := []HandlersRegistry{modeCallbackHandlersRegistry}

	if helper.HasOverlappingKeys(registries...) {
		zap.L().Error("callback handlers registry has overlapping keys")
		return
	}

	c.callbackHandlersRegistry = helper.MergeMaps(registries...)
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

	if err := c.API.Send(ctx, api.MessageConfig{Msg: msg, Markup: markup}); err != nil {
		zap.L().Error("Error sending configure message", zap.Error(err))
	}
}
