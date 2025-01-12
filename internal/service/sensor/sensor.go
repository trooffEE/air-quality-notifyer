package sensor

import (
	"air-quality-notifyer/internal/db/models"
	repo "air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/service/districts"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
)

type Service struct {
	worstAirqualitySensorsChannel chan []AqiSensor
	districts                     *districts.Service
	repo                          repo.SensorRepositoryType
	ctx                           context.Context
}

func NewSensorService(ctx context.Context, repository repo.SensorRepositoryType, districtService *districts.Service) *Service {
	return &Service{
		repo:                          repository,
		districts:                     districtService,
		worstAirqualitySensorsChannel: make(chan []AqiSensor),
		ctx:                           ctx,
	}
}

func (s *Service) ListenChangesInSensors(handler func([]AqiSensor)) {
	for update := range s.worstAirqualitySensorsChannel {
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
	scrappedSensors := scrapSensorData()

	for _, sensor := range scrappedSensors {
		_, err := s.repo.GetSensorByApiId(sensor.Id)
		if errors.Is(err, sql.ErrNoRows) {
			s.saveNewScrappedSensor(sensor)
		} else if err != nil {
			fmt.Printf("Failed to get api_ids of sensors from database: %v\n", err)
		}
	}

	s.invalidateSensors(scrappedSensors)
}

func (s *Service) saveNewScrappedSensor(sensor AirqualitySensorScriptScrapped) {
	districtId := s.districts.GetDistrictByCoords(sensor.Lat, sensor.Lon)
	// TODO Не работаем с датчиками вне районов города
	if districtId == -1 {
		return
	}

	dbModel := models.AirqualitySensor{
		DistrictId: districtId,
		ApiId:      sensor.Id,
		Address:    sensor.Address,
		Lat:        sensor.Lat,
		Lon:        sensor.Lon,
	}
	err := s.repo.SaveSensor(dbModel)
	if err != nil {
		fmt.Printf("Failed to save new scrapped sensor: %v\n", err)
	}
}

func (s *Service) getWorstAirqualitySensors() {
	ctxDistricts := s.ctx.Value("districts").([]models.District)

	respChan := make(chan AqiSensor, len(ctxDistricts))

	for _, district := range ctxDistricts {
		allSensorsInDistrict := s.repo.GetSensorsByDistrictId(district.Id)
		findWorstSensorInDistrict(respChan, allSensorsInDistrict)
	}

	close(respChan)

	var sensors []AqiSensor
	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	s.worstAirqualitySensorsChannel <- sensors
}
