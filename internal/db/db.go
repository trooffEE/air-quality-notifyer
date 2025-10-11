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
)

func New(cfg config.Config) *sqlx.DB {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Name,
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
