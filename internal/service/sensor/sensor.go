package sensor

import (
	"air-quality-notifyer/internal/db/repository/sensor"
	"air-quality-notifyer/internal/service/districts"
	"air-quality-notifyer/internal/service/sensor/model"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo       sensor.Interface
	sDistricts districts.Interface
	cSensors   chan []model.Sensor
	syncCron   chan interface{}
	cache      *redis.Client
}

type Interface interface {
	ListenChangesInSensors(handler func([]model.Sensor))
	StartInvalidatingSensorsPeriodically()
	StartGettingTrustedSensorsEveryHour()
}

func New(
	repo sensor.Interface,
	sDistricts districts.Interface,
	cache *redis.Client,
) Interface {
	return &Service{
		repo:       repo,
		sDistricts: sDistricts,
		cSensors:   make(chan []model.Sensor),
		syncCron:   make(chan interface{}),
		cache:      cache,
	}
}

func (s *Service) ListenChangesInSensors(handler func([]model.Sensor)) {
	for update := range s.cSensors {
		handler(update)
	}
}
