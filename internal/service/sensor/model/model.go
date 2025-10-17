package model

import (
	"air-quality-notifyer/internal/service/sensor/pollution"
	"fmt"
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

func (s *Sensor) GetPollutionLevel() *pollution.Level {
	return pollution.GetPollutionLevelByDangerLevel(s.Level)
}
