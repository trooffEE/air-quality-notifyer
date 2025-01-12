package sensor

import (
	"fmt"
)

var (
	good               = "good"
	moderate           = "moderate"
	unhealthySensitive = "unhealthy_sensitive"
	unhealthy          = "unhealthy"
	unhealthyModerate  = "very_unhealthy"
	hazardous          = "hazardous"
)

type PollutionLevels struct {
	Good               PollutionLevel
	Moderate           PollutionLevel
	UnhealthySensitive PollutionLevel
	Unhealthy          PollutionLevel
	UnhealthyModerate  PollutionLevel
	Hazardous          PollutionLevel
}

type PollutionLevel struct {
	Name                     string
	AqiDescription           string
	AqiSafetyRecommendations string
}

var PollutionLevelsMap = PollutionLevels{
	Good: PollutionLevel{
		Name:                     "Хорошо",
		AqiDescription:           "Нормальный уровень",
		AqiSafetyRecommendations: "Отличный день для активного отдыха на свежем воздухе",
	},
	Moderate: PollutionLevel{
		Name:                     "Приемлемо",
		AqiDescription:           "Нормальный уровень",
		AqiSafetyRecommendations: "Некоторые люди могут быть чувствительны к загрязнению частицами.\n\n<b>Чувствительные люди</b>: попробуйте уменьшить длительные или тяжелые нагрузки. Следите за такими симптомами, как кашель или одышка. Это признаки того, что нужно снизить нагрузку.\n\n<b>Всем остальным</b>: это хороший день для активности на улице.",
	},
	UnhealthySensitive: PollutionLevel{
		Name:                     "Вредно",
		AqiDescription:           "Повышенный уровень - \"плохо\" ⚠️",
		AqiSafetyRecommendations: "К уязвимым группам относятся люди <b>с заболеваниями сердца или легких, пожилые люди, дети и подростки</b>.\n\n<b>Чувствительные группы</b>: уменьшите длительные или тяжелые нагрузки. Активный образ жизни на улице - это нормально, но делайте больше перерывов и делайте менее интенсивные занятия. Следите за такими симптомами, как кашель или одышка.\n\n<b>Люди, страдающие астмой</b>, должны следовать своим планам действий при астме и иметь под рукой лекарства быстрого действия.\n\n<b>Если у вас заболевание сердца</b>: такие симптомы, как учащенное сердцебиение, одышка или необычная усталость, могут указывать на серьезную проблему. Если у вас есть какие-либо из них, обратитесь к своему врачу.",
	},
	Unhealthy: PollutionLevel{
		Name:                     "Вредно",
		AqiDescription:           "Повышенный уровень - \"вредно\" ⚠️",
		AqiSafetyRecommendations: "<b>Касается всех</b>\n\n<b>Чувствительные группы</b>: Избегайте длительных или тяжелых нагрузок. Подумайте о том, чтобы переместиться в помещение или изменить расписание.\n\n<b>Всем остальным</b>: уменьшите длительные или тяжелые нагрузки. Делайте больше перерывов во время активного отдыха.",
	},
	UnhealthyModerate: PollutionLevel{
		Name:                     "Очень вредно",
		AqiDescription:           "Опасный уровень - \"очень вредно\" 💀",
		AqiSafetyRecommendations: "<b>Касается всех</b>\n\n<b>Чувствительные группы</b>: избегайте любых физических нагрузок на открытом воздухе. Переместите занятия в закрытое помещение или перенесите время, когда качество воздуха будет лучше.\n\n<b>Всем остальным</b>: Избегайте длительных или тяжелых нагрузок. Подумайте о том, чтобы переместиться в помещение или перенести время на то время, когда качество воздуха будет лучше.",
	},
	Hazardous: PollutionLevel{
		Name:                     "Чрезвычайно опасно",
		AqiDescription:           "Опасный уровень - \"чрезвычайно опасно\" 💀💀💀",
		AqiSafetyRecommendations: "<b>Для всех</b>: избегайте любых физических нагрузок на открытом воздухе.\n\n<b>Чувствительные группы</b>: оставайтесь в помещении и сохраняйте низкий уровень активности. Следуйте советам по сохранению низкого уровня частиц в помещении.",
	},
}

func (s *AqiSensor) IsDangerousLevelDetected() bool {
	if s.Level == good || s.Level == moderate {
		return false
	}
	return true
}

func (s *AqiSensor) GetExtendedPollutionLevel() *PollutionLevel {
	switch s.Level {
	case good:
		return &PollutionLevelsMap.Good
	case moderate:
		return &PollutionLevelsMap.Moderate
	case unhealthySensitive:
		return &PollutionLevelsMap.UnhealthySensitive
	case unhealthy:
		return &PollutionLevelsMap.Unhealthy
	case unhealthyModerate:
		return &PollutionLevelsMap.UnhealthyModerate
	case hazardous:
		return &PollutionLevelsMap.Hazardous
	}

	return nil
}

type AqiSensor struct {
	Id          int64
	Aqi         int     `json:"aqi"`
	Aqi25       int     `json:"aqi25"`
	Aqi10       int     `json:"aqi10"`
	Level       string  `json:"level"`
	Pm10        float64 `json:"pm10"`
	Pm25        float64 `json:"pm25"`
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
	Date        string  `json:"date"`
	UpdatedAt   string  `json:"updated_at"`
	District    string
	SourceLink  string
}

type AqiSensorResponse struct {
	Id          int         `json:"id"`
	CityId      int         `json:"city_id"`
	Description interface{} `json:"description"`
	Lat         float64     `json:"lat"`
	Lon         float64     `json:"lon"`
	Address     string      `json:"address"`
	Floor       int         `json:"floor"`
	Radius      int         `json:"radius"`
	Source      interface{} `json:"source"`
	Type        string      `json:"type"`
	Last        AqiSensor   `json:"last"`
	Archive     []AqiSensor `json:"archive"`
}

func (s *AqiSensor) withDistrict(districtName string) {
	s.District = districtName
}

func (s *AqiSensor) withApiData(id int64) {
	s.Id = id
	s.SourceLink = fmt.Sprintf("https://airkemerovo.ru/sensor/%d", id)
}
