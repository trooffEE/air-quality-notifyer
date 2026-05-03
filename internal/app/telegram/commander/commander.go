package commander

import (
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/helper"
	"air-quality-notifyer/internal/service/districts"
	"air-quality-notifyer/internal/service/sensor"
	"air-quality-notifyer/internal/service/user"
	"context"

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
		zap.L().Error("message handlers registry already registered, can not override it")
		return
	}

	adminMessageHandlersRegistry := NewAdminMessageHandlersRegistry(c)
	menuMessageHandlersRegistry := NewMenuMessageHandlersRegistry(c)
	coreMessageHandlersRegistry := NewCoreMessageHandlersRegistry(c)

	registries := []HandlersRegistry{adminMessageHandlersRegistry, menuMessageHandlersRegistry, coreMessageHandlersRegistry}

	if helper.HasOverlappingKeys(registries...) {
		zap.L().Error("message handlers registry has overlapping keys")
		return
	}

	c.messageHandlersRegistry = helper.MergeMaps(registries...)
}

func (c *Commander) RegisterCallbackHandlers() {
	if c.callbackHandlersRegistry != nil {
		zap.L().Error("callback handlers registry already registered, can not override it")
		return
	}

	modeCallbackHandlersRegistry := NewModeCallbackHandlersRegistry(c)
	menuCallbackHandlersRegistry := NewMenuCallbackHandlersRegistry(c)

	registries := []HandlersRegistry{modeCallbackHandlersRegistry, menuCallbackHandlersRegistry}

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
