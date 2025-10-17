package pollution

import "go.uber.org/zap"

type DangerLevel = string

type DangerLevelToPollutionLevel = map[DangerLevel]Level
type Level struct {
	Name                     string
	AqiDescription           string
	AqiSafetyRecommendations string
}

const (
	Good               DangerLevel = "good"
	Moderate           DangerLevel = "moderate"
	UnhealthySensitive DangerLevel = "unhealthy_sensitive"
	Unhealthy          DangerLevel = "unhealthy"
	UnhealthyModerate  DangerLevel = "very_unhealthy"
	Hazardous          DangerLevel = "hazardous"
)

var dangerLevelToPollutionLevel = DangerLevelToPollutionLevel{
	Good: Level{
		Name:                     "Хорошо",
		AqiDescription:           "Нормальный уровень",
		AqiSafetyRecommendations: "Отличный день для активного отдыха на свежем воздухе",
	},
	Moderate: Level{
		Name:                     "Приемлемо",
		AqiDescription:           "Нормальный уровень",
		AqiSafetyRecommendations: "Некоторые люди могут быть чувствительны к загрязнению частицами.\n\n<b>Чувствительные люди</b>: попробуйте уменьшить длительные или тяжелые нагрузки. Следите за такими симптомами, как кашель или одышка. Это признаки того, что нужно снизить нагрузку.\n\n<b>Всем остальным</b>: это хороший день для активности на улице.",
	},
	UnhealthySensitive: Level{
		Name:                     "Вредно",
		AqiDescription:           "Повышенный уровень - \"плохо\" ⚠️",
		AqiSafetyRecommendations: "К уязвимым группам относятся люди <b>с заболеваниями сердца или легких, пожилые люди, дети и подростки</b>.\n\n<b>Чувствительные группы</b>: уменьшите длительные или тяжелые нагрузки. Активный образ жизни на улице - это нормально, но делайте больше перерывов и делайте менее интенсивные занятия. Следите за такими симптомами, как кашель или одышка.\n\n<b>Люди, страдающие астмой</b>, должны следовать своим планам действий при астме и иметь под рукой лекарства быстрого действия.\n\n<b>Если у вас заболевание сердца</b>: такие симптомы, как учащенное сердцебиение, одышка или необычная усталость, могут указывать на серьезную проблему. Если у вас есть какие-либо из них, обратитесь к своему врачу.",
	},
	Unhealthy: Level{
		Name:                     "Вредно",
		AqiDescription:           "Повышенный уровень - \"вредно\" ⚠️⚠️",
		AqiSafetyRecommendations: "<b>Касается всех</b>\n\n<b>Чувствительные группы</b>: Избегайте длительных или тяжелых нагрузок. Подумайте о том, чтобы переместиться в помещение или изменить расписание.\n\n<b>Всем остальным</b>: уменьшите длительные или тяжелые нагрузки. Делайте больше перерывов во время активного отдыха.",
	},
	UnhealthyModerate: Level{
		Name:                     "Очень вредно",
		AqiDescription:           "Опасный уровень - \"очень вредно\" ⚠️⚠️⚠️",
		AqiSafetyRecommendations: "<b>Касается всех</b>\n\n<b>Чувствительные группы</b>: избегайте любых физических нагрузок на открытом воздухе. Переместите занятия в закрытое помещение или перенесите время, когда качество воздуха будет лучше.\n\n<b>Всем остальным</b>: Избегайте длительных или тяжелых нагрузок. Подумайте о том, чтобы переместиться в помещение или перенести время на то время, когда качество воздуха будет лучше.",
	},
	Hazardous: Level{
		Name:                     "Чрезвычайно опасно",
		AqiDescription:           "Опасный уровень - \"чрезвычайно опасно\" 💀",
		AqiSafetyRecommendations: "<b>Для всех</b>: избегайте любых физических нагрузок на открытом воздухе.\n\n<b>Чувствительные группы</b>: оставайтесь в помещении и сохраняйте низкий уровень активности. Следуйте советам по сохранению низкого уровня частиц в помещении.",
	},
}

func GetPollutionLevelByDangerLevel(dangerLevel DangerLevel) *Level {
	level, ok := dangerLevelToPollutionLevel[dangerLevel]
	if !ok {
		zap.L().Warn("")
		return nil
	}
	return &level
}
