package db

import (
	"air-quality-notifyer/internal/config"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

var (
	ErrConnectionFailed = errors.New("connection failed")
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
		"host=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.host, cfg.user, cfg.dbname, cfg.password,
	)
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Panicf("Establish failed: %w", ErrConnectionFailed)
	}

	db.MustExec(schema)

	return db
}

func NewConfig(args []func(*Config)) *Config {
	dbConfig := Config{
		host:     "airquality-db-container",
		user:     os.Getenv("DB_USER"),
		dbname:   os.Getenv("DB_NAME"),
		password: os.Getenv("DB_PASSWORD"),
	}

	//TODO
	if config.Cfg.Development {
		dbConfig.host = "localhost"
	}

	for _, fn := range args {
		fn(&dbConfig)
	}

	return &dbConfig
}
