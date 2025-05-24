package telegram

import (
	"air-quality-notifyer/internal/app/commander"
	"air-quality-notifyer/internal/app/keypads"
	"air-quality-notifyer/internal/app/menu"
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
	UserService   *user.Service
	SensorService *sensor.Service
}

func InitTelegramBot(services BotServices, cfg config.ApplicationConfig) *tgBot {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		zap.L().Error("Filed to create new bot api", zap.Error(err))
		panic(err)
	}

	cmder := commander.NewCommander(bot, cfg)
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

	wh, err := tgbotapi.NewWebhook(fmt.Sprintf("https://%s/webhook%s", cfg.WebhookHost, bot.Token))
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

func (t *tgBot) ListenChangesInSensors() {
	t.services.SensorService.ListenChangesInSensors(t.notifyUsersAboutSensors)
}

func (t *tgBot) ListenTelegramUpdates() {
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
				t.Commander.Start(update.Message, t.services.UserService)
			case "users":
				t.Commander.ShowUsers(update.Message, t.services.UserService)
			case menu.FAQ:
				t.Commander.FAQ(update.Message)
			case menu.Setup:
				t.Commander.Setup(update.Message)
			case "ping":
				t.Commander.Pong(update.Message)
			}

			if menu.IsMenuButton(update.Message.Text) {
				t.Commander.Delete(update.Message)
			}
		}

		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := t.bot.Request(callback); err != nil {
				zap.L().Error("Error receiving response from callback with id", zap.Error(err), zap.String("id", update.CallbackQuery.ID))
				continue
			}

			switch update.CallbackQuery.Data {
			case keypads.BackData:
				t.Commander.Back(update.CallbackQuery)
			case keypads.OperationModeFAQData:
				t.Commander.OperatingModeInfo(update.CallbackQuery)
			}
		}
	}
}
