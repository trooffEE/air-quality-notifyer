package main

import (
	"air-quality-notifyer/internal/app/server"
	"air-quality-notifyer/internal/app/telegram"
	"air-quality-notifyer/internal/cache"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/db"
	"air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/service/districts"
	"air-quality-notifyer/internal/service/sensor"
	"air-quality-notifyer/internal/service/user"
	"context"
	_ "database/sql"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.New()
	initLogger(cfg)

	//DB
	database := db.New(cfg)
	districtRepository := repository.NewDistrictRepository(database)
	userRepository := repository.NewUserRepository(database)
	sensorRepository := repository.NewSensorRepository(database)

	//Cache
	cacheClient := cache.New(cfg)

	userService := user.New(userRepository)
	districtService := districts.New(districtRepository)
	sensorService := sensor.New(sensorRepository, districtService, cacheClient)

	httpShutdown := server.Init(ctx, cfg)

	services := telegram.BotServices{
		UserService:   userService,
		SensorService: sensorService,
	}

	bot := telegram.Init(services, cfg)

	go bot.ListenSensorsUpdates()
	go bot.ListenTelegramUpdates()

	sensorService.StartGettingTrustedSensorsEveryHour()
	sensorService.InvalidateSensorsPeriodically()

	<-ctx.Done()
	zap.L().Info("starting application shutdown...")
	httpShutdown()
	zap.L().Info("http server is down")
}

func initLogger(cfg config.Config) {
	var logger *zap.Logger
	if cfg.Development {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
}
