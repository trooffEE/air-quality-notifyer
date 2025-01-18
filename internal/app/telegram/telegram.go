package telegram

import (
	"air-quality-notifyer/internal/app/commands"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/lib"
	"air-quality-notifyer/internal/service/sensor"
	"air-quality-notifyer/internal/service/user"
	"fmt"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"log"
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
	commander := commands.NewCommander(bot, cfg)
	if err != nil {
		log.Panic(err)
	}

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
		log.Panic(err)
	}
	_, err = bot.Request(wh)
	if err != nil {
		log.Panic(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Panic(err)
	}

	if info.LastErrorDate != 0 {
		lib.LogError("InitTelegramBot", "failed to init get info about webhook", err)
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

func (t *tgBot) ListenTelegramUpdates() {
	for update := range t.updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "users":
			t.Commander.ShowUsers(update.Message, t.services.UserService)
		case "help":
			t.Commander.Help(update.Message.Chat.ID)
		case "start":
			t.Commander.Start(update.Message, t.services.UserService)
		}
	}
}
