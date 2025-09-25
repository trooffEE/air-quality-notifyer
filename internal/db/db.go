package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Config struct {
	host     string
	user     string
	dbname   string
	password string
}

func NewDB() *sqlx.DB {
	config := NewConfig()
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		config.user, config.password, config.host, config.dbname,
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

func NewConfig() *Config {
	dbConfig := Config{
		host:     os.Getenv("DB_HOST"),
		user:     os.Getenv("DB_USER"),
		dbname:   os.Getenv("DB_NAME"),
		password: os.Getenv("DB_PASSWORD"),
	}

	return &dbConfig
}
