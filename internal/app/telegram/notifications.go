package telegram

import (
	"air-quality-notifyer/internal/service/sensor"
	"fmt"
	"log"
	"time"
)

func (t *tgBot) notifyUsersAboutSensors(sensors []sensor.AirqualitySensor) {
	var messages []string
	for _, s := range sensors {
		if s.AQIPM10WarningIndex > 1 || s.AQIPM25WarningIndex > 1 {
			t, err := time.Parse("2006-01-02 15", s.Date)
			if err != nil {
				log.Printf("Error parsing date %#v", err)
				return
			}
			loc, _ := time.LoadLocation("Asia/Novosibirsk")
			sDate := t.In(loc)
			message := fmt.Sprintf("<b>В районе - %s</b> 🏠\n\nЗа прошедший час - для времени %s 🕛 \n\nЗафиксировано значительное ухудшение качества воздуха - уровень опасности \"%s\"\n\n<b>AQI(PM10): %.2f  - %s\nAQI(PM2.5): %.2f - %s</b>\n\nПодробнее (отматать 1 час назад): %s",
				s.District, sDate.Format("02.01.2006 15:04"), s.DangerLevel,
				s.AQIPM10, s.AQIPM10Analysis,
				s.AQIPM25, s.AQIPM25Analysis, s.SourceLink,
			)
			messages = append(messages, message)
		}
	}

	userIds := *t.services.UserService.GetUsersIds()
	for _, id := range userIds {
		for _, message := range messages {
			err := t.Commander.DefaultSend(id, message)
			if err != nil && err.Code == 403 {
				t.services.UserService.DeleteUser(id)
				break
			}
		}
	}
}
