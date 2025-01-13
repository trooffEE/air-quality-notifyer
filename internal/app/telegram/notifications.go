package telegram

import (
	s "air-quality-notifyer/internal/service/sensor"
	"fmt"
	"log"
	"time"
)

func prepareDangerousLevelMessage(s s.AqiSensor) string {
	pollutionLevel := s.GetExtendedPollutionLevel()
	if pollutionLevel == nil {
		return ""
	}

	t, err := time.Parse("2006-01-02 15", s.Date)
	if err != nil {
		log.Printf("Error parsing date %#v", err)
		return ""
	}

	loc, err := time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		log.Printf("Error loading timezone %#v", err)
		return ""
	}

	date := t.In(loc).Format("02.01.2006 15:04")
	return fmt.Sprintf(
		"<b>В районе - %s</b> 🏠\n\nЗа прошедший час - для времени %s 🕛 \n\nЗафиксировано значительное ухудшение качества воздуха - уровень опасности \"%s\"\n\n<b>AQI(PM10): %d\nAQI(PM2.5): %d</b>\n\nПодробнее: %s",
		s.District, date, pollutionLevel.Name,
		s.Aqi10, s.Aqi25, s.SourceLink,
	)
}

func getTimezoneHourTime() int {
	loc, err := time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		log.Printf("Error loading timezone %#v", err)
		return -1
	}

	return time.Now().In(loc).Hour()
}

func (t *tgBot) notifyUsersAboutSensors(sensors []s.AqiSensor) {
	var messages []string
	for _, sensor := range sensors {
		if sensor.IsDangerousLevelDetected() {
			msg := prepareDangerousLevelMessage(sensor)
			messages = append(messages, msg)
		}
	}

	hour := getTimezoneHourTime()
	isSilentMessage := false
	if hour < 8 && hour >= 0 {
		isSilentMessage = true
	}
	userIds := *t.services.UserService.GetUsersIds()
	for _, id := range userIds {
		for _, message := range messages {
			err := t.Commander.DefaultSend(id, message, isSilentMessage)
			if err != nil && err.Code == 403 {
				t.services.UserService.DeleteUser(id)
				break
			}
		}
	}
}
