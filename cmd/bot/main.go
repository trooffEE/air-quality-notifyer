package main

import (
	"air-quality-notifyer/internal/app/server"
	"air-quality-notifyer/internal/app/telegram"
	"air-quality-notifyer/internal/app/telegram/commander"
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
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.New()
	logger := initLogger(cfg)
	defer func() {
		if err := logger.Sync(); err != nil {
			zap.L().Debug("failed to sync logger", zap.Error(err))
		}
	}()

	//DB
	database := db.New(ctx, cfg)
	defer func() {
		if err := database.Close(); err != nil {
			zap.L().Error("failed to close database", zap.Error(err))
		}
	}()

	districtRepository := rDistricts.New(database)
	userRepository := rUser.New(database)
	sensorRepository := rSensor.New(database)

	cacheClient := cache.New(cfg)
	defer func() {
		if err := cacheClient.Close(); err != nil {
			zap.L().Error("failed to close cache", zap.Error(err))
		}
	}()

	userService := sUser.New(userRepository)
	districtService := sDistricts.New(districtRepository, cacheClient)
	sensorService := sSensor.New(sensorRepository, districtService, cacheClient)

	services := commander.Services{
		User:     userService,
		Sensor:   sensorService,
		District: districtService,
		Cache:    cacheClient,
	}

	bot := telegram.Init(cfg, &services)
	httpShutdown := server.Init(cfg, server.Services{
		User:   userService,
		Sensor: sensorService,
		Bot:    bot.Commander.API,
	})
	botShutdown := bot.Start(ctx)

	trustedSensorsShutdown := sensorService.StartGettingTrustedSensorsEveryHour(ctx)
	invalidationShutdown := sensorService.StartInvalidatingSensorsPeriodically(ctx)

	<-ctx.Done()
	zap.L().Info("starting application shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	botShutdown(shutdownCtx)
	trustedSensorsShutdown(shutdownCtx)
	invalidationShutdown(shutdownCtx)
	httpShutdown(shutdownCtx)
	zap.L().Info("http server is down")
}

func initLogger(cfg config.Config) *zap.Logger {
	var logger *zap.Logger
	if cfg.Development {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	zap.ReplaceGlobals(logger)
	return logger
}
