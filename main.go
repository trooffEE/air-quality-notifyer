package main

import (
	"air-quality-notifyer/bot"
	"air-quality-notifyer/sensor"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
)

func main() {
	sensors := sensor.NewSensorsData()
	bot.NewTelegramBot()

	c := cron.New()
	sensor.FetchSensorsData(&sensors)
	sensor.NotifyPackages(sensors)
	//c.AddFunc("@every 1m", func() {
	//	sensor.FetchSensorsData(&sensors)
	//})
	c.Start()

	select {}
}
