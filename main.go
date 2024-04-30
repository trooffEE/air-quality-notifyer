package main

import (
	"air-quality-notifyer/bot"
	"air-quality-notifyer/entity"
	sensor "air-quality-notifyer/handlers"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
)

func main() {
	bot.NewTelegramBot()
	sensors := entity.NewSensorsData()
	c := cron.New()
	c.AddFunc("1 * * * * *", func() {
		sensor.FetchSensorsData(&sensors)
	})
	c.Start()
}
