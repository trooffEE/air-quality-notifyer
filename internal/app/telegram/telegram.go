package telegram

import (
	"air-quality-notifyer/internal/app/commander"
	"air-quality-notifyer/internal/app/commander/api"
	"air-quality-notifyer/internal/app/commander/mode"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/service/sensor"
	"air-quality-notifyer/internal/service/user"
	"fmt"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type tgBot struct {
	bot       *tgbotapi.BotAPI
	updates   tgbotapi.UpdatesChannel
	services  BotServices
	Commander *commander.Commander
}

type BotServices struct {
	UserService   user.Interface
	SensorService sensor.Interface
}

func Init(services BotServices, cfg config.Config) *tgBot {
	bot, err := tgbotapi.NewBotAPI(cfg.App.TelegramToken)
	if err != nil {
		zap.L().Error("Filed to create new bot api", zap.Error(err))
		panic(err)
	}

	cmder := commander.New(bot, cfg)
	if cfg.Development {
		bot.Debug = true

		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 30

		return &tgBot{
			bot:       bot,
			updates:   bot.GetUpdatesChan(updateConfig),
			services:  services,
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
		services:  services,
		Commander: cmder,
	}
}

func (t *tgBot) ListenSensors() {
	t.services.SensorService.ListenChangesInSensors(t.notifyUsers)
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

	for update := range t.updates {
		if update.Message != nil {
			zap.L().Info(
				"client message",
				zap.String("msg", update.Message.Text),
				zap.String("username", update.Message.From.UserName),
			)

			switch update.Message.Text {
			case "/start":
				t.Commander.Start(update, t.services.UserService)
			case api.KeypadUsersText:
				t.Commander.Admin.ShowUsers(update, t.services.UserService)
			case api.KeypadFaqText:
				t.Commander.API.MenuFaq(update)
			case api.KeypadSettingsText:
				t.Commander.Settings(update)
			case api.KeypadPingText:
				t.Commander.Admin.Pong(update)
			}

			if api.IsMenuButton(update.Message.Text) {
				t.Commander.API.Delete(update)
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
			case mode.KeypadData:
				t.Commander.Mode.Setup(update)
			}
		}
	}
}
