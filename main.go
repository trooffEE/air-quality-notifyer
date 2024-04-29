package main

import (
	"air-quality-notifyer/bot"
	"air-quality-notifyer/entity"
	sensor "air-quality-notifyer/handlers"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
)

func main() {
	bot.NewTelegramBot()
	sensors := entity.NewSensorsData()
	c := cron.New()
	c.AddFunc("1 * * * * *", func() {
		sensor.FetchSensorsData(&sensors)
	})
	c.Start()

	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
