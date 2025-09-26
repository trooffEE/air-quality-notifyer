package cache

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func NewCacheClient() *redis.Client {
	cfg := newRedisConfig()
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DBIndex,
	})
}

type Config struct {
	Address  string
	Password string
	DBIndex  int
}

func newRedisConfig() Config {
	return Config{
		Address:  os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DBIndex:  0,
	}
}
