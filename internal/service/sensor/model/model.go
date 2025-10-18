package model

import (
	"air-quality-notifyer/internal/service/sensor/pollution"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Sensor struct {
	Id          int64
	Aqi         int                   `json:"aqi"`
	Aqi25       int                   `json:"aqi25"`
	Aqi10       int                   `json:"aqi10"`
	Level       pollution.DangerLevel `json:"level"`
	Pm10        float64               `json:"pm10"`
	Pm25        float64               `json:"pm25"`
	Humidity    float64               `json:"humidity"`
	Temperature float64               `json:"temperature"`
	Date        string                `json:"date"`
	UpdatedAt   string                `json:"updated_at"`
	District    string
	SourceLink  string
}

func (s *Sensor) WithDistrict(districtName string) {
	s.District = districtName
}

func (s *Sensor) WithApiData(id int64) {
	s.Id = id
	s.SourceLink = fmt.Sprintf("https://airkemerovo.ru/sensor/%d", id)
}

func (s *Sensor) IsDangerousLevelDetected() bool {
	return !(s.Level == pollution.Good || s.Level == pollution.Moderate || s.Level == "")
}

func (s *Sensor) DangerLevelText() string {
	pollutionData := s.GetPollutionLevel()
	if pollutionData == nil {
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
		s.District, date, pollutionData.Name,
		s.Aqi10, s.Aqi25, s.SourceLink,
	)
}

func (s *Sensor) GetPollutionLevel() *pollution.Level {
	return pollution.GetPollutionLevelByDangerLevel(s.Level)
}
