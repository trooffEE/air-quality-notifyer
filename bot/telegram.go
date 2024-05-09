package bot

import (
	"air-quality-notifyer/config"
	"air-quality-notifyer/sensor"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"time"
)

var cfg = config.InitConfig()

type TelegramBot struct {
	API *tgbotapi.BotAPI
}

func getUpdatesAboutSensors(tgBot *TelegramBot) {
	for update := range sensor.ChangesInAPIAppearedChannel {
		tgBot.notifyUsersAboutSensors(update)
	}
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

	go getUpdatesAboutSensors(tgBot)

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			if !(update.Message.IsCommand() || IsPublicCommandProvided(update.Message.Text)) {
				fmt.Println(update.Message.Text, update.Message.Chat.ID, update.Message.From.UserName)
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
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := t.API.Send(msg)
	if err != nil {
		log.Print(fmt.Sprintf("Error appeared upon sending message to user %d with message %s", chatID, messagePayload))
	}
}

func (t *TelegramBot) notifyUsersAboutSensors(sensors [][]sensor.Data) {
	var currentSensorsData []sensor.Data
	for _, sensorForDistrict := range sensors {
		currentSensorsData = append(currentSensorsData, sensorForDistrict[len(sensorForDistrict)-1])
	}

	var messages []string
	for _, s := range currentSensorsData {
		if s.AQIPM10WarningIndex > 1 || s.AQIPM25WarningIndex > 1 {
			loc, _ := time.LoadLocation("Asia/Novosibirsk")
			now := time.Now().In(loc)
			message := fmt.Sprintf("<b>В районе - %s</b> 🏠\n\nДля времени %s 🕛\n\nЗафиксировано значительное ухудшение качества воздуха - уровень опасности \"%s\"\n\n<b>AQI(PM10): %.2f  - %s\nAQI(PM2.5): %.2f - %s</b>\n\n%s",
				s.GetFormatedDistrictName(), now.Format("02.01.2006 15:04"), s.DangerLevel,
				s.AQIPM10, s.AQIPM10Analysis,
				s.AQIPM25, s.AQIPM25Analysis, s.AQIAnalysisRecommendations,
			)
			messages = append(messages, message)
		}
	}

	for _, message := range messages {
		// TODO Create notify for each individual user
		t.MessageSend(cfg.GetTestTelegramChatID(), message)
	}
}
