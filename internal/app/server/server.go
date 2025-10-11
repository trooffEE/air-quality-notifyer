package server

import (
	"air-quality-notifyer/internal/config"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

func Init(ctx context.Context, cfg config.Config) func() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.HttpServerPort),
		Handler: nil,
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
