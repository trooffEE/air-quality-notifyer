package sensor

import (
	"air-quality-notifyer/internal/db/models"
	"fmt"

	"github.com/robfig/cron/v3"
)

var (
	InvalidationPeriod = 4
)

func (s *Service) InvalidateSensorsPeriodically() {
	cronCreator := cron.New()
	cronString := fmt.Sprintf("0 */%d * * *", InvalidationPeriod)

	_, err := cronCreator.AddFunc(cronString, func() {
		s.startInvalidation(InvalidationPeriod)
		s.syncCron <- 0
	})
	if err != nil {
		panic(err)
	}

	cronCreator.Start()
}

func (s *Service) startInvalidation(allowedHourDiff int) {
	scrappedSensors := scrapSensorData()
	aliveSensors := filterDeadSensors(scrappedSensors, allowedHourDiff)

	for _, sensor := range aliveSensors {
		s.saveSensor(sensor)
	}
}

func (s *Service) saveSensor(sensor scriptTagScrappedSensor) {
	district := s.sDistricts.GetDistrictByCoords(sensor.Lat, sensor.Lon)
	// TODO Не работаем с датчиками вне районов города
	if district == nil {
		return
	}

	payload := models.AirqualitySensor{
		DistrictId: district.Id,
		ApiId:      sensor.Id,
		Address:    sensor.Address,
		Lat:        sensor.Lat,
		Lon:        sensor.Lon,
		CreatedAt:  sensor.CreatedAt,
		District: models.DistrictSensor{
			Id:   district.Id,
			Name: district.Name,
		},
	}

	s.saveSensorInCache(payload)
}
