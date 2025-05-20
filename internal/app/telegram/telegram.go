package telegram

import (
	"air-quality-notifyer/internal/app/commands"
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
	Commander *commands.Commander
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

	commander := commands.NewCommander(bot, cfg)
	if cfg.Development {
		bot.Debug = true

		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 30

		return &tgBot{
			bot:       bot,
			updates:   bot.GetUpdatesChan(updateConfig),
			services:  services,
			Commander: commander,
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
		Commander: commander,
	}
}

func (t *tgBot) ListenChangesInSensors() {
	t.services.SensorService.ListenChangesInSensors(t.notifyUsersAboutSensors)
}

/*
*

	tgbotapi.BotCommand{
				Command:     "configure",
				Description: "⚙️ Настроить бота",
			},
			tgbotapi.BotCommand{
				Command:     "faq",
				Description: "❓ Ответы на вопросы",
			},
			tgbotapi.BotCommand{
				Command:     "operating_mode",
				Description: "❓ Режимы работы",
			},
*/
func (t *tgBot) ListenTelegramUpdates() {
	cfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "🌀 Перезапустить бота",
		},
	)
	_, err := t.bot.Request(cfg)

	for update := range t.updates {
		if update.Message == nil {
			continue
		}
		if err != nil {
			zap.L().Error("commands request error", zap.Error(err))
			continue
		}
		zap.L().Info(
			"client message",
			zap.String("msg", update.Message.Text),
			zap.String("username", update.Message.From.UserName),
		)

		switch update.Message.Command() {
		case "users":
			t.Commander.ShowUsers(update.Message, t.services.UserService)
		case "faq":
			t.Commander.FAQ(update.Message)
		case "operating_mode":
			t.Commander.OperatingMode(update.Message)
		case "start":
			t.Commander.Start(update.Message, t.services.UserService)
		case "configure":
			t.Commander.Configure(update.Message)
		case "ping":
			t.Commander.Pong(update.Message)
		}
	}
}
