package cache

import (
	"os"
)

type Config struct {
	Address  string
	Password string
	DBIndex  int
}

func NewRedisConfig() Config {
	return Config{
		Address:  os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DBIndex:  0,
	}
}
