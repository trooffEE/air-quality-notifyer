package main

import (
	"air-quality-notifyer/bot"
	"air-quality-notifyer/sensor"
	_ "database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	fmt.Println("TEST")
	connString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal(err, connString)
	}
	result, err := db.Exec("SELECT * from user_table")
	if err != nil {
		log.Fatal(err, 9)
	}
	fmt.Println(result, "test")
	bot.InitTelegramBot().ListenForUpdates()
	sensor.GetSensorsDataOnceIn("0 * * * *")
}
