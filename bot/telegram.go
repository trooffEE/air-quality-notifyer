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

type tgBot struct {
	API     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func InitTelegramBot() *tgBot {
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
		log.Println(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/webhook" + bot.Token)

	go http.ListenAndServe(fmt.Sprintf(":%s", cfg.WebhookPort), nil)

	return &tgBot{bot, updates}
}

// messageSend is basically shortcut to cover potential errors upon sending message + for DRY principle
func (t *tgBot) messageSend(chatID int64, messagePayload string) {
	msg := tgbotapi.NewMessage(chatID, messagePayload)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := t.API.Send(msg)
	if err != nil {
		log.Print(fmt.Sprintf("Error appeared upon sending message to user %d with message %s", chatID, messagePayload))
	}
}

func (t *tgBot) handleSensorsUpdates() {
	sensor.ListenChangesInSensors(t.notifyUsersAboutSensors)
}

func (t *tgBot) handleWebhookUpdates() {
	for update := range t.updates {
		if update.Message == nil {
			continue
		}

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

func (t *tgBot) ListenForUpdates() {
	go t.handleSensorsUpdates()
	go t.handleWebhookUpdates()
}

func (t *tgBot) notifyUsersAboutSensors(sensors []sensor.Data) {
	var messages []string
	for _, s := range sensors {
		if s.AQIPM10WarningIndex > 1 || s.AQIPM25WarningIndex > 1 {
			t, err := time.Parse("2006-01-02 15", s.Date)
			if err != nil {
				log.Printf("Error parsing date %s", s.Date)
				return
			}
			loc, _ := time.LoadLocation("Asia/Novosibirsk")
			sDate := t.In(loc)
			// TODO
			//if !time.Now().In(loc).Equal(sDate) {
			//	fmt.Printf("Sensor with ID - %d - is outdated - TODO Logic remove it from grasp\n", s.Id)
			//	return
			//}
			message := fmt.Sprintf("<b>В районе - %s</b> 🏠\n\nЗа прошедший час - для времени %s 🕛 \n\nЗафиксировано значительное ухудшение качества воздуха - уровень опасности \"%s\"\n\n<b>AQI(PM10): %.2f  - %s\nAQI(PM2.5): %.2f - %s</b>\n\nПодробнее (отматать 1 час назад): %s",
				s.GetFormatedDistrictName(), sDate.Format("02.01.2006 15:04"), s.DangerLevel,
				s.AQIPM10, s.AQIPM10Analysis,
				s.AQIPM25, s.AQIPM25Analysis, s.SourceLink,
			)
			messages = append(messages, message)
		}
	}

	for _, message := range messages {
		// TODO Create notify for each individual user
		t.messageSend(cfg.GetTestTelegramChatID(), message)
		t.messageSend(cfg.GetTestTelegramChatID2(), message)
		t.messageSend(cfg.GetTestTelegramChatID3(), message)
	}
}
