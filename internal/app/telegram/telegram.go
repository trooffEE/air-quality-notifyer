package telegram

import (
	"air-quality-notifyer/internal/app/commands"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/sensor"
	"air-quality-notifyer/internal/service/user"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

type tgBot struct {
	bot       *tgbotapi.BotAPI
	updates   tgbotapi.UpdatesChannel
	services  BotServices
	Commander *commands.Commander
}

type BotServices struct {
	UserService *user.Service
}

func InitTelegramBot(services BotServices) *tgBot {
	bot, err := tgbotapi.NewBotAPI(config.Cfg.TelegramToken)
	commander := commands.NewCommander(bot)
	if err != nil {
		log.Panic(err)
	}

	go http.ListenAndServe(fmt.Sprintf(":%s", config.Cfg.WebhookPort), nil)

	if config.Cfg.Development {
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

	wh, err := tgbotapi.NewWebhook(fmt.Sprintf("https://%s/webhook%s", config.Cfg.WebhookHost, bot.Token))
	if err != nil {
		log.Panic(err)
	}
	_, err = bot.Request(wh)
	if err != nil {
		log.Panic(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Println(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/webhook" + bot.Token)

	return &tgBot{
		bot:       bot,
		updates:   updates,
		services:  services,
		Commander: commander,
	}
}

func (t *tgBot) ListenForUpdates() {
	go sensor.ListenChangesInSensors(t.notifyUsersAboutSensors)
	go t.handleUpdates()
}

func (t *tgBot) handleUpdates() {
	for update := range t.updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "help":
			t.Commander.Help(update.Message.Chat.ID)
		case "start":
			t.Commander.Start(update.Message, t.services.UserService)
		}
	}
}
