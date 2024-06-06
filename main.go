package main

import (
	"air-quality-notifyer/bot"
	"air-quality-notifyer/sensor"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	bot.InitTelegramBot().ListenForUpdates()
	sensor.GetSensorsDataOnceIn("0 * * * *")
}
