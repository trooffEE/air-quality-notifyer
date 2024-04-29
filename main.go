package main

import (
	"air-quality-notifyer/entity"
	sensor "air-quality-notifyer/handlers"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
)

func main() {
	c := cron.New()
	sensors := entity.NewSensorsData()
	sensor.FetchSensorsData(&sensors)
	fmt.Println("MAIN TEST", sensors)
	//c.AddFunc("@every 1m", func() {
	//	sensor.FetchSensorsData(&sensors)
	//	fmt.Println("In Cron!", sensors)
	//})
	c.Start()

	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
