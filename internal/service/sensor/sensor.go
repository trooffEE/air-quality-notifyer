package sensor

import (
	"air-quality-notifyer/internal/db/models"
	repo "air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/service/districts"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	AliveSensorTimeDiff = 4
)

type Service struct {
	repo       repo.SensorRepositoryInterface
	sDistricts *districts.Service
	cSensors   chan []AqiSensor
	syncCron   chan interface{}
	cache      *redis.Client
}

func NewSensorService(
	repo repo.SensorRepositoryInterface,
	sDistricts *districts.Service,
	cache *redis.Client,
) *Service {
	return &Service{
		repo:       repo,
		sDistricts: sDistricts,
		cSensors:   make(chan []AqiSensor),
		syncCron:   make(chan interface{}),
		cache:      cache,
	}
}

func (s *Service) ListenChangesInSensors(handler func([]AqiSensor)) {
	for update := range s.cSensors {
		handler(update)
	}
}

func (s *Service) InvalidateSensorsPeriodically() {
	cronCreator := cron.New()
	cronString := fmt.Sprintf("0 */%d * * *", AliveSensorTimeDiff)

	//s.startInvalidation(AliveSensorTimeDiff)
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
		s.saveSensor(sensor)
	}
}

func (s *Service) saveSensor(sensor AqiSensorScriptScrapped) {
	districtId := s.sDistricts.GetDistrictByCoords(sensor.Lat, sensor.Lon)
	// TODO Не работаем с датчиками вне районов города
	if districtId == -1 {
		return
	}

	payload := models.AirqualitySensor{
		DistrictId: districtId,
		ApiId:      sensor.Id,
		Address:    sensor.Address,
		Lat:        sensor.Lat,
		Lon:        sensor.Lon,
		CreatedAt:  sensor.CreatedAt,
	}

	s.saveSensorInCache(payload)
}

func (s *Service) getTrustedAirqualitySensors() {
	allDistricts := s.sDistricts.GetAllDistricts() // think about it

	respChan := make(chan AqiSensor, len(allDistricts))
	wg := sync.WaitGroup{}
	wg.Add(len(allDistricts))
	for _, district := range allDistricts {
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

	s.cSensors <- sensors
}
