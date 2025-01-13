package db

import (
	"air-quality-notifyer/internal/config"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

type Config struct {
	host     string
	user     string
	dbname   string
	password string
}

func NewDB(args ...func(*Config)) *sqlx.DB {
	cfg := NewConfig(args)

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.user, cfg.password, cfg.host, cfg.dbname,
	)
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatalln("Failed to establish DB connection")
	}

	m, err := migrate.New(
		"file://internal/db/migrations",
		connString,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("🏆 Migrations applied successfully!")

	return db
}

func NewConfig(args []func(*Config)) *Config {
	dbConfig := Config{
		host:     "airquality-db-container",
		user:     os.Getenv("DB_USER"),
		dbname:   os.Getenv("DB_NAME"),
		password: os.Getenv("DB_PASSWORD"),
	}

	if config.Cfg.Development {
		dbConfig.host = "localhost"
	}

	for _, fn := range args {
		fn(&dbConfig)
	}

	return &dbConfig
}
