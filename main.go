package main

import (
	"air-quality-notifyer/bot"
	"air-quality-notifyer/sensor"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
)

func main() {
	sensors := sensor.NewSensorsData()
	tgBot := bot.NewTelegramBot()

	c := cron.New()
	c.AddFunc("@every 1m", func() {
		sensor.FetchSensorsData(&sensors)
		tgBot.ConsumeSensorsData(sensors)
	})

	c.Start()
}
