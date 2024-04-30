package bot

import (
	"air-quality-notifyer/config"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

var cfg = config.InitConfig()

func NewTelegramBot() {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	// Set the webhook for the bot
	wh, err := tgbotapi.NewWebhook(fmt.Sprintf("https://%s/webhook", cfg.WebhookHost))
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

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe(fmt.Sprintf(":%s", cfg.WebhookPort), nil)
	fmt.Println("HUH", updates)
	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
