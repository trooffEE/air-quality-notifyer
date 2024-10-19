package telegram

import (
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/sensor"
	"fmt"
	"log"
	"time"
)

func (t *tgBot) notifyUsersAboutSensors(sensors []sensor.Data) {
	var messages []string
	for _, s := range sensors {
		if s.AQIPM10WarningIndex > 1 || s.AQIPM25WarningIndex > 1 {
			t, err := time.Parse("2006-01-02 15", s.Date)
			if err != nil {
				log.Printf("Error parsing date %s", s.Date)
				return
			}
			loc, _ := time.LoadLocation("Asia/Novosibirsk")
			sDate := t.In(loc)
			// TODO
			//if !time.Now().In(loc).Equal(sDate) {
			//	fmt.Printf("Sensor with ID - %d - is outdated - TODO Logic remove it from grasp\n", s.Id)
			//	return
			//}
			message := fmt.Sprintf("<b>В районе - %s</b> 🏠\n\nЗа прошедший час - для времени %s 🕛 \n\nЗафиксировано значительное ухудшение качества воздуха - уровень опасности \"%s\"\n\n<b>AQI(PM10): %.2f  - %s\nAQI(PM2.5): %.2f - %s</b>\n\nПодробнее (отматать 1 час назад): %s",
				s.GetFormatedDistrictName(), sDate.Format("02.01.2006 15:04"), s.DangerLevel,
				s.AQIPM10, s.AQIPM10Analysis,
				s.AQIPM25, s.AQIPM25Analysis, s.SourceLink,
			)
			messages = append(messages, message)
		}
	}

	for _, message := range messages {
		// TODO Create notify for each individual user
		t.Commander.DefaultSend(config.Cfg.GetTestTelegramChatID2(), message)
		t.Commander.DefaultSend(config.Cfg.GetTestTelegramChatID2(), message)
		t.Commander.DefaultSend(config.Cfg.GetTestTelegramChatID3(), message)
	}
}
