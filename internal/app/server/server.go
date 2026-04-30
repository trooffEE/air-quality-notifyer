package server

import (
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/service/sensor"
	"air-quality-notifyer/internal/service/user"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type Services struct {
	User   user.Interface
	Sensor sensor.Interface
	Bot    *tgbotapi.BotAPI
}

func Init(cfg config.Config, services Services) func(context.Context) {
	mux := http.NewServeMux()
	newMapHandler(cfg, services).Register(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.HttpServerPort),
		Handler: mux,
	}

	var wg sync.WaitGroup
	wg.Go(func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			zap.L().Fatal("http server failed to start on port", zap.String("port", cfg.App.HttpServerPort), zap.Error(err))
		}
	})
	zap.L().Info("🏆 http server started on port", zap.String("port", cfg.App.HttpServerPort))

	return func(ctx context.Context) {
		if err := server.Shutdown(ctx); err != nil {
			zap.L().Error("http server failed to shutdown", zap.Error(err))
			if closeErr := server.Close(); closeErr != nil {
				zap.L().Error("http server failed to close", zap.Error(closeErr))
			}
		}
		wg.Wait()
	}
}
