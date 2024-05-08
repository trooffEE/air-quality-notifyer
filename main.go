package main

import (
	"air-quality-notifyer/bot"
	"air-quality-notifyer/sensor"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
	"log"
)

func main() {
	sensors := sensor.NewSensorsData()
	bot.NewTelegramBot()
	c := cron.New()

	_, err := c.AddFunc("0 * * * *", func() {
		sensor.FetchSensorsData(&sensors)
	})
	if err != nil {
		log.Panic(err)
	}
	c.Start()

	select {}
}
