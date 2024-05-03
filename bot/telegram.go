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

type TelegramBot struct {
	API         *tgbotapi.BotAPI
	sensorsData [][]sensor.SensorData
}

func NewTelegramBot() *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
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
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/webhook" + bot.Token)

	go http.ListenAndServe(fmt.Sprintf(":%s", cfg.WebhookPort), nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !(update.Message.IsCommand() || IsPublicCommandProvided(update.Message.Text)) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, GetMessageByMention(NotCommandMessage))
			fmt.Println("test", update.Message.Chat.ID)
			bot.Send(msg)
			continue
		}

		if isShowAQIForChosenDistrictCommandProvided(update.Message.Text) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, GetMessageWithAQIStatsForChosenDistrict())
			bot.Send(msg)
			continue
		}
	}

	return &TelegramBot{API: bot}
}

func (t *TelegramBot) ConsumeSensorsData(data [][]sensor.SensorData) {
	t.sensorsData = data
	t.notifyUsersAboutSensorConsume()
}

func (t *TelegramBot) notifyUsersAboutSensorConsume() {
	//msg := tgbotapi.NewMessage()
	// TODO In future make all users get notification about districts AQI they subscribed to
	//t.API.Send()
}
