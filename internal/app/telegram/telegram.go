package telegram

import (
	"air-quality-notifyer/internal/app/telegram/commander"
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/app/telegram/commander/mode"
	"air-quality-notifyer/internal/config"
	"fmt"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type tgBot struct {
	bot       *tgbotapi.BotAPI
	updates   tgbotapi.UpdatesChannel
	Commander *commander.Commander
}

func Init(cfg config.Config, services *commander.Services) *tgBot {
	bot, err := tgbotapi.NewBotAPI(cfg.App.TelegramToken)
	if err != nil {
		zap.L().Error("Filed to create new bot api", zap.Error(err))
		panic(err)
	}

	cmder := commander.New(cfg, bot, services)
	if cfg.Development {
		bot.Debug = true

		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 30

		return &tgBot{
			bot:       bot,
			updates:   bot.GetUpdatesChan(updateConfig),
			Commander: cmder,
		}
	}

	wh, err := tgbotapi.NewWebhook(fmt.Sprintf("https://%s/webhook%s", cfg.App.WebhookHost, bot.Token))
	if err != nil {
		zap.L().Panic("Filed to create new webhook", zap.Error(err))
	}
	_, err = bot.Request(wh)
	if err != nil {
		panic(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		zap.L().Panic("Failed to get webhook info", zap.Error(err))
	}

	if info.LastErrorDate != 0 {
		zap.L().Error("failed to init get info about webhook", zap.Error(err))
	}

	updates := bot.ListenForWebhook(fmt.Sprintf("/webhook%s", bot.Token))

	return &tgBot{
		bot:       bot,
		updates:   updates,
		Commander: cmder,
	}
}

func (t *tgBot) Start() {
	go t.ListenUpdates()
	go t.ListenSensors()
}

func (t *tgBot) ListenSensors() {
	t.Commander.Services.Sensor.ListenChanges(t.NotifyUsers)
}

func (t *tgBot) ListenUpdates() {
	cfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "🌀 Перезапустить бота",
		},
	)
	if _, err := t.bot.Request(cfg); err != nil {
		zap.L().Error("commander request error", zap.Error(err))
	}

	//t.Commander.HandleUpdate(t.updates)
	for update := range t.updates {
		if update.Message != nil {
			if !update.Message.From.IsBot {
				zap.L().Info(
					"client message",
					zap.String("msg", update.Message.Text),
					zap.String("username", update.Message.From.UserName),
				)
			}

			switch update.Message.Text {
			case "/start":
				t.Commander.Start(update)
			case api.KeypadUsersText:
				t.Commander.Admin.ShowUsers(update)
			case api.KeypadFaqText:
				t.Commander.API.MenuFaq(update)
			case api.KeypadSettingsText:
				t.Commander.Settings(update)
			case api.KeypadPingText:
				t.Commander.Admin.Pong(update)
			}

			if api.IsMenuButton(update.Message.Text) {
				err := t.Commander.API.Delete(update.Message)
				if err != nil {
					zap.L().Error("failed to delete commander menu item", zap.Error(err))
				}
			}
		}

		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := t.bot.Request(callback); err != nil {
				zap.L().Error("Error receiving response from callback with id", zap.Error(err), zap.String("id", update.CallbackQuery.ID))
				continue
			}

			switch update.CallbackQuery.Data {
			case api.KeypadMenuBackData:
				t.Commander.API.MenuBack(update)
			case api.KeypadModeFaqData, mode.KeypadFaqFromSetupData:
				t.Commander.Mode.Faq(update)
			case mode.KeypadSetupData:
				t.Commander.Mode.Setup(update)
			case mode.KeypadSetCityData:
				t.Commander.Mode.SetCity(update)
			case mode.KeypadSetDistrictData:
				t.Commander.Mode.SetDistrict(update)
				//case mode.KeypadSetDistrictData:
				//	t.Commander.Mode.SetDistrict(update, t.services.UserService, constants.District)
				//case mode.KeypadSetHomeData:
				//	t.Commander.Mode.SetHome(update, t.services.UserService, constants.Home)
			}
		}
	}
}
