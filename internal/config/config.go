package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	App         AppConfig
	DB          DBConfig
	Cache       CacheConfig
	Development bool
}

type DBConfig struct {
	Host     string
	User     string
	Name     string
	Password string
}

type AppConfig struct {
	TelegramToken   string
	HttpServerPort  string
	WebhookHost     string
	AdminTelegramId string
}

type CacheConfig struct {
	Address  string
	Password string
	DBIndex  int
}

func New() Config {
	//TODO Подумать о .env.prod инъекции тут
	if err := godotenv.Load(); err != nil {
		zap.L().Fatal("Error loading environment variables", zap.Error(err))
	}

	var config = Config{
		App: AppConfig{
			TelegramToken:   os.Getenv("TELEGRAM_SECRET"),
			WebhookHost:     os.Getenv("WEBHOOK_HOST"),
			HttpServerPort:  os.Getenv("WEBHOOK_PORT"),
			AdminTelegramId: os.Getenv("ADMIN_TELEGRAM_ID"),
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Name:     os.Getenv("DB_NAME"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Cache: CacheConfig{
			Address:  os.Getenv("REDIS_HOST"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DBIndex:  0,
		},
		Development: os.Getenv("DEVELOPMENT") == "1",
	}

	return config
}
