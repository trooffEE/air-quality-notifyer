package main

import (
	"air-quality-notifyer/internal/app/telegram"
	"air-quality-notifyer/internal/db"
	"air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/sensor"
	"air-quality-notifyer/internal/service/user"
	_ "database/sql"
	_ "github.com/lib/pq"
)

func main() {
	database := db.NewDB()
	psqlRepo := repository.NewUserRepository(database)
	usrService := user.NewUserService(psqlRepo)

	services := telegram.BotServices{
		UserService: usrService,
	}
	telegram.InitTelegramBot(services).ListenForUpdates()
	sensor.GetSensorsDataOnceIn("0 * * * *")
	select {}
}
