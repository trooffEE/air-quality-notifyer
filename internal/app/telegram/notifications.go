package telegram

import (
	"air-quality-notifyer/internal/app/commander/api"
	mSensor "air-quality-notifyer/internal/service/sensor/model"
	"fmt"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

func (t *tgBot) notifyUsers(sensors []mSensor.Sensor) {
	messages := newUserMessages(sensors)

	ids := t.services.UserService.GetUsersIds()
	for _, id := range ids {
		for _, message := range messages {
			msg := tgbotapi.NewMessage(id, message)
			payload := api.MessageConfig{Msg: msg}
			if err := t.Commander.API.Send(payload); err != nil && err.Code == 403 {
				t.services.UserService.DeleteUser(id)
				break
			}
		}
	}
}

func newUserMessages(sensors []mSensor.Sensor) []string {
	var messages []string
	for _, sensor := range sensors {
		if sensor.IsDangerousLevelDetected() {
			msg := prepareDangerousLevelMessage(sensor)
			messages = append(messages, msg)
		}
	}
	return messages
}

func prepareDangerousLevelMessage(s mSensor.Sensor) string {
	pollution := s.GetPollutionData()
	if pollution == nil {
		return ""
	}

	t, err := time.Parse("2006-01-02 15", s.Date)
	if err != nil {
		zap.L().Error("failed to parse time", zap.Error(err))
		return ""
	}

	loc, err := time.LoadLocation("Asia/Novosibirsk") // TODO not good i load it on every message + user (n^2), needs one point of Load such as in commander
	if err != nil {
		zap.L().Error("failed to load timezone", zap.Error(err))
		return ""
	}

	date := t.In(loc).Format("02.01.2006 15:04")
	return fmt.Sprintf(
		"<b>В районе - %s</b> 🏠\n\nЗа прошедший час - для времени %s 🕛 \n\nЗафиксировано значительное ухудшение качества воздуха - уровень опасности \"%s\"\n\n<b>AQI(PM10): %d\nAQI(PM2.5): %d</b>\n\nПодробнее: %s",
		s.District, date, pollution.Name,
		s.Aqi10, s.Aqi25, s.SourceLink,
	)
}
