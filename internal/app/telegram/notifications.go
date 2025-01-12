package telegram

import (
	s "air-quality-notifyer/internal/service/sensor"
	"fmt"
	"log"
	"time"
)

func prepareDangerousLevelMessage(s s.AqiSensor) string {
	t, err := time.Parse("2006-01-02 15", s.Date)
	if err != nil {
		log.Printf("Error parsing date %#v", err)
		return ""
	}
	loc, _ := time.LoadLocation("Asia/Novosibirsk")
	date := t.In(loc).Format("02.01.2006 15:04")
	pollutionLevel := s.GetExtendedPollutionLevel()

	return fmt.Sprintf("<b>В районе - %s</b> 🏠\n\nЗа прошедший час - для времени %s 🕛 \n\nЗафиксировано значительное ухудшение качества воздуха - уровень опасности \"%s\"\n\n<b>AQI(PM10): %d\nAQI(PM2.5): %d</b>\n\nПодробнее: %s",
		s.District, date, pollutionLevel.Name,
		s.Aqi10, s.Aqi25, s.SourceLink,
	)
}

func (t *tgBot) notifyUsersAboutSensors(sensors []s.AqiSensor) {
	var messages []string
	for _, sensor := range sensors {
		if sensor.IsDangerousLevelDetected() {
			msg := prepareDangerousLevelMessage(sensor)
			messages = append(messages, msg)
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
