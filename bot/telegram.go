/*TODO Переписать пакет на поддержку web-hooks*/
package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func NewTelegramBot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_SECRET"))
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true

	cert, err := os.ReadFile("cert.pem")
	if err != nil {
		log.Fatal("Provide cert to establish correct WebSocket connection")
	}
	certFileData := tgbotapi.FileBytes{Name: "cert.pem", Bytes: cert}

	wh, _ := tgbotapi.NewWebhookWithCert("https://t0ffee-dev.ru/"+bot.Token, certFileData)
	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(&err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	//updates := bot.ListenForWebhook("/" + bot.Token)
	//go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

}

func getTelegramUpdates() {

}
