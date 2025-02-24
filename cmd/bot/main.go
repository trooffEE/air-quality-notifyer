package main

import (
	"air-quality-notifyer/internal/app/server"
	"air-quality-notifyer/internal/app/telegram"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/db"
	"air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/service/districts"
	"air-quality-notifyer/internal/service/sensor"
	"air-quality-notifyer/internal/service/user"
	"context"
	_ "database/sql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	cfg := config.NewApplicationConfig()
	database := db.NewDB(cfg)
	districtRepository := repository.NewDistrictRepository(database)
	userRepository := repository.NewUserRepository(database)
	sensorRepository := repository.NewSensorRepository(database)

	districtsList, err := districtRepository.GetAllDistricts()
	if err != nil {
		zap.L().Panic("Failed to get all districts", zap.Error(err))
	}
	ctx = context.WithValue(ctx, "districts", districtsList)

	userService := user.NewUserService(userRepository)
	districtService := districts.NewDistrictService(districtRepository)
	sensorService := sensor.NewSensorService(ctx, sensorRepository, districtService)

	services := telegram.BotServices{
		UserService:   userService,
		SensorService: sensorService,
	}

	httpShutdown := server.InitHttpServer(cfg)
	bot := telegram.InitTelegramBot(services, cfg)

	//TODO: think about wait group wrapping
	go bot.ListenChangesInSensors()
	go bot.ListenTelegramUpdates()

	sensorService.GetTrustedSensorsEveryHour()
	sensorService.InvalidateSensorsPeriodically()

	<-ctx.Done()
	zap.L().Info("starting application shutdown...")
	httpShutdown()
	zap.L().Info("http server is down")
}
