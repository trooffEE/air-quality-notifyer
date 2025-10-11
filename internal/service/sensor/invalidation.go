package sensor

import (
	"air-quality-notifyer/internal/db/repository/sensor"
	"air-quality-notifyer/internal/service/sensor/scrapper"
	"fmt"

	"github.com/robfig/cron/v3"
)

var (
	InvalidationPeriod = 4
)

func (s *Service) StartInvalidatingSensorsPeriodically() {
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
	scrappedSensors := scrapper.Scrap()
	aliveSensors := scrapper.FilterSensorsByHourDiff(scrappedSensors, allowedHourDiff)

	for _, sensor := range aliveSensors {
		s.saveSensor(sensor)
	}
}

func (s *Service) saveSensor(scrappedSensor scrapper.Sensor) {
	district := s.sDistricts.GetDistrictByCoords(scrappedSensor.Lat, scrappedSensor.Lon)
	// TODO Не работаем с датчиками вне районов города
	if district == nil {
		return
	}

	payload := sensor.Sensor{
		DistrictId: district.Id,
		ApiId:      scrappedSensor.Id,
		Address:    scrappedSensor.Address,
		Lat:        scrappedSensor.Lat,
		Lon:        scrappedSensor.Lon,
		CreatedAt:  scrappedSensor.CreatedAt,
		District: sensor.DistrictSensor{
			Id:   district.Id,
			Name: district.Name,
		},
	}

	s.saveSensorInCache(payload)
}
