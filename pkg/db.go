package pkg

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
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

func NewDB(args ...func(*Config)) (*sqlx.DB, error) {
	config := NewConfig(args)

	connString := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=disable",
		config.host, config.user, config.dbname, config.password,
	)
	db, err := sqlx.Connect("postgres", connString)

	if err != nil {
		fmt.Println(fmt.Errorf("Establish failed: %w ", ErrConnectionFailed))
		return nil, err
	}

	return db, nil
}

func NewConfig(args []func(*Config)) *Config {
	dbConfig := Config{
		host:     "airquality-db-container",
		user:     os.Getenv("DB_USER"),
		dbname:   os.Getenv("DB_NAME"),
		password: os.Getenv("DB_PASSWORD"),
	}
	for _, fn := range args {
		fn(&dbConfig)
	}

	return &dbConfig
}

func WithHost(host string) func(config *Config) {
	return func(s *Config) {
		s.host = host
	}
}
