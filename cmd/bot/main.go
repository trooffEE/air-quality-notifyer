package main

import (
	"air-quality-notifyer/internal/app/telegram"
	"air-quality-notifyer/internal/db"
	"air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/sensor"
	"air-quality-notifyer/internal/service/districts"
	"air-quality-notifyer/internal/service/user"
	"context"
	_ "database/sql"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	database := db.NewDB()
	userRepository := repository.NewUserRepository(database)
	districtRepository := repository.NewDistrictRepository(database)
	userService := user.NewUserService(userRepository)
	districtService := districts.NewDistrictService(districtRepository)

	services := telegram.BotServices{UserService: userService}
	bot := telegram.InitTelegramBot(services)
	bot.ListenForUpdates()
	//TODO Rewrite on service
	sensor.GetSensorsDataOnceIn("0 * * * *")
	districtService.Get

	<-ctx.Done()
}
