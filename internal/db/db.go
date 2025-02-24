package db

import (
	"air-quality-notifyer/internal/config"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	host     string
	user     string
	dbname   string
	password string
}

func NewDB(cfg config.ApplicationConfig) *sqlx.DB {
	dbConfig := NewConfig(cfg)

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		dbConfig.user, dbConfig.password, dbConfig.host, dbConfig.dbname,
	)
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		zap.L().Fatal("Failed to establish db connection", zap.Error(err))
	}

	m, err := migrate.New(
		"file://internal/db/migrations",
		connString,
	)
	if err != nil {
		zap.L().Fatal("Failed to create migrate instance", zap.Error(err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		zap.L().Fatal("Failed to run migrations", zap.Error(err))
	}

	zap.L().Info("🏆 Migrations applied successfully!")

	return db
}

func NewConfig(cfg config.ApplicationConfig) *Config {
	dbConfig := Config{
		host:     "airquality-db-container",
		user:     os.Getenv("DB_USER"),
		dbname:   os.Getenv("DB_NAME"),
		password: os.Getenv("DB_PASSWORD"),
	}

	if cfg.Development {
		dbConfig.host = "localhost"
	}

	return &dbConfig
}
