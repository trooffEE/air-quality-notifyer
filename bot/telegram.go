package bot

import (
	"air-quality-notifyer/config"
	"air-quality-notifyer/sensor"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

var cfg = config.InitConfig()

type tgBot struct {
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func InitTelegramBot() *tgBot {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	if cfg.Development {
		bot.Debug = true
		go http.ListenAndServe(fmt.Sprintf(":%s", cfg.WebhookPort), nil)
		return &tgBot{bot, nil}
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
		log.Println(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/webhook" + bot.Token)

	go http.ListenAndServe(fmt.Sprintf(":%s", cfg.WebhookPort), nil)

	return &tgBot{bot, updates}
}

func (t *tgBot) ListenForUpdates() {
	go sensor.ListenChangesInSensors(t.notifyUsersAboutSensors)
	go t.handleUpdates()
}

func (t *tgBot) handleUpdates() {
	if cfg.Development {
		go t.handleUpdatesLocally()
		return
	}
	go t.handleWebhookUpdates()
}

func (t *tgBot) handleUpdatesLocally() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := t.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := t.bot.Send(msg); err != nil {
			panic(err)
		}
	}
}

func (t *tgBot) handleWebhookUpdates() {
	for update := range t.updates {
		if update.Message == nil {
			continue
		}

		fmt.Println("test", update.Message.Text)

		if !(update.Message.IsCommand() || isPublicCommandProvided(update.Message.Text)) {
			t.messageSend(update.Message.Chat.ID, getMessageByMention(notCommandMessage))
			continue
		}

		if isShowAQIForChosenDistrictCommandProvided(update.Message.Text) {
			t.messageSend(update.Message.Chat.ID, getMessageWithAQIStatsForChosenDistrict())
			continue
		}
	}
}
