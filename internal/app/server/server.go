package server

import (
	"air-quality-notifyer/internal/config"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

func InitHttpServer(cfg config.ApplicationConfig) func() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HttpServerPort),
		Handler: nil,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			zap.L().Fatal("http server failed to start on port", zap.String("port", cfg.HttpServerPort), zap.Error(err))
		}
	}()
	zap.L().Info("🏆 http server started on port", zap.String("port", cfg.HttpServerPort))

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			zap.L().Fatal("http server failed to shutdown", zap.Error(err))
		}
		wg.Wait()
	}
}
