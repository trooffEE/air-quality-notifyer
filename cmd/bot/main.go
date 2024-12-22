package main

import (
	"air-quality-notifyer/internal/app/telegram"
	"air-quality-notifyer/internal/db"
	"air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/service/districts"
	"air-quality-notifyer/internal/service/sensor"
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
	sensorRepository := repository.NewSensorRepository(database)

	userService := user.NewUserService(userRepository)
	districtService := districts.NewDistrictService(districtRepository)
	sensorService := sensor.NewSensorService(sensorRepository, districtService)

	services := telegram.BotServices{
		UserService:   userService,
		SensorService: sensorService,
	}
	bot := telegram.InitTelegramBot(services)
	bot.ListenForUpdates()

	sensorService.FetchSensorsEveryHour()
	sensorService.InvalidateSensorsEveryday()

	<-ctx.Done()
}
