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

	cfg := config.NewConfig()
	initLogger(cfg)

	//DB
	database := db.NewDB(cfg)
	districtRepository := repository.NewDistrictRepository(database)
	userRepository := repository.NewUserRepository(database)
	sensorRepository := repository.NewSensorRepository(database)

	//Cache
	cacheClient := cache.NewCacheClient(cfg)

	userService := user.NewUserService(userRepository)
	districtService := districts.NewDistrictService(districtRepository)
	sensorService := sensor.NewSensorService(sensorRepository, districtService, cacheClient)

	httpShutdown := server.InitHttpServer(ctx, cfg)

	services := telegram.BotServices{
		UserService:   userService,
		SensorService: sensorService,
	}

	bot := telegram.InitTelegramBot(services, cfg)

	go bot.ListenSensorsUpdates()
	go bot.ListenTelegramUpdates()

	sensorService.GetTrustedSensorsEveryHour()
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
