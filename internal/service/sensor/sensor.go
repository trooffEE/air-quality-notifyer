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
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	AliveSensorTimeDiff = 4
)

type Service struct {
	trustedSensorAqiChannel chan []AqiSensor
	districts               *districts.Service
	repo                    repo.SensorRepositoryType
	ctx                     context.Context
	syncCron                chan interface{}
}

func NewSensorService(ctx context.Context, repository repo.SensorRepositoryType, districtService *districts.Service) *Service {
	return &Service{
		repo:                    repository,
		districts:               districtService,
		trustedSensorAqiChannel: make(chan []AqiSensor),
		ctx:                     ctx,
		syncCron:                make(chan interface{}),
	}
}

func (s *Service) ListenChangesInSensors(handler func([]AqiSensor)) {
	for update := range s.trustedSensorAqiChannel {
		handler(update)
	}
}

func (s *Service) InvalidateSensorsPeriodically() {
	cronCreator := cron.New()
	cronString := fmt.Sprintf("0 */%d * * *", AliveSensorTimeDiff)

	_, err := cronCreator.AddFunc(cronString, func() {
		s.startInvalidation(AliveSensorTimeDiff)
		s.syncCron <- 0
	})
	if err != nil {
		panic(err)
	}

	cronCreator.Start()
}

func (s *Service) GetTrustedSensorsEveryHour() {
	cronCreator := cron.New()
	cronString := "0 * * * *"

	_, err := cronCreator.AddFunc(cronString, func() {
		if time.Now().UTC().Hour()%AliveSensorTimeDiff == 0 {
			<-s.syncCron
		}
		s.getTrustedAirqualitySensors()
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
		_, err := s.repo.GetSensorByApiId(sensor.Id)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.saveNewScrappedSensor(sensor)
				continue
			}
			zap.L().Error("failed to get api_ids of sensors from database", zap.Error(err))
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
		CreatedAt:  sensor.CreatedAt,
	}
	err := s.repo.SaveSensor(dbModel)
	if err != nil {
		zap.L().Error("failed to save sensor", zap.Error(err), zap.Any("dbModel", dbModel))
	}
}

func (s *Service) getTrustedAirqualitySensors() {
	ctxDistricts := s.ctx.Value("districts").([]models.District)

	respChan := make(chan AqiSensor, len(ctxDistricts))
	wg := sync.WaitGroup{}
	wg.Add(len(ctxDistricts))
	for _, district := range ctxDistricts {
		sensorsInDistrict, err := s.repo.GetSensorsByDistrictId(district.Id)
		if err != nil {
			zap.L().Error("failed to get sensors by districtId", zap.Error(err), zap.Int64("districtId", district.Id))
			continue
		}
		go func() {
			defer wg.Done()
			findTrustedSensor(respChan, sensorsInDistrict)
		}()
	}
	wg.Wait()
	close(respChan)

	var sensors []AqiSensor
	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	s.trustedSensorAqiChannel <- sensors
}
