package main

import (
	"air-quality-notifyer/bot"
	"air-quality-notifyer/pkg"
	"air-quality-notifyer/pkg/repository"
	"air-quality-notifyer/pkg/service"
	"air-quality-notifyer/sensor"
	_ "database/sql"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	db, err := pkg.NewDB()
	if err != nil {
		log.Fatal(err, 9)
	}
	psqlRepo := repository.NewUserRepository(db)
	usrService := service.NewUserService(psqlRepo)
	fmt.Println(usrService)
	bot.InitTelegramBot().ListenForUpdates()
	sensor.GetSensorsDataOnceIn("0 * * * *")

	//c := make(chan os.Signal)
	//signal.Notify(c, os.Kill)
	//
	//go func() {
	//	select {
	//	case <-c:
	//		fmt.Println("test")
	//		os.Exit(1)
	//	}
	//}()
	//select {}
}
