package main

import (
	"air-quality-notifyer/bot"
	"air-quality-notifyer/pkg"
	"air-quality-notifyer/pkg/repository"
	"air-quality-notifyer/pkg/service"
	"air-quality-notifyer/sensor"
	_ "database/sql"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	db, err := pkg.NewDB()
	if err != nil {
		log.Fatal(err, 9)
	}
	psqlRepo := repository.NewUserRepository(db)
	usrService := service.NewUserService(psqlRepo)

	services := bot.BotServices{
		UserService: usrService,
	}
	bot.InitTelegramBot(services).ListenForUpdates()
	sensor.GetSensorsDataOnceIn("0 * * * *")
	select {}
}
