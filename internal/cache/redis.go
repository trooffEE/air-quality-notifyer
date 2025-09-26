package cache

import (
	"github.com/redis/go-redis/v9"
)

func NewCacheClient() *redis.Client {
	cfg := NewRedisConfig()
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DBIndex,
	})
}
