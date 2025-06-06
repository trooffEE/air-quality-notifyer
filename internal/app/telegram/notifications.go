package telegram

import (
	"air-quality-notifyer/internal/app/commander"
	s "air-quality-notifyer/internal/service/sensor"
	"fmt"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
	"time"
)

func prepareDangerousLevelMessage(s s.AqiSensor) string {
	pollutionLevel := s.GetExtendedPollutionLevel()
	if pollutionLevel == nil {
		return ""
	}

	t, err := time.Parse("2006-01-02 15", s.Date)
	if err != nil {
		zap.L().Error("failed to parse time", zap.Error(err))
		return ""
	}

	loc, err := time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		zap.L().Error("failed to load timezone", zap.Error(err))
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
		zap.L().Error("failed to load timezone", zap.Error(err))
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
	userIds := t.services.UserService.GetUsersIds()
	for _, id := range userIds {
		for _, message := range messages {
			msg := tgbotapi.NewMessage(id, message)
			err := t.Commander.Send(commander.SendPayload{Msg: msg, DisableNotification: isSilentMessage})
			if err != nil && err.Code == 403 {
				t.services.UserService.DeleteUser(id)
				break
			}
		}
	}
}
