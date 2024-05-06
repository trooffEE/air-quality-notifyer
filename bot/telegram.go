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
	sensorsData [][]sensor.Data
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

	var tgBot *TelegramBot = &TelegramBot{API: bot}

	//anyNewSensorUpdate := <-sensor.TelegramBotNotifySensorChangeChanel

	//fmt.Println("Как это работает?", anyNewSensorUpdate)

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			if !(update.Message.IsCommand() || IsPublicCommandProvided(update.Message.Text)) {
				tgBot.MessageSend(update.Message.Chat.ID, GetMessageByMention(NotCommandMessage))
				continue
			}

			if isShowAQIForChosenDistrictCommandProvided(update.Message.Text) {
				tgBot.MessageSend(update.Message.Chat.ID, GetMessageWithAQIStatsForChosenDistrict())
				continue
			}
		}
	}()

	return tgBot
}

// MessageSend is basically shortcut to cover potential errors upon sending message + for DRY principle
func (t *TelegramBot) MessageSend(chatID int64, messagePayload string) {
	msg := tgbotapi.NewMessage(chatID, messagePayload)
	_, err := t.API.Send(msg)
	if err != nil {
		log.Print(fmt.Sprintf("Error appeared upon sending message to user %d with message %s", chatID, messagePayload))
	}
}

func (t *TelegramBot) ConsumeSensorsData(data [][]sensor.Data) {
	t.sensorsData = data
	t.notifyUsersAboutSensorConsume()
}

func (t *TelegramBot) notifyUsersAboutSensorConsume() {
	//// TODO In future make all users get notification about districts AQI they subscribed to
	//msg := tgbotapi.NewMessage(cfg.GetTestTelegramChatID(), "Placeholder")
	//t.API.Send(msg)
}
