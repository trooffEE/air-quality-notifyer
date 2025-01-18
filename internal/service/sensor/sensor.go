package sensor

import (
	"air-quality-notifyer/internal/db/models"
	repo "air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/lib"
	"air-quality-notifyer/internal/service/districts"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

var (
	AliveSensorTimeDiff = 4
)

type Service struct {
	worstAirqualitySensorsChannel chan []AqiSensor
	districts                     *districts.Service
	repo                          repo.SensorRepositoryType
	ctx                           context.Context
	syncCron                      chan interface{}
}

func NewSensorService(ctx context.Context, repository repo.SensorRepositoryType, districtService *districts.Service) *Service {
	return &Service{
		repo:                          repository,
		districts:                     districtService,
		worstAirqualitySensorsChannel: make(chan []AqiSensor),
		ctx:                           ctx,
		syncCron:                      make(chan interface{}),
	}
}

func (s *Service) ListenChangesInSensors(handler func([]AqiSensor)) {
	for update := range s.worstAirqualitySensorsChannel {
		handler(update)
	}
}

func (s *Service) InvalidateSensorsPeriodically() {
	cronCreator := cron.New()
	cronString := fmt.Sprintf(fmt.Sprintf("0 */%d * * *", AliveSensorTimeDiff))

	_, err := cronCreator.AddFunc(cronString, func() {
		s.startInvalidation(AliveSensorTimeDiff)
		s.syncCron <- 0
	})
	if err != nil {
		log.Panic(err)
	}

	cronCreator.Start()
}

func (s *Service) FetchSensorsEveryHour() {
	cronCreator := cron.New()
	cronString := "* * * * *"

	_, err := cronCreator.AddFunc(cronString, func() {
		if time.Now().UTC().Hour()%AliveSensorTimeDiff == 0 {
			<-s.syncCron
		}
		s.getWorstAirqualitySensors()
	})
	if err != nil {
		log.Panic(err)
	}

	cronCreator.Start()
}

func (s *Service) startInvalidation(allowedHourDiff int) {
	scrappedSensors := scrapSensorData()

	aliveSensors := filterDeadSensors(scrappedSensors, allowedHourDiff)

	for _, sensor := range aliveSensors {
		_, err := s.repo.GetSensorByApiId(sensor.Id)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.saveNewScrappedSensor(sensor)
				continue
			}
			lib.LogError("startInvalidation", "failed to get api_ids of sensors from database", err)
		}
	}

	s.invalidateSensors(aliveSensors)
}

func (s *Service) saveNewScrappedSensor(sensor AqiSensorScriptScrapped) {
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
		lib.LogError("saveNewScrappedSensor", "failed to save sensor %+v", err, dbModel)
	}
}

func (s *Service) getWorstAirqualitySensors() {
	ctxDistricts := s.ctx.Value("districts").([]models.District)

	respChan := make(chan AqiSensor, len(ctxDistricts))

	for _, district := range ctxDistricts {
		allSensorsInDistrict, err := s.repo.GetSensorsByDistrictId(district.Id)
		if err != nil {
			lib.LogError("getWorstAirqualitySensors", "failed to get sensors by districtId=%d", err, district.Id)
			continue
		}

		findWorstSensorInDistrict(respChan, allSensorsInDistrict)
	}

	close(respChan)

	var sensors []AqiSensor
	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	s.worstAirqualitySensorsChannel <- sensors
}
