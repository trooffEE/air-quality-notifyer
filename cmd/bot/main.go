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
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		zap.L().Fatal("Error loading environment variables", zap.Error(err))
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.NewApplicationConfig()
	initLogger(cfg)

	database := db.NewDB()
	client := cache.NewCacheClient()

	//TODO Test string
	client.Set(ctx, "test", 123, time.Hour)

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

func initLogger(cfg config.ApplicationConfig) {
	var logger *zap.Logger
	if cfg.Development {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
}
