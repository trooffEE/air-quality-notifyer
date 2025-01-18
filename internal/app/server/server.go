package server

import (
	"air-quality-notifyer/internal/config"
	"context"
	"errors"
	"fmt"
	"log"
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
			log.Fatalf("http server failed to start on port: %s %v\n", cfg.HttpServerPort, err)
		}
	}()
	log.Println(fmt.Sprintf("🏆 http server started on port %s", cfg.HttpServerPort))

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("http server shutdown failed: %v", err)
		}
		wg.Wait()
	}
}
