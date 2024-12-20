package sensor

import (
	repo "air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/districts"
	districts2 "air-quality-notifyer/internal/service/districts"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
)

type Service struct {
	worstAirqualitySensorsChannel chan []AirqualitySensor
	districts                     *districts2.Service
	repo                          repo.SensorRepositoryType
}

func NewSensorService(repository repo.SensorRepositoryType, districtService *districts2.Service) *Service {
	return &Service{
		repo:                          repository,
		districts:                     districtService,
		worstAirqualitySensorsChannel: make(chan []AirqualitySensor),
	}
}

func (s *Service) ListenChangesInSensors(handler func([]AirqualitySensor)) {
	for update := range s.worstAirqualitySensorsChannel {
		fmt.Println(update)
		handler(update)
	}
}

func (s *Service) FetchSensorsEveryHour() {
	cronCreator := cron.New()
	cronString := "0 * * * *"

	_, err := cronCreator.AddFunc(cronString, s.getWorstAirqualitySensors)
	if err != nil {
		log.Panic(err)
	}

	cronCreator.Start()
}

func (s *Service) InvalidateSensorsEveryday() {
	cronCreator := cron.New()
	cronString := "0 0 * * *"

	_, err := cronCreator.AddFunc(cronString, s.startInvalidation)
	if err != nil {
		log.Panic(err)
	}

	cronCreator.Start()
}

func (s *Service) startInvalidation() {
	aliveSensors := scrapSensorData()

	for _, sensor := range aliveSensors {
		id := s.districts.GetDistrictByCoords(sensor.Lat, sensor.Lon)
		s.invalidateSensor(sensor, id)
	}
}

func (s *Service) getWorstAirqualitySensors() {
	respChan := make(chan AirqualitySensor, len(districts.Dictionary))

	for _, district := range districts.Dictionary {
		findWorstSensorInDistrict(respChan, district)
	}

	close(respChan)

	var sensors []AirqualitySensor
	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	s.worstAirqualitySensorsChannel <- sensors
}
