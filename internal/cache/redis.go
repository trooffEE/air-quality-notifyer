package cache

import (
	"air-quality-notifyer/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewCacheClient(cfg config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Cache.Address,
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DBIndex,
	})
}
