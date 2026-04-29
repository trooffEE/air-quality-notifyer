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
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type Services struct {
	User   user.Interface
	Sensor sensor.Interface
	Bot    *tgbotapi.BotAPI
}

func Init(ctx context.Context, cfg config.Config, services Services) func() {
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

	return func() {
		_ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := server.Shutdown(_ctx); err != nil {
			zap.L().Fatal("http server failed to shutdown", zap.Error(err))
		}
		wg.Wait()
	}
}
