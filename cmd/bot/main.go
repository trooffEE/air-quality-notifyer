package main

import (
	"air-quality-notifyer/internal/app/server"
	"air-quality-notifyer/internal/app/telegram"
	"air-quality-notifyer/internal/cache"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/db"
	rDistricts "air-quality-notifyer/internal/db/repository/districts"
	rSensor "air-quality-notifyer/internal/db/repository/sensor"
	rUser "air-quality-notifyer/internal/db/repository/user"
	sDistricts "air-quality-notifyer/internal/service/districts"
	sSensor "air-quality-notifyer/internal/service/sensor"
	sUser "air-quality-notifyer/internal/service/user"
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
	districtRepository := rDistricts.New(database)
	userRepository := rUser.New(database)
	sensorRepository := rSensor.New(database)

	//Cache
	cacheClient := cache.New(cfg)

	userService := sUser.New(userRepository)
	districtService := sDistricts.New(districtRepository)
	sensorService := sSensor.New(sensorRepository, districtService, cacheClient)

	httpShutdown := server.Init(ctx, cfg)

	services := telegram.BotServices{
		UserService:   userService,
		SensorService: sensorService,
	}

	bot := telegram.Init(services, cfg)

	go bot.ListenSensors()
	go bot.ListenUpdates()

	sensorService.StartGettingTrustedSensorsEveryHour()
	sensorService.StartInvalidatingSensorsPeriodically()

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
