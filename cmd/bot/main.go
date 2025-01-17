package main

import (
	"air-quality-notifyer/internal/app/telegram"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/db"
	"air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/service/districts"
	"air-quality-notifyer/internal/service/sensor"
	"air-quality-notifyer/internal/service/user"
	"context"
	_ "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := config.NewApplicationConfig()
	database := db.NewDB(cfg)
	districtRepository := repository.NewDistrictRepository(database)
	userRepository := repository.NewUserRepository(database)
	sensorRepository := repository.NewSensorRepository(database)

	districtsList, err := districtRepository.GetAllDistricts()
	if err != nil {
		log.Panicln(fmt.Errorf("panic getting districts list: %w", err))
	}
	ctx = context.WithValue(ctx, "districts", districtsList)

	userService := user.NewUserService(userRepository)
	districtService := districts.NewDistrictService(districtRepository)
	sensorService := sensor.NewSensorService(ctx, sensorRepository, districtService)

	services := telegram.BotServices{
		UserService:   userService,
		SensorService: sensorService,
	}
	bot := telegram.InitTelegramBot(services, cfg)
	bot.ListenForUpdates()
	sensorService.FetchSensorsEveryHour()
	sensorService.InvalidateSensorsEveryday()

	<-ctx.Done()
}
